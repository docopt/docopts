/*
 * This is a lexer code which support state tokenizing.
 *
 * It reuses github.com/alecthomas/participle/lexer code and extend to support in token syntax:
 *	- comments
 *	- states changes
 *
 *
 * Our syntax:
 *
 *  Go Strings are setup with a regexp using named group (?P<NAME>regexp_here)
 *  We parse the string regexp by removing comment and empty lines.
 *  We build the resulting stateRegexpDefinition.
 *  Token cannot contain the string ' => ' in the regexp part.
 *  Comment are introduced with # at column 1, whole line is discarded.
 *  The rest of the sting is a normal Go's regexp input.
 *
 *  And finaly all strings are associated with states name, using a `map[string]string`
 *  Keys of the map[] are the string that can be used on right part of a state change:
 *
 *     (?P<TOKEN_NAME>regexp_here) => new_state
 *
 *  Example:
 *
 *  	All_states = map[string]string{
 *  		"state_Prologue":   State_Prologue,
 *  		"state_Usage":      State_Usage,
 *  		"state_Usage_Line": State_Usage_Line,
 *  		"state_Options":    State_Options,
 *  		"state_Free":       State_Free,
 *  	}
 *
 *  Examples of states:
 *
 * 	   State_Prologue = `
 *      (?P<NEWLINE>\n)
 *      |(?P<SECTION>^Usage:) => state_Usage_Line
 *      |(?P<LINE_OF_TEXT>[^\n]+)
 *      `
 *   	State_Usage = `
 *     (?P<NEWLINE>\n)
 *     |(?P<USAGE>^Usage:)
 *     |(?P<SECTION>^[A-Z][A-Za-z _-]+:) => state_Options
 *     |(?P<LONG_BLANK>[\t ]{2,}) => state_Usage_Line
 *     # skip single blank
 *     |([\t ])
 *     # Match some kind of comment when not preceded by LongBlank
 *     |(?P<LINE_OF_TEXT>[^\n]+)
 *     `
 */
package lexer_state

import (
	"bytes"
	"fmt"
	"github.com/alecthomas/participle/lexer"
	"io"
	"io/ioutil"
	"regexp"
	"strings"
	"unicode/utf8"
)

type stateRegexpDefinition struct {
	State_name  string
	Re          *regexp.Regexp
	Leave_token map[string]string
	Symbols     []string
}

func (def stateRegexpDefinition) String() string {
	return fmt.Sprintf("{ State_name: %s, Re: %v, Leave_token: %v, Symbols: %v }",
		def.State_name,
		def.Re,
		def.Leave_token,
		def.Symbols,
	)
}

type stateLexer struct {
	// lexer.Position from participle
	pos lexer.Position
	b   []byte
	re  *regexp.Regexp
	// TODO: optimize names change
	names []string

	s             []*stateRegexpDefinition
	Current_state *stateRegexpDefinition
	symbols       map[string]rune
}

