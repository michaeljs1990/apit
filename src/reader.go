package apit

import (
	"bytes"
	"encoding/json"
	"fmt"
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

// Run all tests
func RunTests(tests []testCase) {

	color.Blue("Starting test cases...")

	client := &http.Client{}

	for _, test := range tests {
		color.Blue("Test: " + test.Name)
		color.Blue(test.Method + " " + test.Path)

		if data, err := test.Sent.MarshalJSON(); err == nil {

			var contentbody io.Reader = nil

			if string(data) != "null" {
				contentbody = bytes.NewReader(data)
			}

			if req, err := http.NewRequest(test.Method, test.Path, contentbody); err == nil {
				if resp, err := client.Do(req); err == nil {
					if returned, err := ioutil.ReadAll(resp.Body); err == nil {
						// Unmarshal and compare JSON
						if data, err := test.Return.MarshalJSON(); err == nil {

							mashaled_input := make(map[string]interface{})
							fdsa := make(map[string]interface{})

							err := json.Unmarshal(data, &mashaled_input)
							err2 := json.Unmarshal(returned, &fdsa)

							if err == nil && err2 == nil {
								fmt.Println(reflect.DeepEqual(mashaled_input, fdsa))
							} else {
								color.Red(err.Error(), err2.Error())
							}

						} else {
							color.Red(err.Error())
						}
					} else {
						color.Red(err.Error())
					}
				} else {
					color.Red(err.Error())
				}
			} else {
				color.Red(err.Error())
			}
		} else {
			color.Red(err.Error())
		}

	}

}
