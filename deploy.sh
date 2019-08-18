#!/usr/bin/env bash
#
# Tools for deploying our release to github
#
# Usage: ./deploy.sh deploy [-n] [-r REMOTE_REPOS] [--replace] [-u GITHUB_USER] [RELEASE_VERSION]
#        ./deploy.sh build [RELEASE_VERSION]
#        ./deploy.sh delete RELEASE_VERSION
#
# Options:
#   -n                   Dry run, show with version, files and description
#   -r REMOTE_REPOS      Specify a REMOTE_REPOS name [default: origin]
#   --replace            Replace existing release with this one, previous release
#                        will be deleted first.
#   -u GITHUB_USER       force this GITHUB_USER
#
# Arguments:
#   RELEASE_VERSION      a git tag, or current for the local modified version
#
# Actions:
#   build      only build using gox and deployment.yml config
#   deploy     prepare and deploy the release
#   delete     delete the given RELEASE_VERSION from github and all assets
#
# deploy.sh reads description and name for releases in deployment.yml


DEPLOYMENT_FILE=deployment.yml
# change GITHUB_USER + GITHUB_REPO to change repository, it is for building API URL
GITHUB_USER=sylvain303
GITHUB_REPO=docopts
TAG="v0.6.3-alpha2"
BUILD_DEST_DIR=build

create_release()
{
  local release="$1"
  local name="$2"
  local description="$3"

  # detect alpha ==> pre-release
  # match -ending
  local pre_release=""
  if [[ $release =~ -[a-zA-Z0-9_-]$ ]] ; then
    pre_release='--pre-release'
  fi

  gothub release \
      --user $GITHUB_USER \
      --repo $GITHUB_REPO \
      --tag "$release" \
      --name "$name" \
      --description "$description" \
      $pre_release
}

# check the the given release exists, test with $?
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

# after build, generate sha256sum for all file in BUILD_DEST_DIR
# then output all files name from parent directory
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

# read os/arch from $DEPLOYMENT_FILE
# format as list or space separated list for gox
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

  local ldflags

  if [[ $release == current ]] ; then
    cp docopts.go $build_dest_dir
    # will user ./VERSION to get the version
    ldflags="$(govvv -flags)"
  else
    # checkout release version to $BUILD_DEST_DIR and build it from here
    git show $release:./docopts.go > $build_dest_dir/docopts.go
    # we force the version to be $release.
    # before version v0.6.3-rc1. this will overwrite code version string.
    # Which produces a binary diffrent from `go build` (on those older source) for --version
    ldflags="$(govvv -flags -version "$release")"
  fi

	# ldflags need to be synchronised with Makefile
  local go_version="$(go version)"
  ldflags+=" -X 'main.GoBuildVersion=$go_version'"

  local osarch="$(get_arch_build_target gox)"
  # chdir to the build_dest_dir, use the docopts.go copy for source to compile.
  pushd $build_dest_dir > /dev/null
    # -output allow to force generated binaries
    gox -osarch "$osarch" -output="docopts_{{.OS}}_{{.Arch}}" -ldflags "$ldflags"
  popd > /dev/null
}

yaml_keys()
{
  # yq seems to young software, quick fix to get keys
  # https://github.com/mikefarah/yq/issues/20
  yq.v2 r "$1" "$2" | sed -n -e '/^\([^ ]\([^:]\+\)\?\):/  s/:.*// p'
}

show_release_data()
{
  local release=$1
  local name="$2"
  local description="$3"

  local repository=$(git remote -v | grep $ARGS_REMOTE_REPOS | grep push | head -1)

  cat << EOT
GITHUB_REPO: $GITHUB_REPO
GITHUB_USER: $GITHUB_USER
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
  local release=$1
  local name="$2"
  local description="$3"
  if [[ -z $description || $description == null || -z $name || $name == null ]] ; then
    echo "description or name not found for tag '$release' in $DEPLOYMENT_FILE"
    echo "available git tags:"
    indent "$(git tag)"
    echo "available git tags in $DEPLOYMENT_FILE:"
    indent "$(yaml_keys $DEPLOYMENT_FILE releases)"
    echo "VERSION contains"
    indent "$(cat VERSION)"
    return 1
  fi
}

main_deploy()
{
  local release=$1
  local release_version=$release

  # redefine GITHUB_TOKEN to test if exported for strict mode
  GITHUB_TOKEN=${GITHUB_TOKEN:-}

  if [[ $release == current ]] ; then
    release_version=$(cat VERSION)
    echo "using current release in VERSION: $release_version"
  fi

  local description=$(yq.v2 r $DEPLOYMENT_FILE "releases[$release_version].description")
  local name=$(yq.v2 r $DEPLOYMENT_FILE "releases[$release_version].name")

  # will stop the execution (as set -e is enabled)
  check_name_description $release_version "$name" "$description"

  build_binaries $release $BUILD_DEST_DIR
  UPLOAD_FILES=$(prepare_upload $BUILD_DEST_DIR)

  if $ARGS_n ; then
    show_release_data $release_version "$name" "$description"
    exit 0
  else
    if [[ -z $GITHUB_TOKEN ]] ; then
      echo "GITHUB_TOKEN must be exported"
      return 1
    fi

    echo "deploying release $GITHUB_USER/$GITHUB_REPO: $release_version"

    if check_release $release_version ; then
      echo "release already exists: $release_version"
      if $ARGS_replace ; then
        echo "deleting existing release: $release_version"
        delete_release $release_version
        echo "creating release: $release_version"
        create_release $release_version "$name" "$description"
      else
        echo "use --replace to replace the existing release"
        echo "only upload new files"
      fi
    else
      echo "release doesn't exists yet: $release_version"
      echo "creating new release: $release_version"
      create_release $release_version "$name" "$description"
    fi
    upload_binaries $release_version $UPLOAD_FILES
  fi
}

check_env()
{
  local v val
  local error=0

  for v in GOPATH GOBIN
  do
    eval "val=\${$v:-}"
    if [[ -z $val ]] ; then
      echo "$v is undefined, check failed"
      error=$((error+1))
    fi
  done
  return $error
}

if [[ $0 == $BASH_SOURCE ]] ; then
  # bash strict mode
  set -euo pipefail

  check_env

  # we add our repository path to run our local docopts binary
  # you will have to build it first of course.
  PATH=$(dirname $0):$PATH
  source docopts.sh --auto -G "$@"
  docopt_print_ARGS -G

  # fix docopt bug https://github.com/docopt/docopt/issues/386
  ARGS_REMOTE_REPOS=${ARGS_REMOTE_REPOS:-$ARGS_r}
  ARGS_GITHUB_USER=${ARGS_GITHUB_USER:-$ARGS_u}

  if [[ -n $ARGS_GITHUB_USER ]]; then
    GITHUB_USER=$ARGS_GITHUB_USER
  fi

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
  elif $ARGS_deploy ; then
    main_deploy $TAG
  elif $ARGS_delete ; then
    echo "deleting release $GITHUB_USER/$GITHUB_REPO: $TAG"
    delete_release $TAG
  else
    echo "no command found: $*"
    exit 1
  fi
fi
