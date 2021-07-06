#!/bin/bash

set -e

go fmt *.go
golint -set_exit_status *.go
go vet *.go
