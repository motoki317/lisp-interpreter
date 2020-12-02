package token

import (
	"reflect"
	"strings"
	"testing"
)

func readAllTokens(t *testing.T, tokenizer *Tokenizer) []Token {
	t.Helper()

	tokens := make([]Token, 0)
	for {
		token, err := tokenizer.Next()
		if err != nil {
			t.Fatalf("error while reading tokens: %v", err)
		}
		if token == nil {
			break
		}
		tokens = append(tokens, *token)
	}
	return tokens
}

func TestTokenizer(t *testing.T) {
	tests := []struct {
		name   string
		string string
		want   []Token
	}{
		{
			name:   "po",
			string: "po po2 po3 po4",
			want: []Token{
				{
					Type:   Identifier,
					String: "po",
				},
				{
					Type:   Identifier,
					String: "po2",
				},
				{
					Type:   Identifier,
					String: "po3",
				},
				{
					Type:   Identifier,
					String: "po4",
				},
			},
		},
		{
			name:   "numbers",
			string: "po 23 45.6 -78",
			want: []Token{
				{
					Type:   Identifier,
					String: "po",
				},
				{
					Type:   Number,
					String: "23",
				},
				{
					Type:   Number,
					String: "45.6",
				},
				{
					Type:   Number,
					String: "-78",
				},
			},
		},
		{
			name:   "parentheses",
			string: "(define po 42.3)",
			want: []Token{
				{
					Type:   LeftPar,
					String: "",
				},
				{
					Type:   Keyword,
					String: "define",
				},
				{
					Type:   Identifier,
					String: "po",
				},
				{
					Type:   Number,
					String: "42.3",
				},
				{
					Type:   RightPar,
					String: "",
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tokenizer := NewTokenizer(strings.NewReader(tt.string))
			if got := readAllTokens(t, tokenizer); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
