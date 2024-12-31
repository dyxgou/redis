package parser

import (
	"github/dyxgou/redis/pkg/ast"
	"testing"
)

func assertGetCommand(t *testing.T, cmd ast.Command, gc *ast.GetCommand) {
	t.Helper()
	getCmd, ok := cmd.(*ast.GetCommand)

	if !ok {
		t.Errorf("command expected=*ast.GetCommand. got=%T", cmd)
	}

	if getCmd.Key != gc.Key {
		t.Errorf("getCmd key expected=%q. got=%q", gc.Key, getCmd.Key)
	}
}

func assertSetCommand(t *testing.T, cmd ast.Command, sc *ast.SetCommand) {
	t.Helper()
	setCmd, ok := cmd.(*ast.SetCommand)

	if !ok {
		t.Errorf("command expected=*ast.SetCommand. got=%T", cmd)
	}

	if setCmd.Key != sc.Key {
		t.Errorf("setCmd key expected=%q. got=%q", sc.Key, setCmd.Key)
	}

	if setCmd.Value != sc.Value {
		t.Errorf("setCmd value expected=%q. got=%q", sc.Value, setCmd.Value)
	}

	if setCmd.Ex != sc.Ex {
		t.Errorf("setCmd Ex expected=%d. got=%d", sc.Ex, setCmd.Ex)
	}

	if setCmd.Nx != sc.Nx {
		t.Errorf("setCmd Nx expected=%t. got=%t", sc.Nx, setCmd.Nx)
	}

	if setCmd.Xx != sc.Xx {
		t.Errorf("setCmd Xx expected=%t. got=%t", sc.Nx, setCmd.Nx)
	}
}

func assertBoolean(t *testing.T, expr ast.Expression, expected *ast.BooleanExpr) {
	b, ok := expr.(*ast.BooleanExpr)
	if !ok {
		t.Errorf("expr is not *ast.BooleanExpr. got=%T", expr)
	}

	if b.Token.Kind != expected.Token.Kind {
		t.Errorf("boolExpr Kind expected=%d. got=%d", b.Token.Kind, expected.Token.Kind)
	}

	if b.Token.Literal != expected.Token.Literal {
		t.Errorf("boolExpr Literal expected=%q. got=%q", b.Token.Literal, expected.Token.Literal)
	}

	if b.Value != expected.Value {
		t.Errorf("boolExpr Value expected=%t. got=%t", b.Value, expected.Value)
	}
}

func assertString(t *testing.T, expr ast.Expression, expected *ast.StringExpr) {
	s, ok := expr.(*ast.StringExpr)
	if !ok {
		t.Errorf("expr is not *ast.StringExpr. got=%T", expr)
	}

	if s.Token.Kind != expected.Token.Kind {
		t.Errorf("strExpr Kind expected=%d. got=%d", s.Token.Kind, expected.Token.Kind)
	}

	if s.Token.Literal != expected.Token.Literal {
		t.Errorf("strExpr Literal expected=%q. got=%q", s.Token.Literal, expected.Token.Literal)
	}

	if s.Value() != expected.Value() {
		t.Errorf("strExpr Value() expected=%q. got=%q", s.Value(), expected.Value())
	}
}

func assertInteger(t *testing.T, expr ast.Expression, expected *ast.IntegerExpr) {
	i, ok := expr.(*ast.IntegerExpr)
	if !ok {
		t.Errorf("expr is not *ast.IntegerExpr. got=%T", expr)
	}

	if i.Token.Kind != expected.Token.Kind {
		t.Errorf("intExpr Kind expected=%d. got=%d", i.Token.Kind, expected.Token.Kind)
	}

	if i.Token.Literal != expected.Token.Literal {
		t.Errorf("intExpr Literal expected=%q. got=%q", i.Token.Literal, expected.Token.Literal)
	}

	if i.Value != expected.Value {
		t.Errorf("intExpr Value expected=%d. got=%d", i.Value, expected.Value)
	}
}

func assertBigInt(t *testing.T, expr ast.Expression, expected *ast.BigIntegerExpr) {
	bi, ok := expr.(*ast.BigIntegerExpr)
	if !ok {
		t.Errorf("expr is not *ast.BigIntegerExpr. got=%T", expr)
		return
	}

	if bi.Token.Kind != expected.Token.Kind {
		t.Errorf("bigIntExpr Kind expected=%d. got=%d", bi.Token.Kind, expected.Token.Kind)
	}

	if bi.Token.Literal != expected.Token.Literal {
		t.Errorf("bigIntExpr Literal expected=%q. got=%q", bi.Token.Literal, expected.Token.Literal)
	}

	if bi.Value != expected.Value {
		t.Errorf("bigIntExpr Value expected=%d. got=%d", bi.Value, expected.Value)
	}
}

func assertFloat(t *testing.T, expr ast.Expression, expected *ast.FloatExpr) {
	fo, ok := expr.(*ast.FloatExpr)
	if !ok {
		t.Errorf("expr is not *ast.FloatExpr. got=%T", expr)
	}

	if fo.Token.Kind != expected.Token.Kind {
		t.Errorf("floatExpr Kind expected=%d. got=%d", fo.Token.Kind, expected.Token.Kind)
	}

	if fo.Token.Literal != expected.Token.Literal {
		t.Errorf("floatExpr Literal expected=%q. got=%q", fo.Token.Literal, expected.Token.Literal)
	}

	if fo.Value != expected.Value {
		t.Errorf("floatExpr Value expected=%f. got=%f", fo.Value, expected.Value)
	}
}
