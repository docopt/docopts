#!/usr/bin/env bash
#
# functional test for get_doctops.sh download helper
#

_run_get_docopts() {
  # ensure we are still in tests sub-folder
  [[ $(basename $PWD) == 'tests' ]]
  [[ -x docopts ]] && rm docopts
  run ../get_docopts.sh
  echo "$output"
  [[ $status -eq 0 ]]
}

@test "download binary with get_docopts.sh" {
  # cleanup
  unset GIT_USER
  unset RELEASE

  # we are in ./tests for now
  [[ ${PWD##*/} == 'tests' ]]
  # force an existing released version
  version=$(cat ./VERSION)
  export RELEASE=$version
  _run_get_docopts

  # ensure main repository URL
  match_url='Fetching from: https://github.com/docopt/docopts/'
  [[ $output =~ $match_url ]]

  [[ -x docopts ]]
  # test version match
  run ./docopts --version
  [[ $status -eq 0 ]]
  [[ $output =~ $version ]]
}

@test "get_docopts.sh use another repository" {
  export GIT_USER=Sylvain303
  export RELEASE=v0.6.4-alpha1
  _run_get_docopts

  match_url='Fetching from: https://github.com/Sylvain303/docopts/'
  [[ $output =~ $match_url ]]

  [[ -x docopts ]]
  run ./docopts --version
  [[ $status -eq 0 ]]
  # test version match
  version='v0.6.4-alpha1'
  [[ $output =~ $version ]]
}

@test "arch 64bits detection on macos" {
  # https://github.com/docopt/docopts/issues/44
  #export OSTYPE=darwin
  if [[ $OSTYPE =~ ^darwin.* ]] ; then
    ARCH=$(uname -m)
    [[ $ARCH == x86_64 ]]

    version=$(cat ./VERSION)
    export RELEASE=$version
    _run_get_docopts
    expect="I'm on macos"
    [[ $output =~ $expect ]]

    os=darwin
    match_url="Fetching from: https://github.com/docopt/docopts/releases/download/$version/docopts_${os}_amd64"
    [[ $output =~ $match_url ]]

  else
    skip "only on macos, this OS is: $OSTYPE"
  fi
}
