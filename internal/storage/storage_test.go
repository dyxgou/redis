package storage

import (
	"fmt"
	"github/dyxgou/redis/pkg/ast"
	"github/dyxgou/redis/pkg/token"
	"testing"
)

func TestMultipleWrite(t *testing.T) {
	s := New()

	tests := []struct {
		key string
		val ast.Expression
	}{
		{key: "intKey", val: &ast.IntegerLit{Token: token.New(token.INTEGER, ":"), Value: 1}},
		{key: "strKey", val: &ast.StringExpr{Token: token.New(token.BULKSTRING, "string")}},
		{key: "trueKey", val: &ast.BooleanExpr{Token: token.New(token.BOOLEAN, "#"), Value: true}},
		{key: "trueKey", val: &ast.BooleanExpr{Token: token.New(token.BOOLEAN, "#"), Value: false}},
		{key: "bigIntKey", val: &ast.BigIntegerExpr{Token: token.New(token.BIGINT, "("), Value: 123123123123}},
	}

	for i, tt := range tests {
		name := fmt.Sprintf("inserting query=%d", i)

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			s.Set(tt.key, tt.val)

			val, ok := s.Get(tt.key)
			if !ok {
				t.Errorf("val not found. key=%q", tt.key)
			}

			assertValue(t, val, NewValue(tt.val))
		})
	}

}

func assertValue(t *testing.T, v Value, expected Value) {
	switch v := v.(type) {
	case *Int:
		i, ok := expected.(*Int)
		if !ok {
			t.Error("expected is not an int type")
		}
		assertInt(t, v, i.Value)
	case *Int64:
		i, ok := expected.(*Int64)
		if !ok {
			t.Error("expected is not an int64 type")
		}

		assertInt64(t, v, i.Value)
	case *Float:
		f, ok := expected.(*Float)
		if !ok {
			t.Error("expected is not an float type")
		}

		assertFloat(t, v, f.Value)
	case *Bool:
		b, ok := expected.(*Bool)
		if !ok {
			t.Error("expected is not an int type")
		}

		assertBool(t, v, b.Value)
	case *String:
		s, ok := expected.(*String)
		if !ok {
			t.Error("expected is not an int type")
		}

		assertString(t, v, s.Value)
	default:
		t.Errorf("value not supported. got=%+v", v)
	}
}
