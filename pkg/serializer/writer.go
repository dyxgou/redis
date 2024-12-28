package serializer

import (
	"bytes"
	"fmt"
	"github/dyxgou/redis/pkg/token"
	"strconv"
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
	if err := w.writeSymbol(t.Kind); err != nil {
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
