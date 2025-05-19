package main

import (
	"fmt"
	"io"
	"os"

	"github.com/codecrafters-io/grep-starter-go/app/models"
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

// \d apple should match "1 apple", but not "1 orange".
// \d\d\d apple should match "100 apples", but not "1 apple".
func matchLine(line []byte, pattern string) (bool, error) {
	parsedPattern, err := parsePattern(pattern)
	if err != nil {
		return false, err
	}

	if len(parsedPattern) > 0 && parsedPattern[0].ItemType == models.StartOfLineType {
		return match(line, parsedPattern[1:]), nil
	}

	// Try to match from each starting position
	for i := 0; i <= len(line); i++ {
		if match(line[i:], parsedPattern) {
			return true, nil
		}
	}

	return false, nil
}

// "\d apple" -> [{type: digit}, {type: literal, value: ' '}, {type: literal, value: 'a'}, ...]
func parsePattern(pattern string) ([]models.PatternItem, error) {
	var result []models.PatternItem

	// Check the pattern starts with a line start anchor
	if len(pattern) > 0 && pattern[0] == '^' {
		result = append(result, models.PatternItem{ItemType: models.StartOfLineType})
		pattern = pattern[1:]
	}

	for i := 0; i < len(pattern); i++ {
		if pattern[i] == '\\' && i+1 < len(pattern) {
			i++
			switch pattern[i] {
			case 'd':
				result = append(result, models.PatternItem{ItemType: models.DigitType})
			case 'w':
				result = append(result, models.PatternItem{ItemType: models.WordCharType})
			case '\\':
				result = append(result, models.PatternItem{ItemType: models.LiteralType, Value: '\\'})
			default:
				return nil, fmt.Errorf("unsupported escape sequence: \\%c", pattern[i])
			}
		} else if pattern[i] == '[' {
			// Handle character set
			j := i + 1
			isNegated := false
			if j < len(pattern) && pattern[j] == '^' {
				isNegated = true
				j++
			}

			chars := ""
			for j < len(pattern) && pattern[j] != ']' {
				chars += string(pattern[j])
				j++
			}

			if j >= len(pattern) {
				return nil, fmt.Errorf("unclosed character class")
			}

			//  Skip the entire character set. like: [xyz] or [^xyz]
			i = j

			if isNegated {
				result = append(result, models.PatternItem{ItemType: models.NegatedCharSetType, CharSet: chars})
			} else {
				result = append(result, models.PatternItem{ItemType: models.CharSetType, CharSet: chars})
			}
		} else {
			// Handle normal characters
			result = append(result, models.PatternItem{ItemType: models.LiteralType, Value: pattern[i]})
		}
	}
	return result, nil
}

func match(data []byte, pattern []models.PatternItem) bool {
	if len(pattern) == 0 {
		return true
	}

	if len(data) == 0 {
		return false
	}

	switch pattern[0].ItemType {
	case models.LiteralType:
		if len(data) == 0 || data[0] != pattern[0].Value {
			return false
		}
		return match(data[1:], pattern[1:])
	case models.DigitType:
		if len(data) == 0 || data[0] < '0' || data[0] > '9' {
			return false
		}
		return match(data[1:], pattern[1:])
	case models.WordCharType:
		if len(data) == 0 {
			return false
		}
		ch := data[0]
		isWordChar := (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '_'
		if !isWordChar {
			return false
		}
		return match(data[1:], pattern[1:])
	case models.CharSetType:
		if len(data) == 0 {
			return false
		}
		found := false
		for _, c := range pattern[0].CharSet {
			if data[0] == byte(c) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
		return match(data[1:], pattern[1:])
	case models.NegatedCharSetType:
		if len(data) == 0 {
			return false
		}
		for _, c := range pattern[0].CharSet {
			if data[0] == byte(c) {
				return false
			}
		}
		return match(data[1:], pattern[1:])
	}

	return false
}