// helper
func mergeMaps(maps ...map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

// Parse a regexp string removing comment and empty lines
// and build the resulting stateRegexpDefinition
// Token cannot contain the string ' => ' in the regexp part.
// Comment are introduced with # at column 1, whole line is discarded
func Parse_lexer_state(state_name string, pattern string) (*stateRegexpDefinition, error) {
	var final_pat []string
	leave_token := map[string]string{}
	re_extract_rename, _ := regexp.Compile(`\(\?P<([^>]+)`)
	leave_str := " => "
	for i, l := range strings.Split(pattern, "\n") {
		// skip empty line
		if l == "" {
			continue
		}

		l = strings.TrimLeft(l, "\t ")

		// skip empty line after trim
		if l == "" {
			continue
		}

		// skip comment
		if l[0] == '#' {
			continue
		}

		if strings.Count(l, leave_str) > 1 {
			return nil, fmt.Errorf("Parse_lexer_state:%d: error: more than one '%s', in '%s'", i, leave_str, l)
		}

		divide := strings.SplitN(l, leave_str, 2)
		regexp := strings.Trim(divide[0], "\t ")
		//var msg string
		if len(divide) == 2 {
			new_state := strings.Trim(divide[1], "\t ")
			result := re_extract_rename.FindStringSubmatch(regexp)
			if result != nil {
				token := result[1]
				//msg = fmt.Sprintf("'%s' => '%s'", token, new_state)
				leave_token[token] = new_state
			} else {
				return nil, fmt.Errorf("Parse_lexer_state:%d: error: regexp not matched token: %s", i, regexp)
			}
		}

		// fmt.Printf("%d: '%s' : %s\n", i, l, msg)

		final_pat = append(final_pat, regexp)
	}

	// assign result
	re, err := regexp.Compile(strings.Join(final_pat, ""))
	if err != nil {
		return nil, fmt.Errorf("Parse_lexer_state: %v", err)
	}

	// retrieve symbols (token names) from regxep names
	symbols := []string{}
	for _, sym := range re.SubexpNames()[1:] {
		if sym != "" {
			symbols = append(symbols, sym)
		}
	}

	if len(symbols) < 1 {
		return nil, fmt.Errorf("Parse_lexer_state: error: no symbol found after parsing regxep: '%s'", pattern)
	}

	s := stateRegexpDefinition{
		State_name:  state_name,
		Re:          re,
		Leave_token: leave_token,
		Symbols:     symbols,
	}
	return &s, nil
}

var eolBytes = []byte("\n")

// StateLexer creates a lexer definition from a regular expression map[string]string.
//
// Each named sub-expression in the regular expression matches a token. Anonymous sub-expressions
// will be matched and discarded.
//
// Examples:
//
// s2_def_string := `
//   # a state Regexp definition
//   (?P<token1>[a-z]+)
//   # reaching a token2 will change state to s3
//   |(?P<token2>[A-Z]+) => s3
//   `
// states_all := map[string]string{
//    "s1" : `(?P<Ident>[a-z]+)|(\s+)|(?P<Number>\d+)`,
//    "s2" : s2_def_string,
//    "s3" : s3_def_string,
//   }
//
//      def, err := StateLexer(states_all, "s1")
func StateLexer(states_all map[string]string, start_state string) (*stateLexer, error) {
	states := stateLexer{
		s: []*stateRegexpDefinition{},
	}
	for s, p := range states_all {
		def, err := Parse_lexer_state(s, p)
		if err != nil {
			return nil, fmt.Errorf("StateLexer: '%s': %v", s, err)
		}

		states.s = append(states.s, def)
		// initialize state
		if s == start_state {
			states.re = def.Re
			states.Current_state = def
		}
	}

	states.Make_symbols()

	return &states, nil
}

//func (sl *stateLexer) String() string {
//  var out string
//  out += fmt.Sprintf("Current_state: %s\n", s.Current_state)
//  for n, s := range sl.s {
//    out += fmt.Sprintf("%s: %v\n", n, s.Re)
//  }
//}

func (sl *stateLexer) Make_symbols() error {
	// create symbol map common for all state
	symbols := map[string]rune{
		"EOF": lexer.EOF,
	}

	// renumber all symbol
	var tok rune = lexer.EOF - 1
	for _, sdef := range sl.s {
		for _, sym := range sdef.Symbols {
			// skip symbol already known
			if _, ok := symbols[sym]; ok {
				continue
			}
			symbols[sym] = tok
			tok--
		}
	}

	sl.symbols = symbols
	return nil
}

func (sl *stateLexer) ChangeState(new_state string) error {
	for _, def := range sl.s {
		if def.State_name == new_state {
			sl.re = def.Re
			sl.Current_state = def
			sl.names = def.Re.SubexpNames()

			return nil
		}
	}

	return fmt.Errorf("ChangeState: '%s' state_name not found", new_state)
}

func (s *stateLexer) Lex(r io.Reader) (lexer.Lexer, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	s.pos = lexer.Position{
		Filename: lexer.NameOfReader(r),
		Line:     1,
		Column:   1,
	}
	s.b = b
	s.names = s.re.SubexpNames()

	return s, nil
}

func (sl *stateLexer) Symbols() map[string]rune {
	return sl.symbols
}

func (r *stateLexer) Next() (lexer.Token, error) {
nextToken:
	for len(r.b) != 0 {
		matches := r.re.FindSubmatchIndex(r.b)
		if matches == nil || matches[0] != 0 {
			rn, _ := utf8.DecodeRune(r.b)
			return lexer.Token{}, lexer.Errorf(r.pos, "invalid token %q", rn)
		}
		match := r.b[:matches[1]]
		token := lexer.Token{
			Pos:   r.pos,
			Value: string(match),
		}

		// Update lexer state.
		r.pos.Offset += matches[1]
		lines := bytes.Count(match, eolBytes)
		r.pos.Line += lines
		// Update column.
		if lines == 0 {
			r.pos.Column += utf8.RuneCount(match)
		} else {
			r.pos.Column = utf8.RuneCount(match[bytes.LastIndex(match, eolBytes):])
		}
		// Move slice along.
		r.b = r.b[matches[1]:]

		// Finally, assign token type. If it is not a named group, we continue to the next token.
		for i := 2; i < len(matches); i += 2 {
			if matches[i] != -1 {
				tok_name := r.names[i/2]
				// discard unnamed token
				if tok_name == "" {
					continue nextToken
				}

				token.Type = r.symbols[tok_name]

				// if we encounter a leave_token we change our lexer state
				if new_state, ok := r.Current_state.Leave_token[tok_name]; ok {
					r.ChangeState(new_state)
				}
				break
			}
		}

		return token, nil
	}

	return lexer.EOFToken(r.pos), nil
}
