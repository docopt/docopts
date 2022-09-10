# docopts pre-built binaries

`docopts` is a shell helper mainly for bash as now, for parsing command-line arguments using the
[docopt language](https://docopt.org).

This implementation use Go, and we provide pre-built binaries available for download from
[github releases](https://github.com/docopt/docopts/releases).

The sha256sum is provided to ensure the uploaded binaries are conform to the
one built by the uploader.

## get the binary

You can simply download the one you want from [releases](https://github.com/docopt/docopts/releases).

Or we provide a command line helper if you cloned this repository, that will
simply download the good binary format from our repository and provide a `docopts` binary in
your current folder.

```
./get_docopts.sh
```

## Release binaries

This section is for developers. In order to release binaries you will need
some granted access to github API.

You will also need some more developer tools.

All dependencies are installed with:

```
make install_builddep
```

### `goreleaser`

[`goreleaser`](https://github.com/goreleaser/goreleaser) is used to cross-compile binaries for different platforms as well publish GitHub releases.

See the [`.goreleaser.yml`](../goreleaser.yml) for the configuration.

### github API token

You will need a valid github token for the target repository.

https://help.github.com/articles/creating-an-access-token-for-command-line-use

The token needs to have `repo` auth privileges.

Then export it as a bash environment variable:

```
export GITHUB_TOKEN="your token here"
```

### `git tag` a new release

We use [semantic verion tags](https://semver.org/)

Our version is stored in `VERSION` file.

```
echo "v0.6.3-alpha2" > VERSION
git tag -a "$(cat VERSION)" -m "golang 2019"
git push origin "$(cat VERSION)"
```

### deployment steps

In `docopts` project folder.

You need to create the release text in `CHANGELOG.md` first.

Dry-run to create release files and binaries in the `dist/` dir:

```
make snapshot
```

Publish the binaries and release notes:

```
make release
```
