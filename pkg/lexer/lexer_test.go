package lexer

import (
	"github/dyxgou/redis/pkg/token"
	"testing"
)

func FuzzNextChar(f *testing.F) {
	f.Add("*1\r\n$3\r\nGET\r\n")

	f.Fuzz(func(t *testing.T, a string) {
		l := New(a)

		for i := 0; i < len(a); i++ {
			ch := a[i]
			if l.ch != ch {
				t.Errorf("l.ch expected=%q. got=%q", l.ch, ch)
			}

			l.next()
		}
	})
}

func TestTokenizeCRLF(t *testing.T) {
	tt := struct {
		input         string
		expectedToken token.Token
	}{
		input:         "\r\n",
		expectedToken: token.New(token.CRLF, token.EndCRLF),
	}

	l := New(tt.input)

	tok := l.NextToken()

	if tok.Kind != tt.expectedToken.Kind {
		t.Errorf(
			"token kind expected=%d. got=%d (%q)", tt.expectedToken.Kind, tok.Kind, tok.Literal,
		)
	}

	if tok.Literal != tt.expectedToken.Literal {
		t.Errorf(
			"token literal expected=%q. got=%q", tt.expectedToken.Literal, tok.Literal,
		)
	}
}

func TestTokenizeNumber(t *testing.T) {
	tests := []struct {
		input         string
		expectedToken token.Token
	}{
		{"123123", token.New(token.INTEGER, "123123")},
		{"312", token.New(token.INTEGER, "312")},
		{"312.123", token.New(token.FLOAT, "312.123")},
		{"312.", token.New(token.ILLEGAL, "")},
	}

	for _, tt := range tests {
		l := New(tt.input)

		tok := l.NextToken()

		if tok.Kind != tt.expectedToken.Kind {
			t.Errorf(
				"token kind expected=%d. got=%d (%q)", tt.expectedToken.Kind, tok.Kind, tok.Literal,
			)
		}

		if tok.Literal != tt.expectedToken.Literal {
			t.Errorf(
				"token literal expected=%q. got=%q", tt.expectedToken.Literal, tok.Literal,
			)
		}
	}
}

func TestNextToken(t *testing.T) {
	tests := []struct {
		input          string
		expectedTokens []token.Token
	}{
		{
			input: "*2\r\n$3\r\nGET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n",
			expectedTokens: []token.Token{
				token.New(token.ARRAY, "*"),
				token.New(token.INTEGER, "2"),
				token.New(token.CRLF, token.EndCRLF),
				token.New(token.BULKSTRING, "$"),
				token.New(token.INTEGER, "3"),
				token.New(token.CRLF, token.EndCRLF),
				token.New(token.GET, "GET"),
				token.New(token.CRLF, token.EndCRLF),
				token.New(token.BULKSTRING, "$"),
				token.New(token.INTEGER, "3"),
				token.New(token.CRLF, token.EndCRLF),
				token.New(token.IDENT, "key"),
				token.New(token.CRLF, token.EndCRLF),
				token.New(token.BULKSTRING, "$"),
				token.New(token.INTEGER, "5"),
				token.New(token.CRLF, token.EndCRLF),
				token.New(token.IDENT, "value"),
				token.New(token.CRLF, token.EndCRLF),
			},
		},
	}

	for _, tt := range tests {
		l := New(tt.input)

		for _, expTok := range tt.expectedTokens {
			tok := l.NextToken()

			if tok.Kind != expTok.Kind {
				t.Errorf("token kind expected=%d. got=%d (%q)", expTok.Kind, tok.Kind, tok.Literal)
			}

			if tok.Literal != expTok.Literal {
				t.Errorf("token literal expected=%q. got=%q", expTok.Literal, tok.Literal)
			}
		}
	}
}
