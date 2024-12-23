// Package lexer sends a Stream
package lexer

import (
	"fmt"
	"github/dyxgou/redis/pkg/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}

	l.next()

	return l
}

func (l *Lexer) next() {
	if l.readPosition >= len(l.input) {
		l.ch = byte(token.EOF)
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition++
}

// NextToken() provides a stream of tokens from the input given to the lexer
func (l *Lexer) NextToken() token.Token {
	l.skipWhitespace()

	switch l.ch {
	case '\r':
		if ch := l.peekChar(); ch == '\n' {
			return token.New(token.CRLF, "\r\n")
		}
		break
	case '+':
		return token.New(token.STRING, string(l.ch))
	case '-':
		return token.New(token.ERROR, string(l.ch))
	case ':':
		return token.New(token.INTEGER, string(l.ch))
	case '$':
		return token.New(token.BULKSTRING, string(l.ch))
	case '*':
		return token.New(token.ARRAY, string(l.ch))
	case '_':
		return token.New(token.NULL, string(l.ch))
	case '#':
		return token.New(token.BOOLEAN, string(l.ch))
	case ',':
		return token.New(token.FLOAT, string(l.ch))
	case '(':
		return token.New(token.BIGNUMBER, string(l.ch))
	case '!':
		return token.New(token.BULKERROR, string(l.ch))
	case '=':
		return token.New(token.VERTAMINSTRING, string(l.ch))
	case '%':
		return token.New(token.MAPS, string(l.ch))
	case '`':
		return token.New(token.ATTRIBUTES, string(l.ch))
	case '~':
		return token.New(token.SET, string(l.ch))
	case '>':
		return token.New(token.PUSHES, string(l.ch))
	}

	var t token.Token

	if isLetter(l.ch) {
		t.Literal = l.readIdentifier()
		t.Kind = token.LookupIdent(t.Literal)
	} else if isDigit(l.ch) {
		literal, err := l.readNumber()
		t.Literal = literal

		if err != nil {
			t.Kind = token.ILEGAL
		}
	} else {
		t = token.New(token.ILEGAL, string(l.ch))
	}

	l.next()
	return t
}

func isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
}

func isDigit(ch byte) bool {
	return (ch >= '1' && ch <= '9') || ch == '.'
}

func (l *Lexer) skipWhitespace() {
	for l.ch == '_' || l.ch == '\n' || l.ch == '\t' || l.ch == '\r' {
		l.next()
	}
}

func (l *Lexer) readIdentifier() string {
	pos := l.position

	for isLetter(l.ch) {
		l.next()
	}

	return l.input[pos:l.position]
}

func (l *Lexer) readNumber() (string, bool, error) {
	pos := l.position
	var isFloat bool

	for isDigit(l.ch) {
		if l.ch == '.' {
			isFloat = true
		}

		l.next()
	}

	number := l.input[pos:l.position]

	if l.input[l.position-1] == '.' {
		return number, isFloat, fmt.Errorf("unexpected '.' at the end of the number.")
	}

	return number, isFloat, nil
}

func (l *Lexer) peekChar() byte {
	if l.readPosition > len(l.input) {
		return byte(token.EOF)
	}

	return l.input[l.readPosition]
}
