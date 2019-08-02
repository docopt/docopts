#!/usr/bin/env bash
#
# Tools for deploying our release to github
#
# Usage: ./deploy.sh deploy [-n] [-r REMOTE_REPOS] [RELEASE_VERSION]
#        ./deploy.sh build
#
# Options:
#   -n                   Dry run, show with version, files and description
#   -r REMOTE_REPOS      Specify a REMOTE_REPOS name [default: origin]
#
# Arguments:
#   RELEASE_VERSION      a git tag
#
# Actions:
#   build  only build using gox and deployment.yml config
#
# deploy.sh read description in deployment.yml


DEPLOYMENT_FILE=deployment.yml
GITHUB_USER=sylvain303
GITHUB_REPO=docopts
TAG="v0.6.3-alpha2"
BUILD_DEST_DIR=build

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

prepare_upload()
{
  local build_dest_dir=$1
  pushd $build_dest_dir > /dev/null
  rm -f sha256sum.txt
  sha256sum * > sha256sum.txt
  popd > /dev/null
  find $build_dest_dir -type f -a ! -name .\*
}

upload_binaries()
{
  local filenames=$*

  local f
  for f in $filenames
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

indent()
{
  local arg="$1"
  if [[ -f "$arg" ]] ; then
    sed -e 's/^/  /' "$arg"
  else
    sed -e 's/^/  /' <<< "$arg"
  fi
}

get_arch_build_target()
{
  local arch_list=$(yq.v2 r $DEPLOYMENT_FILE build | sed -e 's/^- //')
  if [[ $# -eq 1 && $1 == 'gox' ]] ; then
    # join output and trim
    tr $'\n' ' ' <<< "$arch_list" | sed -e 's/ $//'
  else
    echo "$arch_list"
  fi
}

build_binaries()
{
  gox -osarch "$(get_arch_build_target gox)" -output="$BUILD_DEST_DIR/{{.Dir}}_{{.OS}}_{{.Arch}}"
}

yaml_keys ()
{
  yq.v2 r "$1" "$2" | sed -n -e '/^\([^ ]\([^:]\+\)\?\):/  s/:.*// p'
}

main_deploy()
{
  # redefine GITHUB_TOKEN to test if exported for strict mode
  GITHUB_TOKEN=${GITHUB_TOKEN:-}

  if [[ -n $ARGS_RELEASE_VERSION ]] ; then
    TAG=$ARGS_RELEASE_VERSION
  else
    # fetch last tag from git
    TAG=$(git describe --abbrev=0)
  fi

  repository=$(git remote -v | grep $ARGS_REMOTE_REPOS | grep push | head -1)
  description=$(yq.v2 r $DEPLOYMENT_FILE "releases[$TAG].description")
  if [[ -z $description || $description == null ]] ; then
    echo "description not found for tag '$TAG' in $DEPLOYMENT_FILE"
    echo "available git tags: ($repository)"
    indent "$(git tag)"
    echo "available git tags in $DEPLOYMENT_FILE:"
    indent "$(yaml_keys $DEPLOYMENT_FILE releases)"
    return 1
  fi

  build_binaries
  UPLOAD_FILES=$(prepare_upload $BUILD_DEST_DIR)

  if $ARGS_n ; then
    cat << EOT
GITHUB_TOKEN: $GITHUB_TOKEN
build_dir: $BUILD_DEST_DIR
repository: $repository
tag: $TAG
files: $UPLOAD_FILES
sha256sum.txt:
$(indent $BUILD_DEST_DIR/sha256sum.txt)
description:
$(indent "$description")
EOT
    exit 0
  else
    if [[ -z $GITHUB_TOKEN ]] ; then
      echo "GITHUB_TOKEN must be exported"
      return 1
    fi
    upload_binaries $UPLOAD_FILES
  fi
}

if [[ $0 == $BASH_SOURCE ]] ; then
  # bash strict mode
  set -euo pipefail

  # we add our repository path to run our local docopts binary
  # you will have to build it first of course.
  PATH=$(dirname $0):$PATH
  source docopts.sh --auto -G "$@"
  # fix docopt bug https://github.com/docopt/docopt/issues/386
  ARGS_REMOTE_REPOS=${ARGS_REMOTE_REPOS:-$ARGS_r}

  docopt_print_ARGS -G

  if $ARGS_build ; then
    # build only
    echo "dest build dir: $BUILD_DEST_DIR/"
    build_binaries
    ls -l $BUILD_DEST_DIR
    exit 0
  elif [[ $ARGS_deploy ]] ; then
     main_deploy
  else
    echo "no command found: $*"
    exit 1
  fi
fi
