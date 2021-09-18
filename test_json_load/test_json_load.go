package test_json_loader

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type TestString struct {
	Input                map[string]interface{}
	Expect_args          []string
	Expect_global        []string
	Expect_global_prefix []string
}

func (t TestString) ToString() string {
	return t.toString()
}

// some debug
func (t TestString) toString() (str string) {
	var key string
	var val interface{}

	str = fmt.Sprintf("input: {\n")
	for key, val = range t.Input {
		str += fmt.Sprintf(" key : %s , value : %v\n", key, val)
	}
	str += fmt.Sprintf("}\n")

	str += fmt.Sprintf("Expect_args {\n")
	for _, v := range t.Expect_args {
		str += fmt.Sprintf("%v\n", v)
	}
	str += fmt.Sprintf("}\n")

	str += fmt.Sprintf("Expect_global : %v\n", t.Expect_global)
	str += fmt.Sprintf("Expect_global_prefix : %v\n", t.Expect_global_prefix)

	return str
}

func Load_json(filename string) ([]TestString, error) {
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var c []TestString
	json.Unmarshal(raw, &c)

	// loop over result and convert type when needed to match expected type in tests
	for _, ts := range c {
		for k, v := range ts.Input {
			switch v.(type) {
			case []interface{}:
				// converts []interface{} to []string
				aInterface := v.([]interface{})
				aString := make([]string, len(aInterface))
				// copy and change type
				for i, str_v := range aInterface {
					aString[i] = str_v.(string)
				}
				// overwrite old key with the new converted array
				ts.Input[k] = aString
			case float64:
				// json.Unmarshal read all number to float64
				// See: https://golang.org/pkg/encoding/json/#Unmarshal
				ts.Input[k] = int(v.(float64))
			}
		}
	}
	return c, nil
}
