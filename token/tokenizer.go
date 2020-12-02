package token

import (
	"bufio"
	"bytes"
	"io"
	"unicode"
)

type Tokenizer struct {
	sc *bufio.Scanner
}

// NewTokenizer creates a new tokenizer with the given io.Reader.
func NewTokenizer(r io.Reader) *Tokenizer {
	sc := bufio.NewScanner(r)
	isSpaceParComment := func(r rune) bool {
		return unicode.IsSpace(r) || r == '(' || r == ')' || r == ';'
	}
	isNotSpace := func(r rune) bool {
		return !unicode.IsSpace(r)
	}
	// Using recursive (inner) function to make sure to return token whenever possible
	type splitFunc func(data []byte, atEOF bool, f splitFunc) (advance int, token []byte, err error)
	f := func(data []byte, atEOF bool, f splitFunc) (advance int, token []byte, err error) {
		if len(data) == 0 {
			return 0, nil, nil
		}
		// skip spaces
		if i := bytes.IndexFunc(data, isNotSpace); i > 0 {
			adv, tok, err := f(data[i:], atEOF, f)
			return i + adv, tok, err
		} else if i == -1 {
			// data filled with spaces
			return len(data), nil, nil
		}
		// tokenize (special case)
		switch data[0] {
		case '(':
			return 1, data[0:1], nil
		case ')':
			return 1, data[0:1], nil
		case ';':
			// comment: ignore until next newline
			if i := bytes.IndexByte(data, '\n'); i >= 0 {
				adv, tok, err := f(data[i+1:], atEOF, f)
				return i + 1 + adv, tok, err
			} else {
				return 0, nil, nil
			}
		}
		// tokenize by splitting with spaces, parentheses, or semicolon
		if i := bytes.IndexFunc(data, isSpaceParComment); i >= 0 {
			return i, data[0:i], nil
		}
		if atEOF {
			// We're at EOF, so return the data.
			return len(data), data, err
		} else {
			// Doesn't contain full token, request more data.
			return 0, nil, nil
		}
	}
	sc.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		return f(data, atEOF, f)
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

	return &Token{
		Type:   Word,
		String: str,
	}, nil
}
