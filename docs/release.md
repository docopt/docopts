# Creating a docopts release

Step to produce a release.

## 1. prepare the release

copy the previous `./VERSION` to `./tests/VERSION`
This is for functional test `get_doctops.bats` so it can fetch a
published release version.

```
cp ./VERSION ./tests/VERSION
git add ./tests/VERSION
```

Increment the `./VERSION` number

edit `CHANGELOG.md`:
- add a new heading at the top of the file with `## <insert heading here>` (only the first `##` will be used)
- the description you want, changes etc.

## 2. rebuild all

```
make clean
make test
```

## 3. rebuild README.md

Ensure README.md is commited and rebuild it too

```
touch build_doc.sh
make README.md
```

## 4. commit all the code in master

```
git checkout master
# merge if needed
git commit -a
```

## 5. push on origin for a ci build

```
git push origin master
```

## 6. remove reverse ssh hack macos if any

```
sed -i -e '/ci.reverse_ssh_tunnel.sh/ s/^\([^#]\)/#\1/' .github/workflows/ci.yml
```

## 7. push on docopts for a ci build

```
git commit -a
git push docopts master
```

## 8. tag the new release

remove previous tag test if any

```
git push --delete origin  $(cat ./VERSION)
git tag -d $(cat ./VERSION)
```

create the tag

```
git tag -a $(cat ./VERSION) -m "docopts $(cat ./VERSION)"
git tag
```

## 9. build the release

The following will build all binary version from `.goreleaser.yml`

```
make snapshot
```

## 10. push the tag on docopts

```
git push docopts $(cat ./VERSION)
```

## 11. deploy the release

```
export GITHUB_TOKEN="your token here"
make release
```
