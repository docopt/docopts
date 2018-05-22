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
    "os"
    "io/ioutil"
    "sort"
)

var Version string = "docopts 0.6.3"
var Usage string = `Shell interface for docopt, the CLI description language.

Usage:
  docopts [options] -h <msg> : [<argv>...]

Options:
  -h <msg>, --help=<msg>        The help message in docopt format.
                                Without argument outputs this help.
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
  --no-mangle                   Output parsed option not suitable for bash eval.
                                As without -A but full option names are kept.
                                Rvalue is still shellquoted.
  --debug                       Output extra parsing information for debuging.
                                Output cannot be used in bash eval.

Copyright (C) 2013 Vladimir Keleshev, Lari Rasku.
License MIT <http://opensource.org/licenses/MIT>.
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law.
`

// debug helper
func print_args(args docopt.Opts, message string) {
    // sort keys
    mk := make([]string, len(args))
    i := 0
    for k, _ := range args {
        mk[i] = k
        i++
    }
    sort.Strings(mk)
    fmt.Printf("################## %s ##################\n", message)
    for _, key := range mk {
        fmt.Printf("%20s : %v\n", key, args[key])
    }
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
            //case 1:
            //    // quoting assignment is driven by to_bash()
            //    fmt.Printf("%s['%s']=%s\n", bash_assoc, shellquote(key), to_bash(val_arr[0]))
            default:
                for index, v := range val_arr {
                    fmt.Printf("%s['%s,%d']=%s\n", bash_assoc, shellquote(key), index, to_bash(v))
                }
                // size of the array
                fmt.Printf("%s['%s,#']=%d\n", bash_assoc, shellquote(key), len(val_arr))
            }
        } else {
            // value is not an array
            fmt.Printf("%s['%s']=%s\n", bash_assoc, shellquote(key), to_bash(value))
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
        s = fmt.Sprintf("%v", v.(bool))
    case int:
        s = fmt.Sprintf("%d", v.(int))
    case string:
        s = fmt.Sprintf("'%s'", shellquote(v.(string)))
    case []string:
        // escape all strings
        arr := v.([]string)
        arr_out := make([]string, len(arr))
        for i, e := range arr {
            arr_out[i] = shellquote(e)
        }
        s = fmt.Sprintf("('%s')", strings.Join(arr_out[:],"', '"))
    case nil:
        s = ""
    default:
        panic(fmt.Sprintf("to_bash():unsuported type: %v", reflect.TypeOf(v) ))
    }

    return s
}

func print_bash_global(args docopt.Opts, mangle_key bool) {
    var new_name string
    var err error

    for key, value := range args {
        if mangle_key {
            new_name, err = name_mangle(key)
            if err != nil {
                // skip
                return
            }
        } else {
            new_name = key
        }
        fmt.Printf("%s=%s\n", new_name, to_bash(value))
    }
}

func name_mangle(elem string) (string, error) {
    var v string

    if elem == "-" || elem == "--" {
        return "", fmt.Errorf("not supported")
    }

    if Match(`^<.*>$`, elem) {
        v = elem[1:len(elem)-1]
    } else if Match(`^-[^-]$`, elem) {
        v = fmt.Sprintf("%c", elem[1])
    } else if Match(`^--.+$`, elem) {
        v = elem[2:]
    } else {
        v = elem
    }

    v = strings.Replace(v, "-", "_", -1)

    if ! isbashidentifier(v) {
        return "", fmt.Errorf("cannot transform into a bash identifier: %s", elem)
    }

    return v, nil
}

// helper for lazy typing
func Match(regex string, source string) bool {
    matched, _ := regexp.MatchString(regex, source)
    return matched
}

// our HelpHandler which outputs bash code to be evaled as error and stop or
// display program's help or version
// TODO: handle return or kill instead of exit so it can be launched inside a function
var HelpHandler_for_bash_eval = func(err error, usage string) {
    if err != nil {
        fmt.Printf("echo 'error: %s\n%s' >&2\nexit 64\n", shellquote(err.Error()), shellquote(usage))
        os.Exit(1)
    } else {
        // --help or --version found and --no-help was not given
        fmt.Printf("echo '%s'\nexit 0\n",  shellquote(usage))
        os.Exit(0)
    }
}

