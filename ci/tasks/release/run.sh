#!/bin/bash
set -xe

export GOPATH=$PWD/gopath

echo "Building 64-bit Darwin"
GOARCH=amd64 GOOS=darwin go build -o build/yml2env_osx github.com/EngineerBetter/yml2env
echo "Building 32-bit Linux"
GOARCH=386 GOOS=linux go build -ldflags '-extldflags "-static"' -o build/yml2env_linux_i686 github.com/EngineerBetter/yml2env
echo "Building 64-bit Linux"
GOARCH=amd64 GOOS=linux go build -ldflags '-extldflags "-static"' -o build/yml2env_linux_x86-64 github.com/EngineerBetter/yml2env
echo "Building 32-bit Windows"
GOARCH=386 GOOS=windows go build -o build/yml2env_win32.exe github.com/EngineerBetter/yml2env
echo "Building 64-bit Windows"
GOARCH=amd64 GOOS=windows go build -o build/yml2env_winx64.exe github.com/EngineerBetter/yml2env