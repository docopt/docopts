# Some docker container with docopts installed

No special purpose than to recreate an environment with `docopts`

## Build the docker image

The `--build-arg` is used to pass teh version of `docopts` to download pre-built binary.
If you're publishing `docopts` the `../VERSION` may not yet reflect a published version.
You can find previously published VERSION in `../tests/VERSION`.

```
docker build -f debian-docopts.Dockerfile --build-arg VERSION=$(cat ./VERSION) -t debian-docopts .
```

## Run interactive

You will have latest python version 0.6.1 in PATH as `docopts` and Go version in PATH too as `docopts0`

Golang work space will also be installed if you want to test something.

```
docker run -it debian-docopts:latest
```
