# Developers documenation

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

I rewrote some part of the script so it's now compatible with python3.

This script goal is to parse `testcases.docopt` input format and to send to `testee.sh` to run against our `docopts`
version. The result of `testee.sh` is JSON output of the test case ran with `docopts`.

This script was provided with the original `docopts`. I fixed number/string output parsing failure with an extra function
for bash in [docopts.sh](https://github.com/docopt/docopts/blob/13f0bbcaba5c92deba909139b92fbbf3d768ea1b/docopts.sh#L144-L151)
`docopt_get_raw_value()`. This is a hack to get 100% pass, and it is not very efficient.

This could have been a python's only code without `testee.sh` piping result.

Run these tests from top of repo:
```
time python3 language_agnostic_tester.py ./testee.sh
real  0m3,711s
```

Or a single test by its number (See: [testee.sh top comment documentation](../testee.sh)).

```
python3 language_agnostic_tester.py ./testee.sh 176
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

## Add more test to docopts use-case

### bats (functionnal testing) (bash)

Most `.bats` files in `tests/` folder are test for bash code `docopts.sh` functions,
and are more internal stuff, hack as you wish but it may be more complicated to
handle at first.

Functional tests are in [`tests/functional_tests_docopts.bats`](tests/functional_tests_docopts.bats).

This this `bats` syntax test, this means almost bash.

Add a new functional test:

```
@test "your comment describing the test here" {
    # run is a bats helper that run a command catching all
    # use the $DOCOPTS_BIN to ensure to point to the good binary
    # pass any docopts argument to test, good or wrong
    run $DOCOPTS_BIN -G ARGS -h 'Usage: prog dump [-]' : dump -

    # output is caught by bats (stdout + stderr)
    # and wont be displayed, unless a test fail: any non-0 bash code
    echo "$output"

    # test docopts return value 0 on success
    [[ $status -eq 0 ]]

    regexp='some bash regexp'
    [[ "$output" =~  $regexp ]]
}
```

run only functional test (`bats` in `PATH`):

``` bash
# test are ran from test folder only
cd tests
bats functional_tests_docopts.bats
```

### `docopts_test.go` (unit testing) (golang + JSON)

It uses standard `go test`. Some method are "exposed" with Capital name only for testing purpose.

JSON is not that easy to handle nor really human redable/editable and will probably be dropped.

Options loop test JSON: [`common_input_test.json`](../common_input_test.json)

This is our input file for testing options parsed recieved from docopt lib. We emulate docopt parsed options and give
them as input to our method.

JSON tests are read as a list of test: (others JSON keys should be ignored and are used as comment)
Our parser JSON parser / loader [`test_json_load.go`](../test_json_load/test_json_load.go), only used for testing.

```go
type TestString struct {
    Input                 map[string]interface{}
    Expect_args           []string
    Expect_global         []string
    Expect_global_prefix  []string // optional
}
```

So adding new test case:

add a new bloc of JSON object (dict / hash / map):
```json
  {
    "description" : "convert array into fake nested array",
    "input": {
      "FILE": [
        "pipo",
        "molo",
        "toto"
      ]
    },
    "expect_args": [
      "declare -A args",
      "args['FILE,0']='pipo'",
      "args['FILE,1']='molo'",
      "args['FILE,2']='toto'",
      "args['FILE,#']=3"
    ],
    "expect_global": [
      "FILE=('pipo' 'molo' 'toto')"
    ],
    "expect_global_prefix": [
      "ARGS_FILE=('pipo' 'molo' 'toto')"
    ]
  },
```

* `description` a comment which is ignored
* other extra JSON key that are not the 4 following will be ignored too.
* `input` correspond to the `map[string]interface{}` of docopt parsed options.
* `expect_args` the text rows of the associative array code for bash4 that is outputed by `Print_bash_args()` matched in order.
* `expect_global` the text definition of the bash global vars that is outputed by `Print_bash_global()` matched in order.
* `expect_global_prefix` [optional] if present will be used for testing `Mangle_key` + `Global_prefix` instead of [`rewrite_prefix("ARGS",)`](../docopts_test.go)
  So left hand values in `expect_global_prefix` the prefix must be `ARGS` + `_`.


### testcases.docopt (agnostic test universal to docopt parsing language)

This file is still avaible from python docopt lib original repository
too [testcases.docopt](https://github.com/docopt/docopt/blob/511d1c57b59cd2ed663a9f9e181b5160ce97e728/testcases.docopt)

This is the input file used by `language_agnostic_tester.py`, which is a middleware originaly written to read
`testcases.docopt` and to send it to ~> `testee.sh` ~> `docopts -A` ~> JSON ~> the result is the validated against the
embedded JSON expected result.

Input file format is historically as follow:

* Support comment `#`
* single test definition:
  * `r"""` introduce a new test
  * a docopt definition including `Usage:` (can be multiline or single line)
  * `"""` finish the docopt definition
  * one or more call introduced with `$` + keyword `prog` followed by argument to pass to the program
```
$ prog -a
{"-a": true}
```
  * followed with the exptected output in JSON format (single ligne) (no empty line between `prog` call and expected JSON)
  * `\n` newline separator if some other call are added for the same `Usage:` definition

## Golang debugger

Debugger is a must for any programming language. Go provides an extrenal debugger named [delve](https://github.com/go-delve/delve)

https://github.com/go-delve/delve/tree/master/Documentation

In order to debug `docopts` you obviously need to pass command line argument to our program:
Our arguments start after the first `--` which is `delve` argument stopper.

```
cd path/to/docopts/source
dlv debug -- -h 'Usage: prog dump [--] <unparsed_arguments>...' : dump -- -some -auto-approve FILENAME
```

The as usual in debugger:

```
# put break point in main function
b main.main
continue
n
s # for steping into the current function
# printing some value
p some_varialbles
```

The debugger will then magically bring you, step after step, to the bug!
Enjoy! and promote dubugger in every programming language and every programming course. :wink:

