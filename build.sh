#!/bin/bash
#
# Build embedded docopt.py into docopts.sh

# exit on error
set -e

echo "no build"
exit 0

echo "download: docopt.py form github…"
wget --no-check-certificate -O docopt.py "https://raw.githubusercontent.com/docopt/docopt/master/docopt.py"

[[ -s docopt.py ]]

sourcef=docopts.sh
embedded_mark='### EMBEDDED'

echo "merging docopt.py + docopts.py"
# remove duplicate import and code
sed -i -e '/^import sys/ d' \
    -e '/^import re/ d' \
    -e '/^__all__ =/ d' \
    -e '/^__version__ =/ d' \
    docopt.py

# replace import
# from docopt import docopt, DocoptExit, DocoptLanguageError
sed -e '/^from docopt/ {
   s//# embedded: &/
   r docopt.py
   a\
# ----------------------- end docopt.py ---------------------------
   }' docopts > docopts.py

# check if code is modified
s=$(git status -uno --short $sourcef)
if [[ ! -z "$s" ]] ; then
    echo $s
    echo "$sourcef is modified, commit or revert first"
    exit 1
fi
echo "removing embedded code behind '$embedded_mark'"
sed -i -e "/^$embedded_mark/,\$ d" $sourcef
echo "add marker and reembbed code"
echo "$embedded_mark" >> $sourcef
sed -e "s/^/#> /" docopts.py >> $sourcef
# remove ending blanks there's a TAB + space in the [] below
sed -i -e 's/[ 	]\+$//g' $sourcef

