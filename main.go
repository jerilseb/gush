package main

/*
extern void disableRawMode();
extern void enableRawMode();
*/
import "C"

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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
	replLoop()
}

func initShell() {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf(ERRFORMAT, err.Error())
	}

	wdSlice := strings.Split(wd, "/")
	os.Setenv("CWD", wdSlice[len(wdSlice)-1])
}

func replLoop() {
	status := 1
	reader := bufio.NewReader(os.Stdin)

	for status != 0 {
		C.enableRawMode()
		symbol := "\u2713"
		fmt.Printf("\033[36mgush \033[33m%s \033[36m%s \033[m", os.Getenv("CWD"), symbol)
		line, shellEditor := "", false

		for {
			c, _ := reader.ReadByte()

			if shellEditor && c == 13 {
				line = line[:len(line)-1]
				fmt.Println()
				shellEditor = false
				continue
			}
			shellEditor = false

			if c == 27 {
				fmt.Println("Exiting...")
				exit()
			}

			// the enter key was pressed
			if c == 13 {
				fmt.Println()
				break
			}

			// Enter shell editor
			if c == '\\' {
				shellEditor = true
			}
		}
	}
}

func exit() {
	C.disableRawMode()
	os.Exit(0)
}
