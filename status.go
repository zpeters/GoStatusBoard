package main

import (
	"os"
	"encoding/json"
	"io/ioutil"
)

const jsonFile = "statuses.json"

func updateStatus(item string, status string) {
	r, err := ioutil.ReadFile(jsonFile)
	if err != nil { panic(err) }	
	var m map[string]string
	err2 := json.Unmarshal(r, &m)
	if err2 != nil { panic(err2) }
	m[item] = status
	b, err3 := json.Marshal(m)
	if err3 != nil { panic(err3) }
	err4 := ioutil.WriteFile(jsonFile, b, 0644)
	if err4 !=nil { panic(err4) }
	
}

func main () {
	args := os.Args
	if len(args) == 3 {
		updateStatus(args[1], args[2])
	}
}