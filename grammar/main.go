package main

import (
	"github.com/alecthomas/participle/lexer"
  "fmt"
  "strings"
)


var Usage string = `Shell interface for docopt, the CLI description language.

Usage:
  docopts [options] -h <msg> : [<argv>...]
  docopts [options] [--no-declare] -A <name>   -h <msg> : [<argv>...]
  docopts [options] -G <prefix>  -h <msg> : [<argv>...]
  docopts [options] --no-mangle  -h <msg> : [<argv>...]

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
  -G <prefix>                   Don't use associative array but output
                                Bash 3.2 compatible GLOBAL variables assignment:
                                  <prefix>_{mangled_args}={parsed_value}
                                Can be used with numeric incompatible options
                                as well.  See also: --no-mangle
  --no-mangle                   Output parsed option not suitable for bash eval.
                                Full option names are kept. Rvalue is still
                                shellquoted. Extra parsing is required.
  --no-declare                  Don't output 'declare -A <name>', used only
                                with -A argument.
  --debug                       Output extra parsing information for debugging.
                                Output cannot be used in bash eval.
`

// A custom lexer for docopt input
var doctop_Lexer = lexer.Must(lexer.Regexp(
    `(?P<NewLine>\n)` +
  // catch multiple blank
    `|(?P<LongBlank>\s{2,})` +
  // skip single blank
		`|(\s)` +
    `|(?P<Section>^(Usage|Options):)` +
    `|(?P<Option>-[A-Za-z0-9]|--[A-Za-z][A-Za-z0-9_-]+)` +
    `|(?P<Keyword>default:)` +
		`|(?P<Argument><[a-z]+>|[A-Z_]{2,})` +
		`|(?P<Punct>\]|[[=,]|\.{3})` +
		`|(?P<String>[^ \t\n\]=,[]+)`,
))

func main() {
  // extract symbols from the Lexer
  sym := lexer.SymbolsByRune(doctop_Lexer)
  // initialize the Lexer with a string
	lex, _ := doctop_Lexer.Lex(strings.NewReader(Usage))
  // extract all token
  tok, _ := lexer.ConsumeAll(lex)
  // display all tokens with their type
  for _, v := range tok {
    fmt.Printf("Token{%s, %q}\n", sym[v.Type], v.Value)
  }
}
