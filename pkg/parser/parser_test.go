package parser

import (
	"github/dyxgou/redis/pkg/ast"
	"github/dyxgou/redis/pkg/lexer"
	"github/dyxgou/redis/pkg/token"
	"testing"
)

func TestParseGetCommand(t *testing.T) {
	tt := struct {
		input    string
		expected *ast.GetCommand
	}{
		input: "*2\r\n$3\r\nGET\r\n$5\r\nmykey\r\n",
		expected: &ast.GetCommand{
			Token: token.New(token.GET, "GET"),
			Key:   "mykey",
		},
	}

	l := lexer.New(tt.input)
	p := New(l)

	cmd, err := p.Parse()
	if err != nil {
		t.Error(err)
		return
	}

	assertGetCommand(t, cmd, tt.expected)
}

func TestParseSetCommand(t *testing.T) {
	tests := []struct {
		input    string
		expected *ast.SetCommand
	}{
		{"SET mykey myvalue", &ast.SetCommand{
			Key:   "mykey",
			Value: "myvalue",
		}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		cmd, err := p.Parse()
		if err != nil {
			t.Error(err)
		}

		assertSetCommand(t, cmd, tt.expected)
	}
}
