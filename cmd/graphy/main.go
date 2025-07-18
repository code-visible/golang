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
		excludes  string
		types     string
		test      string
	)

	// set up command line arguments
	flag.StringVar(&project, "project", ".", "path of the project")
	flag.StringVar(&directory, "directory", ".", "directory of the project to parse")
	flag.StringVar(&dump, "dump", "parsed.json", "dump path of the project")
	flag.StringVar(&module, "module", "", "module name of the project, it will search go.mod if not provided")
	flag.StringVar(&minify, "minify", "", "keep only the core informations to minimize the output, default minify=0")
	flag.StringVar(&excludes, "excludes", golang.DEFAULT_EXCLUDE_DIRECTORIES, "exclude the given directories, for example: `excludes=test,vendor,docs`")
	flag.StringVar(&types, "types", golang.DEFAULT_INCLUDE_FILE_TYPES, "parse the given types, for example: `types=*.yml,*.json,,README.md,LICENSE,go.mod,go.sum`")
	flag.StringVar(&test, "test", "", "allow test, default test=0")
	flag.Parse()

	if module != "" {
		fmt.Printf("gopher: parsing project (%s) with directory (%s), with customed module (%s)\n", project, directory, module)
	} else {
		fmt.Printf("gopher: parsing project (%s) with directory (%s)\n", project, directory)
	}

	currentPath, err := os.Getwd()
	if err != nil {
		fmt.Println("fail to get current path (pwd)")
		panic(err)
	}
	dumpPath := path.Join(currentPath, dump)

	allowTest := false
	if test != "" && test != "false" && test != "0" {
		allowTest = true
	}

	// enter the parse progress
	p := golang.NewProject(project, directory, excludes, module, types, allowTest)
	p.Initialize()
	p.Parse()

	var d []byte
	if minify != "" && minify != "false" && minify != "0" {
		minifiedProject := &golang.ProjectMinify{
			Name:       p.Name,
			Lang:       p.Lang,
			Parser:     p.Parser,
			Protocol:   p.Protocol,
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
	fmt.Printf("gopher: parse project successfully, dump to %s\n", dumpPath)
}
