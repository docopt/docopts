#
# Makefile for managing docopts build
#
# See also: deploy.sh

PREFIX ?= /usr/local

# keep docopts: as first target for development

# govvv define main.Version with the contents of ./VERSION file, if exists
BUILD_FLAGS=$(shell ./get_ldflags.sh)
docopts: docopts.go Makefile docopt-go/docopt.go
	go build -o $@ -ldflags "${BUILD_FLAGS} ${LDFLAGS}"

# dependancies
install_builddep:
	go install github.com/docopt/docopts
	go install github.com/docopt/docopt-go
	go install github.com/mitchellh/gox
	go install github.com/itchio/gothub
	go install gopkg.in/yaml.v2
	go install github.com/ahmetb/govvv

parser: parser/docoptnodetype_string.go parser/docopt_language.go parser/docopt_ast.go
	go build cmd/docopt-analyze/main.go

parser/docoptnodetype_string.go: parser/docopt_ast.go
	go generate parser/docopt_ast.go

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
	python3 language_agnostic_tester.py ./testee.sh
	cd ./tests/ && bats .

# README.md is composed with external source too
# Markdown hidden markup are used to insert some text form the dependancies
README.md: examples/legacy_bash/rock_hello_world.sh examples/legacy_bash/rock_hello_world_with_grep.sh docopts build_doc.sh
	./build_doc.sh README.md > README.tmp
	mv README.tmp README.md

clean:
	rm -f docopts-* docopts README.tmp build/*
