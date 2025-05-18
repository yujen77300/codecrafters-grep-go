package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"
)

// Ensures gofmt doesn't remove the "bytes" import above (feel free to remove this!)
// var _ = bytes.ContainsAnyc

// Usage: echo <input_text> | your_program.sh -E <pattern>
func main() {
	if len(os.Args) < 3 || os.Args[1] != "-E" {
		fmt.Fprintf(os.Stderr, "usage: mygrep -E <pattern>\n")
		os.Exit(2) // 1 means no lines were selected, >1 means error
	}

	pattern := os.Args[2]

	line, err := io.ReadAll(os.Stdin) // assume we're only dealing with a single line
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: read input text: %v\n", err)
		os.Exit(2)
	}

	ok, err := matchLine(line, pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	if !ok {
		os.Exit(1)
	}

	// default exit code is 0 which means success
}

func matchLine(line []byte, pattern string) (bool, error) {
	var ok bool
	if pattern == "\\d" {
		pattern = strings.Replace(pattern, "\\d", "0123456789", -1)
		ok = bytes.ContainsAny(line, pattern)
	} else if pattern == "\\w" {
		pattern = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
		ok = bytes.ContainsAny(line, pattern)
	} else if len(pattern) >= 2 && pattern[0] == '[' && pattern[1] == '^' && pattern[len(pattern)-1] == ']' {
		pattern := pattern[2 : len(pattern)-1]
		for _, b := range line {
			if !bytes.ContainsAny([]byte{b}, pattern) {
				return true, nil
			}
		}
		return false, nil
	} else if len(pattern) >= 2 && pattern[0] == '[' && pattern[len(pattern)-1] == ']' {
		pattern = pattern[1 : len(pattern)-1]
		ok = bytes.ContainsAny(line, pattern)
	} else if utf8.RuneCountInString(strings.Trim(pattern, "\\n")) == 1 {
		ok = bytes.Contains(line, []byte(pattern))
	} else {
		return false, fmt.Errorf("unsupported pattern: %q", pattern)
	}

	return ok, nil
}
