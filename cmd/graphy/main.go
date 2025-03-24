package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/code-visible/golang"
)

func main() {
	var (
		project   string
		directory string
		dump      string
		module    string
		minify    string
		// TODO: do we need excludes feature ?
		// excludes  string
	)

	// set up command line arguments
	flag.StringVar(&project, "project", ".", "path of the project")
	flag.StringVar(&directory, "directory", ".", "directory of the project to parse")
	flag.StringVar(&dump, "dump", "parsed.json", "dump path of the project")
	flag.StringVar(&module, "module", "", "module name of the project, it will search go.mod if not provided")
	flag.StringVar(&minify, "minify", "", "keep only the core informations to minimize the output")
	// flag.StringVar(&excludes, "excludes", "", "exclude the given directories")
	flag.Parse()

	fmt.Printf("graphy: try to parse project (%s) with folder (%s), dump to (%s)\n", project, directory, dump)

	currentPath, err := os.Getwd()
	if err != nil {
		fmt.Println("fail to get current path (pwd)")
		panic(err)
	}
	dumpPath := path.Join(currentPath, dump)

	// enter the parse progress
	p := golang.NewProject(project, directory, module)
	p.Initialize()
	p.Parse()

	var d []byte
	if minify != "" {
		minifiedProject := &golang.ProjectMinify{
			Name:       p.Name,
			Lang:       p.Lang,
			Parser:     p.Parser,
			Timestamp:  p.Timestamp,
			Repository: p.Repository,
			Version:    p.Version,
			Typ:        golang.PARSE_TYPE_MINIFY,
			Abstracts:  uint64(len(p.Abstracts)),
			Callables:  uint64(len(p.Callables)),
			Calls:      uint64(len(p.Calls)),
			References: uint64(len(p.References)),
			Deps:       uint64(len(p.Deps)),
			Pkgs:       p.Pkgs,
			Files:      p.Files,
		}
		// marshal the whole project into a json file
		d, err = json.Marshal(minifiedProject)
	} else {
		// marshal the whole project into a json file
		d, err = json.Marshal(p)
	}
	if err != nil {
		panic(err)
	}
	// dump out the json file
	err = os.WriteFile(dumpPath, d, os.ModePerm)
	if err != nil {
		fmt.Println("fail to dump result to given dump path")
		panic(err)
	}

}
