# docopts (docopt for bash) TODO list

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

