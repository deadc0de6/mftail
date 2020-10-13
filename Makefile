SRC = mftail.go
BIN = mftail
OS = linux

all: build

build:
	go build -o $(BIN) $(SRC)

build-all:
	GOOS=$(OS) GOARCH=arm 	go build -v -o $(BIN)-$(OS)-arm $(SRC)
	GOOS=$(OS) GOARCH=arm64 go build -v -o $(BIN)-$(OS)-arm64 $(SRC)
	GOOS=$(OS) GOARCH=386 	go build -v -o $(BIN)-$(OS)-386 $(SRC)
	GOOS=$(OS) GOARCH=amd64 go build -v -o $(BIN)-$(OS)-amd64 $(SRC)

clean:
	rm -f $(BIN) \
		$(BIN)-$(OS)-arm \
		$(BIN)-$(OS)-arm64 \
		$(BIN)-$(OS)-386 \
		$(BIN)-$(OS)-amd64
