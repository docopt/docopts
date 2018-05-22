// vim: set ts=4 sw=4 sts=4 et:
package main

import "testing"

func TestShellquote(t *testing.T) {
    tables := []struct {
        input string
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
        input string
        expect bool
    }{
        {"pipo", true},
        {"i''i", false},
        {"'\\''pipo'\\''", false},
        {"OK", true},
        {"123", false},
        {"var%%", false},
        {"varname ", false},
        {"var name", false},
    }

    for _, table := range tables {
      res := IsBashIdentifier(table.input)
      if res != table.expect {
         t.Errorf("IsBashIdentifier for '%s', got: %v, want: %v.", table.input, res, table.expect)
      }
    }
}
