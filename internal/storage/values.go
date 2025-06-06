package storage

import (
	"errors"
	"github/dyxgou/redis/pkg/ast"
	"strconv"
)

var NotSupportedValErr = errors.New("value provided not supported")
var Nil = &NilVal{}

type values interface {
	int | int64 | float64 | bool | string
}

type valueKind byte

const (
	IntKind valueKind = iota
	Int64Kind
	BoolKind
	StringKind
	FloatKind
	NilKind
)

type Value interface {
	Kind() valueKind
	String() string
}

func NewValue(v ast.Expression) Value {
	switch v := v.(type) {
	case *ast.IntegerLit:
		return &Int{v.Value}
	case *ast.BigIntegerExpr:
		return &Int64{v.Value}
	case *ast.BooleanExpr:
		return &Bool{v.Value}
	case *ast.StringExpr:
		return &String{v.Value()}
	case *ast.FloatExpr:
		return &Float{v.Value}
	}

	return nil
}

type (
	Bool struct {
		Value bool
	}

	Float struct {
		Value float64
	}

	Int struct {
		Value int
	}

	Int64 struct {
		Value int64
	}

	String struct {
		Value string
	}

	NilVal struct{}
)

func (n *NilVal) Kind() valueKind { return NilKind }
func (i *Int) Kind() valueKind    { return IntKind }
func (b *Bool) Kind() valueKind   { return BoolKind }
func (f *Float) Kind() valueKind  { return FloatKind }
func (i *Int64) Kind() valueKind  { return Int64Kind }
func (s *String) Kind() valueKind { return StringKind }

func (n *NilVal) String() string { return "(nil)" }
func (i *Int) String() string    { return strconv.Itoa(i.Value) }
func (f *Float) String() string  { return strconv.FormatFloat(f.Value, 'E', -1, 64) }
func (i *Int64) String() string  { return strconv.FormatInt(i.Value, 10) }
func (s *String) String() string { return s.Value }
func (b *Bool) String() string {
	if b.Value {
		return "true"
	} else {
		return "false"
	}
}
