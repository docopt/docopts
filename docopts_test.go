// vim: set ts=4 sw=4 sts=4 noet:
//
// unit test for docopts.go
//
package main

import (
	"bytes"
	"errors"
	"reflect"
	"strings"
	"testing"
	// our json loader for common_input_test.json
	"fmt"
	"github.com/docopt/docopts/test_json_load"
)

func TestShellquote(t *testing.T) {
	tables := []struct {
		input  string
		expect string
	}{
		{"pipo", "pipo"},
		{"i''i", "i'\\'''\\''i"},
		{"'pipo'", "'\\''pipo'\\''"},
	}

	for _, table := range tables {
		str := Shellquote(table.input)
		if str != table.expect {
			t.Errorf("Shellquote error, got: %s, want: %s.", str, table.expect)
		}
	}
}

func TestIsBashIdentifier(t *testing.T) {
	tables := []struct {
		input  string
		expect bool
	}{
		{"pipo", true},
		{"i''i", false},
		{"'\\''pipo'\\''", false},
		{"OK", true},
		{"ARGS", true},
		// unsecable space at first char
		{" ARGS", false},
		{"123", false},
		{"var%%", false},
		{"varname ", false},
		{"var name", false},
		{"", false},
		{"--", false},
	}

	for _, table := range tables {
		res := IsBashIdentifier(table.input)
		if res != table.expect {
			t.Errorf("IsBashIdentifier for '%s', got: %v, want: %v.", table.input, res, table.expect)
		}
	}
}

func TestIsArray(t *testing.T) {
	tables := []struct {
		input  interface{}
		expect bool
	}{
		{[]string{"pipo", "molo", "--clip"}, true},
		{"pipo", false},
		{42, false},
		{[3]int{1, 2, 3}, true},
	}

	for _, table := range tables {
		rt := reflect.TypeOf(table.input)
		res := IsArray(rt)
		if res != table.expect {
			t.Errorf("IsArray for '%v', got: %v, want: %v.", table.input, res, table.expect)
		}
	}
}

func TestPrint_bash_args(t *testing.T) {
	// replace out (os.Stdout) by a buffer
	bak := out
	out = new(bytes.Buffer)
	defer func() { out = bak }()

	//tables := []struct{
	//    input map[string]interface{}
	//    expect []string
	//}{
	//    {
	//     map[string]interface{}{ "FILE" : []string{"pipo", "molo", "toto"} },
	//     []string{
	//      "declare -A args",
	//      "args['FILE,0']='pipo'",
	//      "args['FILE,1']='molo'",
	//      "args['FILE,2']='toto'",
	//      "args['FILE,#']=3",
	//   },
	//  },
	//    {
	//     map[string]interface{}{ "--counter" : 2 },
	//     []string{
	//      "declare -A args",
	//      "args['--counter']=2",
	//   },
	//  },
	//    {
	//     map[string]interface{}{ "--counter" : "2" },
	//     []string{
	//      "declare -A args",
	//      "args['--counter']='2'",
	//   },
	//  },
	//    {
	//     map[string]interface{}{ "bool" : true },
	//     []string{
	//      "declare -A args",
	//      "args['bool']=true",
	//   },
	//  },
	//}

	d := &Docopts{
		Global_prefix:  "",
		Mangle_key:     true,
		Output_declare: true,
	}

	tables, _ := test_json_loader.Load_json("./common_input_test.json")
	for _, table := range tables {
		d.Print_bash_args("args", table.Input)
		res := out.(*bytes.Buffer).String()
		expect := strings.Join(table.Expect_args[:], "\n") + "\n"
		if res != expect {
			t.Errorf("Print_bash_args for '%v'\ngot: '%v'\nwant: '%v'\n", table.Input, res, expect)
		}
		out.(*bytes.Buffer).Reset()
	}
}

func TestTo_bash(t *testing.T) {
	tables := []struct {
		input  interface{}
		expect string
	}{
		{"pipo", "'pipo'"},
		{"i''i", "'i'\\'''\\''i'"},
		{123, "123"},
		{nil, ""},
		{"", "''"},
		{[]string{"pipo", "molo"}, "('pipo' 'molo')"},
		{true, "true"},
	}

	for _, table := range tables {
		res := To_bash(table.input)
		if res != table.expect {
			t.Errorf("To_bash for '%s', got: %v, want: %v.", table.input, res, table.expect)
		}
	}
}

// helpers compose no-mangle output for matching test
func rewrite_not_mangled(input map[string]interface{}) string {
	var out string
	for _, k := range Sort_args_keys(input) {
		v := input[k]
		out += fmt.Sprintf("%s=%s\n", k, To_bash(v))
	}
	return out
}

func rewrite_prefix(prefix string, expected []string) string {
	var out string
	for _, l := range expected {
		out += fmt.Sprintf("%s_%s\n", prefix, l)
	}
	return out
}

