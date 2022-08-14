#!/usr/bin/env bash

git clone https://github.com/bats-core/bats-core.git --depth=1
cd bats-core
sudo ./install.sh /usr/local

pwd
type bats
bats --version
