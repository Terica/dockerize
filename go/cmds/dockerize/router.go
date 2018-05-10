package main

import (
	"fmt"
	"local/dockerize/docker"
	"local/dockerize/progress"
	"os"
	"path/filepath"
)

// programMode is used to store a lookup table of command names we process in particular functions
// _ represents a default value

type modeFunc func() int

var (
	// This means we have an initialized map before the init() functions run.
	programMode = map[string]modeFunc{}
	programName string
	programPath string
)

func init() {
	// Each dockerize module should add the "modes" aka dockerize <foo> that it supports.
	// That's how we can make sure that stuff only intended for a particular OS can only work on that OS.
	// If you +build !linux then the init() for that .go source won't be run either.
	programMode["_"] = modeIndirect
	programPath, programName = ParseProgram(os.Args[0])
}

// ParseProgram converts an arg0 to directory and program name. Stripping .exe if necessary
func ParseProgram(arg string) (string, string) {
	abs, _ := filepath.Abs(arg)
	dir := filepath.Dir(abs)
	name := filepath.Base(abs)
	if ext := filepath.Ext(abs); ext == ".exe" {
		name = name[0 : len(name)-len(ext)]
	}
	return dir, name
}

func main() {
	var ret int
	if fn, ok := programMode[programName]; ok {
		ret = fn()
	} else {
		ret = programMode["_"]()
	}
	os.Exit(ret)
}

func modeIndirect() int {
	cli := docker.Connect()
	prog := progress.New(fmt.Sprintf("Starting Container for %s", programName), 0, 0)
	//fmt.Printf("It was %s\n", programName)
	cli.Pull(programName,"latest", prog)
	containers := cli.PStat(map[string][]string{"name": []string{programName+"_latest"}})
	var id string
	if len(containers)>0 {
		id = containers[0].ID
	} else {
		if ret, err := cli.Run(programName+":latest", programName+"_latest", []string{"/c/Users/faye:/c/Users/faye"},nil,[]string{"cat"}); err != nil {
			fmt.Printf("%s\n", error.Error(err))
		} else {
			id = ret
		}
	}
	//fmt.Printf("running %s\n", id)
	wd, _ := os.Getwd()
	exit, _ := cli.Exec(id, nil, wd, os.Args[1:])
	return exit
}