package main

import (
	"os"
	"encoding/json"
	"io/ioutil"
)

const jsonFile = "statuses.json"
const debug = true

func updateStatus(item string, status string) {
	_, err := os.Stat(jsonFile)
	if err != nil {
		// file doesn't exist, create an empty json file
		// create an empty map
			m := make(map[string]string)
		// marshall to a empty json object
		b, json_err := json.Marshal(m)
		if json_err != nil { panic(json_err) }
		// write out to file
		write_err := ioutil.WriteFile(jsonFile, b, 0644)
		if write_err != nil { panic(write_err) }
	}	
	// read current statues file
	r, read_err := ioutil.ReadFile(jsonFile)
	if read_err != nil { panic(read_err) }	
	// create any empty map to hold json contents
	m := make(map[string]string)
	// unmarshall data read in from file
	unmarshall_error:= json.Unmarshal(r, &m)
	if unmarshall_error != nil { panic(unmarshall_error) }
	// set map with the passed in data
	m[item] = status
	// marshall data
	b, marshall_error := json.Marshal(m)
	if marshall_error != nil { panic(marshall_error) }
	// write json data tile file
	write_error := ioutil.WriteFile(jsonFile, b, 0644)
	if write_error !=nil { panic(write_error) }
	
}

func main () {
	args := os.Args
	if len(args) == 3 {
		updateStatus(args[1], args[2])
	} else {
		panic("Specify Item and Status")
	}
}