package apit

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
)

// Struct to hold all passed in json data
type testCase struct {
	Name     string
	Path     string
	Response []int
	Object   []byte
}

// Grab json from passed in file
func ReadJSON(file string) ([]testCase, bool) {

	var tests []testCase

	if input, err := ioutil.ReadFile(file); err == nil {
		// Check if file contains valid json
		decoder := json.NewDecoder(bytes.NewReader(input))
		// Decode and check for error
		if err := decoder.Decode(tests); err == nil {
			return tests, true
		}
	}

	return nil, false
}
