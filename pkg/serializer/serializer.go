package serializer

import (
	"github/dyxgou/redis/pkg/lexer"
	"github/dyxgou/redis/pkg/token"
)

type Serializer struct {
	l      *lexer.Lexer
	curTok token.Token

	w *writer
}

func New(l *lexer.Lexer) *Serializer {
	s := &Serializer{
		l: l,
		w: newWriter(),
	}

	s.next()

	return s
}

func (s *Serializer) next() {
	s.curTok = s.l.NextToken()
}

func (s *Serializer) Serialize() (string, error) {
	for s.curTok.Kind != token.EOF {
		if token.IsKeyword(s.curTok.Kind) {
			err := s.w.writeKeyword(s.curTok)

			if err != nil {
				return "", err
			}
		}

		s.next()
	}

	return s.w.string(), nil
}
