package serializer

import (
	"fmt"
	"github/dyxgou/redis/pkg/lexer"
	"github/dyxgou/redis/pkg/token"
)

type Serializer struct {
	l      *lexer.Lexer
	curTok token.Token

	w *writer
}

func New(input string) *Serializer {
	s := &Serializer{
		l: lexer.New(input),
		w: newWriter(),
	}

	s.next()
	return s
}

func (s *Serializer) next() {
	s.curTok = s.l.NextToken()
}

func (s *Serializer) done() bool {
	return s.curTokIs(token.EOF)
}

func (s *Serializer) Serialize() (string, error) {
	for !s.done() {
		if s.curTokIs(token.ILLEGAL) {
			return "", fmt.Errorf("token ilegal")
		}

		if s.curTokIs(token.STRING) {
			if err := s.w.writeSimpleString(s.curTok); err != nil {
				return "", err
			}
		}

		if s.curTokIs(token.BOOLEAN) {
			cur := s.curTok
			s.next()

			if err := s.w.writeBool(cur, s.curTok); err != nil {
				return "", err
			}

			s.next()
		}

		if token.IsKeyword(s.curTok.Kind) || s.curTokIs(token.BULKSTRING) {
			if err := s.w.writeWord(s.curTok); err != nil {
				return "", err
			}
		}

		if s.curTokIs(token.NUMBER) {
			if err := s.w.writeNumber(s.curTok); err != nil {
				return "", err
			}
		}

		s.next()
	}

	return s.w.string(), nil
}

func (s *Serializer) curTokIs(k token.TokenKind) bool {
	return s.curTok.Kind == k
}
