package main

import (
	"fmt"
	"github.com/docopt/docopts/scanner"
	"github.com/docopt/docopts/token"
	"log"
	"os"
)

func main() {
	var s scanner.Scanner

	// error handler
	eh := func(_ token.Position, msg string) {
		fmt.Errorf("error handler called (msg = %s)", msg)
	}

	filname := os.Args[1]
	data, err := os.ReadFile(filname)
	// Open file for reading
	if err != nil {
		log.Fatal(err)
	}
	fset := token.NewFileSet()
	s.Init(fset.AddFile("", fset.Base(), len(data)), data, eh, scanner.ScanComments)
	for {
		pos, tok, lit := s.Scan()

		fmt.Printf("pos: %d, %s, %s\n", pos, tok.String(), lit)

		if tok == token.EOF {
			break
		}
	}
	fmt.Println("OK")
}
