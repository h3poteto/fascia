#!/bin/sh

gom run main.go
ret=$?

if [ $ret -ne 0 ]; then
    gom install
fi
