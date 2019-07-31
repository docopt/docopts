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
  local filename=$1
  sha256sum $filename >> sha256sum.txt
  gothub upload \
      --user $GITHUB_USER \
      --repo $GITHUB_REPO \
      --tag $TAG \
      --name "$filename" \
      --file $filename

  gothub upload \
      --user $GITHUB_USER \
      --repo $GITHUB_REPO \
      --tag $TAG \
      --name "sha256sum.txt" \
      --file  sha256sum.txt
}

rm sha256sum.txt
upload_binaries $1
