package token_docopt

// ================================ lexer for docopt language ===========
var (
	// State_Prologue can only leave if we encounter "usage:" section
	State_Prologue = `
  (?P<NEWLINE>\n)
  |(?P<USAGE>[Uu][Ss][Aa][Gg][Ee]:) => state_First_Program_Usage
  |(?P<WORD>\S+)
  |(?P<BLANK>\s+)
  `

	State_Usage = `
  (?P<NEWLINE>\n)
  |(?P<USAGE>[Uu][Ss][Aa][Gg][Ee]:)
	#|(?P<PROG_NAME>[^ :]+)
  |(?P<SECTION>^[A-Z][A-Za-z _-]+:) => state_Options
  |(?P<LONG_BLANK>\t|[\t ]{2,}) => state_Usage_Line
  # skip single blank
  |([\t ])
  # Match some kind of comment when not preceded by LongBlank
  |(?P<LINE_OF_TEXT>[^\n]+)
  `

	State_First_Program_Usage = `
  (?P<NEWLINE>\n)
  |(?P<BLANK>[\t ]+)
	|(?P<PROG_NAME>\S+) => state_Usage_Line
	`

	State_Usage_Line = `
  (?P<NEWLINE>\n) => state_Usage
  |(?P<LONG_BLANK>\t|[\t ]{2,})
  |([\t ])
	|(?P<@PROG_NAME>@PROG_NAME)
  |(?P<SHORT>-[A-Za-z0-9?])
  |(?P<LONG>--[A-Za-z][A-Za-z0-9_-]+|^--$)
	# argument are free text colonn is an ARGUMENT
	|(?P<ARGUMENT><[A-Za-z][A-Za-z0-9_-]*>|[A-Z_][A-Z0-9_-]+|[:])
	# Punctuation doesn't accept comma, accepts elipsis ...
	|(?P<PUNCT>[\][=()|-]|\.{3})
  |(?P<IDENT>[A-Za-z][A-Za-z0-9_-]+)
  `

	State_Options = `
  (?P<NEWLINE>\n)
  # Options: is matched by state_Usage
  |(?P<SECTION>^[A-Z][A-Za-z _-]+:) => state_Free
  |(?P<LONG_BLANK>\t|[\t ]{2,})
	|(?P<DEFAULT>^default: )
  # skip single blank
  |([\t ])
  |(?P<SHORT>-[A-Za-z0-9?])
  |(?P<LONG>--[A-Za-z][A-Za-z0-9_-]+|^--$)
  |(?P<ARGUMENT><[A-Za-z][A-Za-z0-9_-]+>|[A-Z_][A-Z0-9_-]+)
  # Punctuation differe from state_Usage accepts comma and dot
  |(?P<PUNCT>[=,()|.[\]])
	 # LINE_OF_TEXT not some PUNCT []
  |(?P<LINE_OF_TEXT>[^\n[\]]+)
  `

	State_Free = `
  (?P<NEWLINE>\n)
  |(?P<SECTION>^[A-Z][A-Za-z _-]+:)
  |(?P<LINE_OF_TEXT>[^\n]+)
  `

	All_states = map[string]string{
		"state_Prologue":            State_Prologue,
		"state_First_Program_Usage": State_First_Program_Usage,
		"state_Usage":               State_Usage,
		"state_Usage_Line":          State_Usage_Line,
		"state_Options":             State_Options,
		"state_Free":                State_Free,
	}
)
