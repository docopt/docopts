#
# Makefile for managing docopts build
#

PREFIX ?= /usr/local
GOVERSION := $$(go version)

# keep docopts: as first target for development

docopts: docopts.go Makefile
	GOVERSION=$(GOVERSION) goreleaser build --rm-dist --snapshot --single-target -o $@

install_builddep:
	go mod tidy

all: install_builddep docopts README.md

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
	rm -f docopts-* docopts README.tmp dist/*

release: clean all test
	GOVERSION=$(GOVERSION) goreleaser release
