package node

import (
	"errors"
	"fmt"
	"github.com/motoki317/lisp-interpreter/token"
	"regexp"
	"strconv"
)

var (
	EOF      = errors.New("end of input")
	keywords = []string{
		"define",
		"lambda",
	}
	numRegexp = regexp.MustCompile("^-?[0-9]+?(\\.[0-9]*)?$")
)

type Parser struct {
	t   *token.Tokenizer
	buf *token.Token
}

func NewParser(t *token.Tokenizer) *Parser {
	return &Parser{t: t}
}

func (p *Parser) read() error {
	if p.buf != nil {
		return nil
	}
	t, err := p.t.Next()
	if err != nil {
		return err
	}
	if t == nil {
		return EOF
	}
	p.buf = t
	return nil
}

func (p *Parser) consume(tokenType token.Type) (*token.Token, bool, error) {
	err := p.read()
	if err != nil {
		return nil, false, err
	}
	if p.buf.Type == tokenType {
		t := p.buf
		p.buf = nil
		return t, true, nil
	}
	return nil, false, nil
}

// Next parses tokens from the tokenizer, and returns the next node.
// Returns nil and EOF error on end of input.
func (p *Parser) Next() (*Node, error) {
	err := p.read()
	if err != nil {
		return nil, err
	}

	t := p.buf
	p.buf = nil
	switch t.Type {
	case token.RightPar:
		return nil, errors.New("unexpected right parenthesis")
	case token.Word:
		s := t.String

		// Boolean
		if s == "#t" || s == "#f" {
			return &Node{
				Type: Boolean,
				B:    s == "#t",
			}, nil
		}

		// Number
		if numRegexp.MatchString(s) {
			num, err := strconv.ParseFloat(t.String, 64)
			if err != nil {
				return nil, fmt.Errorf("error while parsing number: %w", err)
			}
			return &Node{
				Type: Number,
				Num:  num,
			}, nil
		}

		// Reserved keywords
		for _, keyword := range keywords {
			if s == keyword {
				return &Node{
					Type: Keyword,
					Str:  s,
				}, nil
			}
		}

		// Other words -> identifier
		return &Node{
			Type: Identifier,
			Str:  s,
		}, nil
	case token.LeftPar:
		node := &Node{
			Type:     Branch,
			Children: make([]*Node, 0),
		}
		for {
			_, stop, err := p.consume(token.RightPar)
			if err != nil {
				return nil, err
			}
			if stop {
				return node, nil
			}

			child, err := p.Next()
			if err != nil {
				return nil, err
			}
			node.Children = append(node.Children, child)
		}
	}

	return nil, errors.New(fmt.Sprintf("parser internal error: unexpected token: %v", t.Type))
}
