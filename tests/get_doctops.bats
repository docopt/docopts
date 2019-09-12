#!/usr/bin/env bash
#
# functional test for get_doctops.sh download helper
#

@test "download binary with get_docopts.sh" {
  # go to repository basepath
  cd ..
  rm docopts
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
