[![Build Status](https://travis-ci.org/deadc0de6/mftail.svg?branch=master)](https://travis-ci.org/deadc0de6/mftail)
[![License: GPL v3](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](http://www.gnu.org/licenses/gpl-3.0)

# mftail

[mftail](https://github.com/deadc0de6/mftail) is a multiple files follower (similar to `tail -f`) with color.

# Usage

usage
```bash
mftail v0.1

Usage: mftail <path>...
  -v	Print version.
  -version
    	Print version.
```

# Install

Quick start:
```bash
## You need at least golang 1.14
$ go install -v github.com/deadc0de6/mftail@latest
```

Or pick a release from [the release page](https://github.com/deadc0de6/mftail/releases) and install it in your `$PATH`

Or compile it from source
```bash
$ go mod tidy
$ make
$ mftail -help
```

# Contribution

If you are having trouble using mftail, open an issue.

If you want to contribute, feel free to do a PR.

# License

This project is licensed under the terms of the GPLv3 license.