func TestPrint_bash_global(t *testing.T) {
	// replace out (os.Stdout) by a buffer
	bak := out
	out = new(bytes.Buffer)
	defer func() { out = bak }()

	// now loads test from a JSON file
	tables, _ := test_json_loader.Load_json("./common_input_test.json")

	// static tables format
	//tables := []struct{
	//    input map[string]interface{}
	//    expect []string
	//}{
	//    {
	//     map[string]interface{}{ "FILE" : []string{"pipo", "molo", "toto"} },
	//     []string{
	//      "FILE=('pipo' 'molo' 'toto')",
	//   },
	//  },
	//    {
	//     map[string]interface{}{ "--counter" : 2 },
	//     []string{
	//      "counter=2",
	//   },
	//  },
	//    {
	//     map[string]interface{}{ "--counter" : "2" },
	//     []string{
	//      "counter='2'",
	//   },
	//  },
	//    {
	//     map[string]interface{}{ "bool" : true },
	//     []string{
	//      "bool=true",
	//   },
	//  },
	//}

	var err error
	// with Mangle_key
	d := &Docopts{
		Global_prefix: "",
		Mangle_key:    true,
	}
	for _, table := range tables {
		err = d.Print_bash_global(table.Input)
		if err != nil {
			t.Errorf("Print_bash_global doesn't return nil for err: %v\n", err)
		}
		res := out.(*bytes.Buffer).String()
		expect := strings.Join(table.Expect_global[:], "\n") + "\n"
		if res != expect {
			t.Errorf("Print_bash_global for '%v'\ngot: '%v'\nwant: '%v'\n", table.Input, res, expect)
		}
		out.(*bytes.Buffer).Reset()
	}

	// without Mangle_key --no-mangle
	d = &Docopts{
		Global_prefix: "",
		Mangle_key:    false,
	}
	for _, table := range tables {
		err = d.Print_bash_global(table.Input)
		if err != nil {
			t.Errorf("Print_bash_global doesn't return nil for err: %v\n", err)
		}
		res := out.(*bytes.Buffer).String()
		expect := rewrite_not_mangled(table.Input)
		if res != expect {
			t.Errorf("Mangle_key false: Print_bash_global for '%v'\ngot: '%v'\nwant: '%v'\n", table.Input, res, expect)
		}
		out.(*bytes.Buffer).Reset()
	}

	// with Mangle_key and Global_prefix
	d = &Docopts{
		Global_prefix: "ARGS",
		Mangle_key:    true,
	}
	for _, table := range tables {
		err = d.Print_bash_global(table.Input)
		if err != nil {
			t.Errorf("Print_bash_global doesn't return nil for err: %v\n", err)
		}
		res := out.(*bytes.Buffer).String()

		var expect string
		if len(table.Expect_global_prefix) > 0 {
			// if special case was provided
			expect = strings.Join(table.Expect_global_prefix[:], "\n") + "\n"
		} else {
			expect = rewrite_prefix("ARGS", table.Expect_global)
		}
		if res != expect {
			t.Errorf("with prefix: Print_bash_global for '%v'\ngot: '%v'\nwant: '%v'\n", table.Input, res, expect)
		}
		out.(*bytes.Buffer).Reset()
	}

	// with Mangle_key plus name collision
	input_args := make(map[string]interface{})
	input_args["--long-option"] = true
	input_args["<long-option>"] = "dummy_value"

	err = d.Print_bash_global(input_args)
	if err == nil {
		t.Errorf("Print_bash_global expecting err on duplicate Mangle_key options")
	}
	out.(*bytes.Buffer).Reset()
}

type Expected struct {
	s string
	e error
}

func TestName_mangle(t *testing.T) {
	tables := []struct {
		input  string
		expect Expected
	}{
		{
			"FILE",
			Expected{s: "FILE", e: nil},
		},
		{
			"--counter",
			Expected{s: "counter", e: nil},
		},
		{
			"--counter-strike",
			Expected{s: "counter_strike", e: nil},
		},
		{
			"--",
			Expected{s: "", e: errors.New("fail")},
		},
		{
			"<key_word>",
			Expected{s: "key_word", e: nil},
		},
		{
			"<key-word>",
			Expected{s: "key_word", e: nil},
		},
		{
			"-A",
			Expected{s: "A", e: nil},
		},
		{
			"-9",
			Expected{s: "", e: errors.New("fail")},
		},
		{
			"CamelCase",
			Expected{s: "CamelCase", e: nil},
		},
		{
			"SCREAMING_SNAKE_CASE",
			Expected{s: "SCREAMING_SNAKE_CASE", e: nil},
		},
		{
			"l-i-s-p",
			Expected{s: "l_i_s_p", e: nil},
		},
	}

	d := &Docopts{
		Global_prefix: "",
		Mangle_key:    true,
	}

	for _, table := range tables {
		res, err := d.Name_mangle(table.input)
		if table.expect.e != nil && err == nil {
			t.Errorf("Name_mangle for '%v'\ngot: '%v'\nwant: '%v'\n", table.input, err, table.expect.e)
		}
		if res != table.expect.s {
			t.Errorf("Name_mangle for '%v'\ngot: '%v'\nwant: '%v'\n", table.input, res, table.expect.s)
		}
	}
}
