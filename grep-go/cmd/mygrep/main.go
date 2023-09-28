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
	m := newMatcher(line, pattern)
	return m.Match(), nil
}

type matcher struct {
	line          []byte
	pattern       string
	lineOffset    int
	patternOffset int
}

func newMatcher(line []byte, pattern string) matcher {
	return matcher{
		line:    line,
		pattern: pattern,
	}
}

func (m *matcher) Match() bool {
	m.startOfStringAnchor()
	m.endOfStringAnchor()

	for m.lineOffset < len(m.line) && m.patternOffset < len(m.pattern) {
		var ok bool
		switch m.pattern[m.patternOffset] {
		case '.':
			ok = true
		case '\\':
			m.patternOffset++
			switch m.pattern[m.patternOffset] {
			case 'd':
				ok = isNumeric(m.line[m.lineOffset])
			case 'w':
				ok = m.line[m.lineOffset] == '_' || isAlpha(m.line[m.lineOffset]) || isNumeric(m.line[m.lineOffset])
			}
		case '[':
			ok = m.group()
		case '(':
			ok = m.alternation()
		default:
			ok = m.line[m.lineOffset] == m.pattern[m.patternOffset]
			if m.patternOffset+1 < len(m.pattern) {
				switch m.pattern[m.patternOffset+1] {
				case '+':
					ok = m.oneOrMoreTimes()
				case '?':
					ok = m.zeroOrOneTimes()
				}
			}
		}

		m.lineOffset++
		m.patternOffset++
		if !ok {
			m.patternOffset = 0
		}
	}

	return m.patternOffset == len(m.pattern)
}

func (m *matcher) startOfStringAnchor() bool {
	if m.pattern[0] == '^' {
		if m.line[0] != m.pattern[1] {
			return false
		}
		m.pattern = m.pattern[1:]
	}
	return true
}

func (m *matcher) endOfStringAnchor() bool {
	if m.pattern[len(m.pattern)-1] == '$' {
		if m.line[len(m.line)-1] != m.pattern[len(m.pattern)-2] {
			return false
		}
		m.pattern = m.pattern[:len(m.pattern)-1]
	}
	return true
}

func (m *matcher) group() (ok bool) {
	k := strings.IndexByte(m.pattern[m.patternOffset:], ']')
	if k > -1 {
		ok = matchGroup(m.line[m.lineOffset], m.pattern[m.patternOffset:k+1])
		m.patternOffset = k
	}
	return
}

func matchGroup(c byte, pattern string) (ok bool) {
	isGroup := pattern[0] == '[' && pattern[len(pattern)-1] == ']' && len(pattern) > 2
	if !isGroup {
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

func (m *matcher) alternation() bool {
	k := strings.IndexByte(m.pattern, ')')
	if k > -1 {
		patterns := strings.Split(m.pattern[m.patternOffset+1:k], "|")
		for _, pattern := range patterns {
			l, p := m.lineOffset, 0
			for l < len(m.line) && p < len(pattern) {
				if m.line[l] != pattern[p] {
					break
				}
				l++
				p++
			}

			if p == len(pattern) {
				m.patternOffset = k
				m.lineOffset = l
				return true
			}
		}
	}
	return false
}

func (m *matcher) oneOrMoreTimes() bool {
	ok := m.line[m.lineOffset] == m.pattern[m.patternOffset]
	m.patternOffset++
	current := m.line[m.lineOffset]
	for ; m.lineOffset+1 < len(m.line); m.lineOffset++ {
		if m.line[m.lineOffset+1] != current {
			break
		}
	}
	return ok
}

func (m *matcher) zeroOrOneTimes() bool {
	ok := m.line[m.lineOffset] == m.pattern[m.patternOffset]
	if !ok {
		m.pattern = m.pattern[:m.patternOffset] + m.pattern[m.patternOffset+3:]
		m.lineOffset = -1
	} else {
		m.patternOffset++
	}
	return ok
}

func isNumeric(c byte) bool {
	return '0' <= c && c <= '9'
}

func isAlpha(c byte) bool {
	return 'A' <= c && c <= 'Z' || 'a' <= c && c <= 'z'
}
