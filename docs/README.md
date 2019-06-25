# `docopts` Documentation

`docopts` is a standalone binary program that generate bash code suitable for `eval`.

See: [README.md](../README.md)

## `doctops.sh` bash library helper

`doctops.sh` is a library companion of `doctops` binary.

It may become the main entry point using `doctops`, but keep in mind that `doctops` is also a standalone binary.

`doctops.sh` is intensively used in [../examples/](../examples/).

It has 2 behaviors:

* a library `source doctops.sh` (found in your `$PATH`)
* a wrapper that auto parse option from a dedicated `Usage:` comment from your header script

### `doctops.sh` as a library

No action is performed the functions and globals are now part as your bash environment.

You can manualy call the lib functions.

[make README.md]: # (sed -n -e '/^# Doc:$/,/^docopt_/ p' docopts.sh)

```
docopt_get_help_string()
docopt_get_version_string()
docopt_get_values()
docopt_get_eval_array()
docopt_get_raw_value()
docopt_print_ARGS()
```

### `docopts.sh` as a wrapper

Automaticaly parse and eval the `Usage:` header of your script and build `docopts` call for you.
The script will fail if option parse error are encountered.

Generate Bash 4.0 associative array `${ARGS[--option]}`:

```
source docopts.sh --auto "$@"
```

Generate Bash mangled GLOBALS `$ARGS_option`:

```
source docopts.sh --auto -G "$@"
```
