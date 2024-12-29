package ast

// Command nodes represent each command read from a connection
type Node interface {
	String() string
	TokenLiteral() string
}

type Command interface {
	Node

	cmdNode()
}

type Expression interface {
	Node

	exprNode()
}
