package ast

// Node represent either a command or an expression
type Node interface {
	String() string
	TokenLiteral() string
}

// Command reprensets a Redis Command
type Command interface {
	Node

	cmdNode()
}

// Expression represents a value of any command
type Expression interface {
	Node

	exprNode()
}
