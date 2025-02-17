package evaluator

import (
	"fmt"
	"github/dyxgou/redis/internal/storage"
	"github/dyxgou/redis/pkg/ast"
	"github/dyxgou/redis/pkg/token"
	"os"
	"testing"
)

var e *Evaluator

func TestMain(m *testing.M) {
	e = New()

	code := m.Run()

	os.Exit(code)
}

func TestEvalGetNilKey(t *testing.T) {
	tt := struct {
		cmd         *ast.GetCommand
		expectedVal string
	}{
		cmd: &ast.GetCommand{
			Token: token.New(token.GET, "GET"),
			Key:   "nilKey",
		},
		expectedVal: storage.Nil.String(),
	}

	res, err := e.Eval(tt.cmd)
	if err != nil {
		t.Error(err)
		return
	}

	if res != tt.expectedVal {
		t.Errorf("result expected=%q. got=%q", tt.expectedVal, res)
	}
}

func TestEvalSet(t *testing.T) {
	tt := struct {
		cmd      *ast.SetCommand
		expected *storage.Int
	}{
		cmd: &ast.SetCommand{
			Token: token.New(token.SET, "SET"),
			Key:   "keyInt",
			Value: &ast.IntegerLit{Token: token.New(token.INTEGER, ":"), Value: 1},
		},
		expected: &storage.Int{Value: 1},
	}

	res, err := e.Eval(tt.cmd)
	if err != nil {
		t.Error(err)
		return
	}

	if res != opSuccesful {
		t.Errorf("operation SET was not succesful. res=%q", res)
		return
	}

	val, ok := e.s.Get(tt.cmd.Key)
	if !ok {
		t.Errorf("key=%q not found", tt.cmd.Key)
	}

	if val.String() != tt.expected.String() {
		t.Errorf("value expected=%q. got=%q", tt.expected.String(), val.String())
	}
}

func TestEvalGetNotNil(t *testing.T) {
	tests := []struct {
		key string
		val ast.Expression
	}{
		{
			key: "stringKey",
			val: &ast.StringExpr{
				Token: token.New(token.BULKSTRING, "string"),
			},
		},
		{
			key: "trueKey",
			val: &ast.BooleanExpr{
				Token: token.New(token.BOOLEAN, "#"),
				Value: true,
			},
		},
		{
			key: "falseKey",
			val: &ast.BooleanExpr{
				Token: token.New(token.BOOLEAN, "#"),
				Value: false,
			},
		},
		{
			key: "intKey",
			val: &ast.IntegerLit{
				Token: token.New(token.INTEGER, ":"), Value: 1,
			},
		},
		{
			key: "floatKey",
			val: &ast.FloatExpr{
				Token: token.New(token.FLOAT, ","), Value: 1.1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Creating and Checking key=%q", tt.key), func(t *testing.T) {
			t.Parallel()

			if err := e.s.Set(tt.key, tt.val); err != nil {
				t.Error(err)
			}

			val, ok := e.s.Get(tt.key)
			if !ok {
				t.Errorf("value is nil. key=%q", tt.key)
			}

			expected := storage.NewValue(tt.val)
			if val.String() != expected.String() {
				t.Errorf("val expected=%q. got=%q", val.String(), expected.String())
			}
		})
	}
}
