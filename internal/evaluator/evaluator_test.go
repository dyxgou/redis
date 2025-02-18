package evaluator

import (
	"context"
	"fmt"
	"github/dyxgou/redis/internal/storage"
	"github/dyxgou/redis/pkg/ast"
	"github/dyxgou/redis/pkg/token"
	"os"
	"testing"
	"time"
)

var e *Evaluator

func TestMain(m *testing.M) {
	e = New(context.Background())

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

func TestEvalSetEx(t *testing.T) {
	tt := struct {
		cmd      *ast.SetCommand
		expected *storage.Int
	}{
		cmd: &ast.SetCommand{
			Token: token.New(token.SET, "SET"),
			Key:   "keyInt",
			Value: &ast.IntegerLit{Token: token.New(token.INTEGER, ":"), Value: 1},
			Ex:    1,
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
	t.Run("check deleted key", func(t *testing.T) {
		t.Parallel()
		time.Sleep(2 * time.Second)
		if ok = e.s.Exists(tt.cmd.Key); ok {
			t.Errorf("key=%q still exists after EX setted", tt.cmd.Key)
		}
	})
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

func TestEvalGetSet(t *testing.T) {
	tt := struct {
		cmd      *ast.GetSetCommand
		expected *storage.Int
	}{
		cmd: &ast.GetSetCommand{
			Token: token.New(token.GETSET, "GETSET"),
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

	if res != tt.expected.String() {
		t.Errorf("GETSET value expected=%q. got=%q", tt.expected.String(), res)
	}
}

func TestEvalGetDel(t *testing.T) {
	tt := struct {
		cmd      *ast.GetDelCommand
		expected *storage.String
	}{
		cmd: &ast.GetDelCommand{
			Token: token.New(token.GETDEL, "GELDEl"),
			Key:   "valkey",
		},
		expected: &storage.String{Value: "value deleted"},
	}
	e.s.Set("valkey", &ast.StringExpr{Token: token.New(token.BULKSTRING, "value deleted")})

	res, err := e.Eval(tt.cmd)
	if err != nil {
		t.Error(err)
		return
	}

	if res != tt.expected.String() {
		t.Errorf("res expected=%q. got=%q", res, tt.expected.String())
	}

	ok := e.s.Exists(tt.cmd.Key)
	if ok {
		t.Errorf("key=%q still exists after GETDEL command", tt.cmd.Key)
	}
}
