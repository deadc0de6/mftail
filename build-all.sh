#!/bin/bash

name="mftail"
goos="linux"
for goarch in arm arm64 386 amd64; do
  echo "building ${goos}-${goarch}"
  env GOOS="${goos}" GOARCH="${goarch}" go build -v -o ${name}-${goos}-${goarch}
done

