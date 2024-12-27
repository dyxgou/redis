package lexer

import (
	"fmt"
	"github/dyxgou/redis/pkg/token"
	"log/slog"
)

type Lexer struct {
	input   string
	pos     int
	readPos int

	ch byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input, pos: 0, readPos: 0}
	l.next()

	return l
}

func (l *Lexer) next() {
	if l.readPos >= len(l.input) {
		l.ch = byte(token.EOF)
		return
	} else {
		l.ch = l.input[l.readPos]
	}

	l.pos = l.readPos
	l.readPos++
}

func (l *Lexer) NextToken() token.Token {
	l.skipWhitespaces()

	if l.ch == '\r' && l.peekChar() == '\n' {
		l.next()
		return token.New(token.CRLF, token.EndCRLF)
	}

	if k, ok := token.GetKindWithSymbol(l.ch); ok {
		t := token.New(k, string(l.ch))
		l.next()
		return t
	}

	var t token.Token

	if isLetter(l.ch) {
		t.Literal = l.readIdent()
		t.Kind = token.LookupIdent(t.Literal)
		return t
	} else if isDigit(l.ch) {
		num, isFloat, err := l.readNumber()

		if err != nil {
			slog.Error("tokenizing number failed", "err", err)
			t.Literal = ""
			t.Kind = token.ILLEGAL
			return t
		}

		if isFloat {
			t.Kind = token.FLOAT
		} else {
			t.Kind = token.INTEGER
		}

		t.Literal = num
		return t
	}

	t.Literal = ""
	t.Kind = token.ILLEGAL
	return t
}

func isLetter(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9' || ch == '.'
}

func (l *Lexer) readIdent() string {
	pos := l.pos

	for isLetter(l.ch) {
		l.next()
	}

	return l.input[pos:l.pos]
}

func (l *Lexer) readNumber() (string, bool, error) {
	pos := l.pos
	var isFloat bool

	for isDigit(l.ch) {
		slog.Info("digit", "ch", string(l.ch))
		if l.ch == '.' {
			isFloat = true
		}

		l.next()
	}

	var offset int
	if l.pos == len(l.input)-1 {
		offset = 1
	}

	number := l.input[pos : l.pos+offset]
	slog.Info("num", "num", number)

	if number[len(number)-1] == '.' {
		return "", false, fmt.Errorf("number invalid. last char cannot be '.'")
	}

	return number, isFloat, nil
}

func (l *Lexer) peekChar() byte {
	if l.readPos >= len(l.input) {
		return byte(token.EOF)
	}

	return l.input[l.readPos]
}

func (l *Lexer) skipWhitespaces() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' {
		l.next()
	}
}
