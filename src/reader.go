package apit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/fatih/color"
)

// Struct to hold all passed in json data
type testCase struct {
	Name     string          `json:"name"`
	Path     string          `json:"path"`
	Response []int           `json:"response"`
	Sent     json.RawMessage `json:"sent"`
	Return   json.RawMessage `json:"return"`
}

// Grab json from passed in file
func ReadJSON(file string) ([]testCase, bool) {

	var tests []testCase

	// Get the filepath of the current directory
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))

	if err != nil {
		log.Fatal(err)
	} else {
		dir = dir + "/"
	}

	if input, err := ioutil.ReadFile(dir + file); err == nil {
		// Check if file contains valid json
		decoder := json.NewDecoder(bytes.NewReader(input))
		// Decode and check for error
		if err := decoder.Decode(&tests); err == nil {
			return tests, true
		} else {
			color.Red(err.Error())
		}
	} else {
		color.Red(err.Error())
	}

	return nil, false
}

// Run all tests and log to file
func RunTests(tests []testCase) {

	color.Blue("Starting test casses...")

	for test := range tests {
		color.Green("run")
		fmt.Println(tests[test])
	}

}
