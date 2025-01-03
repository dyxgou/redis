package ast

import (
	"github/dyxgou/redis/pkg/token"
	"strconv"
	"strings"
	"time"
)

type (
	// A GetCommand represents the command "GET <key>"
	GetCommand struct {
		Token token.Token
		Key   string
	}

	// A GetSetCommand represents the "GETSET <key> <value>" and acts like a SET followed by a GET
	GetSetCommand struct {
		token token.Token
		key   string
		value string
	}

	// A GetExCommand represents the "GETEX <key> <value> <EX> <time>" if the Ex is not provided, the parser throws an error
	GetExCommand struct {
		token token.Token
		key   string
		value string
		Ex    time.Duration
	}

	// A GetDelCommand represents the "GETDEL <key>"
	GetDelCommand struct {
		token token.Token
		key   string
	}

	// A SetCommand represents the command "SET <key> <value>" if Ex is not provided it is assumed that the key won't be expired
	SetCommand struct {
		Token token.Token
		Key   string
		Value Expression
		Ex    int

		// Sets the key if not exists
		Nx bool
		// Sets the key if it exists
		Xx bool
	}

	// A IncrCommand represents "INCR <key>" and increments the key by 1.
	IncrCommand struct {
		token token.Token
		key   string
	}

	// A IncrByCommand represents "INCRBY <key> <increment>" and increments the key by increment.
	IncrByCommand struct {
		token     token.Token
		key       string
		increment int
	}

	// A DercCommand represents "DECR <key>" and decrements the key by 1.
	DercCommand struct {
		token token.Token
		key   string
	}

	// A DercByCommand represents "DECRBY <key> <decrement>" and decrements the key by decrement.
	DercByCommand struct {
		token     token.Token
		key       string
		decrement int
	}

	// A MGetCommand stands for Multiple Get and represents "MGET <key> [<key> ...]"
	MGetCommand struct {
		token token.Token
		keys  []string
	}

	// MSetCommand stands for Multiple Set and represents "MSET <key> <value> [<key> <value> ...]"
	MSetCommand struct {
		token token.Token
		pairs []struct {
			key   string
			value string
		}
	}

	// AppendCommand represents "APPEND <key> <value>"
	AppendCommand struct {
		token token.Token
		key   string
		value string
	}

	// ExistsCommand represents "EXISTS <key> <value>"
	ExistsCommand struct {
		token token.Token
		key   string
	}
)

func (gc *GetCommand) cmdNode()             {}
func (gc *GetCommand) TokenLiteral() string { return gc.Token.Literal }
func (gc *GetCommand) String() string {
	var sb strings.Builder

	sb.WriteString(gc.Token.Literal)
	sb.WriteByte(' ')

	sb.WriteString(gc.Key)

	return sb.String()
}

func (sc *SetCommand) cmdNode()             {}
func (sc *SetCommand) TokenLiteral() string { return sc.Token.Literal }
func (sc *SetCommand) String() string {
	var sb strings.Builder

	sb.WriteString(sc.Token.Literal)
	sb.WriteByte(' ')

	sb.WriteString(sc.Key)
	sb.WriteByte(' ')

	sb.WriteString(sc.Value)
	sb.WriteByte(' ')

	writeNumArg(&sb, "Ex", sc.Ex)
	writeBoolArg(&sb, "NX", sc.Nx)
	writeBoolArg(&sb, "XX", sc.Xx)

	return sb.String()
}

type (
	BooleanExpr struct {
		Token token.Token
		Value bool
	}

	StringExpr struct {
		Token token.Token
	}

	IntegerExpr struct {
		Token token.Token
		Value int
	}

	BigIntegerExpr struct {
		Token token.Token
		Value int64
	}

	FloatExpr struct {
		Token token.Token
		Value float64
	}
)

func (be *BooleanExpr) exprNode()    {}
func (se *StringExpr) exprNode()     {}
func (ie *IntegerExpr) exprNode()    {}
func (bi *BigIntegerExpr) exprNode() {}
func (fo *FloatExpr) exprNode()      {}

func (be *BooleanExpr) TokenLiteral() string    { return be.Token.Literal }
func (se *StringExpr) TokenLiteral() string     { return se.Token.Literal }
func (ie *IntegerExpr) TokenLiteral() string    { return ie.Token.Literal }
func (bi *BigIntegerExpr) TokenLiteral() string { return bi.Token.Literal }
func (fo *FloatExpr) TokenLiteral() string      { return fo.Token.Literal }

func (be *BooleanExpr) String() string    { return be.Token.Literal }
func (se *StringExpr) String() string     { return se.Token.Literal }
func (ie *IntegerExpr) String() string    { return ie.Token.Literal }
func (bi *BigIntegerExpr) String() string { return bi.Token.Literal }
func (fo *FloatExpr) String() string      { return fo.Token.Literal }

func (se *StringExpr) Value() string { return se.Token.Literal }

func writeNumArg(sb *strings.Builder, arg string, n int) {
	if n == 0 {
		return
	}

	sb.WriteString(arg)
	sb.WriteString(strconv.Itoa(n))
	sb.WriteByte(' ')
}

func writeBoolArg(sb *strings.Builder, arg string, val bool) {
	if !val {
		return
	}

	sb.WriteString(arg)
	sb.WriteByte(' ')
}
