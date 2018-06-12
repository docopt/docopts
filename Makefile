
DOCOPT_GO=${GOPATH}/linux_amd64/github.com/docopt/docopt-go.a

# build 64 bits version
docopts: docopts.go
	go build docopts.go

docopt-go:
	go get github.com/docopt/docopt-go

all: docopt-go docopts docopts-arm docopts-32bits docopts-OSX

# build 32 bits version too
docopts-32bits: docopts.go
	env GOOS=linux GOARCH=386 go build -o docopts-32bits docopts.go

# build for OSX
docopts-OSX: docopts.go
	env GOOS=darwin go build -o docopts-OSX docopts.go

# build 32 bits version too
docopts-arm: docopts.go
	env GOOS=linux GOARCH=arm go build -o docopts-arm docopts.go
