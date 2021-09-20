package main

import (
	"fmt"
	//	"github.com/alecthomas/repr"
	"github.com/docopt/docopts/docopt-go"
	"os"
)

func main() {
	section := os.Args[1]
	filename := os.Args[2]
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("error: fail to open file: %s\n", filename)
		return
	} else {
		fmt.Printf("parsing: %s\n", filename)
	}

	doc := string(data)
	// extract the Usage: section (case-insensitive) until the next empty line
	usageSections := docopt.ParseSection(section, doc)
	for i, s := range usageSections {
		fmt.Printf("%d: %s\n", i, s)
	}
}
