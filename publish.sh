#!/bin/bash
while getopts v:m: option
do 
case "${option}"
in 
v) VERSION=${OPTARG};;
m) MESSAGE=${OPTARG};;
esac
done

hub release create -a ./out/tldr.linux.x64.tgz -a ./out/tldr.macos.x64.tgz -a ./out/tldr/win.x64.zip -m \"$MESSAGE\" $VERSION