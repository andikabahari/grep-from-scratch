package main

import (
	"fmt"
	"io"
	"os"
	"strings"
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
	if pattern[0] == '^' {
		if line[0] != pattern[1] {
			return false, nil
		}
		pattern = pattern[1:]
	}

	if pattern[len(pattern)-1] == '$' {
		if line[len(line)-1] != pattern[len(pattern)-2] {
			return false, nil
		}
		pattern = pattern[:len(pattern)-1]
	}

	l, p := 0, 0
	for l < len(line) && p < len(pattern) {
		var ok bool
		if pattern[p] == '\\' {
			p++
			switch pattern[p] {
			case 'd':
				ok = isNumeric(line[l])
			case 'w':
				ok = line[l] == '_' || isAlpha(line[l]) || isNumeric(line[l])
			}
		} else if pattern[p] == '[' {
			k := strings.IndexByte(pattern[p:], ']')
			if k > -1 {
				ok = matchGroup(line[l], pattern[p:k+1])
				p = k
			}
		} else {
			if pattern[p] == '.' {
				ok = true
			} else {
				ok = line[l] == pattern[p]
			}
			if p+1 < len(pattern) {
				if pattern[p+1] == '+' {
					p++
					current := line[l]
					for ; l+1 < len(line); l++ {
						if line[l+1] != current {
							break
						}
					}
				} else if pattern[p+1] == '?' {
					if !ok {
						pattern = pattern[:p] + pattern[p+3:]
						l = -1
					} else {
						p++
					}
				}
			}
		}

		l++
		p++
		if !ok {
			p = 0
		}
	}

	return p == len(pattern), nil
}

func isNumeric(c byte) bool {
	return '0' <= c && c <= '9'
}

func isAlpha(c byte) bool {
	return 'A' <= c && c <= 'Z' || 'a' <= c && c <= 'z'
}

func matchGroup(c byte, pattern string) (ok bool) {
	if !isGroup(pattern) {
		return
	}
	for p := 0; p < len(pattern); p++ {
		if c == pattern[p] {
			ok = true
			break
		}
	}
	if pattern[1] == '^' {
		return !ok
	}
	return
}

func isGroup(pattern string) bool {
	return pattern[0] == '[' && pattern[len(pattern)-1] == ']' && len(pattern) > 2
}
