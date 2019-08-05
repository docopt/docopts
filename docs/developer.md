## Developers

All python related stuff has been removed, excepted `language_agnostic_tester.py`.

If you want to clone this repository and hack `docopts`:

Use `git clone --recursive`, to get submodules - these are only required for testing with `bats`.
If [`bats-core`](https://github.com/bats-core/bats-core) installed in your PATH should work too.

Fetch the extra golang version of `docopt-go` (required for building `docopts`)

```
go get github.com/docopt/docopt-go
```

If you forgot `--recursive`, you can also run it afterwards:

~~~bash
git submodule init
git submodule update
~~~

Current folder structure:

~~~
.
├── docopts.go - main source code
├── docopts_test.go - go unit tests
├── docopts.sh - library wrapper and helpers
├── examples - many ported examples in bash, all must be working
├── language_agnostic_tester.py - old python JSON tester still used with testee.sh
├── LICENSE-MIT
├── PROGRESS.md - what I'm working on
├── README.md
├── testcases.docopt - agnostic testcases copied from python's docopt
├── testee.sh - bash wrapper to convert docopts output to JSON (now uses docopts.sh)
├── tests - unit and functional testing written in bats (requires submodule)
└── TODO.md - Some todo list on this golang version of docopts
~~~

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

`bats.alias` modify your current environment to define a alias on the submodule of `bats` installed (if you did it).

```bash
cd ./tests
. bats.alias
bats .
```

#### `language_agnostic_tester`

This script was provided with the original `docopts`. I fixed number/string output parsing failure with an extra function
for bash in [docopts.sh](https://github.com/docopt/docopts/blob/docopts-go/docopts.sh#L108)
`docopt_get_raw_value()`. This is a hack to get 100% pass, and it is not very efficient.

Run these tests from top of repo:
```
python language_agnostic_tester.py ./testee.sh
```

#### golang docopt.go (golang parser lib)

This lib is outside this project, but it is the base of the `docopt` language parsing for this wrapper.

```
cd PATH/to/go/src/github.com/docopt/docopt-go/
go test -v .
```

#### golang docopts (our bash wrapper)

```
cd PATH/to/go/src/github.com/docopt/docopts
go test -v
```
