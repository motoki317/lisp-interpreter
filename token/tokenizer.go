package token

import (
	"bufio"
	"bytes"
	"io"
	"regexp"
	"unicode"
)

var (
	keywords = []string{
		"define",
		"lambda",
	}
	numRegexp = regexp.MustCompile("^-?[0-9]+?(\\.[0-9]*)?$")
)

type Tokenizer struct {
	sc *bufio.Scanner
}

// NewTokenizer creates a new tokenizer with the given io.Reader.
func NewTokenizer(r io.Reader) *Tokenizer {
	sc := bufio.NewScanner(r)
	isSpaceOrPar := func(r rune) bool {
		return unicode.IsSpace(r) || r == '(' || r == ')'
	}
	isNotSpace := func(r rune) bool {
		return !unicode.IsSpace(r)
	}
	sc.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if len(data) == 0 {
			return 0, nil, nil
		}
		// number of leading spaces
		leadingSpaces := bytes.IndexFunc(data, isNotSpace)
		if leadingSpaces == -1 {
			// data filled with spaces
			return len(data), nil, nil
		}
		// tokenize (special case)
		switch data[leadingSpaces] {
		case '(':
			return leadingSpaces + 1, data[leadingSpaces : leadingSpaces+1], nil
		case ')':
			return leadingSpaces + 1, data[leadingSpaces : leadingSpaces+1], nil
		}
		// tokenize by splitting with spaces or parentheses
		if i := bytes.IndexFunc(data[leadingSpaces:], isSpaceOrPar); i >= 0 {
			return leadingSpaces + i, data[leadingSpaces : leadingSpaces+i], nil
		}
		if atEOF {
			// We're at EOF, so return the data.
			return len(data), data[leadingSpaces:], err
		} else {
			// Doesn't contain full token, request more data.
			return leadingSpaces, nil, nil
		}
	})
	return &Tokenizer{
		sc: sc,
	}
}

// Next returns the next token, or if any, errors.
// Returns nil, nil on end of the input.
func (t *Tokenizer) Next() (*Token, error) {
	if !t.sc.Scan() {
		return nil, t.sc.Err()
	}

	str := t.sc.Text()
	switch str {
	case "(":
		return &Token{
			Type:   LeftPar,
			String: "",
		}, nil
	case ")":
		return &Token{
			Type:   RightPar,
			String: "",
		}, nil
	}

	if numRegexp.MatchString(str) {
		return &Token{
			Type:   Number,
			String: str,
		}, nil
	}

	for _, keyword := range keywords {
		if str == keyword {
			return &Token{
				Type:   Keyword,
				String: str,
			}, nil
		}
	}

	return &Token{
		Type:   Identifier,
		String: str,
	}, nil
}
