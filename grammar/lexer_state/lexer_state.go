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
	// lexer.Position copied from participle
	pos lexer.Position
	// content to scan
	b  []byte
	re *regexp.Regexp

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
func (s *StateLexer) Lex(r io.Reader) (lexer.Lexer, error) {
	// read all bytes
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

func (sl *StateLexer) InitSource(source []byte) error {
	sl.pos = lexer.Position{
		Filename: "source",
		Line:     1,
		Column:   1,
	}
	sl.b = source

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

func (r *StateLexer) Next() (lexer.Token, error) {
nextToken:
	for len(r.b) != 0 {
		matches := r.re.FindSubmatchIndex(r.b)
		if matches == nil || matches[0] != 0 {
			rn, _ := utf8.DecodeRune(r.b)
			return lexer.Token{}, lexer.Errorf(r.pos, "invalid token %q state: %s", rn, r.Current_state)
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
				token.Regex_name = tok_name

				if r.State_auto_change {
					// if we encounter a leave_token we change our lexer state
					if new_state, ok := r.Current_state.Leave_token[tok_name]; ok {
						r.ChangeState(new_state)
					}
				}

				break
			}
		}

		return token, nil
	}

	return lexer.EOFToken(r.pos), nil
}

func (r *StateLexer) Discard(pos lexer.Position, nb_uchar int) {
	nb_byte_to_move := 0
	for i := 0; i < nb_uchar; i++ {
		_, nb_byte := utf8.DecodeRune(r.b[nb_byte_to_move:])
		nb_byte_to_move += nb_byte
	}
	match := r.b[:nb_byte_to_move]
	// Update lexer state.
	r.pos.Offset += nb_byte_to_move
	lines := bytes.Count(match, eolBytes)
	r.pos.Line += lines
	// Update column.
	if lines == 0 {
		r.pos.Column += utf8.RuneCount(match)
	} else {
		r.pos.Column = utf8.RuneCount(match[bytes.LastIndex(match, eolBytes):])
	}
	// update bytes array
	r.b = r.b[nb_byte_to_move:]
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
