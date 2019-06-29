#!/bin/sh
rm -rf out
mkdir out
go build -o out/tldr
tar -czf out/tldr.tgz out/tldr