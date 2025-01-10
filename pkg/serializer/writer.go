package serializer

import (
	"bytes"
	"fmt"
	"github/dyxgou/redis/pkg/token"
	"strconv"
	"strings"
)

type writer struct {
	head *bytes.Buffer
	body *bytes.Buffer

	len int
}

func newWriter() *writer {
	return &writer{
		head: new(bytes.Buffer),
		body: new(bytes.Buffer),
	}
}

func (w *writer) writeCRLF() {
	w.body.WriteString(token.EndCRLF)
}

func (w *writer) writeSymbol(k token.TokenKind) error {
	sym, ok := token.GetSymbolWithKind(k)

	if !ok {
		return fmt.Errorf("not symbol found. token=%d", k)
	}

	w.body.WriteByte(sym)
	return nil
}

func (w *writer) writeSymbolWithAmount(k token.TokenKind, n int) error {
	err := w.writeSymbol(k)

	if err != nil {
		return err
	}

	w.body.WriteString(strconv.Itoa(n))
	w.writeCRLF()

	return nil
}

const TRUE = "t"
const FALSE = "f"

func (w *writer) writeBool(cur, next token.Token) error {
	if cur.Kind != token.BOOLEAN {
		return fmt.Errorf("curToken expected=%d ('BOOLEAN'). got=%d (%q)", token.BOOLEAN, cur.Kind, cur.Literal)
	}

	if next.Kind != token.IDENT {
		return fmt.Errorf("nextToken expected=%d ('IDENT'). got=%d (%q)", token.IDENT, next.Kind, next.Literal)
	}

	if next.Literal != TRUE && next.Literal != FALSE {
		return fmt.Errorf("nextToken expected 't' or 'f'. got=%q", next.Literal)
	}

	if err := w.writeSymbol(token.BOOLEAN); err != nil {
		return err
	}

	w.body.WriteString(next.Literal)
	w.writeCRLF()
	w.len++

	return nil
}

func (w *writer) writeWord(t token.Token) error {
	if err := w.writeSymbolWithAmount(token.BULKSTRING, len(t.Literal)); err != nil {
		return err
	}

	w.body.WriteString(t.Literal)
	w.writeCRLF()
	w.len++

	return nil
}

func (w *writer) writeNumber(t token.Token) error {
	var sym token.TokenKind

	if len(t.Literal) >= 10 {
		sym = token.BIGINT
	} else if isFloat(t.Literal) {
		sym = token.FLOAT
	} else {
		sym = token.INTEGER
	}

	if err := w.writeSymbol(sym); err != nil {
		return err
	}

	w.body.WriteString(t.Literal)
	w.writeCRLF()
	w.len++

	return nil
}

func (w *writer) writeLen() {
	sym, _ := token.GetSymbolWithKind(token.ARRAY)

	w.head.WriteByte(sym)
	w.head.WriteString(strconv.Itoa(w.len))
	w.head.WriteString(token.EndCRLF)
}

func (w *writer) string() string {
	w.writeLen()

	w.head.Write(w.body.Bytes())

	return w.head.String()
}

func isFloat(n string) bool {
	return strings.Contains(n, ".")
}
