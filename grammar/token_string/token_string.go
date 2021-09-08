package token_string

// ================================ lexer definition for string sub-lexer ===========
var (
  State_Base = `
  (?P<NEWLINE>\n)
  # skip all blank
  |([\t ]+)
  |(?P<WORD>[A-Za-z][^ \t]+)
  |(?P<DQUOTE>") => State_String
  `

  State_String = `
  (?P<NEWLINE>\n)
  |(?P<DOLLAR_OPEN>\$\{) => State_Code_embbeded
  |(?P<ESCAPE_DQUOTE>\\")
  |(?P<DQUOTE>")
  |(?P<STRING_CHAR>.)
  `

  State_Code_embbeded = `
  (?P<NEWLINE>\n)
  |(?P<CLOSING_BRASSE>\}) => State_String
  # skip all blank
  |([\t ]+)
  |(?P<WORD>[A-Za-z][^ \t]+)
  |(?P<DQUOTE>") => State_String
  |([\t ]+)
  |(?P<SHORT>-[A-Za-z0-9?])
  |(?P<LONG>--[A-Za-z][A-Za-z0-9_-]+|^--$)
  |(?P<ARGUMENT><[A-Za-z][A-Za-z0-9_-]*>|[A-Z_][A-Z0-9_-]+)
  # Punctuation doesn't accept comma but elipsis ...
  |(?P<PUNCT>[\][=()|]|\.{3})
  |(?P<IDENT>[A-Za-z][A-Za-z0-9_-]+)
  `

  State_Options = `
  (?P<NEWLINE>\n)
  # Options: is matched by state_Usage
  |(?P<SECTION>^[A-Z][A-Za-z _-]+:) => state_Free
  |(?P<DEFAULT>^default: )
  |(?P<LONG_BLANK>[\t ]{2,})
  # skip single blank
  |([\t ])
  |(?P<SHORT>-[A-Za-z0-9?])
  |(?P<LONG>--[A-Za-z][A-Za-z0-9_-]+|^--$)
  |(?P<ARGUMENT><[A-Za-z][A-Za-z0-9_-]+>|[A-Z_][A-Z0-9_-]+)
  # Punctuation differe from state_Usage accepts comma
  |(?P<PUNCT>[\][=,()|])
  # LINE_OF_TEXT not matching []
  |(?P<LINE_OF_TEXT>[^\n[\]]+)
  `

  State_Free = `
  (?P<NEWLINE>\n)
  |(?P<SECTION>^[A-Z][A-Za-z _-]+:)
  |(?P<LINE_OF_TEXT>[^\n]+)
  `

  All_states = map[string]string{
    "state_Prologue" : State_Prologue,
    "state_Usage" : State_Usage,
    "state_Usage_Line" : State_Usage_Line,
    "state_Options" : State_Options,
    "state_Free" : State_Free,
  }
)


import "fmt"

func main() {
	fmt.Println("vim-go")
}
