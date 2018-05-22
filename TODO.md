# docopts (docopt for bash) TODO list

## global prefix

-G PREFIX

this option must be incompatible with -A (see also: --no-mangle)

outputs

`PREFIX_MANGLED_NAME=value`

## embeded JSON

```
DOCOPTS_JSON=$(docopts --json --help-mesg="Usage: mystuff [--code] INFILE [--out=OUTFILE]" -- "$@")

# automaticly use  $DOCOPTS_JSON
if [[ $(docopts get -- --code) == checkit ]]
then
  action
fi

# or
docopts get --env SOME_DOCOPTS_JSON -- --code

# or
DOCOPTS_JSON_VAR=SOME_DOCOPTS_JSON
docopts get --code
```

## --no-declare

remove the output of `declare -A hash_name`

## build and publish binary

reuse build.sh to build golang binary and pubilsh it as a new release too.
