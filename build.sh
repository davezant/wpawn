#!/bin/bash
temp=$(realpath "$0"| sed 's|\(.*\)/.*|\1|')

go build main.go
cp $temp/main $temp/wpawn
