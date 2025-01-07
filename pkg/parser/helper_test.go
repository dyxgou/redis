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

	if setCmd.Value.String() != sc.Value.String() {
		t.Errorf("setCmd value expected=%q. got=%q", sc.Value, setCmd.Value)
	}

	if setCmd.Ex != sc.Ex {
		t.Errorf("setCmd Ex expected=%d. got=%d", sc.Ex, setCmd.Ex)
	}

	if setCmd.Nx != sc.Nx {
		t.Errorf("setCmd Nx expected=%t. got=%t", sc.Nx, setCmd.Nx)
	}

	if setCmd.Xx != sc.Xx {
		t.Errorf("setCmd Xx expected=%t. got=%t", sc.Xx, setCmd.Xx)
	}
}

func assertGetSetCommand(t *testing.T, expr ast.Command, expected *ast.GetSetCommand) {
	t.Helper()
	gsc, ok := expr.(*ast.GetSetCommand)
	if !ok {
		t.Errorf("expr is not *ast.GetSetCommand. got=%T", expr)
	}

	if gsc.Token.Kind != expected.Token.Kind {
		t.Errorf("gscCmd Kind expected=%d. got=%d", expected.Token.Kind, gsc.Token.Kind)
	}

	if gsc.Token.Literal != expected.Token.Literal {
		t.Errorf("gscCmd Literal expected=%q. got=%q", expected.Token.Literal, gsc.Token.Literal)
	}

	if gsc.Value.String() != expected.Value.String() {
		t.Errorf("gscCmd Value expected=%q. got=%q", expected.Value, gsc.Value)
	}
}

func assertBoolean(t *testing.T, expr ast.Expression, expected *ast.BooleanExpr) {
	b, ok := expr.(*ast.BooleanExpr)
	if !ok {
		t.Errorf("expr is not *ast.BooleanExpr. got=%T", expr)
	}

	if b.Token.Kind != expected.Token.Kind {
		t.Errorf("boolExpr Kind expected=%d. got=%d", expected.Token.Kind, b.Token.Kind)
	}

	if b.Token.Literal != expected.Token.Literal {
		t.Errorf("boolExpr Literal expected=%q. got=%q", expected.Token.Literal, b.Token.Literal)
	}

	if b.Value != expected.Value {
		t.Errorf("boolExpr Value expected=%t. got=%t", expected.Value, b.Value)
	}
}

func assertString(t *testing.T, expr ast.Expression, expected *ast.StringExpr) {
	s, ok := expr.(*ast.StringExpr)
	if !ok {
		t.Errorf("expr is not *ast.StringExpr. got=%T", expr)
	}

	if s.Token.Kind != expected.Token.Kind {
		t.Errorf("strExpr Kind expected=%d. got=%d", expected.Token.Kind, s.Token.Kind)
	}

	if s.Token.Literal != expected.Token.Literal {
		t.Errorf("strExpr Literal expected=%q. got=%q", expected.Token.Literal, s.Token.Literal)
	}

	if s.Value() != expected.Value() {
		t.Errorf("strExpr Value() expected=%q. got=%q", expected.Value(), s.Value())
	}
}

func assertInteger(t *testing.T, expr ast.Expression, expected *ast.IntegerLit) {
	i, ok := expr.(*ast.IntegerLit)
	if !ok {
		t.Errorf("expr is not *ast.IntegerExpr. got=%T", expr)
	}

	if i.Token.Kind != expected.Token.Kind {
		t.Errorf("intExpr Kind expected=%d. got=%d", expected.Token.Kind, i.Token.Kind)
	}

	if i.Token.Literal != expected.Token.Literal {
		t.Errorf("intExpr Literal expected=%q. got=%q", expected.Token.Literal, i.Token.Literal)
	}

	if i.Value != expected.Value {
		t.Errorf("intExpr Value expected=%d. got=%d", expected.Value, i.Value)
	}
}

func assertBigInt(t *testing.T, expr ast.Expression, expected *ast.BigIntegerExpr) {
	bi, ok := expr.(*ast.BigIntegerExpr)
	if !ok {
		t.Errorf("expr is not *ast.BigIntegerExpr. got=%T", expr)
		return
	}

	if bi.Token.Kind != expected.Token.Kind {
		t.Errorf("bigIntExpr Kind expected=%d. got=%d", expected.Token.Kind, bi.Token.Kind)
	}

	if bi.Token.Literal != expected.Token.Literal {
		t.Errorf("bigIntExpr Literal expected=%q. got=%q", expected.Token.Literal, bi.Token.Literal)
	}

	if bi.Value != expected.Value {
		t.Errorf("bigIntExpr Value expected=%d. got=%d", expected.Value, bi.Value)
	}
}

func assertFloat(t *testing.T, expr ast.Expression, expected *ast.FloatExpr) {
	fo, ok := expr.(*ast.FloatExpr)
	if !ok {
		t.Errorf("expr is not *ast.FloatExpr. got=%T", expr)
	}

	if fo.Token.Kind != expected.Token.Kind {
		t.Errorf("floatExpr Kind expected=%d. got=%d", expected.Token.Kind, fo.Token.Kind)
	}

	if fo.Token.Literal != expected.Token.Literal {
		t.Errorf("floatExpr Literal expected=%q. got=%q", expected.Token.Literal, fo.Token.Literal)
	}

	if fo.Value != expected.Value {
		t.Errorf("floatExpr Value expected=%f. got=%f", expected.Value, fo.Value)
	}
}
