package main

import (
	"os"
	"encoding/json"
	"io/ioutil"
	"fmt"
)

const jsonFile = "statuses.json"
const htmlFile = "statuses.html"

func printHelp() {
	fmt.Printf("USAGE:\n")
	fmt.Printf("statusboard ACTION OPTIONS\n")
	fmt.Printf("\n")
	fmt.Printf("\t statusboard update ITEM STATUS\n")
	fmt.Printf("\t\t ITEM = name of service\n")
	fmt.Printf("\t\t STATUS = status of service\n")
	fmt.Printf("\n")
	fmt.Printf("\t statusboard dump\n")
	fmt.Printf("\t\t dumps the current statues to html\n")
	fmt.Printf("\n")
	fmt.Printf("\t statusboard help\n")
	fmt.Printf("\t\t This file\n")
}

func createEmptyDb() {
	m := make(map[string]string)
	b, json_err := json.Marshal(m)
	if json_err != nil { panic(json_err) }
	write_err := ioutil.WriteFile(jsonFile, b, 0644)
	if write_err != nil { panic(write_err) }
}

func dumpStatus() {
	// check if db exists - if not create a blank one
	_, err := os.Stat(jsonFile)
	if err != nil { createEmptyDb() }	

	// create any empty map to hold json contents
	m := make(map[string]string)

	// read current statues file
	r, read_err := ioutil.ReadFile(jsonFile)
	if read_err != nil { panic(read_err) }	

	// unmarshall data read in from file
	unmarshall_error:= json.Unmarshal(r, &m)
	if unmarshall_error != nil { panic(unmarshall_error) }

	// create the header
	d := fmt.Sprintf("<html>\n\t<head>\n\t\t<title>Title</title>\n\t</head>\n\t<body>\n\t\t<table>\n")
	// dump data
	for key := range m {
		d += fmt.Sprintf("\t\t\t<tr><td>%s</td><td>%s</td></tr>\n", key, m[key])
	}	
	// footer
	d += fmt.Sprintf("\t\t</table>\n\t</body>\n</html>\n")

	// write file
	write_error := ioutil.WriteFile(htmlFile, []byte(d), 0644)
	if write_error !=nil { panic(write_error) }
}

func updateStatus(item string, status string) {
	// check if db exists - if not create a blank one
	_, err := os.Stat(jsonFile)
	if err != nil { createEmptyDb() }	

	// create any empty map to hold json contents
	m := make(map[string]string)

	// read current statues file
	r, read_err := ioutil.ReadFile(jsonFile)
	if read_err != nil { panic(read_err) }	

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
	if len(args) > 1  {
		action := args[1]
		switch action {
		case "update":
			if len(args) == 4 {
				item, status := args[2], args[3]
				updateStatus(item, status)
			} else {
				printHelp()
			}
		case "dump":
			dumpStatus()
		case "help":
			printHelp()
		default:
			printHelp()
		}

	} else {
		printHelp()
	}
}