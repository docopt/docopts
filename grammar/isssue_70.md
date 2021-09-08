Hi,

I'm a bit confused, `Doc()` would not be called by the parser, as `Capture()` would, right?

I suppose, your hint represents a method that I could call after parsing, in order to rebuild the Documentation string, right?

So in this spirit, I coded [this version](https://github.com/docopt/docopts/commit/590e559fdc0f9900f96dd9095ffc77363575ece6#diff-2991ad015b6cdc2ba88467b1806edfa4)

Also, I have the opportunity to not match `"[", "default: ", "value_here_is_parsed", "]", "."` during parsing, but to match it during a post parsing phase.

So I [tested this way too](https://github.com/docopt/docopts/commit/3fce67be4932d25653de5c75b3f0de468655c91b), and rewrote the struct with an unassigned pointer:

```go
type Options_line struct {
	Pos lexer.Position

	Option_def     Option_def `( LONG_BLANK @@`
	Option_doc     Option_doc `  @@`
	Option_default *string
	Comment        []string   `| @( LINE_OF_TEXT "\n" | "\n"+ ) )`
}
```

Removed all token extraction from the lexer and the grammar and call:

```go
if err = parser.Parse(f, ast); err == nil {
  ast.Options.Post_parsing_Options_extract_default()
  // [...]
```

Which ends with an ast I prefere:

```go
[...]
      main.Options_line{
        Option_def: main.Option_def{
          Options: []main.Option{
            main.Option{
              Long: &"--which-support",
              Argument: &"<argument>",
            },
          },
        },
        Option_doc: main.Option_doc{
          Option_doc: "The <argument> for this option has a [default: value_here_is_parsed].",
        },
        Option_default: &"value_here_is_parsed",
      },
```


