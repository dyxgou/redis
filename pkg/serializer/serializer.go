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
	return s.curTok.Kind == token.EOF
}

func (s *Serializer) Serialize() (string, error) {
	for !s.done() {
		if token.IsKeyword(s.curTok.Kind) || token.IsArg(s.curTok.Kind) {
			if err := s.w.writeWord(s.curTok); err != nil {
				return "", err
			}
		} else if token.IsNumber(s.curTok.Kind) {
			if err := s.w.writeNumber(s.curTok); err != nil {
				return "", err
			}
		} else {
			return "", fmt.Errorf("token not supported. got=%d (%q)", s.curTok.Kind, s.curTok.Literal)
		}

		s.next()
	}

	return s.w.string(), nil
}
