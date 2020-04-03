# current work in PROGRRESS

## correct macos arch geting 32bits

https://github.com/docopt/docopts/issues/44
OK - done

* speedup travis build on macos
  * provide our own bash?
  * caching?

## install bats as a travis dependancies

- in .travis.yml - OK
- remove bats subproject
- update developpers documentation

actually bats is a submodule.
and it's run from `make test`

```
cd tests/ && ./bats/bin/bats .
```

## grammar programming

Testing https://github.com/alecthomas/participle

## provide test on old environment

docker?
32bits
bash 3

