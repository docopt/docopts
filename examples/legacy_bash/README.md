# docopts legacy examples

These examples work on bash 3.x by not using `docopts.sh` `--auto` mode.
And without using associative arrays (bash 4.x only).

## Example of conversion

`sshdiff_legacy.sh` doesn't use `docopts` at all. It has been converted by hand to `sshdiff_with_docopts.sh` as an
example of rewrite legacy option parsing not using associative arrays.
