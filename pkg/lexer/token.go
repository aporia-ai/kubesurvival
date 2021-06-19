package lexer

// Token represents a lexical token.
type TokenType int

const (
	// Special tokens
	ILLEGAL TokenType = iota
	EOF

	// Symbols
	LPAREN // (
	RPAREN // )
	COMMA  // ,
	COLON  // :

	// Keywords
	POD    // pod
	CPU    // cpu
	MEMORY // memory
	GPU    // gpu

	// Operators
	ADD // +
	MUL // *

	// Literals
	INTEGER // 5
	STRING  // "100m"

	// Errors
	BADSTRING // "abc
	BADESCAPE // \q
)

// Position specifies the line and character position of a token.
// The Column and Line are both zero-based indexes.
type Position struct {
	Line   int
	Column int
}

type Token struct {
	TokenType TokenType
	Lexeme    string
	Position  Position
}

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",

	// Symbols
	LPAREN: "(",
	RPAREN: ")",
	COMMA:  ",",
	COLON:  ":",

	// Keywords
	POD:    "pod",
	CPU:    "cpu",
	MEMORY: "memory",
	GPU:    "gpu",

	// Operators
	ADD: "+",
	MUL: "*",

	// Literals
	INTEGER: "INTEGER",
	STRING:  "STRING",

	// Errors
	BADSTRING: "BADSTRING",
}

// String returns the string representation of the token.
func (tok TokenType) String() string {
	if tok >= 0 && tok < TokenType(len(tokens)) {
		return tokens[tok]
	}
	return ""
}
