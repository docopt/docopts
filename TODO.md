# docopts (docopt for bash) TODO list



## better error handling

https://github.com/docopt/docopts/issues/17

See also:

PR: https://github.com/docopt/docopt.go/pull/65

## --json output

same as `--no-mangle` but json formated

Somewhat discussed here: https://github.com/docopt/docopt/issues/174

## functional testing for all options

`./docopts --help`

## return or kill for function instead of exit

Add a parameter to handle return or kill instead of exit so it can be launched inside a function.

## embeded JSON

See [API_proposal.md](API_proposal.md)

## generate bash completion from usage

```
docopts -h "$help" --generate-completion
```

## embed test routine (validation)?

May we can interract with the caller to eval some validation…
It is needed? Is it our goal?

2019-06-07: I think it's ouside `docopts` goal to perform validation. It requires extra language to validate data and it
will pollute bash own programming role.


```bash
# with tests
# pass value to parent: JSON or some_thing_else
eval $(docopts --eval --json --help="Usage: mystuff [--code] INFILE [--out=OUTFILE]" -- "$@")

# docopts test would perform some check based on our own testing language
if docopts test -- file_exists:INFILE !file_exists:OUTFILE
then
  # normal action INFILE exists and OUTFILE will not be ovrerwritten
else
  # some error
fi

eval $(docopts --eval --json --help="Usage: prog [--count=NUM] INFILE..."  -- "$@")
if docopts test -- num:gt:1:NUM file_exists:INFILE
then
  # normal action can be performed
fi
```

## config file parse config to option format

À la nslcd… ?

* json merge
* toml merge (ini)
