package serializer

import (
	"github/dyxgou/redis/pkg/token"
	"testing"
)

func TestWriteCRLF(t *testing.T) {
	w := newWriter()

	w.writeCRLF()

	if b := w.body.String(); b != token.EndCRLF {
		t.Errorf("writer expected=%q. got=%q", b, token.EndCRLF)
	}
}

func TestWriteWord(t *testing.T) {
	tt := struct {
		input    token.Token
		expected string
	}{
		input:    token.New(token.GET, "GET"),
		expected: "$3\r\nGET\r\n",
	}

	w := newWriter()
	if err := w.writeWord(tt.input); err != nil {
		t.Error(err)
	}

	if b := w.body.String(); b != tt.expected {
		t.Errorf("writer expected=%q. got=%q", tt.expected, b)
	}
}

func TestWriteLen(t *testing.T) {
	w := newWriter()
	expected := "*0\r\n"

	w.writeLen()

	if w.head.String() != expected {
		t.Errorf("head expected=%q. got=%q", expected, w.head.String())
	}
}

func TestWriteNumber(t *testing.T) {
	tests := []struct {
		input    token.Token
		expected string
	}{
		{token.New(token.INTEGER, "123"), ":123\r\n"},
		{token.New(token.FLOAT, "123.123"), ",123.123\r\n"},
	}

	for _, tt := range tests {
		w := newWriter()
		w.writeNumber(tt.input)

		if w.body.String() != tt.expected {
			t.Errorf("body expected=%q. got=%q", tt.expected, w.body.String())
		}
	}
}

func TestWriteBoolean(t *testing.T) {
	tests := []struct {
		cur, next token.Token
		expected  string
	}{
		{
			cur:      token.New(token.BOOLEAN, "#"),
			next:     token.New(token.IDENT, "t"),
			expected: "#t\r\n",
		},
		{
			cur:      token.New(token.BOOLEAN, "#"),
			next:     token.New(token.IDENT, "f"),
			expected: "#f\r\n",
		},
	}

	for _, tt := range tests {
		w := newWriter()
		if err := w.writeBool(tt.cur, tt.next); err != nil {
			t.Error(err)
			return
		}

		if w.body.String() != tt.expected {
			t.Errorf("body expected=%q. got=%q", tt.expected, w.body.String())
		}
	}
}
