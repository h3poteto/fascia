#!/bin/sh

set -x
go run main.go
ret=$?

if [ $ret -ne 0 ]; then
    dep ensure
fi
