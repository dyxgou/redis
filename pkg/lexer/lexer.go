package lexer

import (
	"fmt"
	"github/dyxgou/redis/pkg/token"
	"log/slog"
)

const qoute = '"'

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

func (l *Lexer) Reset(input string) {
	l.input = input
	l.pos = 0
	l.readPos = 0
	l.next()
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
	if l.ch == byte(token.EOF) {
		return token.New(token.EOF, "")
	}

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

	if l.ch == qoute {
		t.Kind = token.BULKSTRING
		t.Literal = l.readString()
		l.next()
		return t
	}

	if isLetter(l.ch) {
		t.Literal = l.readIdent()
		t.Kind = token.LookupIdent(t.Literal)
		return t
	} else if isDigit(l.ch) {
		num, err := l.readNumber()

		if err != nil {
			slog.Error("tokenizing number failed", "err", err)
			t.Literal = ""
			t.Kind = token.ILLEGAL
			return t
		}

		t.Kind = token.NUMBER
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

	for isLetter(l.ch) || isDigit(l.ch) {
		l.next()
	}

	offset := l.getReadOffset()

	return l.input[pos : l.pos+offset]
}

func (l *Lexer) readNumber() (string, error) {
	pos := l.pos

	for isDigit(l.ch) {
		l.next()
	}

	offset := l.getReadOffset()

	number := l.input[pos : l.pos+offset]

	if number[len(number)-1] == '.' {
		return "", fmt.Errorf("number invalid. last char cannot be '.'")
	}

	return number, nil
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

// getReadOffset checks if thue current possition is equal to the String input, if so, it returns the offset of 1 to
func (l *Lexer) getReadOffset() int {
	if l.pos == len(l.input)-1 {
		return 1
	}

	return 0
}

func (l *Lexer) readString() string {
	pos := l.pos + 1

	for l.ch != byte(token.EOF) {
		l.next()

		if l.ch == qoute {
			break
		}
	}

	return l.input[pos:l.pos]
}
