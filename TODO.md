# docopts (docopt for bash) TODO list or questions

## better error handling

https://github.com/docopt/docopts/issues/17

See also:

PR: https://github.com/docopt/docopt.go/pull/65

It probably needs to rewrite the docopt parser.

## --json output

same as `--no-mangle` but json formated

Somewhat discussed here: https://github.com/docopt/docopt/issues/174

Trivial, could be implemented, even without embbeding JSON lib.
See branch `json-api` too.

## functional testing for all options

`./docopts --help`
* `tests/functional_tests_docopts.bats` was introduced in PR #52

## return or kill for function instead of exit

Add a parameter to handle return or kill instead of exit so it can be launched inside a function.

See also: https://github.com/docopt/docopts/issues/43

## verb action

```
docopts parse "$usage" : [<args>...]
```

## generate bash completion from usage

Would probably need a new docopt parser too.

```
docopts completion "$usage"
```

## config file parse config to option format

À la nslcd… ?

* json merge
* toml merge (ini)
