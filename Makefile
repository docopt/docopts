#
# Makefile for building docopts
#

# local path for dependent golang lib (not used)
DOCOPT_GO=${GOPATH}/linux_amd64/github.com/docopt/docopt-go.a

PREFIX ?= /usr/local
VERSION=0.6.3
LOCAL_DIR=./local

# build 64 bits version (default target)
docopts: docopts.go
	go build docopts.go

# simply get the library
docopt-go:
	go get github.com/docopt/docopt-go

all: docopts docopts-arm docopts-32bits docopts-OSX sha256sum.txt

# compile 32 bits version too
docopts-32bits: docopts.go
	env GOOS=linux GOARCH=386 go build -o docopts-32bits docopts.go

# compile for OSX (not tested)
docopts-OSX: docopts.go
	env GOOS=darwin go build -o docopts-OSX docopts.go

# compile version for arm architecture (not suitable for android)
docopts-arm: docopts.go
	env GOOS=linux GOARCH=arm go build -o docopts-arm docopts.go

test: docopts
	go test -v
	python language_agnostic_tester.py ./testee.sh
	cd tests/ && ./bats/bin/bats .

# for building debian package
build: all build_orig_tgz

build_orig_tgz: docopts docopts-32bits docopts.sh
	tar czf ./docopts_${VERSION}.tar.gz docopts docopts-32bits docopts.sh

sha256sum.txt: docopts docopts-32bits docopts.sh
	sha256sum docopts docopts-32bits docopts.sh > sha256sum.txt

# some fake local dir for testing install
${LOCAL_DIR}:
	mkdir -p ${LOCAL_DIR}
	mkdir -p ${LOCAL_DIR}/bin
	mkdir -p ${LOCAL_DIR}/lib/docopts

# https://www.gnu.org/savannah-checkouts/gnu/make/manual/html_node/Prerequisite-Types.html#Prerequisite-Types
# create local dir for test install
create_local: | ${LOCAL_DIR}

install: docopts
	strip docopts
	install -m 755 docopts $(PREFIX)/bin
	install -m 644 docopts.sh $(PREFIX)/lib/docopts

clean:
	rm -f docopts-* docopts
	rm -rf ${LOCAL_DIR}
