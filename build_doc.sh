#!/bin/bash
#
# helper that insert external source into README.md content
#
# Usage: build_doc.sh [-d] [-i] README_FILE
#
# Options:
#   -i         Edit in place README_FILE is modified
#   -d         diff: README_FILE is unmodified diff is outputed
#
# The generated text from README_FILE is outputed on stdout by default.
#
# Behavior:
#   parse README_FILE and look for marker matching markdown link identifier.
#   https://stackoverflow.com/questions/4823468/comments-in-markdown#20885980
#
#   our markup format: (with and empty line above)
#
#   [make README.md]: # (PARSED_CODE_HERE)
#
# PARSED_CODE_HERE can be:
#   bash command         : will be executed and inserted verbatim
#   include local/file   : imported verbatim
#
# The markdown following code bloc ``` [...] ``` will be remplaced with
# the new content.
#
# For include, a markdown link to the local/file will also be remplaced
# and inserted above the code block.
#
# Our markdown comment markup is conserved.
#
# some temporary files: /tmp/content.* are created.

# blank line above ^^

################################ our functions

# Input parser:
# return a @ separated $line_num@$edit_line@$shell_cmd
parse_input() {
  local input=$1
  local cmd line_num edit_line src
  local oldIFS=$IFS
  IFS=$'\n\t '
  # read the output of the extracted content at the end of the while
  while read line_num src
  do
    # ${var:0:7} extract first 7 chars of a string in bash
    if [[ ${src:0:7} == 'include' ]] ; then
      # local file link
      # format: include local/path/to/filename

      local src=${src:8}
      cmd="include '$src'"

      # line_num+4 is supposed to be the starting ```
      # because to maintain a source link in markdown too
      edit_line=$((line_num + 4))
    else
      # free bash command

      cmd="$src"
      # line_num+2 is supposed to be the starting ```
      edit_line=$((line_num + 2))
    fi

    # display result
    echo "$line_num@$edit_line@$cmd"

  done < <(extract_markup $input)
  IFS=$oldIFS
}

# fetch line with our markup
# for each line:
#
#   [make README.md]: # (PARSED_CODE_HERE)
#
# we output:
#
#   line_num PARSED_CODE_HERE
extract_markup() {
  local input=$1
  local our_markup='[make README.md]: # ('
  grep -n -F "$our_markup" $input | sed -e 's/^\([0-9]\+\):[^(]\+(/\1 /' -e 's/)$//'
}

# transform a free string to a valid filename identifier
# remove quote, space and so on
to_filename() {
  echo "$1" | tr -d '/. `"'"'"'$,^:|\\{}[]()&#\\'
}

# given $1 contains a markdown link: [text](path/to/local/file)
# return the path/to/local/file
# or no_match is not found
markdown_extract_link() {
  local regexp='\[[^]]+\]\(([^)]+)\)'
  if [[ $1 =~ $regexp ]] ; then
    echo "${BASH_REMATCH[1]}"
  else
    echo "no_match"
  fi
}

# search position of $2 in $1
# return -1 if not found
# 0 to n-1
strpos() {
	# remove $2 from $1 and any character following to the end
  local x="${1%%$2*}"
  [[ "$x" = "$1" ]] && echo -1 || echo "${#x}"
}

find_end_content() {
  local start_line=$1
  local input=$2
  # find closing ``` : first matching ``` after $start_line
  gawk "NR > $start_line && /^\`{3}/ {print NR; exit}" $input
}

eval_wrapper() {
  local line_num=$1
  local start=$2
  local mycmd="$3"

  # used by printf -v var_name to return a modified value
  # See man bash for printf
  local return_begin_var=$4
  local return_content_filename=$5

  # build temporary filename used later by sed
  local content_filename=/tmp/content.$(to_filename "$mycmd")
  # merge multiple output
  if [[ $mycmd =~ ^include ]] ; then
    # we shift back to also replace markdown source link
    printf -v $return_begin_var "%d" $((start - 2))
    # pass the line number
    eval "$mycmd $start" > $content_filename
  else
    # default
    printf -v $return_begin_var "%d" $start
    echo '```'       >  $content_filename
    eval "$mycmd"    >> $content_filename
    echo -e '```\n'  >> $content_filename
  fi

  printf -v $return_content_filename "%s" $content_filename
}

build_sed_cmd() {
  local input=$1
  local build=$(parse_input $input)

  local oldIFS=$IFS
  IFS=$'@\n'
  # the loop combine the parsed input into a valid sed command
  # building a oneliner sed allow us to keep lines number while deleting old content for replacing it.
  local sed_cmd=""
  local filename start begin_line end_line mycmd line_num
  while read line_num start mycmd
  do
    eval_wrapper $line_num $start "$mycmd" begin_line filename
    end_line=$(find_end_content $start $input)
    # The begin_line,end_line computation must keep the insert line itself or no content will be outputed.
    # begin_line is the index on the first line to replace (we remove the blank line above)
    # end_line   is the index of the last line.
    # We insert at end_line+1, reusing the blank line after the bloc, which will become the new blank
    # line on top of the inserted block. Inserted block also append a new extra blank line to keep our
    # logic.
    sed_cmd="$sed_cmd -e '$((begin_line-1)),$((end_line)) d' -e '$((end_line+1)) r $filename'"
  done <<<"$build"
  echo "$sed_cmd"
  IFS=$oldIFS
}

##################################### formating helpers

# extract text bloc starting at Usage to the first blank line.
# remove Usage and blank line.
get_usage() {
  sed -n '/^Usage:/,/^$/ p' | sed -e '1d' -e '$d'
}

# extract docopts version
get_version() {
  local text="$1"
  local version=$(gawk 'NR == 1 { print $2 }')
  echo "$text $version"
}

# include the given file
include() {
  local local_file="$1"
  local num_line=$2
  if [[ ! -f $local_file ]] ; then
    echo "include: file not found '$local_file' at '$num_line'"
    return 1
  else
    local format_type=""
    # detect file type
    if grep -q 'shell script' <<< "$(file $local_file)" ; then
      format_type="bash"
    fi

    echo -e "[source $local_file]($local_file)\n"
    echo '```'"$format_type"
    cat "$local_file"
    echo -e '```\n'
  fi
}

################## main

main_build_doc() {
  README_FILE=$1

  local sed_cmd=$(build_sed_cmd $README_FILE)

  # eval is required to preserve multiple args quoting and passing multiple actions to one sed process.
  if $ARGS_i ; then
    # edit inplace
    eval "sed -i $sed_cmd $README_FILE"
  elif $ARGS_d ; then
    diff -u $README_FILE <(eval "sed $sed_cmd $README_FILE")
  else
    # echo "$sed_cmd"
    eval "sed $sed_cmd $README_FILE"
  fi
}

if [[ $0 == $BASH_SOURCE ]] ; then
  # we add our repository path to run our local docopts binary
  # you will have to build it first of course.
  PATH=$(dirname $0):$PATH
  source docopts.sh --auto -G "$@"
  main_build_doc "$ARGS_README_FILE"
fi
