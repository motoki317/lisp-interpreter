package node

import (
	"errors"
	"fmt"
	"github.com/motoki317/lisp-interpreter/token"
	"regexp"
	"strconv"
)

var (
	EOF       = errors.New("end of input")
	keywords  map[string]bool
	numRegexp = regexp.MustCompile("^-?[0-9]+?(\\.[0-9]*)?$")
)

func init() {
	keywordsList := []string{
		"define",
		"lambda",
		"and",
		"or",
		"if",
		"cond",
		"else",
		"let",
		"let*",
		"quote",
		"set!",
		"begin",
		".",
		"_",
		"define-syntax",
		"syntax-rules",
		"...",
		"delay",
	}
	keywords = make(map[string]bool, len(keywordsList))
	for _, keyword := range keywordsList {
		keywords[keyword] = true
	}
}

type Parser struct {
	t   *token.Tokenizer
	buf *token.Token
}

func NewParser(t *token.Tokenizer) *Parser {
	return &Parser{t: t}
}

// SetTokenizer sets internal tokenizer used by this parser, to start using from the next Next() call.
func (p *Parser) SetTokenizer(t *token.Tokenizer) {
	p.t = t
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

		// quote
		if s == "'" {
			next, err := p.Next()
			if err != nil {
				return nil, fmt.Errorf("an error occurred while parsing quote: %v", err)
			}
			return &Node{
				Type: Branch,
				Children: []*Node{
					{Type: Keyword, Str: "quote"},
					next,
				},
			}, nil
		}

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
		if keywords[s] {
			return &Node{
				Type: Keyword,
				Str:  s,
			}, nil
		}

		// String
		if s[0] == '"' && s[len(s)-1] == '"' {
			return &Node{
				Type: String,
				Str:  s[1 : len(s)-1],
			}, nil
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
				return nil, fmt.Errorf("an error occurred while parsing node: %v", err)
			}
			if stop {
				return node, nil
			}

			child, err := p.Next()
			if err != nil {
				return nil, fmt.Errorf("an error occurred while parsing node: %v", err)
			}
			node.Children = append(node.Children, child)
		}
	}

	return nil, errors.New(fmt.Sprintf("parser internal error: unexpected token: %v", t.Type))
}
