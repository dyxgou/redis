package parser

import (
	"fmt"
	"github/dyxgou/redis/pkg/ast"
	"github/dyxgou/redis/pkg/lexer"
	"github/dyxgou/redis/pkg/token"
	"strconv"
)

const trueStr = "t"
const falseStr = "f"

type Parser struct {
	l *lexer.Lexer

	curTok  token.Token
	readTok token.Token

	len int
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
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

func (p *Parser) Reset(input string) {
	p.l.Reset(input)
	p.next()
	p.next()
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
	case token.SET:
		return p.parseSetCommand()
	case token.GETSET:
		return p.parseGetSetCommand()
	case token.GETEX:
		return p.parseGetExCommand()
	case token.GETDEL:
		return p.parseGetDelCommand()
	case token.INCR:
		return p.parseIncrCommand()
	case token.INCRBY:
		return p.parseIncrByCommand()
	case token.DECR:
		return p.parseDecrCommand()
	case token.DECRBY:
		return p.parseDecrByCommand()
	}

	return nil, fmt.Errorf("command not supported. got=%d (%q)", p.curTok.Kind, p.curTok.Literal)
}

func (p *Parser) skipBulkString() error {
	if !p.curTokIs(token.BULKSTRING) {
		return p.isNotBulkStringErr()
	}

	p.next()
	if !p.curTokIs(token.NUMBER) {
		return p.isNotNumberErr()
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
	if !p.curTokIs(token.NUMBER) {
		return p.isNotNumberErr()
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

func (p *Parser) readTokIs(k token.TokenKind) bool {
	return p.readTok.Kind == k
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

func (p *Parser) parseSetCommand() (*ast.SetCommand, error) {
	sc := &ast.SetCommand{Token: p.curTok}
	p.next()
	if err := p.checkCRLF(); err != nil {
		return nil, err
	}

	key, err := p.parseIdent()
	if err != nil {
		return nil, err
	}

	sc.Key = key

	p.next()
	if err := p.checkCRLF(); err != nil {
		return nil, err
	}

	v, err := p.parseValue()
	if err != nil {
		return nil, err
	}

	sc.Value = v

	p.next()
	if err := p.checkCRLF(); err != nil {
		return nil, err
	}

	if err := p.parseSetArgs(sc); err != nil {
		if !p.readTokIs(token.EOF) {
			return nil, err
		}
	}

	return sc, nil
}

func (p *Parser) parseGetSetCommand() (*ast.GetSetCommand, error) {
	gsc := &ast.GetSetCommand{Token: p.curTok}

	p.next()
	if err := p.checkCRLF(); err != nil {
		return nil, err
	}

	key, err := p.parseIdent()
	if err != nil {
		return nil, err
	}

	gsc.Key = key

	p.next()
	if err := p.checkCRLF(); err != nil {
		return nil, err
	}

	v, err := p.parseValue()
	if err != nil {
		return nil, err
	}

	gsc.Value = v

	return gsc, nil
}

func (p *Parser) parseGetExCommand() (*ast.GetExCommand, error) {
	ge := &ast.GetExCommand{Token: p.curTok}

	p.next()
	if err := p.checkCRLF(); err != nil {
		return nil, err
	}

	key, err := p.parseIdent()
	if err != nil {
		return nil, err
	}

	ge.Key = key

	p.next()
	if err := p.checkCRLF(); err != nil {
		return nil, err
	}

	if err := p.skipBulkString(); err != nil {
		return nil, err
	}

	if !p.curTokIs(token.EX) {
		return nil, fmt.Errorf(
			"token expected=%d ('EX'). got=%d (%q)", token.EX, p.curTok.Kind, p.curTok.Literal,
		)
	}

	ex, err := p.parseExArg()
	if err != nil {
		return nil, err
	}
	ge.Ex = ex

	return ge, nil
}

func (p *Parser) parseGetDelCommand() (*ast.GetDelCommand, error) {
	gd := &ast.GetDelCommand{Token: p.curTok}

	p.next()
	if err := p.checkCRLF(); err != nil {
		return nil, err
	}

	key, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	gd.Key = key

	p.next()
	if err := p.checkCRLF(); err != nil {
		return nil, err
	}

	return gd, nil
}

func (p *Parser) parseIncrCommand() (*ast.IncrCommand, error) {
	inc := &ast.IncrCommand{Token: p.curTok}

	p.next()
	if err := p.checkCRLF(); err != nil {
		return nil, err
	}

	key, err := p.parseIdent()
	if err != nil {
		return nil, err
	}

	inc.Key = key

	return inc, nil
}

func (p *Parser) parseDecrCommand() (*ast.DecrCommand, error) {
	dec := &ast.DecrCommand{Token: p.curTok}

	p.next()
	if err := p.checkCRLF(); err != nil {
		return nil, err
	}

	key, err := p.parseIdent()
	if err != nil {
		return nil, err
	}

	dec.Key = key

	return dec, nil
}

func (p *Parser) parseDecrByCommand() (*ast.DecrByCommand, error) {
	dec := &ast.DecrByCommand{Token: p.curTok}

	p.next()
	if err := p.checkCRLF(); err != nil {
		return nil, err
	}

	key, err := p.parseIdent()
	if err != nil {
		return nil, err
	}

	dec.Key = key

	p.next()
	if err := p.checkCRLF(); err != nil {
		return nil, err
	}

	if !p.curTokIs(token.INTEGER) {
		return nil, p.isNotIntegerErr()
	}
	p.next()

	i, err := strconv.Atoi(p.curTok.Literal)
	if err != nil {
		return nil, err
	}

	dec.Decrement = i

	return dec, nil
}

func (p *Parser) parseIncrByCommand() (*ast.IncrByCommand, error) {
	inc := &ast.IncrByCommand{Token: p.curTok}

	p.next()
	if err := p.checkCRLF(); err != nil {
		return nil, err
	}

	key, err := p.parseIdent()
	if err != nil {
		return nil, err
	}

	inc.Key = key

	p.next()
	if err := p.checkCRLF(); err != nil {
		return nil, err
	}

	if !p.curTokIs(token.INTEGER) {
		return nil, p.isNotIntegerErr()
	}
	p.next()

	i, err := strconv.Atoi(p.curTok.Literal)
	if err != nil {
		return nil, err
	}

	inc.Increment = i

	return inc, nil
}

func (p *Parser) parseSetArgs(sc *ast.SetCommand) error {
	if err := p.skipBulkString(); err != nil {
		return err
	}

	switch p.curTok.Kind {
	case token.EX:
		n, err := p.parseExArg()
		if err != nil {
			return err
		}
		sc.Ex = n
	case token.NX:
		if sc.Xx {
			return fmt.Errorf("NX and XX cannot be used together")
		}

		sc.Nx = true
	case token.XX:
		if sc.Nx {
			return fmt.Errorf("NX and XX cannot be used together")
		}

		sc.Xx = true
	default:
		return fmt.Errorf("token invalid=%d (%q). SET doesn't expect this argument", p.curTok.Kind, p.curTok.Literal)
	}

	p.next()
	if err := p.checkCRLF(); err != nil {
		return err
	}

	if !p.readTokIs(token.EOF) {
		return p.parseSetArgs(sc)
	}

	return nil
}

func (p *Parser) parseExArg() (int64, error) {
	p.next()
	if err := p.checkCRLF(); err != nil {
		return 0, err
	}

	if !p.curTokIs(token.INTEGER) {
		return 0, p.isNotIntegerErr()
	}

	p.next()
	n, err := strconv.ParseInt(p.curTok.Literal, 10, 64)
	if n < 1 {
		return 0, fmt.Errorf("EX argument cannot be less than 1. got=%d", n)
	}

	if err != nil {
		return 0, err
	}

	return n, nil
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
	case token.BULKSTRING:
		return p.parseString()
	case token.INTEGER:
		return p.parseInteger()
	case token.BIGINT:
		return p.parseBigInt()
	case token.FLOAT:
		return p.parseFloat()
	}

	return nil, fmt.Errorf("curTok kind is not a value. got=%d (%q)", p.curTok.Kind, p.curTok.Literal)
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

func (p *Parser) parseInteger() (*ast.IntegerLit, error) {
	ie := &ast.IntegerLit{Token: p.curTok}
	p.next()

	v, err := strconv.Atoi(p.curTok.Literal)
	if err != nil {
		return nil, err
	}

	ie.Value = v

	return ie, nil
}

func (p *Parser) parseBigInt() (*ast.BigIntegerExpr, error) {
	bi := &ast.BigIntegerExpr{Token: p.curTok}
	p.next()

	v, err := strconv.ParseInt(p.curTok.Literal, 10, 64)
	if err != nil {
		return nil, err
	}

	bi.Value = v
	return bi, nil
}

func (p *Parser) parseFloat() (*ast.FloatExpr, error) {
	fo := &ast.FloatExpr{Token: p.curTok}
	p.next()

	v, err := strconv.ParseFloat(p.curTok.Literal, 64)
	if err != nil {
		return nil, err
	}

	fo.Value = v
	return fo, nil
}

func (p *Parser) isNotArrayErr() error {
	return fmt.Errorf(
		"token expected=%d ('ARRAY'). got=%d (%s)",
		token.ARRAY,
		p.curTok.Kind,
		p.curTok.Literal,
	)
}

func (p *Parser) isNotNumberErr() error {
	return fmt.Errorf(
		"token expected=%d ('NUMBER'). got=%d (%s)",
		token.NUMBER,
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
