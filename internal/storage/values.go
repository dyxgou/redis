package storage

import (
	"errors"
)

var NotSupportedValErr = errors.New("value provided not supported")

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
)

type val interface {
	kind() valueKind
}

func NewValue[V values | any](v V) val {
	switch v := any(v).(type) {
	case int:
		return &Int{v}
	case int64:
		return &Int64{v}
	case bool:
		return &Bool{v}
	case string:
		return &String{v}
	case float64:
		return &Float{v}
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
)

func (i *Int) kind() valueKind    { return intKind }
func (b *Bool) kind() valueKind   { return boolKind }
func (f *Float) kind() valueKind  { return floatKind }
func (i *Int64) kind() valueKind  { return int64Kind }
func (s *String) kind() valueKind { return stringKind }
