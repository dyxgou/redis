package parser

import (
	"errors"
	"fmt"
	"github/dyxgou/redis/pkg/ast"
	"github/dyxgou/redis/pkg/lexer"
	"github/dyxgou/redis/pkg/token"
	"strconv"
)

const trueStr = "t"
const falseStr = "f"

var errNilPtr = errors.New("destination pointer is nil")

type Parser struct {
	l *lexer.Lexer

	curTok  token.Token
	readTok token.Token

	len int
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, len: 100}
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

func (p *Parser) parseCommand() (ast.Command, error) {
	switch p.curTok.Kind {
	case token.GET:
		return p.parseGetCommand()
		// case token.SET:
		// 	return p.parseSetCommand()
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
	p.len = n

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

	key, err := p.parseIdent()
	if err != nil {
		return nil, err
	}

	getCmd.Key = key

	return getCmd, nil
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

func (p *Parser) isNotIdentErr() error {
	return fmt.Errorf(
		"token expected=%d ('IDENT'). got=%d (%q)",
		token.IDENT,
		p.curTok.Kind,
		p.curTok.Literal,
	)
}

func (p *Parser) parseIdent() (string, error) {
	if err := p.skipBulkString(); err != nil {
		return "", err
	}

	if !p.curTokIs(token.IDENT) {
		return "", p.isNotIdentErr()
	}

	return p.curTok.Literal, nil
}

func (p *Parser) parseValue() (ast.Expression, error) {
	switch p.curTok.Kind {
	case token.BOOLEAN:
		return p.parseBoolean()
	case token.STRING:
	case token.BULKSTRING:
		return p.parseString()
	case token.INTEGER:
		return p.parseInteger()
	case token.BIGNUMBER:
		return p.parseBigInt()
	case token.FLOAT:
		return p.parseFloat()
	}

	return nil, fmt.Errorf("curTok is not a value. got=%d (%q)", p.curTok.Kind, p.curTok.Literal)
}

func (p *Parser) parseBoolean() (*ast.BooleanExpr, error) {
	be := &ast.BooleanExpr{Token: p.curTok}
	p.next()

	if !p.curTokIs(token.IDENT) {
		return nil, p.isNotIdentErr()
	}

	if p.curTok.Literal == trueStr {
		be.Value = true
	} else if p.curTok.Literal == falseStr {
		be.Value = false
	} else {
		return nil, fmt.Errorf("curTok invalid expected='t' or 'f'. got=%q", p.curTok.Literal)
	}

	return be, nil
}

func (p *Parser) parseString() (*ast.StringExpr, error) {
	if err := p.skipBulkString(); err != nil {
		return nil, err
	}
	se := &ast.StringExpr{Token: p.curTok}

	return se, nil
}

func (p *Parser) parseInteger() (*ast.IntegerExpr, error) {
	ie := &ast.IntegerExpr{Token: p.curTok}

	v, err := strconv.Atoi(p.curTok.Literal)
	if err != nil {
		return nil, err
	}

	ie.Value = v

	return ie, nil
}

func (p *Parser) parseBigInt() (*ast.BigIntegerExpr, error) {
	bi := &ast.BigIntegerExpr{Token: p.curTok}

	v, err := strconv.ParseInt(p.curTok.Literal, 10, 64)
	if err != nil {
		return nil, err
	}

	bi.Value = v
	return bi, nil
}

func (p *Parser) parseFloat() (*ast.FloatExpr, error) {
	fo := &ast.FloatExpr{Token: p.curTok}
	v, err := strconv.ParseFloat(p.curTok.Literal, 64)
	if err != nil {
		return nil, err
	}

	fo.Value = v
	return fo, nil
}
