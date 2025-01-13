package token

const EndCRLF = "\r\n"

var symbols = map[TokenKind]byte{
	STRING:         '+',
	BULKSTRING:     '$',
	VERTAMINSTRING: '=',
	ERROR:          '-',
	BULKERROR:      '!',
	INTEGER:        ':',
	FLOAT:          ',',
	BIGINT:         '(',
	BOOLEAN:        '#',
	ARRAY:          '*',
	NIL:            '_',
	MAPS:           '%',
	ATTRIBUTES:     '`',
	SETS:           '~',
	PUSHES:         '>',
}

var chars = map[byte]TokenKind{
	'+': STRING,
	'$': BULKSTRING,
	'=': VERTAMINSTRING,
	'-': ERROR,
	'!': BULKERROR,
	':': INTEGER,
	',': FLOAT,
	'(': BIGINT,
	'#': BOOLEAN,
	'*': ARRAY,
	'_': NIL,
	'%': MAPS,
	'`': ATTRIBUTES,
	'~': SETS,
	'>': PUSHES,
}

func GetKindWithSymbol(ch byte) (TokenKind, bool) {
	k, ok := chars[ch]

	return k, ok
}

func GetSymbolWithKind(kind TokenKind) (byte, bool) {
	sym, ok := symbols[kind]

	return sym, ok
}
