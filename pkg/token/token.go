package token

type TokenKind byte

const (
	ILEGAL TokenKind = iota
	EOF

	// Special
	CRLF
	IDENT

	// String operations
	keyword_beg
	SET
	GET
	GETSET
	GETDEL
	GETEX
	INCR
	INCRBY
	DECR
	DECRBY
	MGET
	MSET
	APPEND
	EXISTS
	STRLEN
	SUBSTR

	// CONFIG
	CONFIG
	keyword_end

	args_beg
	XX
	NX
	EX
	args_end

	// STRINGS
	types_beg
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
	types_end
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
	"GET":    GET,
	"GETSET": GETSET,
	"GETEX":  GETEX,
	"GETDEL": GETDEL,
	"SET":    SET,
	"INCR":   INCR,
	"INCRBY": INCRBY,
	"DECR":   DECR,
	"DECRBY": DECRBY,
	"MGET":   MGET,
	"MSET":   MSET,
	"APPEND": APPEND,
	"EXISTS": EXISTS,
	"STRLEN": STRLEN,
	"SUBSTR": SUBSTR,
	"XX":     XX,
	"NX":     NX,
	"EX":     EX,
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
