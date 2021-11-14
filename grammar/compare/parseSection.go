package main

import (
	"fmt"
	//	"github.com/alecthomas/repr"
	"github.com/docopt/docopts/docopt-go"
	"os"
)

func main() {
	call_method := os.Args[1]
	argument := os.Args[2]
	filename := os.Args[3]
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("error: fail to open file: %s\n", filename)
		return
	} else {
		fmt.Printf("parsing: %s\n", filename)
	}

	doc := string(data)

	switch call_method {
	case "ParseSection":
		// extract the Usage: section (case-insensitive) until the next empty line
		usageSections := docopt.ParseSection(argument, doc)
		for i, s := range usageSections {
			fmt.Printf("%d: %s\n", i, s)
		}
	case "FormalUsage":
		usageSections := docopt.ParseSection(argument, doc)
		formal, err := docopt.FormalUsage(usageSections[0])
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("extrated Usage:\n%s\n", usageSections[0])
		fmt.Printf("FormalUsage: %s\n", formal)
	default:
		fmt.Printf("unknown method: %s\n", call_method)
	}

}
