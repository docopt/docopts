#
# Makefile for managing docopts build
#
# See also: deploy.sh

PREFIX ?= /usr/local

# dependancies
GOVVV=${GOPATH}/bin/govvv
DOCTOP_LIB=${GOPATH}/src/github.com/docopt/docopt-go/docopt.go

# keep docopts: as first target for development

# govvv define main.Version with the contents of ./VERSION file, if exists
BUILD_FLAGS=$(shell ./get_ldflags.sh)
docopts: docopts.go Makefile ${GOVVV} ${DOCTOP_LIB}
	go build -o $@ -ldflags "${BUILD_FLAGS} ${LDFLAGS}"

install_builddep: ${GOVVV} ${DOCTOP_LIB}
	go get github.com/mitchellh/gox
	go get github.com/itchio/gothub
	go get gopkg.in/mikefarah/yq.v2

${DOCTOP_LIB}:
	go get github.com/docopt/docopt-go

${GOVVV}:
	go get github.com/ahmetb/govvv

all: install_builddep docopts README.md
	./deploy.sh build current

############################ cross compile, we use gox now inside deploy.sh

## build 32 bits version too
#docopts-32bits: docopts.go
#	env GOOS=linux GOARCH=386 go build -o docopts-32bits docopts.go
#
## build for OSX
#docopts-OSX: docopts.go
#	env GOOS=darwin go build -o docopts-OSX docopts.go
#
## build 32 bits version too
#docopts-arm: docopts.go
#	env GOOS=linux GOARCH=arm go build -o docopts-arm docopts.go

###########################

# requires write access to $PREFIX
install: all
	install -m 755 docopts    $(PREFIX)/bin
	install -m 755 docopts.sh $(PREFIX)/bin

test: docopts
	./docopts --version
	go test -v
	python language_agnostic_tester.py ./testee.sh
	cd tests/ && bats .

# README.md is composed with external source too
# Markdown hidden markup are used to insert some text form the dependancies
README.md: examples/legacy_bash/rock_hello_world.sh examples/legacy_bash/rock_hello_world_with_grep.sh docopts build_doc.sh
	./build_doc.sh README.md > README.tmp
	mv README.tmp README.md

clean:
	rm -f docopts-* docopts README.tmp build/*
