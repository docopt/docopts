# docopts pre-built binaries

`docopts` is a shell helper mainly for bash as now, for parsing command-line arguments using the
[docopt language](https://docopt.org).

This implementation use Go, and we provide pre-built binaries on available for download from
github releases.

The sha256sum is provided to ensure the uploaded binaries are conform to the
one built by the uploader.

## get the binary

You can simply download the one you want from
[releases](https://github.com/docopt/docopts/releases).

Or we provide a command line helper if you cloned this repository, that will
simply download the good binary format from our repository and provide a `docopts` binary in
your current folder.


```
./get_docopts.sh
```

## Release binaries

This section is for developper. In order to release binaries you will need
some granted access to github API.

You will also need some more developper tools.

Most of the tools require a working [Go developper
environment](https://golang.org/doc/code.html#Organization). Which should not be too
complicated to setup.

Go for it:

### gox

We provide binaries cross-compiled from GNU/Linux using [gox](https://github.com/mitchellh/gox)

```
go get github.com/mitchellh/gox
```

### github release uploader

We publish release using a command line tools and using github API [gothub](https://github.com/itchio/gothub)

```
go get github.com/itchio/gothub
```

### github API token

You will need a valid gitub token for the target repository.

https://help.github.com/articles/creating-an-access-token-for-command-line-use

The token needs to have `repos` auth priviledges.

then export it as a bash environment variable:

```
export GITHUB_TOKEN="you token here"
```

### git tag a new release

We use [semantic verion tags](https://semver.org/)

```
git tag -a v0.6.3-alpha2 -m "golang 2019"
git push origin v0.6.3-alpha2
```

### yaml command-line tool

See: http://mikefarah.github.io/yq/

For extracting yaml data from `deployment.yml`

```
go get gopkg.in/mikefarah/yq.v2
```

### build docopts

Our script uses `docopts` for parsing our command line option.

```
make
```

### deploy

We provide a deploy script, which will take the last git tag, and a deployment
message written in a yaml file `deployment.yml`.

So you need to create the release text in `deployment.yml` before you run
`deploy.sh`.

See what will going on:

```
./deploy.sh deploy -n
```

Deploy and replace existing binaries for this release.

```
./deploy.sh deploy --replace
```

