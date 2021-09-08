
I've a similar problem for sub-something, that could help us  is it `Capture()` or `Parseable()` role?

Here is an input ended by an extra dot ".", it will complexify my grammar if I try to parse this corner
case:

```
  --which-support=<argument>  The <argument> for this option has a [default: value_here_is_parsed].
```

My actual lexer is extracting it as expected:

```
LONG_BLANK, "  "
LONG, "--which-support"
PUNCT, "="
ARGUMENT, "<argument>"
LONG_BLANK, "  "
LINE_OF_TEXT, "The <argument> for this option has a "
PUNCT, "["
DEFAULT, "default: "
LINE_OF_TEXT, "value_here_is_parsed"
PUNCT, "]"
".", LINE_OF_TEXT
```

My idea here, is to remove the `"[",  "default:" ,  "value", "]"` token extraction from the lexer
and to delegate the `default:` extraction during parsing with a simple regexp on `LINE_OF_TEXT`, in order to simplify the grammar.

How can I extract the `value_here_is_parsed` and put it back in the AST?



