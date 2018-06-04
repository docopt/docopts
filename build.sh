#!/bin/bash
#
# Build docopts go version
#

# exit on error
set -e

go get github.com/docopt/docopt-go

# build 64 bits version
go build docopts.go

# build 32 bits version too
env GOOS=linux GOARCH=386 go build -o docopts-32bits docopts.go
