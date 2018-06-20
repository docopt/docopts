package main

import (
    "fmt"
    "github.com/docopt/docopts/test_json_load"
)

func main() {
  t, err := test_json_loader.Load_json("../common_input_test.json")
  if err != nil {
    fmt.Println(err.Error())
  } else {
    for _, e := range t {
      fmt.Printf("%v\n", e.ToString())
    }
  }

}
