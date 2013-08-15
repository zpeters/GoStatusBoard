package main

import (
	"os"
	"encoding/json"
	"io/ioutil"
	"log"
	"fmt"
	"time"
	"flag"
	"os/exec"
	"strings"
	"syscall"
)

var jsonFile = os.Getenv("HOME") + "/" + ".status/statuses.json"
var debug = false
const Version = "0.2"

type Record struct {
	Object string
	Status string
	Timestamp time.Time
}

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\t [-d] Turn on debugging\n")
	fmt.Fprintf(os.Stderr, "\t update OBJECT STATUS - Set objects status\n")
	fmt.Fprintf(os.Stderr, "\t output - Dump current statuses\n")
	fmt.Fprintf(os.Stderr, "\t test COMMAND OBJECT SUCCESS_STATUS FAIL_STATUS - Run command COMMAND if exit code is 0 update OBJECT with SUCCESS_STATUS, otherwise OBJECT with FAIL_STATUS\n")
}

func createEmptyDb() {
	if _, err := os.Stat(os.Getenv("HOME") + "/.status"); err != nil {
		if os.IsNotExist(err) {
			mkdir_err := os.Mkdir(os.Getenv("HOME") + "/.status", 0700)
			if mkdir_err !=nil {panic(mkdir_err)}
		}
	}
	m := make(map[string]Record)
	b, json_err := json.Marshal(m)
	if json_err != nil { panic(json_err) }
	write_err := ioutil.WriteFile(jsonFile, b, 0700)
	if write_err != nil { panic(write_err) }
}

func outputStatus() {
	// check if db exists - if not create a blank one
	_, err := os.Stat(jsonFile)
	if err != nil { createEmptyDb() }	

	// create any empty map to hold json contents
	m := make(map[string]Record)

	// read current statues file
	r, read_err := ioutil.ReadFile(jsonFile)
	if read_err != nil { panic(read_err) }	

	// unmarshall data read in from file
	unmarshall_error:= json.Unmarshal(r, &m)
	if unmarshall_error != nil { panic(unmarshall_error) }

	// create the header
	d := fmt.Sprintf("<!DOCTYPE html>\n<html lang=\"en\">\n\t<head>\n\t\t<title>StatusBoard</title>\n\t\t<link href=\"css/bootstrap.min.css\" rel=\"stylesheet\" media=\"screen\"\n\t</head>\n\t<body style=\"background-color: #dddddd;\">\n\t<script=\"js/jquery-1.8.3.min.js\"></script>\n\t<script src=\"js/bootstrap.min.js\"></script>\n\t<div class=\"row\"><div class=\"span8\"><h1>StatusBoard</h1></div></div>\n\t<div class=\"row\">\n\t<div class=\"span8\">\n\t\t<table class=\"table table-bordered table-condensed table-striped table-hover\">\n")
	// dump data
	for key := range m {
		rec := m[key]
		if rec.Status == "Success" {
			d += fmt.Sprintf("\t\t\t<tr><td>%s</td><td><span class=\"label label-success\">%s</span></td><td>%s</td></tr>\n", rec.Object, rec.Status, rec.Timestamp)
		} else if rec.Status == "Fail" {
			d += fmt.Sprintf("\t\t\t<tr><td>%s</td><td><span class=\"label label-important\">%s</span></td><td>%s</td></tr>\n", rec.Object, rec.Status, rec.Timestamp)
		} else {
			d += fmt.Sprintf("\t\t\t<tr><td>%s</td><td><span class=\"label\">%s</span></td><td>%s</td></tr>\n", rec.Object, rec.Status, rec.Timestamp)
		}
	}	
	// footer
	t := time.Now()
	d += fmt.Sprintf("\t\t</table>\n\t<hr>\n\r<i>Last Updated: %s</i>\t\n</div>\n\t</div>\n\t</body>\n</html>\n", t)

	fmt.Printf(d)
}

