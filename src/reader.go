package apit

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"

	"github.com/fatih/color"
)

// Struct to hold all passed in json data
type testCase struct {
	Name     string          `json:"name"`
	Method   string          `json:"method"`
	Path     string          `json:"path"`
	Response []int           `json:"response"`
	Sent     json.RawMessage `json:"sent"`
	Return   json.RawMessage `json:"return"`
}

// Grab JSON from passed in file
func ReadJSON(file string) ([]testCase, bool) {

	var tests []testCase
	// Get the filepath of the current directory
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))

	if err != nil {
		log.Fatal(err)
	}

	if input, err := ioutil.ReadFile(dir + "/" + file); err == nil {
		// Check if file contains valid JSON
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

// Run all tests and keep track of results.
func Execute(tests []testCase) {

	color.Blue("Starting test cases...")

	client := &http.Client{}

	var passed int = 0

	for _, test := range tests {
		color.Blue("Test: " + test.Name)
		color.Blue(test.Method + " " + test.Path)

		if data, err := test.Sent.MarshalJSON(); err == nil {

			var contentbody io.Reader = nil

			if string(data) != "null" {
				contentbody = bytes.NewReader(data)
			}

			if success := makeRequest(test, contentbody, client); success {
				passed++
			}

		} else {
			color.Red(err.Error())
		}
	}

	color.Blue("Passed: %v/%v", passed, len(tests))

	if passed == len(tests) {
		color.Green("SUCCESS")
	} else {
		color.Red("FAILED")
	}

}

// Run for every test case we loop through
func makeRequest(test testCase, body io.Reader, client *http.Client) bool {

	req, err := http.NewRequest(test.Method, test.Path, body)

	if err != nil {
		color.Red("NewRequest: %s", err.Error())
		return false
	}

	resp, err := client.Do(req)

	if err != nil {
		color.Red("Do: %s", err.Error())
		return false
	}

	returned, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		color.Red("ReadAll: %s", err.Error())
		return false
	}

	// Unmarshal and compare JSON
	data, err := test.Return.MarshalJSON()

	if err != nil {
		color.Red("MarshalJSON: %s", err.Error())
		return false
	}

	json_input := make(map[string]interface{})
	web_output := make(map[string]interface{})

	err = json.Unmarshal(data, &json_input)

	if err != nil {
		color.Red("Unmarshal File: %s", err.Error())
		return false
	}

	err = json.Unmarshal(returned, &web_output)

	if err == nil {
		truthy := reflect.DeepEqual(json_input, web_output)

		switch truthy {
		case true:
			color.Green("Result: %v", truthy)
			return true
		case false:
			color.Red("Result: %v", truthy)
			return false
		}
	} else {
		color.Red("Unmarshal Web: %s, Check for proper method and path.", err.Error())
		return false
	}

	return false
}
