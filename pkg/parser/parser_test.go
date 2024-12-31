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

func TestParseBoolean(t *testing.T) {
	tests := []struct {
		input    string
		expected ast.BooleanExpr
	}{
		{
			input: "#t",
			expected: ast.BooleanExpr{
				Token: token.New(token.BOOLEAN, "#"),
				Value: true,
			},
		},
		{
			input: "#f",
			expected: ast.BooleanExpr{
				Token: token.New(token.BOOLEAN, "#"),
				Value: false,
			},
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		val, err := p.parseValue()
		if err != nil {
			t.Error(err)
			return
		}

		assertBoolean(t, val, &tt.expected)
	}
}

func TestParseString(t *testing.T) {
	tests := []struct {
		input    string
		expected ast.StringExpr
	}{
		{
			input: "\"string 1\"",
			expected: ast.StringExpr{
				Token: token.New(token.STRING, "string 1"),
			},
		},
		{
			input: "\"string 2\"",
			expected: ast.StringExpr{
				Token: token.New(token.STRING, "string 2"),
			},
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		val, err := p.parseValue()
		if err != nil {
			t.Error(err)
			return
		}

		assertString(t, val, &tt.expected)
	}
}

func TestParseInteger(t *testing.T) {
	tests := []struct {
		input    string
		expected ast.IntegerExpr
	}{
		{
			input: "5",
			expected: ast.IntegerExpr{
				Token: token.New(token.INTEGER, "5"),
				Value: 5,
			},
		},
		{
			input: "5614",
			expected: ast.IntegerExpr{
				Token: token.New(token.INTEGER, "5614"),
				Value: 5614,
			},
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		val, err := p.parseValue()
		if err != nil {
			t.Error(err)
			return
		}

		assertInteger(t, val, &tt.expected)
	}
}

func TestParseBigNumber(t *testing.T) {
	tests := []struct {
		input    string
		expected ast.BigIntegerExpr
	}{
		{
			input: "5123123123123123",
			expected: ast.BigIntegerExpr{
				Token: token.New(token.BIGNUMBER, "5123123123123123"),
				Value: 5123123123123123,
			},
		},
		{
			input: "5614546145614456123",
			expected: ast.BigIntegerExpr{
				Token: token.New(token.BIGNUMBER, "5614546145614456123"),
				Value: 5614546145614456123,
			},
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		val, err := p.parseValue()
		if err != nil {
			t.Error(err)
			return
		}

		assertBigInt(t, val, &tt.expected)
	}
}

func TestParseFloat(t *testing.T) {
	tests := []struct {
		input    string
		expected ast.FloatExpr
	}{
		{
			input: "123.123",
			expected: ast.FloatExpr{
				Token: token.New(token.FLOAT, "123.123"),
				Value: 123.123,
			},
		},
		{
			input: "321.321",
			expected: ast.FloatExpr{
				Token: token.New(token.FLOAT, "321.321"),
				Value: 321.321,
			},
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		val, err := p.parseValue()
		if err != nil {
			t.Error(err)
			return
		}

		assertFloat(t, val, &tt.expected)
	}
}
