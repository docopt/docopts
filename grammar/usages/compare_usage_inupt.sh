#!/usr/bin/env bash

for d in $(ls valid/*.docopt)
 do f=$(basename $d .docopt)
 echo "========= $d $f"
 a=ast/${f}_ast.yaml
 if [[ ! -e $a ]]
 then
   echo "not found $a"
   continue
  fi
  # side by side formating
  #paste -d @ <(sed -n -e '/Usage:/, /^$/ p' $d) <(grep usage_line_input: $a) | column -t -s@
  sed -n -e '/Usage:/, /^$/ p' $d
  grep usage_line_input: $a
done
