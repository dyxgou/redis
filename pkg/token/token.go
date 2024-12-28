package token

type TokenKind byte

const (
	EOF TokenKind = iota
	ILLEGAL

	// Special
	CRLF
	TEXT

	// String operations
	keyword_beg
	IDENT
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

func IsKeyword(kind TokenKind) bool {
	return keyword_beg < kind && kind < keyword_end
}

func IsArg(kind TokenKind) bool {
	return args_beg < kind && kind < args_end
}

func IsNumber(kind TokenKind) bool {
	return kind == INTEGER || kind == FLOAT
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
