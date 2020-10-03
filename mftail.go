/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2020, deadc0de6
*/

// a multiple file follower tool similar to `tail -f`
package main

import (
	"bufio"
	"flag"
	"fmt"
	"golang.org/x/sys/unix"
	"io"
	"math/rand"
	"os"
	"path"
	"time"
	"unsafe"
)

const (
	version = "0.2"
	reset   = "\033[0m"
	color   = "\033[%dm"
)

var (
	// FileMatcher match fd to channel
	FileMatcher map[int]chan int // wd -> channel
	// Colors all colors
	Colors = []int{31, 32, 33, 34, 35, 36, 37, 90, 91, 92, 93, 94, 95, 96}
	// EventCounts total event to watch for
	EventCounts int = 128
)

const (
	evModify = iota
	evDelete
	evAttrib
)

type fileEvent struct {
	Fd    int // inotify file descriptor
	Wd    int // watch descriptor
	Path  string
	Color string
	Chan  chan int
}

func init() {
	rand.Seed(time.Now().UnixNano())
	FileMatcher = make(map[int]chan int)
}

// check file exists and is not dir
func isFile(path string) bool {
	s, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !s.IsDir()
}

// watch a file for change
func addWatch(fd int, path string) (int, error) {
	var flag uint32 = unix.IN_MODIFY | unix.IN_ATTRIB | unix.IN_DELETE_SELF
	wd, err := unix.InotifyAddWatch(fd, path, flag)
	if err != nil {
		return 0, err
	}
	return wd, nil
}

// unwatch a file
func rmWatch(fd int, wd int) {
	unix.InotifyRmWatch(fd, uint32(wd))
}

// follow file
func follow(fe fileEvent) {
	f := fileopen(nil, fe)
	if f == nil {
		return
	}
	bname := path.Base(fe.Path)

	// go to end of file
	prevSize, _ := f.Seek(0, os.SEEK_END)

	r := bufio.NewReader(f)
	for {
		// wait on channel for trigger
		b, ok := <-fe.Chan
		if !ok {
			// something's wrong
			rmWatch(fe.Fd, fe.Wd)
			return
		}

		switch b {
		case evModify:
			fileInfo, _ := f.Stat()
			curSize := fileInfo.Size()
			if curSize < prevSize {
				// truncated
				fmt.Printf("mftail: %s: file truncated\n", fe.Path)
				f = fileopen(f, fe)
				if f == nil {
					return
				}
				r = bufio.NewReader(f)
			}
			prevSize = curSize
		//case evAttrib:
		case evDelete:
			f = fileopen(f, fe)
			if f == nil {
				return
			}
			r = bufio.NewReader(f)
		}

		// tail last lines
		err := freadlines(r, bname, fe.Color)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
	}
}

// open/re-open file
func fileopen(cur *os.File, fe fileEvent) *os.File {
	if !isFile(fe.Path) {
		// file inexistant or removed
		rmWatch(fe.Fd, fe.Wd)
		return nil
	}

	// open/re-open
	if cur != nil {
		cur.Close()
	}
	f, err := os.Open(fe.Path)
	if err != nil {
		rmWatch(fe.Fd, fe.Wd)
		fmt.Fprintln(os.Stderr, err)
		return nil
	}
	return f
}

// read all lines from current pos and print to output
func freadlines(r *bufio.Reader, header string, color string) error {
	// read lines
	for {
		line, err := r.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return err
		}
		// print the line
		if len(line) > 0 {
			fmt.Printf("%s:%s%s%s", header, color, string(line), reset)
		}
		if err == io.EOF {
			break
		}
	}
	return nil
}

// read the events
func readEvents(buf []byte, bytesRead int) {
	var offset int
	for offset <= bytesRead {
		ev := *(*unix.InotifyEvent)(unsafe.Pointer(&buf[offset]))
		m := ev.Mask
		ch := FileMatcher[int(ev.Wd)]
		if m&unix.IN_ATTRIB == unix.IN_ATTRIB {
			// unlink (inode count changes)
			ch <- evAttrib
		}
		if m&unix.IN_DELETE_SELF == unix.IN_DELETE_SELF {
			// file removed
			ch <- evDelete
		}
		if m&unix.IN_MODIFY == unix.IN_MODIFY {
			// modified
			ch <- evModify
		}
		offset += unix.SizeofInotifyEvent
	}
}

// monitor for events
func waitForNotif(fd int) error {
	// buffer to store events
	// field "Name" is not present/filled when monitoring files
	// see https://linux.die.net/man/7/inotify
	size := unix.SizeofInotifyEvent * EventCounts
	var buf = make([]byte, size)
	for {
		// read the event
		bytesRead, err := unix.Read(fd, buf[:])
		if err != nil || bytesRead < unix.SizeofInotifyEvent {
			continue
		}
		readEvents(buf[:], bytesRead)
	}
}

// print version
func printVersion() {
	name := path.Base(os.Args[0])
	fmt.Fprintf(os.Stdout, "%s v%s\n", name, version)
}

// print usage and exit
func usage() {
	printVersion()
	name := path.Base(os.Args[0])
	fmt.Fprintf(os.Stdout, "\nUsage: %s <path>...\n", name)
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	// parse args
	vshortArg := flag.Bool("v", false, "Print version.")
	vlongArg := flag.Bool("version", false, "Print version.")
	flag.Parse()

	if *vshortArg || *vlongArg {
		printVersion()
		os.Exit(0)
	}

	if len(flag.Args()) < 1 {
		usage()
	}

	// randomize colors
	rand.Shuffle(len(Colors), func(i, j int) { Colors[i], Colors[j] = Colors[j], Colors[i] })

	// init inotify
	fd, err := unix.InotifyInit()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer unix.Close(fd)

	// parse file paths to monitor
	var cnt int
	for idx, path := range flag.Args() {
		if !isFile(path) {
			fmt.Fprintf(os.Stderr, "\"%s\" not found!\n", path)
			continue
		}

		// add file to monitor
		wd, err := addWatch(fd, path)
		if err != nil {
			fmt.Fprintln(os.Stderr, fmt.Sprintf("\"%s\": %s", path, err.Error()))
			continue
		}
		defer rmWatch(fd, wd)

		// create a new channel for this file
		ch := make(chan int)
		FileMatcher[wd] = ch

		col := Colors[idx%len(Colors)]
		fevent := fileEvent{
			Fd:    fd,
			Wd:    wd,
			Path:  path,
			Color: fmt.Sprintf(color, col),
			Chan:  ch,
		}

		// follow the file
		go follow(fevent)
		cnt++
	}

	if cnt < 1 {
		// no file to follow
		fmt.Fprintf(os.Stderr, "no file to follow\n")
		os.Exit(1)
	}

	err = waitForNotif(fd)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
