// https://github.com/golang/go/issues/25922#issuecomment-1038394599
//go:build tools
// +build tools

package main

import (
	_ "github.com/goreleaser/goreleaser"
)
