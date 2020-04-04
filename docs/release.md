# Creating a docopts release

Step to produce a release.

## 1. prepare the release

copy the previous ./VERSION to ./tests/VERSION

```
cp ./VERSION ./tests/VERSION
git add ./tests/VERSION
```

Increment the ./VERSION number

edit ./deployment.yml

- add the new tagget release by copying ./VERSION
- the description you want, changes etc.

## 2. rebuild all

```
make clean
make test
```

## 3. rebuild README.md

Ensure README.md is commited and rebuild it too

```
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

## 6. tag the new release

```
git tag -a $(cat ./VERSION) -m "docopts $(cat ./VERSION)"
```

## 7. push on docopts for a travis build

```
git push docopts master
```

## 8. build the release

With [deploy.sh](https://github.com/opensource-expert/deploy.sh) installed in our PATH

```
deploy.sh build
```

## 9. deploy the release

load github env

```
deploy.sh deploy
```


