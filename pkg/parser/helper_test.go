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
