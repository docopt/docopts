# current work in PROGESS

## `mangle_name` error handling

-9 ??

is it the same as "TODO.md":(TODO.md) `-G` for global prefix.

```
sylvain@lap40:~/code/go/src/github.com/Sylvain303/docopts$ ./docopts  -h "usage: prog [-9] FILE..." --debug : pipo molo toto
################## golang ##################
             --debug : true
              --help : usage: prog [-9] FILE...
           --no-help : false
         --no-mangle : false
     --options-first : false
         --separator : ----
           --version : <nil>
                  -A : <nil>
                   : : true
              <argv> : [pipo molo toto]
                 doc : usage: prog [-9] FILE...
        bash_version :
################## bash ##################
                  -9 : false
                FILE : [pipo molo toto]
----------------------------------------
# name_mangle:error:cannot transform into a bash identifier: -9:-9
FILE=('pipo', 'molo', 'toto')

```


```
sylvain@lap40:~/code/docopt/docopts$ ./docopts  -h "usage: prog [-9] FILE..." : pipo molo toto
./docopts: name could not be mangled into a valid Bash identifier: -9
```


# next

from `testee.sh` merge `get_raw_value` into `docopts.sh`

## provide test on old environment

docker?
32bist
bash 3
