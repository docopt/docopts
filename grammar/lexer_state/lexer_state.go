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
  State_name   string
  Re           *regexp.Regexp
  Leave_token  map[string]string
	Symbols map[string]rune
}

func (def stateRegexpDefinition) String() string {
  return fmt.Sprintf("{ State_name: %s, Re: %v, Leave_token: %v, Symbols: %v",
    def.State_name,
    def.Re,
    def.Leave_token,
    def.Symbols,
  )
}

type stateLexer struct {
	pos   lexer.Position
	b     []byte
	re    *regexp.Regexp
	names []string

  s []*stateRegexpDefinition
  current_state string
	symbols map[string]rune
}

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
// Token cannot contain the string ' => ' in the regex part.
// Comment are introduced with # at column 1, whole line is discarded
func Parse_lexer_state(state_name string, pattern string) (*stateRegexpDefinition, error) {
  var final_pat  []string
  leave_token := make(map[string]string)
  re_extract_rename, _ := regexp.Compile(`\(\?P<([^>]+)`)
  leave_str := " => "
  for i, l := range strings.Split(pattern, "\n") {
    // skip comment and blank
    if l == "" || l[0] == "#"[0] {
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
    return nil, err
  }

	symbols := map[string]rune{
		"EOF": lexer.EOF,
	}
	for i, sym := range re.SubexpNames()[1:] {
		if sym != "" {
			symbols[sym] = lexer.EOF - 1 - rune(i)
		}
	}

  if len(symbols) < 2 {
    return nil, fmt.Errorf("Parse_lexer_state: error: no symbol found after parsing regxep: '%s'", pattern)
  }

  s := stateRegexpDefinition {
    State_name: state_name,
    Re: re,
    Leave_token: leave_token,
    Symbols: symbols,
  }
  return &s, nil
}

var eolBytes = []byte("\n")

// StateLexer creates a lexer definition from a regular expression map[string]string.
//
// Each named sub-expression in the regular expression matches a token. Anonymous sub-expressions
// will be matched and discarded.
//
// eg.
//
// s2_def_string := `
// # a state Regexp definition
// (?P<token1>[a-z]+)
// # reaching a token2 will change state to s3
// |(?P<token2>[A-Z]+) => s3
// `
// states_all := map[string]string{
//  "s1" : `(?P<Ident>[a-z]+)|(\s+)|(?P<Number>\d+)`,
//  "s2" : s2_def_string,
//  "s3" : s3_def_string,
// }
//
//     	def, err := StateLexer(states_all, "s1")
func StateLexer(states_all map[string]string, start_state string) (lexer.Definition, error) {
  states := stateLexer{
    current_state: start_state,
    s: make([]*stateRegexpDefinition, len(states_all)),
  }
  for s, p := range states_all {
    lex, err := Parse_lexer_state(s, p)
    if err != nil {
      return nil, fmt.Errorf("StateLexer: '%s': %v", s, err)
    }

    states.s = append(states.s, lex)
    // initialize state
    if s == start_state {
      states.re = lex.Re
    }
  }

  states.Make_symbols()

	return &states, nil
}

//func (sl *stateLexer) String() string {
//  var out string
//  out += fmt.Sprintf("current_state: %s\n", s.current_state)
//  for n, s := range sl.s {
//    out += fmt.Sprintf("%s: %v\n", n, s.Re)
//  }
//}

func (sl *stateLexer) Make_symbols() error {
  // create symbol map
	symbols := map[string]rune{
		"EOF": lexer.EOF,
	}

  i := 0
  for _, s := range sl.s {
    for _, sym := range s.Re.SubexpNames()[1:] {
      if sym != "" {
        symbols[sym] = lexer.EOF - 1 - rune(i)
        i ++
      }
    }
  }

  sl.symbols = symbols
  return nil
}

func (sl *stateLexer) ChangeState(new_state string) error {
  for _, def := range sl.s {
    if def.State_name == new_state {
      sl.re = def.Re
      sl.current_state = new_state
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

func (s *stateLexer) Symbols() map[string]rune {
	return s.symbols
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
				if r.names[i/2] == "" {
					continue nextToken
				}
				token.Type = lexer.EOF - rune(i/2)
				break
			}
		}

		return token, nil
	}

	return lexer.EOFToken(r.pos), nil
}
