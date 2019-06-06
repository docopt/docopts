#!/bin/bash
#
# helper that insert external source into README.md content
#
# Usage: ./build_doc.sh README_FILE > README.tmp
#
# Behavior:
#   parse README_FILE and look for marker matching: ^make build_doc
#
#   make build_doc: `markdown bash command`
#   make build_doc: `markdown link to a local file`
#
# bash command : will be executed and inserted verbatim
# local file   : imported verbatim
#
# The README_FILE modifed text flow is outputed on stdout.
#
# The markup "make build_doc: *" is conserved.
#
# README_FILE is unmodified.
#
# some files: /tmp/content.* are created

README_FILE=$1

# return a colon separated $start_line:$end_line:$type:$shell_cmd
get_build_doc() {
  local cmd num_line start_line end_line type
  # read the output of the grep at the end of the while
  # which is a splited grep -n output
  while read num_line source 
  do
    # ${var:0:1} extract first char of a string in bash
    if [[ ${source:0:1} == '`' ]] ; then
      # free bash command

      # extract inner text, removing first char and last char
      cmd=${source:1:$((${#source} - 2))}
      type=""
    elif [[ ${source:0:1} == '[' ]] ; then
      # local file link

      # remove prefix including openinng '('
      cmd=${source#*(}
      # remove last char, supposed to be the closing ')'
      cmd=${cmd:0:$((${#cmd} - 1))}
      cmd="cat $cmd"

      # detect file type
      if grep -q 'shell script' <<< "$(file $cmd)" ; then
        type="bash"
      fi
    fi

    # num_line+2 is supposed to be the starting ```
    start_line=$(($num_line + 2))

    # find closing ``` : first matching ``` after $start_line
    end_line=$(awk "NR > $start_line && /^\`{3}/ {print NR; exit}" $README_FILE)

    # display result
    echo "$start_line:$end_line:$type:$cmd"

  done < <(grep -n '^make build_doc' $README_FILE | awk -F':' '{$2=""; print $0}')
}

# transform a string to a valid filename identified
# remove quote, space and so on
to_filename() {
  echo "$1" | tr -d '/. `"'"'"'$,^:|\\'
}


################## main

build=$(get_build_doc)

# NOTE: changing ISF wil alter get_build_doc parsing, it must be done after.
IFS=$':\n'
# the loop combine the parsed input into a valid sed command
sed_cmd=""
while read start end format_type mycmd
do
  content=/tmp/content.$(to_filename "$mycmd")
  # merge multiple output
  { echo '```'"$format_type" ; eval "$mycmd" ; echo -e '```\n'; }  > $content
  # the start end computation must keep the line itself or no content will be ouptputed
  sed_cmd="$sed_cmd -e '$((start -1)),$end d' -e '$((end +1)) r $content'"
done <<<"$build"

tmp=README_FILE.tmp
cp $README_FILE $tmp
echo "sed -i $sed_cmd $tmp"
# eval is required to preserve multiple args quoting and passing multiple actions to one sed process.
eval "sed -i $sed_cmd $tmp"
diff -u $README_FILE $tmp
