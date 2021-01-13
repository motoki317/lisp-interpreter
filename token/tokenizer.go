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

// isSpaceParCommentQuote returns true if r is one of: space, (, ), ;, ', or "
func isSpaceParCommentQuote(r rune) bool {
	return unicode.IsSpace(r) || r == '(' || r == ')' || r == ';' || r == '\'' || r == '"'
}

// isNotSpace returns true if r is not space.
func isNotSpace(r rune) bool {
	return !unicode.IsSpace(r)
}

// splitFunc parses the possibly incomplete input and advances to the next token if possible.
// See: bufio.SplitFunc
func splitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if len(data) == 0 {
		return 0, nil, nil
	}
	// skip spaces if any
	if i := bytes.IndexFunc(data, isNotSpace); i > 0 {
		adv, tok, err := splitFunc(data[i:], atEOF)
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
			// newline exists, parse next token if any.
			adv, tok, err := splitFunc(data[i+1:], atEOF)
			return i + 1 + adv, tok, err
		} else {
			// no newline, request more data.
			return 0, nil, nil
		}
	case '\'':
		return 1, data[0:1], nil
	case '"':
		// string: read till next double quote
		if i := bytes.IndexByte(data[1:], '"'); i >= 0 {
			return i + 2, data[0 : i+2], nil
		} else {
			return 0, nil, nil
		}
	}
	// tokenize by splitting with spaces, parentheses, or semicolon
	if i := bytes.IndexFunc(data, isSpaceParCommentQuote); i >= 0 {
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

// NewTokenizer creates a new tokenizer with the given io.Reader.
func NewTokenizer(r io.Reader) *Tokenizer {
	sc := bufio.NewScanner(r)
	sc.Split(splitFunc)
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
