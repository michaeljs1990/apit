package main

import (
	"flag"
	"terame.com/michaeljs1990/apit/src"
)

// Get all arguments and call out to the api
// to complete the request.
func main() {
	// flag variables
	var file string

	flag.StringVar(&file, "file", "", "json file that contains all your test infromation")

	// Define all flags above this call
	flag.Parse()

	// Check if file contains valid json and can be open
	if data, valid := apit.ReadJSON(file); valid {
		// Run tests against API
		apit.Execute(data)
	}

}
