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
		{"GETEX", token.New(token.GETEX, "GETEX")},
		{"GETDEL", token.New(token.GETDEL, "GETDEL")},
		{"SET", token.New(token.SET, "SET")},
		{"INCR", token.New(token.INCR, "INCR")},
		{"INCRBY", token.New(token.INCRBY, "INCRBY")},
		{"DECR", token.New(token.DECR, "DECR")},
		{"DECRBY", token.New(token.DECRBY, "DECRBY")},
		{"MGET", token.New(token.MGET, "MGET")},
		{"MSET", token.New(token.MSET, "MSET")},
		{"APPEND", token.New(token.APPEND, "APPEND")},
		{"EXISTS", token.New(token.EXISTS, "EXISTS")},
		{"STRLEN", token.New(token.STRLEN, "STRLEN")},
		{"SUBSTR", token.New(token.SUBSTR, "SUBSTR")},
		{"XX", token.New(token.XX, "XX")},
		{"NX", token.New(token.NX, "NX")},
		{"EX", token.New(token.EX, "EX")},
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
		input: "SET my_key my_value EX 20 NX XX",
		expectedTokens: []token.Token{
			{Kind: token.SET, Literal: "SET"},
			{Kind: token.IDENT, Literal: "my_key"},
			{Kind: token.IDENT, Literal: "my_value"},
			{Kind: token.EX, Literal: "EX"},
			{Kind: token.INTEGER, Literal: "20"},
			{Kind: token.NX, Literal: "NX"},
			{Kind: token.XX, Literal: "XX"},
		},
	}

	l := New(test.input)

	for _, tok := range test.expectedTokens {
		cur := l.NextToken()

		if cur.Kind != tok.Kind {
			t.Errorf("tok kind expected=%d. got=%d (%s)", tok.Kind, cur.Kind, cur.Literal)
		}

		if cur.Literal != tok.Literal {
			t.Errorf("tok literal expected=%s. got=%s", tok.Literal, cur.Literal)
		}
	}
}
