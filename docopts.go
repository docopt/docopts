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
    "io"
    "io/ioutil"
    "sort"
)

var Version string = "docopts 0.6.3"
var Usage string = `Shell interface for docopt, the CLI description language.

Usage:
  docopts [options] -h <msg> : [<argv>...]
  docopts [options] [--no-declare] -A <name>   -h <msg> : [<argv>...]
  docopts [options] -G <prefix> -h <msg> : [<argv>...]
  docopts [options] --no-mangle -h <msg> : [<argv>...]

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
  -s <str>, --separator=<str>   The string to use to separate the help message
                                from the version message when both are given
                                via standard input. [default: ----]
  -O, --options-first           Disallow interspersing options and positional
                                arguments: all arguments starting from the
                                first one that does not begin with a dash will
                                be treated as positional arguments.
  -H, --no-help                 Don't handle --help and --version specially.
  -A <name>                     Export the arguments as a Bash 4.x associative
                                array called <name>.
  -G <prefix>                   As without -A, but outputs Bash compatible
                                GLOBAL varibles assignment, uses the given
                                <prefix>_{option}={parsed_option}. Can be used
                                with numerical incompatible option as well.
                                See also: --no-mangle
  --no-mangle                   Output parsed option not suitable for bash eval.
                                As without -A but full option names are kept.
                                Rvalue is still shellquoted.
  --no-declare                  Don't output 'declare -A <name>', used only
                                with -A argument.
  --debug                       Output extra parsing information for debuging.
                                Output cannot be used in bash eval.

Copyright (C) 2013 Vladimir Keleshev, Lari Rasku.
License MIT <http://opensource.org/licenses/MIT>.
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law.
`

// testing trick, out can be mocked to catch stdout and validate
// https://stackoverflow.com/questions/34462355/how-to-deal-with-the-fmt-golang-library-package-for-cli-testing
var out io.Writer = os.Stdout

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

// store global behavior to not pass to method as optional arguments
type Docopts struct {
    Global_prefix string
    Mangle_key bool
    Output_declare bool
}

// output assoc array output
func (d *Docopts) Print_bash_args(bash_assoc string, args docopt.Opts) {
    // fake nested Bash arrays for repeatable arguments with values
    // structure is:
    // bash_assoc[key,#]=length
    // bash_assoc[key,i]=value
    // i is an integer from 0 to length-1

    if d.Output_declare {
        fmt.Fprintf(out, "declare -A %s\n" ,bash_assoc)
    }

    for key, value := range args {
        // some golang tricks here using reflection to loop over the map[]
        rt := reflect.TypeOf(value)
        if IsArray(rt) {
            // all array is outputed even 0 size
            val_arr := value.([]string)
            for index, v := range val_arr {
                fmt.Fprintf(out, "%s['%s,%d']=%s\n", bash_assoc, Shellquote(key), index, To_bash(v))
            }
            // size of the array
            fmt.Fprintf(out, "%s['%s,#']=%d\n", bash_assoc, Shellquote(key), len(val_arr))
        } else {
            // value is not an array
            fmt.Fprintf(out, "%s['%s']=%s\n", bash_assoc, Shellquote(key), To_bash(value))
        }
    }
}

// test if a value is an array
func IsArray(rt reflect.Type) bool {
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

func Shellquote(s string) string {
    return strings.Replace(s, "'", `'\''`, -1)
}

func IsBashIdentifier(s string) bool {
    identifier := regexp.MustCompile(`^([A-Za-z]|[A-Za-z_][0-9A-Za-z_]+)$`)
    return identifier.MatchString(s)
}

// convert a parsed type to a text string suitable for bash eval
// as a right-hand side of an assignment
func To_bash(v interface{}) string {
    var s string
    switch v.(type) {
    case bool:
        s = fmt.Sprintf("%v", v.(bool))
    case int:
        s = fmt.Sprintf("%d", v.(int))
    case string:
        s = fmt.Sprintf("'%s'", Shellquote(v.(string)))
    case []string:
        // escape all strings
        arr := v.([]string)
        arr_out := make([]string, len(arr))
        for i, e := range arr {
            arr_out[i] = Shellquote(e)
        }
        s = fmt.Sprintf("('%s')", strings.Join(arr_out[:],"' '"))
    case nil:
        s = ""
    default:
        panic(fmt.Sprintf("To_bash():unsuported type: %v for '%v'", reflect.TypeOf(v), v ))
    }

    return s
}

func (d *Docopts) Print_bash_global(args docopt.Opts) {
    var new_name string
    var err error
    var out_buf string

    // value is an interface{}
    for key, value := range args {
        if d.Mangle_key {
            new_name, err = d.Name_mangle(key)
            if err != nil {
                docopts_error("%v", err)
            }
        } else {
            new_name = key
        }

        out_buf += fmt.Sprintf("%s=%s\n", new_name, To_bash(value))
    }

    // final output
    fmt.Fprintf(out, "%s", out_buf)
}

func (d *Docopts) Name_mangle(elem string) (string, error) {
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

    // alter output if we have a prefix
    key_fmt := "%s"
    if d.Global_prefix != "" {
        key_fmt = fmt.Sprintf("%s_%%s", d.Global_prefix)
    }

    v = fmt.Sprintf(key_fmt, strings.Replace(v, "-", "_", -1))

    if ! IsBashIdentifier(v) {
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
        fmt.Printf("echo 'error: %s\n%s' >&2\nexit 64\n", Shellquote(err.Error()), Shellquote(usage))
        os.Exit(1)
    } else {
        // --help or --version found and --no-help was not given
        fmt.Printf("echo '%s'\nexit 0\n",  Shellquote(usage))
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
            // no arg at all, display small usage, also exits 1
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

func docopts_error(msg string, err error) {
    if err != nil {
        msg = fmt.Sprintf(msg, err)
    }
    fmt.Fprintf(os.Stderr, "docopts:error: %s\n", msg)
    os.Exit(1)
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

    // create our Docopts struct
    d := &Docopts{
        Global_prefix: "",
        Mangle_key: true,
        Output_declare: true,
    }

    // parse docopts's own arguments
    argv := arguments["<argv>"].([]string)
    doc := arguments["--help"].(string)
    bash_version, _ := arguments.String("--version")
    options_first := arguments["--options-first"].(bool)
    no_help :=  arguments["--no-help"].(bool)
    separator := arguments["--separator"].(string)
    d.Mangle_key = ! arguments["--no-mangle"].(bool)
    d.Output_declare = ! arguments["--no-declare"].(bool)
    global_prefix, err := arguments.String("-G")
    if err == nil {
        d.Global_prefix = global_prefix
    }

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
            if ! IsBashIdentifier(name) {
                fmt.Printf("-A: not a valid Bash identifier: '%s'", name)
                return
            }
            d.Print_bash_args(name, bash_args)
        } else {
            d.Print_bash_global(bash_args)
        }
    } else {
        panic(err)
    }
}
