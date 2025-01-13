package storage

import (
	"errors"
	"github/dyxgou/redis/pkg/ast"
)

var NotSupportedValErr = errors.New("value provided not supported")
var Nil = &NilVal{}

type values interface {
	int | int64 | float64 | bool | string
}

type valueKind byte

const (
	intKind valueKind = iota
	int64Kind
	boolKind
	stringKind
	floatKind
	nilKind
)

type Value interface {
	kind() valueKind
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

func (n *NilVal) kind() valueKind { return nilKind }
func (i *Int) kind() valueKind    { return intKind }
func (b *Bool) kind() valueKind   { return boolKind }
func (f *Float) kind() valueKind  { return floatKind }
func (i *Int64) kind() valueKind  { return int64Kind }
func (s *String) kind() valueKind { return stringKind }
