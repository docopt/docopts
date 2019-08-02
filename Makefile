#
# Makefile for managing docopts build
#
# See also: deploy.sh


PREFIX ?= /usr/local

# keep this as first target for development
# build 64 bits version
docopts: docopts.go
	go build docopts.go

docopt-go:
	go get github.com/docopt/docopt-go

all: docopt-go docopts docopts-arm docopts-32bits docopts-OSX README.md

# build 32 bits version too
docopts-32bits: docopts.go
	env GOOS=linux GOARCH=386 go build -o docopts-32bits docopts.go

# build for OSX
docopts-OSX: docopts.go
	env GOOS=darwin go build -o docopts-OSX docopts.go

# build 32 bits version too
docopts-arm: docopts.go
	env GOOS=linux GOARCH=arm go build -o docopts-arm docopts.go

# requires write access to $PREFIX
install: all
	cp docopts docopts.sh $(PREFIX)/bin

test: docopts
	go test -v
	python language_agnostic_tester.py ./testee.sh
	cd tests/ && ./bats/bin/bats .

# README.md is composed with external source too
# Markdown hidden markup are used to insert some text form the dependancies
README.md: examples/legacy_bash/rock_hello_world.sh examples/legacy_bash/rock_hello_world_with_grep.sh docopts build_doc.sh
	./build_doc.sh README.md > README.tmp
	mv README.tmp README.md

clean:
	rm -f docopts-* docopts README.tmp build/docopts_*
