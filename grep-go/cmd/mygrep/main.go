package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

// Usage: echo <input_text> | your_grep.sh -E <pattern>
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
}

func matchLine(line []byte, pattern string) (bool, error) {
	if pattern == "\\d" {
		return isNumeric(line), nil
	}

	if pattern == "\\w" {
		return isAlphanumeric(line), nil
	}

	if isPositiveCharacterGroup(pattern) {
		return bytes.ContainsAny(line, pattern[1:len(pattern)-1]), nil
	}

	if isNegativeCharacterGroup(pattern) {
		return !bytes.ContainsAny(line, pattern[2:len(pattern)-1]), nil
	}

	return bytes.ContainsAny(line, pattern), nil
}

func isNumeric(line []byte) bool {
	for _, c := range line {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func isAlphanumeric(line []byte) bool {
	for _, c := range line {
		if !(c == '_' ||
			'0' <= c && c <= '9' ||
			'A' <= c && c <= 'Z' ||
			'a' <= c && c <= 'z') {
			return false
		}
	}
	return true
}

func isPositiveCharacterGroup(pattern string) bool {
	return isCharacterGroup(pattern) && pattern[1] != '^'
}

func isNegativeCharacterGroup(pattern string) bool {
	return isCharacterGroup(pattern) && pattern[1] == '^'
}

func isCharacterGroup(pattern string) bool {
	return pattern[0] == '[' && pattern[len(pattern)-1] == ']' && len(pattern) > 2
}
