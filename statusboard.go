package main

import (
	"os"
	"encoding/json"
	"io/ioutil"
	"log"
	"fmt"
	"time"
	"flag"
)

/***
 - last updated field
***/

var jsonFile = os.Getenv("HOME") + "/" + ".status/statuses.json"
const Version = "0.1"

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\t [-d] Turn on debugging\n")
	fmt.Fprintf(os.Stderr, "\t update OBJECT STATUS - Set objects status\n")
	fmt.Fprintf(os.Stderr, "\t output - Dump current statuses\n")
}

func createEmptyDb() {
	if _, err := os.Stat(os.Getenv("HOME") + "/.status"); err != nil {
		if os.IsNotExist(err) {
			mkdir_err := os.Mkdir(os.Getenv("HOME") + "/.status", 0700)
			if mkdir_err !=nil {panic(mkdir_err)}
		}
	}
	m := make(map[string]string)
	b, json_err := json.Marshal(m)
	if json_err != nil { panic(json_err) }
	write_err := ioutil.WriteFile(jsonFile, b, 0700)
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
	t := time.Now()
	d += fmt.Sprintf("\t\t</table>\n\t<hr>\n\r<i>Last Updated: %s</i>\t\n</body>\n</html>\n", t)

	// write file
	write_error := ioutil.WriteFile(htmlFile, []byte(d), 0700)
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
	write_error := ioutil.WriteFile(jsonFile, b, 0700)
	if write_error !=nil { panic(write_error) }	
}

func main () {
	var debug = flag.Bool("d", false, "Turn debugging on")

	flag.Parse()
	
	if *debug == true { log.Printf("Flags: %v\n", flag.NFlag()) }
	if *debug == true { log.Printf("Args: %v\n", flag.NArg()) }


	if (flag.NFlag() == 0) && (flag.NArg() == 0) {
		Usage()
		return
	} else {
		args := flag.Args()
		action := args[0]
		if *debug == true { log.Printf("Processing Action '%s'", action) }
		switch action {
		case "update":
			if len(args) == 3 {
				object := args[1]
				status := args[2]
				if *debug == true { log.Printf("Action '%s', Object '%s', Status '%s'", action, object, status) }
				updateStatus(object, status)
			} else {
				Usage()
				return
			}
		case "output":
			if len(args) == 1 {
				if *debug == true { log.Printf("Outputting data") }
				outputStatus()
			} else {
				Usage()
				return
			}
		default:
			Usage()
			return
		}
	}
}
