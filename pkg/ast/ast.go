package ast

import "strings"

type Node interface {
	String() string
	TokenLiteral() string
}

// Command nodes represent each command read from a connection
type Command interface {
	Node

	commandNode()
}

// CommandReader is the reader we are gonna use in each connection which is gonna transfrom the string input into an AST.
type CommandReader struct {
	Commands []Command
}

// NewCommandReader() creates a new *CommandReader with a default lenght of 20
func NewCommandReader() *CommandReader {
	return &CommandReader{
		Commands: make([]Command, 0, 20),
	}
}

func (cr *CommandReader) String() string {
	var sb strings.Builder

	for _, cmd := range cr.Commands {
		sb.WriteString(cmd.String())
	}

	return sb.String()
}

func (cr *CommandReader) TokenLiteral() string {
	if len(cr.Commands) == 0 {
		return ""
	}

	return cr.Commands[0].TokenLiteral()
}
