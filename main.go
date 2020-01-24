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
)

func main() {
	for {
		fmt.Printf("\033[36mgush \033[36m\u2713 \033[m")
		reader := bufio.NewReader(os.Stdin)
		C.enableRawMode()

		line, cursorPos := "", 0

		for {
			c, _ := reader.ReadByte()

			// Ctrl+C is pressed
			if c == 3 {
				fmt.Println("Exiting...")
				exit()
			}

			// the enter key was pressed
			if c == 10 {
				fmt.Println()
				break
			}

			// Special control key was pressed
			if c == 27 {
				c1, _ := reader.ReadByte()
				if c1 == '[' {
					c2, _ := reader.ReadByte()
					switch c2 {
					case 'C':
						if cursorPos < len(line) {
							fmt.Printf("\033[C")
							cursorPos++
						}
					case 'D':
						if cursorPos > 0 {
							fmt.Printf("\033[D")
							cursorPos--
						}
					}
				}
				continue
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
		}
	}
}

func exit() {
	C.disableRawMode()
	os.Exit(0)
}
