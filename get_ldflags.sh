#!/usr/bin/env bash
#
# Usage: ./get_ldflags.sh [BUILD_FLAGS]
#
# This file has been created by: deploy.sh init -- 2020-04-04 08:27:32
#
# This script is an helper for both Makefile + deploy.sh
#
# You can reuse it in your Makefile:
# BUILD_FLAGS=$(shell ./get_ldflags.sh)
# your_target: your_target.go Makefile ${OTHER_DEP}
# 	go build -o $@ -ldflags "${BUILD_FLAGS} ${LDFLAGS}"

set -eu
build_flags=${1:-}
if [[ -z $build_flags ]]; then
  # govvv define main.Version with the contents of ./VERSION file, if exists
  build_flags=$(govvv -flags)
fi

# you can add more flags here:
build_flags+=" -X 'main.GoBuildVersion=$(go version)' -X 'main.ByUser=${USER}'"

# the last command MUST display all build_flags
echo "$build_flags"
