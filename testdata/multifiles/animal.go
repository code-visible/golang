package main

import (
	alias "fmt"
)

type Animal struct {
	Name string
}

func (ani Animal) Run() {
	alias.Printf("%s is Running...", ani.Name)
}
