# docopts (docopt for bash) TODO list

## functional testing for all option

## return or kill for function instead of exit

Add a parameter to handle return or kill instead of exit so it can be launched inside a function.

## embeded JSON

See [API_proposal.md](API_proposal.md)

## build and publish binary

Reuse build.sh to build golang binary and pubilsh it as a new release too.

## generate bash completion from usage

```
docopts -h "$help" --generate-completion
```

## embed test routine (validation)?

May we can interract with the caller to eval some validation…
It is needed? Is it our goal?

```
# with tests
# pass value to parent: JSON or some_thing_else
eval $(docopts --eval --json --help="Usage: mystuff [--code] INFILE [--out=OUTFILE]" -- "$@")
if docopts test -- file_exists:--code !file_exists:--out

eval $(docopts --eval --json --help="Usage: prog [--count=NUM] INFILE..."  -- "$@")
if docopts test -- num:gt:1:--count file_exists:INFILE
```

## config file parse config to option format

À la nslcd… ?

* json merge
* toml merge (ini)
