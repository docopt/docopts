// vim: set ts=4 sw=4 sts=4 et:
//
// docopts.go is a command line wrapper for docopt.go to be used by bash scripts.
//
package main

import (
	"fmt"
	"github.com/docopt/docopt-go"
    "regexp"
    "strings"
	"reflect"
)

var Version string = "docopts 0.6.3"
var Usage string = `Shell interface for docopt, the CLI description language.

Usage:
  docopts [options] --help-mesg=<msg> -- [<argv>...]

Options:
  --help-mesg=<msg>             The help message in docopt format.
                                If - is given, read the help message from
                                standard input.
                                If no argument is given, print docopts's own
                                help message and quit.
  -V <msg>, --version=<msg>     A version message.
                                If - is given, read the version message from
                                standard input.  If the help message is also
                                read from standard input, it is read first.
                                If no argument is given, print docopts's own
                                version message and quit.
  -O, --options-first           Disallow interspersing options and positional
                                arguments: all arguments starting from the
                                first one that does not begin with a dash will
                                be treated as positional arguments.
  -H, --no-help                 Don't handle --help and --version specially.
  -A <name>                     Export the arguments as a Bash 4.x associative
                                array called <name>.
  -s <str>, --separator=<str>   The string to use to separate the help message
                                from the version message when both are given
                                via standard input. [default: ----]

Copyright (C) 2013 Vladimir Keleshev, Lari Rasku.
License MIT <http://opensource.org/licenses/MIT>.
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law.
`

// debug helper
func print_args(args docopt.Opts) {
    for key, value := range args {
        fmt.Printf("%20s : %s\n", key, value)
	}
	fmt.Println("----------------------------------------")
}

// output assoc array output
func print_bash_args(bash_assoc string, args docopt.Opts) {
    // fake nested Bash arrays for repeatable arguments with values
    // structure is:
    // bash_assoc[key,#]=length
    // bash_assoc[key,i]=value
    // i is an integer from 0 to length-1

    fmt.Printf("declare -A %s\n" ,bash_assoc)

    for key, value := range args {
        // some golang tricks here using reflection to loop over the map[]
        rt := reflect.TypeOf(value)
		if isArray(rt) {
            val_arr := value.([]string)
            switch len(val_arr) {
            case 0:
                // empty
                fmt.Printf("%s['%s']=''\n", bash_assoc, shellquote(key))
            case 1:
                fmt.Printf("%s['%s']='%s'\n", bash_assoc, shellquote(key), to_bash(val_arr[0]))
            default:
                for index, v := range val_arr {
                    fmt.Printf("%s['%s,%d']='%s'\n", bash_assoc, shellquote(key), index, to_bash(v))
                }
                // size of the array
                fmt.Printf("%s['%s,#']=%d\n", bash_assoc, shellquote(key), len(val_arr))
            }
        } else {
            // value is not an array
            fmt.Printf("%s['%s']='%s'\n", bash_assoc, shellquote(key), value)
        }
	}
}

// test if a value is an array
func isArray(rt reflect.Type) bool {
    if rt == nil {
        return false
    }
	switch rt.Kind() {
	case reflect.Slice:
		return true
	case reflect.Array:
		return true
	default:
		return false
	}
}

func shellquote(s string) string {
    return strings.Replace(s, "'", `'\''`, -1)
}

func isbashidentifier(s string) bool {
    identifier := regexp.MustCompile(`^([A-Za-z]|[A-Za-z_][0-9A-Za-z_]+)$`)
    return identifier.MatchString(s)
}

func to_bash(v interface{}) string {
    var s string
    switch v.(type) {
        case bool:
            s = fmt.Sprintf("%b", v.(bool))
        case int:
            s = fmt.Sprintf("%d", v.(int))
        case string:
            s = shellquote(v.(string))
        case []string:
            // escape all strings
            arr := v.([]string)
            arr_out := make([]string, len(arr))
            for i, e := range arr {
                arr_out[i] = shellquote(e)
            }
            s = fmt.Sprintf("('%s')", strings.Join(arr_out[:],"', '"))
        default:
            panic("unsported type")
    }

    return s
}

func print_bash_global(args docopt.Opts) {
    for key, value := range args {
        new_name, err := name_mangle(key)
        if err == nil {
            fmt.Printf("%s='%s'\n", new_name, to_bash(value))
        }
	}
}

func name_mangle(elem string) (string, error) {
    var v string

    if elem == "-" || elem == "--" {
        return "", fmt.Errorf("not supported")
    }

    if Match(`^<.*>$`, elem) {
        v = elem[1:len(v)-1]
    } else if Match(`^-[^-]$`, elem) {
        v = fmt.Sprintf("%c", elem[1])
    } else if Match(`^--.+$`, elem) {
        v = elem[2:]
    } else {
        v = elem
    }

    v = strings.Replace(v, "-", "_", -1)

    if ! isbashidentifier(v) {
        return "", fmt.Errorf("not bash identifier: %s", elem)
    }

    return v, nil
}

func Match(regex string, source string) bool {
	matched, _ := regexp.MatchString(regex, source)
	return matched
}

func main() {
	arguments, err := docopt.ParseArgs(Usage, nil, Version)
    if err == nil {
        print_args(arguments)
    } else {
        fmt.Println(err)
        return
    }

    // parse docopts's own arguments

    argv := arguments["<argv>"].([]string)
    doc := arguments["--help-mesg"].(string)
    bash_version, _ := arguments.String("--version")

    // available in go, with "OptionsFirst: true,"
    //
    // options_first := arguments["--options-first"].(bool)
    //no_help := ! arguments["--no-help"].(bool)
    // separator := arguments["--separator"].(string)

    // read from stdin
    //if doc == "-" and version == "-":
    //    doc, version = (page.strip() for page in
    //                    sys.stdin.read().split(separator, 1))
    //elif doc == "-":
    //    doc = sys.stdin.read().strip()
    //elif version == "-":
    //    version = sys.stdin.read().strip()

    bash_args, err := docopt.ParseArgs(doc, argv, bash_version)
    if err == nil {
        print_args(bash_args)
    }

    name, err := arguments.String("-A")
    if err == nil {
        if ! isbashidentifier(name) {
            fmt.Printf("-A: not a valid Bash identifier: %s", name)
            return
        }

        print_bash_args(name, bash_args)
    } else {
        print_bash_global(bash_args)
    }
}



