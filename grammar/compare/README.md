# study docopt-go

This code is for studiying go lib docopt-go (modified as internal copy)

We Capilatize the lib method and call them with some source.

## `ParseSection(name, source string) []string`

Note: In vim use ctags to jump to function definition (pwd where tags file live)

Extract a section from the whole doc containing all the docopt source.

possible call look like:

```go
// extract Usage section (start a Usage: case insensitive end with an empty line)
usageSections := ParseSection("usage:", doc)

// same, extract options: section then loop on the result
[...] range ParseSection("options:", doc)
```

## `FormalUsage(section string) (string, error)`

Transforms extracted usage string into another string with group of alternative usages.

Partition `Usage: .*` on `:` in three partion and split usage lines on blank with `strings.Fields()`
progname is `pu[0]` and it's stripped from FormalUsage result. It serves as delimiter to each multiple
possible usage.

```go
formal, err := FormalUsage(usage)
```

## `TokenListFromPattern(source string) *tokenList`

called from `ParsePattern()`.
Return a tokenList pointer which contains the extacted token in form of a list on string.
composed option with named argument are not yet extracted (splited into token)

```go
tokens := TokenListFromPattern(source)
```

## `ParsePattern(source string, options *patternList) (*pattern, error)`

Parser entry point will call: TokenListFromPattern + ParseExpr

