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
				{Type: Word, String: "po"},
				{Type: Word, String: "po2"},
				{Type: Word, String: "po3"},
				{Type: Word, String: "po4"},
			},
		},
		{
			name:   "numbers",
			string: "po 23 45.6 -78",
			want: []Token{
				{Type: Word, String: "po"},
				{Type: Word, String: "23"},
				{Type: Word, String: "45.6"},
				{Type: Word, String: "-78"},
			},
		},
		{
			name:   "parentheses",
			string: "(define po 42.3)",
			want: []Token{
				{Type: LeftPar},
				{Type: Word, String: "define"},
				{Type: Word, String: "po"},
				{Type: Word, String: "42.3"},
				{Type: RightPar},
			},
		},
		{
			name: "comment",
			string: "po ; this is a comment\n" +
				"; another comment\n" +
				"123 ; last comment",
			want: []Token{
				{Type: Word, String: "po"},
				{Type: Word, String: "123"},
			},
		},
		{
			name:   "quote",
			string: "'() '1",
			want: []Token{
				{Type: Word, String: "'"},
				{Type: LeftPar},
				{Type: RightPar},
				{Type: Word, String: "'"},
				{Type: Word, String: "1"},
			},
		},
		{
			name:   "string",
			string: "po \"po po\" po",
			want: []Token{
				{Type: Word, String: "po"},
				{Type: Word, String: "\"po po\""},
				{Type: Word, String: "po"},
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
