package token

type TokenKind byte

const (
	ILEGAL TokenKind = iota
	EOF

	// Special
	CRLF
	IDENT

	// String operations
	SET
	SETNX
	GET
	GETSET
	INCR
	INCRBY
	DECR
	DECRBY
	MGET
	MSET
	APPEND

	// CONFIG
	CONFIG

	// STRINGS
	STRING
	BULKSTRING
	VERTAMINSTRING
	ERROR
	BULKERROR

	// INTEGERS
	INTEGER
	FLOAT
	BIGNUMBER

	// DATA
	BOOLEAN
	ARRAY
	NULL
	MAPS
	ATTRIBUTES
	SETS
	PUSHES
)

type Token struct {
	Kind    TokenKind
	Literal string
}

func New(kind TokenKind, literal string) Token {
	return Token{
		Kind:    kind,
		Literal: literal,
	}
}

var keywords = map[string]TokenKind{
	"SET":    SET,
	"SETNX":  SETNX,
	"GET":    GET,
	"GETSET": GETSET,
	"INCR":   INCR,
	"INCRBY": INCRBY,
	"DECR":   DECR,
	"DECRBY": DECRBY,
	"MGET":   MGET,
	"MSET":   MSET,
	"APPEND": APPEND,
}

func LookupIdent(ident string) TokenKind {
	kw, ok := keywords[ident]

	if ok {
		return kw
	}

	return IDENT

}

func (t *Token) String() string {
	return t.Literal
}
