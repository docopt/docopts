# docopts legacy examples

These examples work on bash 3.x by avoiding `docopts.sh` code
that uses associative arrays (bash 4.x only).

## conversion

`sshdiff_legacy.sh` doesn't use `docopts` at all it is converted to `sshdiff_with_docopts.sh` as an example of rewrite
legacy option parsing not using associative arrays.
