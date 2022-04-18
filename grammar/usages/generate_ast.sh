#!/usr/bin/env bash
#
# Generate the parsed AST.yaml from docopt input file
#
# Usage: generate_ast.sh [--rebuild] DIRNAME... --out-dir=DEST_DIR
#
# Arguments:
#   DIRNAME    a foldername to look for .docopt file
#   DEST_DIR   a foldername into which to generate AST.yaml parsed file
#
# Options:
#   --realpath   force rebuild of docopt-analyze (previous binary will be removed)
#

me=$(realpath  $0)
my_dir=$(dirname $me)
source $my_dir/../../docopts.sh --auto -G "$@"

if [[ ! -d $ARGS_DIRNAME ]]
then
  echo "DIRNAME, not found: '$DIRNAME'"
  exit 1
fi

usage_folder=$(realpath $ARGS_DIRNAME)

if [[ ! -d $ARGS_out_dir ]]
then
  mkdir -p $ARGS_out_dir || { echo "error: creating folder '$ARGS_out_dir'"; exit 1; }
fi

dest_dir=$(realpath $ARGS_out_dir)

parser=$(realpath $my_dir/../../cmd/docopt-analyze/docopt-analyze)
if [[ ! -x $parser || $ARGS_rebuild == 'true' ]]
then
  echo "try building parser '$parser'"
  cd $(dirname $parser)
  go build
  if [[ ! $? -eq 0 ]]
  then
    echo "building error"
    exit 1
  fi
fi

cd $usage_folder
for u in *.docopt
do
  echo "==================== $u"
  ast_fname="$dest_dir/$(basename $u .docopt)_ast.yaml"
  $parser -y $u > $ast_fname
done

