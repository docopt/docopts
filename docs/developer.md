## Developers

All python related stuff has been removed, excepted `language_agnostic_tester.py`.

## Hacking docopts

In order to  hack `docopts` you will need:

- A working [Go developper environment](https://golang.org/doc/code.html#Organization)
- [`bats-core`](https://github.com/bats-core/bats-core) in your PATH. (optional: for testing with `make test`)
- GNU awk (for `build_doc.sh`)

Fetch all dependancies:

```
make install_builddep
```

## Current folder structure

It should looks like:

```
.
├── docopts.go                   - main source code
├── docopts_test.go              - go unit tests
├── docopts.sh                   - library wrapper and helpers
├── examples                     - many ported examples in bash, all must be working
├── language_agnostic_tester.py  - old python JSON tester still used with testee.sh
├── LICENSE-MIT                  - original docopts license
├── README.md                    - our main documentation
├── testcases.docopt             - agnostic testcases copied from python's docopt
├── testee.sh                    - bash wrapper to convert docopts output to JSON (now uses docopts.sh)
├── tests                        - unit and functional testing written in bats
└── TODO.md                      - Some todo list on this golang version of docopts
[...]                            - other file are helpers or current hack, not documented
```

## Tests

Some tests are coded along with this code base:

* bats - bash unit tests and functional testing.
* `language_agnostic_tester.py` - old python wrapper, full docopt compatibility tests.
* See also: [docopt.go](https://github.com/docopt/docopt.go) has its own tests in golang.
* `docopts_test.go` - go unit test for `docopts.go`

### Running tests

```
make test
```

#### bats

Bats is a unittest / functional testing framework for bash.

All our tests expect to run from `./tests/` directory and the `docopts` binary in `../docopts`.
Ensure bats is in your PATH.

```bash
cd ./tests
bats .
```

#### `language_agnostic_tester`

This script was provided with the original `docopts`. I fixed number/string output parsing failure with an extra function
for bash in [docopts.sh](https://github.com/docopt/docopts/blob/13f0bbcaba5c92deba909139b92fbbf3d768ea1b/docopts.sh#L144-L151)
`docopt_get_raw_value()`. This is a hack to get 100% pass, and it is not very efficient.

Run these tests from top of repo:
```
python language_agnostic_tester.py ./testee.sh
```

#### golang docopt.go (golang parser lib)

This lib is outside this project, but it is the base of the `docopt` language parsing for this wrapper.

```
cd $GOPATH/src/github.com/docopt/docopt-go
go test -v .
```

#### golang docopts (our bash wrapper)

```
cd $GOPATH/src/github.com/docopt/docopts
go test -v
```
