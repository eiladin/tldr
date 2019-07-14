#!/bin/sh
rm -rf out
mkdir out
go build -o out/linux/tldr
tar -czf out/tldr.linux.x64.tgz out/linux/tldr

export GOOS="darwin"
export GOARCH="amd64"

go build -o out/macos/tldr
tar -czf out/tldr.macos.x64.tgz out/macos/tldr

export GOOS="windows"
export GOARCH="amd64"

go build -o out/windows/tldr.exe
zip -r out/tldr.win.x64.zip out/windows/tldr.exe
