## during grammar design

Missing: Which state parser we where (

```
go run docopt_grammar.go docopt_language.docopt
```

Output

```
  Options: &main.Options{
    Pos: Position{Filename: "docopt_language.docopt", Offset: 1153, Line: 22, Column: 1},
    Options_lines: []main.Options_line{
      main.Options_line{
        Pos: Position{Filename: "docopt_language.docopt", Offset: 1162, Line: 23, Column: 1},
        Option_def: main.Option_def{
          Options: []main.Option{
            main.Option{
              Pos: Position{Filename: "docopt_language.docopt", Offset: 1165, Line: 23, Column: 4},
              Short: &"-h",
            },
            main.Option{
              Pos: Position{Filename: "docopt_language.docopt", Offset: 1169, Line: 23, Column: 8},
              Long: &"--help",
            },
          },
        },
        Option_doc: main.Option_doc{
          Option_doc: "Indent is important before options, and between option and help message.",
        },
        Option_default: main.Option_default{
        },
      },
    },
  },
}
```

error


```
docopt_language.docopt:24:3: unexpected "--optional-long-argument" (expected <line_of_text>)
```


## Lexer get all token first

All token are eaten first which report some invalid token at this phase (develpement phase only?)


```
sylvain@lap42:~/code/go/src/github.com/docopt/docopts/grammar$ go run docopt_grammar.go docopt_language.docopt
parsing: docopt_language.docopt
Parse error
======================= partial AST ==========================
&main.Docopt{
}
================ end of partial AST ==========================
docopt_language.docopt:26:68: invalid token '['
```


