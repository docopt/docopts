#!/bin/bash
#
# helper that insert external source into README.md
#

README=$1

# return a colon separated $start_line:$end_line:$shell_cmd
get_build_doc() {
  local cmd num_line start_line end_line
  # read the output of the grep at the end of the while
  # which is a splited grep -n output
  while read num_line source 
  do
    # ${var:0:1} extract first char of a string in bash
    if [[ ${source:0:1} == '`' ]] ; then
      # extract inner text, removing first char and last char
      cmd=${source:1:$((${#source} - 2))}
    elif [[ ${source:0:1} == '[' ]] ; then
      # remove prefix including openinng '('
      cmd=${source#*(}
      # remove last char, supposed to be the closing ')'
      cmd=${cmd:0:$((${#cmd} - 1))}
      cmd="cat $cmd"
    fi

    # num_line+2 is supposed to be the starting ```
    start_line=$(($num_line + 2))

    # find closing ``` : first matching ``` after $start_line
    end_line=$(awk "NR > $start_line && /^\`{3}/ {print NR; exit}" $README)

    # display result
    echo "$start_line:$end_line:$cmd"

  done < <(grep -n '^make build_doc' $README | awk -F':' '{$2=""; print $0}')
}

to_filename() {
  echo "$1" | tr -d '/. `"'"'"'$,^:|\\'
}

build=$(get_build_doc)
IFS=$':\n'
tmp=README.tmp
cp $README $tmp
cmd=""
while read start end mycmd
do
  content=/tmp/content.$(to_filename "$mycmd")
  { echo '```bash' ; eval "$mycmd" ; echo -e '```\n'; }  > $content
    cmd="$cmd -e '$((start -1)),$end d' -e '$((end +1)) r $content'"
done <<<"$build"

echo "sed -i $cmd $tmp"
eval "sed -i $cmd $tmp"
diff -u $README $tmp
