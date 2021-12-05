/*
 * This is a lexer code which support state tokenizing.
 *
 * It reuses lexer code from participle and extend to support in token syntax:
 *	- comments
 *	- states changes
 *
 * See Also: github.com/alecthomas/participle/lexer
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
	"github.com/docopt/docopts/grammar/lexer"
	//"github.com/alecthomas/participle/lexer"
	"io"
	"io/ioutil"
	"regexp"
	"strings"
	"unicode/utf8"
)

type dynamicRegexp struct {
	re_name     string
	re_template string
	re_string   string
	is_dynamic  bool
	resolved    bool
	want        string
}

type stateRegexpDefinition struct {
	State_name string
	All_regexp []*dynamicRegexp
	Re         *regexp.Regexp
	// map a named lexer rule to a new stateRegexpDefinition's name
	Leave_token map[string]string
	// list of Symbols name
	Symbols     []string
	DynamicRule bool
}

func (def stateRegexpDefinition) String() string {
	return fmt.Sprintf("{ State_name: %s, Re: %v, Leave_token: %v, Symbols: %v }",
		def.State_name,
		def.Re,
		def.Leave_token,
		def.Symbols,
	)
}

type StateLexer struct {
	// current position in the buffer, lexer.Position copied from participle
	pos lexer.Position

	// content to scan
	b         []byte
	re        *regexp.Regexp
	byte_left int

	// TODO: optimize names change
	names             []string
	State_auto_change bool

	// states
	s             []*stateRegexpDefinition
	Current_state *stateRegexpDefinition
	// map lexer named pattern name to rune of symbols
	symbols map[string]rune
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
// Support dynamic regexp replace with syntax:
// |(?P<@PROG_NAME>@PROG_NAME)
// will be resolved with ==> |(?P<PROG_NAME>%s) where %s will be replaced with the value of PROG_NAME
func Parse_lexer_state(state_name string, pattern string) (*stateRegexpDefinition, error) {
	var all_regexp []*dynamicRegexp
	leave_token := map[string]string{}
	// our regexp to extract the regxep's name from parsed input (valid go named regxep)
	re_extract_rename, _ := regexp.Compile(`\(\?P<([^>]+)`)
	leave_str := " => "
	var new_regexp *dynamicRegexp
	dynamic_rule := 0
	var symbols []string
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

		new_regexp = nil

		// handle leaving state definition
		divide := strings.SplitN(l, leave_str, 2)
		regexp := strings.Trim(divide[0], "\t ")

		// extract regexp name
		re_name_result := re_extract_rename.FindStringSubmatch(regexp)
		re_name := ""
		if re_name_result != nil {
			re_name = re_name_result[1]
		}

		if len(divide) == 2 {
			new_state := strings.Trim(divide[1], "\t ")
			if re_name != "" {
				token := re_name
				//msg = fmt.Sprintf("'%s' => '%s'", token, new_state)
				leave_token[token] = new_state
			} else {
				return nil, fmt.Errorf("Parse_lexer_state:%d: error: regexp not matched token: %s", i, regexp)
			}
		}

		// handle dynamic regexp definition
		if re_name != "" && strings.Index(regexp, "(?P<@") != -1 {
			// remove leading @
			re_name = re_name[1:]
			new_regexp = &dynamicRegexp{
				re_name:     re_name,
				re_template: "|(?P<%s>%s)",
				re_string:   "",
				is_dynamic:  true,
				resolved:    false,
				want:        re_name,
			}
			dynamic_rule++
		} else {
			// also handle anonymous regxep
			new_regexp = &dynamicRegexp{
				re_name:     re_name,
				re_template: "",
				re_string:   regexp,
				is_dynamic:  false,
				resolved:    true,
				want:        "",
			}
		}

		if re_name != "" {
			symbols = append(symbols, re_name)
		}
		all_regexp = append(all_regexp, new_regexp)
	}

	if len(symbols) < 1 {
		return nil, fmt.Errorf("Parse_lexer_state: error: no symbol found after parsing regxep: '%s'", pattern)
	}

	s := stateRegexpDefinition{
		State_name:  state_name,
		All_regexp:  all_regexp,
		Re:          nil,
		Leave_token: leave_token,
		Symbols:     symbols,
		DynamicRule: dynamic_rule > 0,
	}
	if dynamic_rule == 0 {
		err := s.compile_regexp()
		if err != nil {
			return nil, err
		}
	}
	return &s, nil
}

func (s *stateRegexpDefinition) compile_regexp() error {
	var re *regexp.Regexp
	var err error
	var final_pat []string
	for _, dr := range s.All_regexp {
		if !dr.resolved {
			return fmt.Errorf("compile_regexp: dynamicRegexp not resolved: %s missing %s", dr.re_name, dr.want)
		}
		final_pat = append(final_pat, dr.re_string)
	}

	// assign result
	re, err = regexp.Compile(strings.Join(final_pat, ""))
	if err != nil {
		s.Re = nil
		return fmt.Errorf("compile_regexp: %v", err)
	}

	s.Re = re
	return nil
}

var eolBytes = []byte("\n")

// CreateStateLexer creates a lexer definition from a regular expression map[string]string.
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
//      def, err := CreateStateLexer(states_all, "s1")
func CreateStateLexer(states_all map[string]string, start_state string) (*StateLexer, error) {
	states := StateLexer{
		s:                 []*stateRegexpDefinition{},
		State_auto_change: true,
	}
	for s, p := range states_all {
		def, err := Parse_lexer_state(s, p)
		if err != nil {
			return nil, fmt.Errorf("CreateStateLexer: '%s': %v", s, err)
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

//func (sl *StateLexer) String() string {
//  var out string
//  out += fmt.Sprintf("Current_state: %s\n", s.Current_state)
//  for n, s := range sl.s {
//    out += fmt.Sprintf("%s: %v\n", n, s.Re)
//  }
//}

func (sl *StateLexer) Make_symbols() error {
	// create symbol map common for all state
	symbols := map[string]rune{
		"EOF": lexer.EOF,
	}

	// renumber all symbol (symbols are associated with negative rune)
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

func (sl *StateLexer) ChangeState(new_state string) error {
	// search in our states name and update our regexp
	for _, def := range sl.s {
		if def.State_name == new_state {
			// on the fly compile_regexp
			if def.Re == nil {
				if err := def.compile_regexp(); err != nil {
					return fmt.Errorf("ChangeState: '%s' dynamicRegexp error %v", new_state, err)
				}
			}
			// assign the new regexp to the StateLexer
			sl.re = def.Re
			sl.Current_state = def
			sl.names = def.Re.SubexpNames()

			return nil
		}
	}

	return fmt.Errorf("ChangeState: '%s' state_name not found", new_state)
}

// Initialize the Lexer with an io.Reader
// return: a participle lexer.Lexer
func (sl *StateLexer) Lex(r io.Reader) (lexer.Lexer, error) {
	// read all bytes
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	sl.pos = lexer.Position{
		Filename: lexer.NameOfReader(r),
		Line:     1,
		Column:   1,
		Offset:   0,
	}
	sl.b = b
	sl.byte_left = len(b)
	sl.names = sl.re.SubexpNames()

	return sl, nil
}

func (sl *StateLexer) InitSource(source []byte) error {
	sl.pos = lexer.Position{
		Filename: "string",
		Line:     1,
		Column:   1,
		Offset:   0,
	}
	sl.b = source
	sl.byte_left = len(source)

	if sl.Current_state.Re == nil {
		if err := sl.Current_state.compile_regexp(); err != nil {
			return fmt.Errorf("InitSource: '%s'regexp compile error %v", sl.Current_state.State_name, err)
		}
	}

	sl.names = sl.Current_state.Re.SubexpNames()
	if len(sl.names) == 0 {
		return fmt.Errorf("InitSource: start_state '%s' regexp has no names", sl.Current_state.State_name)
	}

	return nil
}

func (sl *StateLexer) Symbols() map[string]rune {
	return sl.symbols
}

func (sl *StateLexer) Next() (lexer.Token, error) {
nextToken:
	for sl.byte_left > 0 {
		matches := sl.re.FindSubmatchIndex(sl.b[sl.pos.Offset:])
		if matches == nil || matches[0] != 0 {
			rn, _ := utf8.DecodeRune(sl.b[sl.pos.Offset:])
			return lexer.Token{}, lexer.Errorf(sl.pos, "invalid token %q state: %s", rn, sl.Current_state)
		}

		// matched_pos is always 0 (start of the pos.Offset)
		// matches[1] is the end position of the match, but as match always happen
		// at pos:0 it's also the length of the match
		matched_len := matches[1]

		match := sl.b[sl.pos.Offset : sl.pos.Offset+matched_len]
		token := lexer.Token{
			Pos:        sl.pos,
			Value:      string(match),
			State_name: sl.Current_state.State_name,
		}

		// Update lexer state.
		sl.pos.Offset += matched_len
		nb_line := bytes.Count(match, eolBytes)
		sl.pos.Line += nb_line
		// Update column.
		if nb_line == 0 {
			sl.pos.Column += utf8.RuneCount(match)
		} else {
			sl.pos.Column = utf8.RuneCount(match[bytes.LastIndex(match, eolBytes):])
		}
		sl.byte_left -= matched_len

		// Finally, assign token type. If it is not a named group, we continue to the next token.
		for i := 2; i < len(matches); i += 2 {
			if matches[i] != -1 {
				tok_name := sl.names[i/2]
				// discard unnamed token
				if tok_name == "" {
					continue nextToken
				}

				token.Type = sl.symbols[tok_name]
				token.Regex_name = tok_name

				if sl.State_auto_change {
					// if we encounter a leave_token we change our lexer state
					if new_state, ok := sl.Current_state.Leave_token[tok_name]; ok {
						sl.ChangeState(new_state)
					}
				}

				break
			}
		}

		return token, nil
	}

	return lexer.EOFToken(sl.pos), nil
}

// Discard some utf-8 char from the input buffer at the current position
// so move forward the requested number of mbchar
func (sl *StateLexer) Discard(nb_mbchar int) string {
	// compute the number of bytes from pos.Offset to discard
	nb_byte_to_move := 0
	for i := 0; i < nb_mbchar; i++ {
		_, nb_byte := utf8.DecodeRune(sl.b[sl.pos.Offset+nb_byte_to_move:])
		nb_byte_to_move += nb_byte
	}
	match := sl.b[sl.pos.Offset : sl.pos.Offset+nb_byte_to_move]
	// Update lexer state.
	sl.pos.Offset += nb_byte_to_move
	nb_line := bytes.Count(match, eolBytes)
	sl.pos.Line += nb_line
	// Update column.
	if nb_line == 0 {
		sl.pos.Column += utf8.RuneCount(match)
	} else {
		sl.pos.Column = utf8.RuneCount(match[bytes.LastIndex(match, eolBytes):])
	}
	// update bytes array
	sl.byte_left -= nb_byte_to_move

	return string(match)
}

func (sl *StateLexer) DynamicRuleUpdate(variable string, value string) error {
	found := false
	for _, s := range sl.s {
		if rd, ok := s.has_dynamic_regexp(variable); ok {
			rd.update_regexp(value)
			found = true
		}
	}

	if !found {
		return fmt.Errorf("DynamicRuleUpdate: variable not found: '%s'", variable)
	}

	return sl.compile_regexp()
}

func (sl *StateLexer) compile_regexp() error {
	var err error = nil
	for _, s := range sl.s {
		err = s.compile_regexp()
	}
	return err
}

func (s *stateRegexpDefinition) has_dynamic_regexp(variable string) (*dynamicRegexp, bool) {
	for _, dr := range s.All_regexp {
		if dr.want == variable && dr.is_dynamic {
			return dr, true
		}
	}
	return nil, false
}

func (dr *dynamicRegexp) update_regexp(value string) error {
	dr.re_string = fmt.Sprintf(dr.re_template, dr.re_name, value)
	dr.resolved = true

	return nil
}

func (sl *StateLexer) Reject(tok *lexer.Token) {
	nb_byte_diff := sl.pos.Offset - tok.Pos.Offset

	// Update lexer state
	sl.pos.Offset = tok.Pos.Offset
	sl.pos.Line = tok.Pos.Line
	sl.pos.Column = tok.Pos.Column

	sl.byte_left += nb_byte_diff
}
