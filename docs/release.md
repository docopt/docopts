# Creating a docopts release

Step to produce a release.

## 1. prepare the release

copy the previous `./VERSION` to `./tests/VERSION`

```
cp ./VERSION ./tests/VERSION
git add ./tests/VERSION
```

Increment the `./VERSION` number

edit `./deployment.yml`

- add the new tag name release matching ./VERSION
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
git commit -a
```

## 5. push on origin for a travis build

```
git push origin master
```

## 6. remove travis hack macos if any

```
sed -i -e '/travis.reverse_ssh_tunnel.sh/ s/^\([^#]\)/#\1/' .travis.yml
```

## 7. push on docopts for a travis build

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

With [deploy.sh](https://github.com/opensource-expert/deploy.sh) installed in our PATH


```
deploy.sh build
```

## 10. push the tag on docopts

```
git push docopts $(cat ./VERSION)
```

## 11. deploy the release

load github env

```
. ./env
export GITHUB_USER=docopt
export GITHUB_REPO=docopts
deploy.sh deploy
```
