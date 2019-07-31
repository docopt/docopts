#!/usr/bin/env bash
#export GITHUB_TOKEN=...

set -euo pipefail

GITHUB_USER=sylvain303
GITHUB_REPO=docopts
TAG="v0.6.3-alpha2"

create_release()
{
  gothub release \
      --user $GITHUB_USER \
      --repo $GITHUB_REPO \
      --tag $TAG \
      --name "docopt for shell Bash" \
      --description "Written in Go. This binaries is for GNU/Linux 32bits and 64bits" \
      --pre-release
}

upload_binaries()
{
  local filenames=$*
  sha256sum $filenames > sha256sum.txt

  local f
  for f in $filenames sha256sum.txt
  do
    echo "uploading '$f' ..."
    gothub upload \
        --user $GITHUB_USER \
        --repo $GITHUB_REPO \
        --tag $TAG \
        --name "$f" \
        --file $f \
        --replace
  done
}

upload_binaries docopts docopts-32bits docopts-OSX
