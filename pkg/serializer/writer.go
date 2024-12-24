package serializer

import (
	"bytes"
	"fmt"
	"github/dyxgou/redis/pkg/token"
	"strconv"
)

// headSize represents the size of the array beggining of any command like "$NUM\r\n"
const headSize = 6

type writer struct {
	len int
	sb  *bytes.Buffer

	head *bytes.Buffer
}

func newWriter() *writer {
	return &writer{
		len: 0,
		sb:  new(bytes.Buffer),

		head: new(bytes.Buffer),
	}
}

func (w *writer) writeLen() {
	if w.len == 0 {
		return
	}

	symbol := token.Symbols[token.ARRAY]

	w.head.WriteByte(symbol)
	w.head.WriteString(strconv.Itoa(w.len))
}

func (w *writer) writeCRLF() {
	w.sb.WriteString("\r\n")
}

func (w *writer) writeSymbol(kind token.TokenKind, n int) error {
	symbol, ok := token.Symbols[kind]

	if !ok {
		return fmt.Errorf("symbol not found for kind=%d", kind)
	}

	w.sb.WriteByte(symbol)
	w.sb.WriteString(strconv.Itoa(n))

	w.writeCRLF()

	return nil
}

func (w *writer) writeKeyword(kw token.Token) error {
	n := len(kw.Literal)

	if err := w.writeSymbol(token.BULKSTRING, n); err != nil {
		return err
	}
	w.sb.WriteString(kw.Literal)

	w.writeCRLF()

	w.len++

	return nil
}

func (w *writer) string() string {
	symbol := token.Symbols[token.ARRAY]
	w.head.WriteByte(symbol)
	w.head.WriteString(strconv.Itoa(w.len))
	w.head.WriteString("\r\n")

	w.head.Write(w.sb.Bytes())

	return w.head.String()
}
