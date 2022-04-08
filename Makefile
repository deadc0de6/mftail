SRC = mftail.go
BIN = mftail
OS = linux

all: build

build:
	CGO_ENABLED=0 GO111MODULE=on go build -o $(BIN) $(SRC)

build-all:
	CGO_ENABLED=0 GO111MODULE=on GOOS=$(OS) GOARCH=arm go build -v -o $(BIN)-$(OS)-arm $(SRC)
	CGO_ENABLED=0 GO111MODULE=on GOOS=$(OS) GOARCH=arm64 go build -v -o $(BIN)-$(OS)-arm64 $(SRC)
	CGO_ENABLED=0 GO111MODULE=on GOOS=$(OS) GOARCH=386 go build -v -o $(BIN)-$(OS)-386 $(SRC)
	CGO_ENABLED=0 GO111MODULE=on GOOS=$(OS) GOARCH=amd64 go build -v -o $(BIN)-$(OS)-amd64 $(SRC)

clean:
	rm -f $(BIN) $(BIN)-*
