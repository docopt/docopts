# WARNING: a new shell API is comming. 

See Proposal lib API: [The Wiki](https://github.com/docopt/docopts/wiki)
The [develop branch](https://github.com/docopt/docopts/tree/develop) is abandoned

## docopts

See the Reference manual for command line `docopts` (old
[README.rst](https://github.com/docopt/docopts/blob/master/old_README.rst))

## Features

The current command line tools `docopts`  written in python is maintained. A new
shell lib is added. 

The `docopts.sh` is a bash library you need to source in your CLI script. It
automaticaly embed docopts.py and docopt.py, and is standalone. Just drop it
and source it.

Of course, it needs a python interperter in the $PATH.

See examples/ for details.

## Developpers

If you want to clone this repository and hack docopts.

Use `git clone --recursive`, to get submodules.

If you forgot, you can also run:

~~~bash
git submodule init
git submodule update
~~~

Folder structure:

~~~
.
├── API_proposal.md - doc See wiki - to be removed
├── build.sh        - build the embedded docopts.py into docopts.sh
├── docopt.py       - original copy of docopt.py
├── docopts         - current python wrapper - almost unmodified
├── docopts.py      - copy of docopts, See build.sh
├── docopts.sh      - bash lib - already embed both docopt.py + docopts.py
├── examples
│   ├── calculator_example.sh
│   ├── cat-n_wrapper_example.sh
│   ├── docopts_auto_examples.sh
│   └── quick_example.sh
├── language_agnostic_tester.py
├── setup.py
├── testcases.docopt
├── testee.sh
└── tests
    ├── bats/                - git submodules
    ├── bats.alias           - source it to have bats working
    ├── docopts-auto.bats    - unit test for --auto
    ├── docopts.bats         - unit test docopts.sh
    └── exit_handler.sh      - helper
~~~
