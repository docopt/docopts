#!/usr/bin/env bash

git clone https://github.com/bats-core/bats-core.git
cd bats-core
./install.sh /usr/local

pwd
type bats
bats --version
