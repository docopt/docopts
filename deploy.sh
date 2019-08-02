#!/usr/bin/env bash
#
# Tools for deploying our release to github
#
# Usage: ./deploy.sh deploy [-n] [-r REMOTE_REPOS] [--replace] [RELEASE_VERSION]
#        ./deploy.sh build [RELEASE_VERSION]
#
# Options:
#   -n                   Dry run, show with version, files and description
#   -r REMOTE_REPOS      Specify a REMOTE_REPOS name [default: origin]
#   --replace            Replace existing release with this one, previous release
#                        will be deleted first.
#
# Arguments:
#   RELEASE_VERSION      a git tag
#
# Actions:
#   build      only build using gox and deployment.yml config
#   deploy     prepare and deploy the release
#
# deploy.sh read description in deployment.yml


DEPLOYMENT_FILE=deployment.yml
GITHUB_USER=sylvain303
GITHUB_REPO=docopts
TAG="v0.6.3-alpha2"
BUILD_DEST_DIR=build

create_release()
{
  local release="$1"
  local name="$2"
  local description="$3"

  # TODO: detect alpha ==> pre-release
  gothub release \
      --user $GITHUB_USER \
      --repo $GITHUB_REPO \
      --tag "$release" \
      --name "$name" \
      --description "$description" \
      --pre-release
}

check_release()
{
  local release=$1
  gothub info \
      --user $GITHUB_USER \
      --repo $GITHUB_REPO \
      --tag "$release" > /dev/null 2>&1
}

delete_release()
{
  local release=$1
  gothub delete \
      --user $GITHUB_USER \
      --repo $GITHUB_REPO \
      --tag "$release"
}

prepare_upload()
{
  local build_dest_dir=$1
  pushd $build_dest_dir > /dev/null
  # remove docopts source used for build
  rm -f sha256sum.txt docopts.go
  sha256sum * > sha256sum.txt
  popd > /dev/null
  find $build_dest_dir -type f -a ! -name .\*
}

upload_binaries()
{
  local release=$1
  shift
  local filenames=$*

  local f
  for f in $filenames
  do
    echo "uploading '$f' ..."
    gothub upload \
        --user $GITHUB_USER \
        --repo $GITHUB_REPO \
        --tag "$release" \
        --name "$(basename $f)" \
        --file "$f" \
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
  local release=$1
  local build_dest_dir=$2
  # checkout release version
  git show $release:./docopts.go > $build_dest_dir/docopts.go
  local osarch="$(get_arch_build_target gox)"
  pushd $build_dest_dir > /dev/null
  gox -osarch "$osarch" -output="docopts_{{.OS}}_{{.Arch}}"
  popd > /dev/null
}

yaml_keys ()
{
  yq.v2 r "$1" "$2" | sed -n -e '/^\([^ ]\([^:]\+\)\?\):/  s/:.*// p'
}

show_release_data()
{
  local release=$1
  local name="$2"
  local description="$3"

  local repository=$(git remote -v | grep $ARGS_REMOTE_REPOS | grep push | head -1)

  cat << EOT
GITHUB_TOKEN: $GITHUB_TOKEN
build_dir: $BUILD_DEST_DIR
repository: $repository
name: $name
tag: $release
files: $UPLOAD_FILES
sha256sum.txt:
$(indent $BUILD_DEST_DIR/sha256sum.txt)
description:
$(indent "$description")
EOT
}

check_name_description()
{
  local name="$1"
  local description="$2"
  if [[ -z $description || $description == null || -z $name || $name == null ]] ; then
    echo "description or name not found for tag '$TAG' in $DEPLOYMENT_FILE"
    echo "available git tags:"
    indent "$(git tag)"
    echo "available git tags in $DEPLOYMENT_FILE:"
    indent "$(yaml_keys $DEPLOYMENT_FILE releases)"
    return 1
  fi
}

main_deploy()
{
  # redefine GITHUB_TOKEN to test if exported for strict mode
  GITHUB_TOKEN=${GITHUB_TOKEN:-}

  local description=$(yq.v2 r $DEPLOYMENT_FILE "releases[$TAG].description")
  local name=$(yq.v2 r $DEPLOYMENT_FILE "releases[$TAG].name")

  check_name_description "$name" "$description"

  build_binaries $TAG $BUILD_DEST_DIR
  UPLOAD_FILES=$(prepare_upload $BUILD_DEST_DIR)

  if $ARGS_n ; then
    show_release_data $TAG "$name" "$description"
    exit 0
  else
    if [[ -z $GITHUB_TOKEN ]] ; then
      echo "GITHUB_TOKEN must be exported"
      return 1
    fi

    if check_release $TAG ; then
      echo "release: already $TAG exists"
      if $ARGS_replace ; then
        echo "deleting existing release $TAG"
        delete_release $TAG
        echo "creating release $TAG"
        create_release $TAG "$name" "$description"
      else
        echo "use --replace to replace the existing release"
        echo "only upload new files"
      fi
    else
      echo "release: $TAG doesn't exists yet"
      echo "creating release $TAG ..."
      create_release $TAG "$name" "$description"
    fi
    upload_binaries $TAG $UPLOAD_FILES
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

  if [[ -n $ARGS_RELEASE_VERSION ]] ; then
    TAG=$ARGS_RELEASE_VERSION
  else
    # fetch last tag from git
    TAG=$(git describe --abbrev=0)
  fi

  if $ARGS_build ; then
    echo "build only ..."
    echo "dest build dir: $BUILD_DEST_DIR/"
    build_binaries $TAG $BUILD_DEST_DIR
    ls -l $BUILD_DEST_DIR
    exit 0
  elif [[ $ARGS_deploy ]] ; then
     main_deploy
  else
    echo "no command found: $*"
    exit 1
  fi
fi
