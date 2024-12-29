package parser

import (
	"errors"
	"fmt"
	"github/dyxgou/redis/pkg/ast"
	"github/dyxgou/redis/pkg/lexer"
	"github/dyxgou/redis/pkg/token"
	"strconv"
)

const setMinArgs int = 3

type Parser struct {
	l *lexer.Lexer

	curTok  token.Token
	readTok token.Token

	maxLen int
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, maxLen: 100}
	p.next()
	p.next()

	return p
}

func (p *Parser) next() {
	p.curTok = p.readTok
	p.readTok = p.l.NextToken()
}

func (p *Parser) done() bool {
	return p.curTok.Kind == token.EOF
}

func (p *Parser) Parse() (ast.Command, error) {
	var n int
	for !p.done() && n < p.maxLen {
		if err := p.parseBegCommand(); err != nil {
			return nil, err
		}

		if err := p.skipBulkString(); err != nil {
			return nil, err
		}

		cmd, err := p.parseCommand()
		if err != nil {
			return nil, err
		}

		return cmd, nil
	}

	return nil, errors.New("could not parse input")
}

func (p *Parser) parseCommand() (ast.Command, error) {
	switch p.curTok.Kind {
	case token.GET:
		return p.parseGetCommand()
	case token.SET:
		return p.parseSetCommand()
	}

	return nil, fmt.Errorf("command not supported. got=%d (%q)", p.curTok.Kind, p.curTok.Literal)
}

func (p *Parser) skipBulkString() error {
	if !p.curTokIs(token.BULKSTRING) {
		return p.isNotBulkStringErr()
	}

	p.next()
	if !p.curTokIs(token.INTEGER) {
		return p.isNotIntegerErr()
	}

	p.next()
	if err := p.checkCRLF(); err != nil {
		return err
	}

	if !p.curTokIs(token.IDENT) {
		return fmt.Errorf(
			"token expected=%d ('IDENT'). got=%d (%s)",
			token.IDENT,
			p.curTok.Kind,
			p.curTok.Literal,
		)
	}

	return nil
}

func (p *Parser) parseBegCommand() error {
	if !p.curTokIs(token.ARRAY) {
		return p.isNotArrayErr()
	}

	p.next()
	if !p.curTokIs(token.INTEGER) {
		return p.isNotIntegerErr()
	}
	n, err := strconv.Atoi(p.curTok.Literal)

	if err != nil {
		return err
	}
	p.maxLen = n

	p.next()

	if err := p.checkCRLF(); err != nil {
		return err
	}

	return nil
}

func (p *Parser) curTokIs(k token.TokenKind) bool {
	return p.curTok.Kind == k
}

func (p *Parser) checkCRLF() error {
	if !p.curTokIs(token.CRLF) {
		return fmt.Errorf(
			"token expected=%d ('CRLF'). got=%d (%s)",
			token.CRLF,
			p.curTok.Kind,
			p.curTok.Literal,
		)
	}

	p.next()
	return nil
}

func (p *Parser) parseGetCommand() (*ast.GetCommand, error) {
	getCmd := &ast.GetCommand{Token: p.curTok}

	p.next()
	if err := p.checkCRLF(); err != nil {
		return nil, err
	}

	if err := p.skipBulkString(); err != nil {
		return nil, err
	}

	getCmd.Key = p.curTok.Literal

	return getCmd, nil
}

func (p *Parser) parseSetCommand() (*ast.SetCommand, error) {
	setCmd := &ast.SetCommand{Token: p.curTok}

	p.next()
	if err := p.checkCRLF(); err != nil {
		return nil, err
	}

	if err := p.skipBulkString(); err != nil {
		return nil, err
	}

	setCmd.Key = p.curTok.Literal

	if err := p.checkCRLF(); err != nil {
		return nil, err
	}
	if err := p.skipBulkString(); err != nil {
		return nil, err
	}

	setCmd.Value = p.curTok.Literal

	if p.maxLen == setMinArgs {
		return setCmd, nil
	}

	return setCmd, nil
}

func (p *Parser) isNotArrayErr() error {
	return fmt.Errorf(
		"token expected=%d ('ARRAY'). got=%d (%s)",
		token.ARRAY,
		p.curTok.Kind,
		p.curTok.Literal,
	)
}

func (p *Parser) isNotIntegerErr() error {
	return fmt.Errorf(
		"token expected=%d ('INTEGER'). got=%d (%s)",
		token.INTEGER,
		p.curTok.Kind,
		p.curTok.Literal,
	)
}

func (p *Parser) isNotBulkStringErr() error {
	return fmt.Errorf(
		"token expected=%d ('BULKSTRING'). got=%d (%s)",
		token.BULKSTRING,
		p.curTok.Kind,
		p.curTok.Literal,
	)
}
