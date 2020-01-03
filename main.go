package main

import (
	"fmt"
	"os"
)

var (
	HISTSIZE  = 100
	HISTFILE  string
	HISTMEM   []string
	HISTCOUNT int
	HISTLINE  string
	CONFIG    string
	aliases   map[string]string
)

const (
	TOKDELIM  = " \t\r\n\a"
	ERRFORMAT = "gush: %s\n"
)

func main() {
	initShell()
}

func initShell() {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf(ERRFORMAT, err.Error())
	}

	println(wd)
}
