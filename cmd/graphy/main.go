package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/code-visible/golang"
)

func main() {
	var (
		project   string
		directory string
	)

	// set up command line arguments
	flag.StringVar(&project, "project", ".", "path of the project")
	flag.StringVar(&directory, "directory", ".", "directory of the project to parse")
	flag.Parse()

	fmt.Printf("graphy: try to parse project (%s) with folder (%s)\n", project, directory)

	// enter the parse progress
	p := golang.NewProject(project, directory)
	p.Initialize()
	p.Parse()

	// marshal the whole project into a json file
	d, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}

	// dump out the json file
	err = os.WriteFile("compiled.json", d, os.ModePerm)
	if err != nil {
		panic(err)
	}
}