func updateStatus(object string, status string) {
	// check if db exists - if not create a blank one
	_, err := os.Stat(jsonFile)
	if err != nil { createEmptyDb() }	

	// create any empty map to hold json contents
	m := make(map[string]Record)

	// read current statues file
	r, read_err := ioutil.ReadFile(jsonFile)
	if read_err != nil { panic(read_err) }	

	// unmarshall data read in from file
	unmarshall_error:= json.Unmarshal(r, &m)
	if unmarshall_error != nil { panic(unmarshall_error) }

	// create a record from the passed in data
	rec := Record{
		Object: object,
		Status: status,
		Timestamp: time.Now(),
	}

	// set map with the passed in data
	m[rec.Object] = rec

	// marshall data
	b, marshall_error := json.Marshal(m)
	if marshall_error != nil { panic(marshall_error) }

	// write json data tile file
	write_error := ioutil.WriteFile(jsonFile, b, 0700)
	if write_error !=nil { panic(write_error) }	
}

func runCommandAndReturnStatus(command string)  int {
	var returnCode = 0
	
	cmdString := strings.Fields(command)
	if debug { log.Printf("Cmd: '%s'\n", cmdString) }

	out, err := exec.Command(cmdString[0], cmdString[1:]...).Output()
	if err != nil {
		fmt.Printf("Err: %s\n", err)
		if debug { log.Printf("Err: %s\n", err) }
		// Couldn't run the command, returning our own error code 1
		if debug { log.Printf("Setting error code to 1\n") }
		returnCode = 1
		return returnCode
	}


	// get the OS exit code
	if msg, ok := err.(*exec.ExitError); ok {
		code := msg.Sys().(syscall.WaitStatus).ExitStatus()
		fmt.Printf("Output: %s\n", out)
		if debug { log.Printf("Output: %s\n", out) }
		if debug { log.Printf("Error code not 0: %s\n", code) }
		returnCode = code
	} else {
		if debug { log.Printf("Output: %s\n", out) }
		if debug { log.Printf("Error code is 0\n") }
		returnCode = 0
	}
	return returnCode
}

func main () {
	flag.BoolVar(&debug, "d", false, "Turn debugging on")

	flag.Parse()
	
	if debug { log.Printf("Flags: %v\n", flag.NFlag()) }
	if debug { log.Printf("Args: %v\n", flag.NArg()) }


	if (flag.NFlag() == 0) && (flag.NArg() == 0) {
		Usage()
		return
	} else {
		args := flag.Args()
		action := args[0]
		if debug { log.Printf("Processing Action '%s'", action) }
		switch action {
		case "update":
			if len(args) == 3 {
				object := args[1]
				status := args[2]
				if debug { log.Printf("Action '%s', Object '%s', Status '%s'", action, object, status) }
				updateStatus(object, status)
			} else {
				Usage()
				return
			}
		case "output":
			if len(args) == 1 {
				if debug { log.Printf("Outputting data") }
				outputStatus()
			} else {
				Usage()
				return
			}

		case "test":
			if len(args) == 5 {
				command := args[1]
				object := args[2]
				success_status := args[3]
				fail_status := args[4]
				if debug { log.Printf("Command: '%s'\n", command) }
				if debug { log.Printf("Object: '%s'\n", object) }
				if debug { log.Printf("Success: '%s'\n", success_status) }
				if debug { log.Printf("Fail: '%s'\n", fail_status) }
				if debug { log.Printf("Getting results for '%s'\b", command) }
				returnCode := runCommandAndReturnStatus(command)
				if debug { log.Printf("Return code: %d\n", returnCode) }
				if returnCode == 0 {
					if debug { log.Printf("Updating '%s' with '%s'\n", object, success_status) }
					updateStatus(object, success_status)
				} else {
					if debug { log.Printf("Updating '%s' with '%s'\n", object, fail_status) }
					updateStatus(object, fail_status)
				}
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
