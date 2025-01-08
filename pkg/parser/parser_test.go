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

func TestSetCommand(t *testing.T) {
	tests := []struct {
		input    string
		expected ast.SetCommand
	}{
		{"*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n:1\r\n", ast.SetCommand{
			Token: token.New(token.SET, "SET"),
			Key:   "key",
			Value: &ast.IntegerLit{Token: token.New(token.INTEGER, ":"), Value: 1},
		}},
		{"*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$8\r\n\"my string\"\r\n", ast.SetCommand{
			Token: token.New(token.SET, "SET"),
			Key:   "key",
			Value: &ast.StringExpr{Token: token.New(token.BULKSTRING, "my string")},
		}},
		{"*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n#t\r\n", ast.SetCommand{
			Token: token.New(token.SET, "SET"),
			Key:   "key",
			Value: &ast.BooleanExpr{Token: token.New(token.BOOLEAN, "#"), Value: true},
		}},
		{"*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n#f\r\n", ast.SetCommand{
			Token: token.New(token.SET, "SET"),
			Key:   "key",
			Value: &ast.BooleanExpr{Token: token.New(token.BOOLEAN, "#"), Value: false},
		}},
		{"*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n:1\r\n", ast.SetCommand{
			Token: token.New(token.SET, "SET"),
			Key:   "key",
			Value: &ast.IntegerLit{Token: token.New(token.INTEGER, ":"), Value: 1},
		}},
		{"*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n:1\r\nEX\r\n:3600\r\n", ast.SetCommand{
			Token: token.New(token.SET, "SET"),
			Key:   "key",
			Value: &ast.IntegerLit{Token: token.New(token.INTEGER, ":"), Value: 1},
			Ex:    3600,
		}},
		{"*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n:1\r\nEX\r\n:3600\r\nXX\r\n", ast.SetCommand{
			Token: token.New(token.SET, "SET"),
			Key:   "key",
			Value: &ast.IntegerLit{Token: token.New(token.INTEGER, ":"), Value: 1},
			Ex:    3600,
			Xx:    true,
		}},
		{"*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n:1\r\nEX\r\n:3600\r\nNX\r\n", ast.SetCommand{
			Token: token.New(token.SET, "SET"),
			Key:   "key",
			Value: &ast.IntegerLit{Token: token.New(token.INTEGER, ":"), Value: 1},
			Ex:    3600,
			Nx:    true,
		}},
		{"*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n:1\r\nNX\r\n", ast.SetCommand{
			Token: token.New(token.SET, "SET"),
			Key:   "key",
			Value: &ast.IntegerLit{Token: token.New(token.INTEGER, ":"), Value: 1},
			Nx:    true,
		}},
		{"*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n:1\r\nXX\r\n", ast.SetCommand{
			Token: token.New(token.SET, "SET"),
			Key:   "key",
			Value: &ast.IntegerLit{Token: token.New(token.INTEGER, ":"), Value: 1},
			Xx:    true,
		}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		cmd, err := p.Parse()
		if err != nil {
			t.Error(err)
			return
		}

		assertSetCommand(t, cmd, &tt.expected)
	}
}

func TestParseGetSetCommand(t *testing.T) {
	tt := struct {
		input    string
		expected ast.GetSetCommand
	}{
		input: "*3\r\n$6\r\nGETSET\r\n$3\r\nkey\r\n#t\r\n",
		expected: ast.GetSetCommand{
			Token: token.New(token.GETSET, "GETSET"),
			Key:   "key",
			Value: &ast.BooleanExpr{Token: token.New(token.BOOLEAN, "#"), Value: true},
		},
	}

	l := lexer.New(tt.input)
	p := New(l)

	cmd, err := p.Parse()
	if err != nil {
		t.Error(err)
		return
	}

	assertGetSetCommand(t, cmd, &tt.expected)
}

func TestParseGetExCommand(t *testing.T) {
	tt := struct {
		input    string
		expected ast.GetExCommand
	}{
		input: "*4\r\n$6\r\nGETEX\r\n$3\r\nkey\r\n$2\r\nEX\r\n:20\r\n",
		expected: ast.GetExCommand{
			Token: token.New(token.GETEX, "GETEX"),
			Key:   "key",
			Ex:    20,
		},
	}

	l := lexer.New(tt.input)
	p := New(l)

	cmd, err := p.Parse()
	if err != nil {
		t.Error(err)
		return
	}

	assertGetExCommand(t, cmd, &tt.expected)
}

func TestParseGetDel(t *testing.T) {
	tt := struct {
		input    string
		expected ast.GetDelCommand
	}{
		input: "*2\r\n$6\r\nGETDEL\r\n$3\r\nkey\r\n",
		expected: ast.GetDelCommand{
			Token: token.New(token.GETDEL, "GETDEL"),
			Key:   "key",
		},
	}

	l := lexer.New(tt.input)
	p := New(l)

	cmd, err := p.Parse()
	if err != nil {
		t.Error(err)
		return
	}

	gdCmd, ok := cmd.(*ast.GetDelCommand)

	if !ok {
		t.Errorf("command expected=*ast.GetDelCommand. got=%T", cmd)
	}

	if gdCmd.Key != tt.expected.Key {
		t.Errorf("gdCmd key expected=%q. got=%q", tt.expected.Key, gdCmd.Key)
	}
}

func TestParseIncrByCommand(t *testing.T) {
	tt := struct {
		input    string
		expected ast.IncrByCommand
	}{
		input: "*2\r\n$6\r\nINCRBY\r\n$3\r\nkey\r\n:123\r\n",
		expected: ast.IncrByCommand{
			Token:     token.New(token.INCRBY, "INCRBY"),
			Key:       "key",
			Increment: 123,
		},
	}

	l := lexer.New(tt.input)
	p := New(l)

	cmd, err := p.Parse()
	if err != nil {
		t.Error(err)
		return
	}

	incCmd, ok := cmd.(*ast.IncrByCommand)

	if !ok {
		t.Errorf("command expected=*ast.IncrByCommand. got=%T", cmd)
	}

	if incCmd.Key != tt.expected.Key {
		t.Errorf("incCmd key expected=%d. got=%d", tt.expected.Increment, incCmd.Increment)
	}

	if incCmd.Increment != tt.expected.Increment {
		t.Errorf("incCmd increment expected=%d. got=%d", tt.expected.Increment, incCmd.Increment)
	}
}

func TestParseIncrCommand(t *testing.T) {
	tt := struct {
		input    string
		expected ast.IncrCommand
	}{
		input: "*2\r\n$2\r\nINCR\r\n$3\r\nkey\r\n",
		expected: ast.IncrCommand{
			Token: token.New(token.INCR, "INCR"),
			Key:   "key",
		},
	}

	l := lexer.New(tt.input)
	p := New(l)

	cmd, err := p.Parse()
	if err != nil {
		t.Error(err)
		return
	}

	gdCmd, ok := cmd.(*ast.IncrCommand)

	if !ok {
		t.Errorf("command expected=*ast.IncrCommand. got=%T", cmd)
	}

	if gdCmd.Key != tt.expected.Key {
		t.Errorf("gdCmd key expected=%q. got=%q", tt.expected.Key, gdCmd.Key)
	}
}

func TestParseDecrByCommand(t *testing.T) {
	tt := struct {
		input    string
		expected ast.IncrByCommand
	}{
		input: "*2\r\n$6\r\nINCRBY\r\n$3\r\nkey\r\n:123\r\n",
		expected: ast.IncrByCommand{
			Token:     token.New(token.INCRBY, "INCRBY"),
			Key:       "key",
			Increment: 123,
		},
	}

	l := lexer.New(tt.input)
	p := New(l)

	cmd, err := p.Parse()
	if err != nil {
		t.Error(err)
		return
	}

	incCmd, ok := cmd.(*ast.IncrByCommand)

	if !ok {
		t.Errorf("command expected=*ast.IncrByCommand. got=%T", cmd)
	}

	if incCmd.Key != tt.expected.Key {
		t.Errorf("incCmd key expected=%d. got=%d", tt.expected.Increment, incCmd.Increment)
	}

	if incCmd.Increment != tt.expected.Increment {
		t.Errorf("incCmd increment expected=%d. got=%d", tt.expected.Increment, incCmd.Increment)
	}
}

func TestParseDecrCommand(t *testing.T) {
	tt := struct {
		input    string
		expected ast.DecrCommand
	}{
		input: "*2\r\n$4\r\nDECR\r\n$3\r\nkey\r\n",
		expected: ast.DecrCommand{
			Token: token.New(token.DECR, "DECR"),
			Key:   "key",
		},
	}

	l := lexer.New(tt.input)
	p := New(l)

	cmd, err := p.Parse()
	if err != nil {
		t.Error(err)
		return
	}

	decCmd, ok := cmd.(*ast.DecrCommand)

	if !ok {
		t.Errorf("command expected=*ast.DecrCommand. got=%T", cmd)
	}

	if decCmd.Key != tt.expected.Key {
		t.Errorf("decCmd key expected=%q. got=%q", tt.expected.Key, decCmd.Key)
	}
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
			input: "$7\r\n\"string 1\"",
			expected: ast.StringExpr{
				Token: token.New(token.BULKSTRING, "string 1"),
			},
		},
		{
			input: "$7\r\n\"string 2\"",
			expected: ast.StringExpr{
				Token: token.New(token.BULKSTRING, "string 2"),
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
		expected ast.IntegerLit
	}{
		{
			input: ":5",
			expected: ast.IntegerLit{
				Token: token.New(token.INTEGER, ":"),
				Value: 5,
			},
		},
		{
			input: ":5614",
			expected: ast.IntegerLit{
				Token: token.New(token.INTEGER, ":"),
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
			input: "(5123123123123123",
			expected: ast.BigIntegerExpr{
				Token: token.New(token.BIGINT, "("),
				Value: 5123123123123123,
			},
		},
		{
			input: "(5614546145614456123",
			expected: ast.BigIntegerExpr{
				Token: token.New(token.BIGINT, "("),
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
			input: ",123.123",
			expected: ast.FloatExpr{
				Token: token.New(token.FLOAT, ","),
				Value: 123.123,
			},
		},
		{
			input: ",321.321",
			expected: ast.FloatExpr{
				Token: token.New(token.FLOAT, ","),
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
