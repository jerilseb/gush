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
		line, cursorPos, shellEditor := "", 0, false

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

			// backspace was pressed
			if c == 127 {
				if cursorPos > 0 {
					if cursorPos != len(line) {
						temp, oldLength := line[cursorPos:], len(line)
						fmt.Printf("\b\033[K%s", temp)
						for oldLength != cursorPos {
							fmt.Printf("\033[D")
							oldLength--
						}
						line = line[:cursorPos-1] + temp
						cursorPos--
					} else {
						fmt.Print("\b\033[K")
						line = line[:len(line)-1]
						cursorPos--
					}
				}
				continue
			}

			// Any normal character
			if cursorPos == len(line) {
				fmt.Printf("%c", c)
				line += string(c)
				cursorPos = len(line)
			} else {
				temp, oldLength := line[cursorPos:], len(line)
				fmt.Printf("\033[K%c%s", c, temp)
				for oldLength != cursorPos {
					fmt.Printf("\033[D")
					oldLength--
				}
				line = line[:cursorPos] + string(c) + temp
				cursorPos++
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