var HelpHandler_golang = func(err error, usage string) {
    if err != nil {
        err_str := err.Error()
        if len(err_str) >= 9 {
            // we hack for our polymorphic argument -h or -V
            // it was the same hack in python version
            if err_str[0:2] == "-h" || err_str[0:6] == "--help" {
                // print full usage message (global var)
                fmt.Println(strings.TrimSpace(Usage))
                os.Exit(0)
            }
            if err_str[0:2] == "-V" || err_str[0:9] == "--version" {
                fmt.Println(strings.TrimSpace(Version))
                os.Exit(0)
            }
        }

        if len(err_str) == 0 {
            // no arg at all, display small usage, also exits
            HelpHandler_for_bash_eval(fmt.Errorf("no argument"), usage)
        }

        // real error
        fmt.Fprintf(os.Stderr, "my error: %v, %v\n", err, usage)
        os.Exit(1)
    } else {
        // no error, never reached?
        fmt.Println(usage)
        os.Exit(0)
    }
}

func main() {
    golang_parser := &docopt.Parser{
      OptionsFirst: true,
      SkipHelpFlags: true,
      HelpHandler: HelpHandler_golang,
    }
    arguments, err := golang_parser.ParseArgs(Usage, nil, Version)

    if err != nil {
        msg := fmt.Sprintf("mypanic: %v\n", err)
        panic(msg)
    }

    debug := arguments["--debug"].(bool)
    if debug {
        print_args(arguments, "golang")
    }

    // parse docopts's own arguments
    argv := arguments["<argv>"].([]string)
    doc := arguments["--help"].(string)
    bash_version, _ := arguments.String("--version")
    options_first := arguments["--options-first"].(bool)
    no_help :=  arguments["--no-help"].(bool)
    separator := arguments["--separator"].(string)
    mangle_key := ! arguments["--no-mangle"].(bool)

    // read from stdin
    if doc == "-" && bash_version == "-" {
        bytes, _ := ioutil.ReadAll(os.Stdin)
        arr := strings.Split(string(bytes), separator)
        if len(arr) == 2 {
            doc, bash_version = arr[0], arr[1]
        } else {
            msg := "error: help + version stdin, not found"
            if debug {
                msg += fmt.Sprintf("\nseparator is: '%s'\n", separator)
                msg += fmt.Sprintf("spliting has given %d blocs, exactly 2 are expected\n", len(arr))
            }
            panic(msg)
        }
    } else if doc == "-" {
        bytes, _ := ioutil.ReadAll(os.Stdin)
        doc = string(bytes)
    } else if bash_version == "-" {
        bytes, _ := ioutil.ReadAll(os.Stdin)
        bash_version = string(bytes)
    }

    doc = strings.TrimSpace(doc)
    bash_version = strings.TrimSpace(bash_version)
    if debug {
        fmt.Printf("%20s : %v\n", "doc", doc)
        fmt.Printf("%20s : %v\n", "bash_version", bash_version)
    }

    // now parse bash program's arguments
    parser := &docopt.Parser{
      HelpHandler: HelpHandler_for_bash_eval,
      OptionsFirst: options_first,
      SkipHelpFlags: no_help,
    }
    bash_args, err := parser.ParseArgs(doc, argv, bash_version)
    if err == nil {
        if debug {
            print_args(bash_args, "bash")
            fmt.Println("----------------------------------------")
        }
        name, err := arguments.String("-A")
        if err == nil {
            if ! isbashidentifier(name) {
                fmt.Printf("-A: not a valid Bash identifier: %s", name)
                return
            }
            print_bash_args(name, bash_args)
        } else {
            // TODO: add global prefix
            print_bash_global(bash_args, mangle_key)
        }
    } else {
        panic(err)
    }
}
