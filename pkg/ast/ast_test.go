package ast

import (
	"github/dyxgou/redis/pkg/token"
	"testing"
)

func TestAstString(t *testing.T) {
	tests := []struct {
		cmd      Command
		expected string
	}{
		{
			&GetCommand{
				Token: token.Token{
					Kind:    token.GET,
					Literal: "GET",
				},
				Key: "mykey",
			},
			"GET mykey",
		},
	}

	for _, tt := range tests {
		s := tt.cmd.String()

		if s != tt.expected {
			t.Errorf("command expected=%q. got=%q", tt.expected, s)
		}
	}
}
