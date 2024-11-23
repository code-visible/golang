package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/code-visible/golang/parser"
)

func main() {
	var project string
	var directory string
	flag.StringVar(&project, "project", ".", "path of the project")
	flag.StringVar(&directory, "directory", ".", "directory of the project to parse")
	flag.Parse()
	fmt.Printf("graphy: try to parse project (%s) with folder (%s)\n", project, directory)
	p := parser.NewProject(project, directory)
	p.Initialize()
	p.Parse()
	d, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("compiled.json", d, os.ModePerm)
	if err != nil {
		panic(err)
	}
}
