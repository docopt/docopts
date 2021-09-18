# New API proposal for docopts - docopt on shell (bash)

## Update: 2021-09-16

[docopts.sh](docopts.sh) has been publised as a `docopts` companion its documentation live here:
[docs/README.md](docs/README.md).

It offer some API alternatives and allow to build external wrapper.

## new cli API proposal for 2021

* use verb for action
* stop using `-h` `<msg>` for reading help
* keep compatibily with old syntax trough a new verb: `compat`
* introduce more useful error message


### new verb

### replace `-h`

```
Usage:
   docopts [options] parse USAGE : [<argv>...]
   docopts [options] parse [--no-declare] -A <name> USAGE : [<argv>...]
   docopts [options] parse -G <prefix> USAGE : [<argv>...]
```

Example:

```
docopts parse "$USAGE" : <option passed here>
```

`-h` still available trough `compat`

```
docopts compat -h "$USAGE" : <option passed here>
```

### generate usage completion

generate usage completion (as `kubect completion` done by [cobra](https://github.com/spf13/cobra/blob/master/shell_completions.md))

```
docopts completion "$usage"
```

### debug mode:

```
Usage:
   docopts debug [--explain-parsed] USAGE : [<argv>...]
   docopts debug [--show-argument-match] USAGE : [<argv>...]
```

### `docopt.sh` compatibile generator

```
Usage:
   docopts generate USAGE
   docopts generate -f FILENAME
```


## Dropped JSON API proposal

old [API_proposal.md](https://github.com/docopt/docopts/blob/031ceb2f0700ac0e40303d72167b91586a6d60da/API_proposal.md)
