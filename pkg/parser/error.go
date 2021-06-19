package parser

import (
	"fmt"
	"strings"

	"github.com/aporia-ai/kubesurvival/v2/pkg/lexer"
)

// ParseError represents an error that occurred during parsing.
type ParseError struct {
	Message  string
	Found    string
	Expected []string
	Pos      lexer.Position
}

// newParseError returns a new instance of ParseError.
func newParseError(found string, expected []string, pos lexer.Position) ParseError {
	return ParseError{Found: found, Expected: expected, Pos: pos}
}

// Error returns the string representation of the error.
func (e *ParseError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("%s at line %d, char %d", e.Message, e.Pos.Line+1, e.Pos.Column+1)
	}
	return fmt.Sprintf("found %s, expected %s at line %d, char %d", e.Found,
		strings.Join(e.Expected, ", "), e.Pos.Line+1, e.Pos.Column+1)
}
