package lexer

import (
	"github/dyxgou/redis/pkg/token"
	"testing"
)

func TestInputToken(t *testing.T) {
	tests := []struct {
		input    string
		expected token.Token
	}{
		{"+", token.New(token.STRING, "+")},
		{"-", token.New(token.ERROR, "-")},
		{":", token.New(token.INTEGER, ":")},
		{"GET", token.New(token.GET, "GET")},
		{"GETSET", token.New(token.GETSET, "GETSET")},
		{"GETDEL", token.New(token.GETDEL, "GETDEL")},
		{"SET", token.New(token.SET, "SET")},
	}

	for _, tt := range tests {
		l := New(tt.input)

		tok := l.NextToken()

		if tok.Literal != tt.expected.Literal {
			t.Errorf("tok literal expected=%s. got=%s", tt.expected.Literal, tok.Literal)
		}

		if tok.Kind != tt.expected.Kind {
			t.Errorf("tok kind expected=%d. got=%d", tt.expected.Kind, tok.Kind)
		}
	}
}

func TestIlegalToken(t *testing.T) {
	input := "@^&)'"

	l := New(input)

	for l.ch != byte(token.EOF) {
		tok := l.NextToken()

		if tok.Kind != token.ILEGAL {
			t.Errorf("token kind ILEGAL expected=%d, got=%d (%s)", token.ILEGAL, tok.Kind, tok.Literal)
		}
	}
}

func TestCommandTokenized(t *testing.T) {
	test := struct {
		input          string
		expectedTokens []token.Token
	}{
		input: "GET my_key my_value",
		expectedTokens: []token.Token{
			{Kind: token.GET, Literal: "GET"},
			{Kind: token.IDENT, Literal: "my_key"},
			{Kind: token.IDENT, Literal: "my_value"},
		},
	}

	l := New(test.input)

	for _, tok := range test.expectedTokens {
		cur := l.NextToken()

		if cur.Kind != tok.Kind {
			t.Errorf("tok kind expected=%d. got=%d", tok.Kind, cur.Kind)
		}

		if cur.Literal != tok.Literal {
			t.Errorf("tok literal expected=%s. got=%s", tok.Literal, cur.Literal)
		}
	}
}
