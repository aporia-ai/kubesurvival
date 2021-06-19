package parser

import "github.com/aporia-ai/kubesurvival/v2/pkg/lexer"

// DataType represents the primitive data types.
type DataType int

const (
	// Unknown primitive data type.
	Unknown DataType = iota
	// Float means the data type is a float.
	Float DataType = 1
	// Integer means the data type is an integer.
	Integer DataType = 2
)

// Operator represents an arithmatic operator.
type Operator int

// Types of operators
const (
	Add      Operator = iota // +
	Multiply                 // *
)

// Node represents a node in the abstract syntax tree.
type Node interface {
	// node is unexported to ensure implementations of Node
	// can only originate in this package.
	node()
}

// Expression is a combination of numbers, variables and operators that
// can be evaluated to a value.
type Expression interface {
	Node
	// expression is unexported to ensure implementations of Expression
	// can only originate in this package.
	expression()
}

// IntLiteral is an expression that contains a single constant integer number.
type IntLiteral struct {
	Value    int64
	Position lexer.Position
}

// String is an expression that contains a string.
type StringLiteral struct {
	Value    string
	Position lexer.Position
}

// ArithmeticExpression is an expression that contains a +, * operator.
type ArithmeticExpression struct {
	LHS      Expression
	Operator Operator
	RHS      Expression
	Position lexer.Position
}

// PodExpression is an expression that represents a pod.
type PodExpression struct {
	CPU      Expression
	Memory   Expression
	GPU      Expression
	Position lexer.Position
}

func (*IntLiteral) node()           {}
func (*StringLiteral) node()        {}
func (*ArithmeticExpression) node() {}
func (*PodExpression) node()        {}

func (*IntLiteral) expression()           {}
func (*StringLiteral) expression()        {}
func (*ArithmeticExpression) expression() {}
func (*PodExpression) expression()        {}
