package token

const EndCRLF = "\r\n"

var Symbols = map[TokenKind]byte{
	STRING:         '+',
	BULKSTRING:     '$',
	VERTAMINSTRING: '=',
	ERROR:          '-',
	BULKERROR:      '!',
	INTEGER:        ':',
	FLOAT:          ',',
	BIGNUMBER:      '(',
	BOOLEAN:        '#',
	ARRAY:          '*',
	NULL:           '_',
	MAPS:           '%',
	ATTRIBUTES:     '`',
	SETS:           '~',
	PUSHES:         '>',
}
