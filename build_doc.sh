#!/bin/bash
#
# helper that insert external source into README.md content
#
# Usage: ./build_doc.sh [-d] [-i] README_FILE
#
# Options:
#   -i         Edit in place README_FILE is modified
#   -d         diff: README_FILE is unmodified diff is outputed
#
# Behavior:
#   parse README_FILE and look for marker matching: ^make build_doc
#
#   make build_doc: `markdown bash command`
#   make build_doc: include [markdown link](to/a/local file)
#
# bash command : will be executed and inserted verbatim
# local file   : imported verbatim
#
# The generated text from README_FILE is outputed on stdout.
#
# The markup "make build_doc: *" is conserved.
#
# some files: /tmp/content.* are created

# blank line above ^^

# Input parser:
# return a @ separated $start_line@$end_line@$type@$shell_cmd
get_build_doc() {
  local input=$1
  local cmd num_line start_line end_line type
  # read the output of the grep at the end of the while
  # which is a splited grep -n output
  while read num_line source 
  do
    type=""
    # ${var:0:1} extract first char of a string in bash
    if [[ ${source:0:1} == '`' ]] ; then
      # free bash command

      # extract inner text, removing first char and last char
      cmd=${source:1:$((${#source} - 2))}
    elif [[ ${source:0:7} == 'include' ]] ; then
      # local file link

      local src=$(markdown_extract_link "$source")
      cmd="include '$src' $num_line"

      # detect file type
      if grep -q 'shell script' <<< "$(file $src)" ; then
        type="bash"
      fi
    else
      # unknown command
      cmd="echo '$num_line: not recognized command: $source'"
    fi

    # num_line+2 is supposed to be the starting ```
    start_line=$(($num_line + 2))

    # find closing ``` : first matching ``` after $start_line
    end_line=$(awk "NR > $start_line && /^\`{3}/ {print NR; exit}" $input)

    # display result
    echo "$start_line@$end_line@$type@$cmd"

  done < <(grep -n '^make build_doc' $input | awk -F':' '{$2=""; print $0}')
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

# formating helpers

# extract text bloc starting at Usage to the first blank line.
# remove Usage and blank line.
get_usage() {
  sed -n '/^Usage:/,/^$/ p' | sed -e '1d' -e '$d'
}

# include the given file
include() {
  local local_file="$1"
  local num_line=$2
  if [[ $local_file == no_match ]] ; then
    echo "include: file not matched at line '$num_line'"
  else
    cat "$local_file"
  fi
}

################## main

main_build_doc() {
  README_FILE=$1

  build=$(get_build_doc $README_FILE)

  # NOTE: changing IFS will alter get_build_doc parsing, it must be done after.
  IFS=$'@\n'
  # the loop combine the parsed input into a valid sed command
  # building a oneliner sed allow us to keep lines number while deleting old content for replacing it.
  sed_cmd=""
  while read start end format_type mycmd
  do
    content=/tmp/content.$(to_filename "$mycmd")
    # merge multiple output
    { echo '```'"$format_type" ; eval "$mycmd" ; echo -e '```\n'; }  > $content
    # the start end computation must keep the line itself or no content will be outputed
    sed_cmd="$sed_cmd -e '$((start -1)),$end d' -e '$((end +1)) r $content'"
  done <<<"$build"

  # eval is required to preserve multiple args quoting and passing multiple actions to one sed process.
  if $ARGS_i ; then
    # edit inplace
    eval "sed -i $sed_cmd $README_FILE"
  elif $ARGS_d ; then
    diff -u $README_FILE <(eval "sed $sed_cmd $README_FILE")
  else
    eval "sed $sed_cmd $README_FILE"
  fi
}

if [[ $0 == $BASH_SOURCE ]] ; then
  source ./docopts.sh --auto -G "$@"
  main_build_doc "$ARGS_README_FILE"
fi
