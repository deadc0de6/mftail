SRC = mftail.go
BIN = mftail
OS = linux

all: build

build:
	GO111MODULE=on go build -o $(BIN) $(SRC)

build-all:
	GO111MODULE=on GOOS=$(OS) GOARCH=arm go build -v -o $(BIN)-$(OS)-arm $(SRC)
	GO111MODULE=on GOOS=$(OS) GOARCH=arm64 go build -v -o $(BIN)-$(OS)-arm64 $(SRC)
	GO111MODULE=on GOOS=$(OS) GOARCH=386 go build -v -o $(BIN)-$(OS)-386 $(SRC)
	GO111MODULE=on GOOS=$(OS) GOARCH=amd64 go build -v -o $(BIN)-$(OS)-amd64 $(SRC)

clean:
	rm -f $(BIN) \
		$(BIN)-$(OS)-arm \
		$(BIN)-$(OS)-arm64 \
		$(BIN)-$(OS)-386 \
		$(BIN)-$(OS)-amd64
