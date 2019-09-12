#!/usr/bin/env bash
#
# functional test for get_doctops.sh download helper
#

@test "download binary with get_docopts.sh" {
  # go to repository basepath
  cd ..
  [[ -x docopts ]] && rm docopts
  run ./get_docopts.sh
  echo "$output"
  [[ $status -eq 0 ]]

  # ensure main repository URL
  match_url='Fetching from: https://github.com/docopt/docopts/'
  [[ $output =~ $match_url ]]

  [[ -x docopts ]]
  run ./docopts --version
  [[ $status -eq 0 ]]
  # test version match
  version=$(cat VERSION)
  [[ $output =~ $version ]]
}


@test "get_docopts.sh use another repository" {
  # go to repository basepath
  cd ..
  [[ -x docopts ]] && rm docopts
  export GIT_USER=Sylvain303
  export RELEASE=v0.6.4-alpha1
  run ./get_docopts.sh
  echo "$output"
  [[ $status -eq 0 ]]

  match_url='Fetching from: https://github.com/Sylvain303/docopts/'
  [[ $output =~ $match_url ]]

  [[ -x docopts ]]
  run ./docopts --version
  [[ $status -eq 0 ]]
  # test version match
  version='v0.6.4-alpha1'
  [[ $output =~ $version ]]
}
