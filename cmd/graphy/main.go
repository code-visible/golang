package main

import (
	"fmt"

	"github.com/code-visible/golang/parser"
)

func main() {
	const project = "/root/go/src/github.com/challenai/golang"
	const direcotry = "testdata/hierarchy"
	p := parser.NewProject(project, direcotry)
	p.Initialize()
	p.Parse()
	fmt.Println(p)
}
