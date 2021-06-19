package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/aporia-ai/kubesurvival/v2/pkg/lexer"
)

type Parser struct {
	Errors    []ParseError
	scanner   *lexer.Scanner
	lookahead lexer.Token
}

// NewParser returns a new instance of Parser.
func NewParser(scanner *lexer.Scanner) *Parser {
	return &Parser{
		Errors:    []ParseError{},
		scanner:   scanner,
		lookahead: scanner.Scan(),
	}
}

// Parse parses an expression and returns its AST representation.
func Parse(s string) (Expression, []ParseError) {
	parser := NewParser(lexer.NewScanner(strings.NewReader(s)))
	return parser.ParseExpression(), parser.Errors
}

func (p *Parser) matchToken(tokenTypes ...lexer.TokenType) (*lexer.Token, bool) {
	for _, tokType := range tokenTypes {
		if tokType == p.lookahead.TokenType {
			token := p.lookahead
			p.lookahead = p.scanner.Scan()
			return &token, true
		}
	}

	return &p.lookahead, false
}

func (p *Parser) match(tokenTypes ...lexer.TokenType) (*lexer.Token, bool) {
	// Try to find the requested token.
	if token, ok := p.matchToken(tokenTypes...); ok {
		return token, true
	}

	return &p.lookahead, false
}

func (p *Parser) skip() {
	p.lookahead = p.scanner.Scan()
}

// ParseExpression parses expressions that might contain any arthimatic operator.
func (p *Parser) ParseExpression() Expression {
	result := p.ParseTerm()
	for p.lookahead.TokenType == lexer.ADD {
		position := p.lookahead.Position
		p.match(lexer.ADD)

		rhs := p.ParseTerm()
		result = &ArithmeticExpression{
			Position: position,
			LHS:      result,
			RHS:      rhs,
			Operator: Add,
		}
	}

	return result
}

// ParseTerm parses expressions that might contain multipications.
func (p *Parser) ParseTerm() Expression {
	var result Expression

	var isLHSInteger = (p.lookahead.TokenType == lexer.INTEGER)
	if isLHSInteger {
		result = p.ParseInteger()
		if p.lookahead.TokenType != lexer.MUL {
			p.addError(newParseError(p.lookahead.Lexeme, []string{"*"}, p.lookahead.Position))
		}
	} else {
		result = p.ParseFactor()
	}

	for p.lookahead.TokenType == lexer.MUL {
		position := p.lookahead.Position
		p.match(lexer.MUL)

		var rhs Expression
		if isLHSInteger {
			rhs = p.ParseFactor()
		} else {
			rhs = p.ParseInteger()
		}

		result = &ArithmeticExpression{
			Position: position,
			LHS:      result,
			RHS:      rhs,
			Operator: Multiply,
		}
	}

	return result
}

// ParseFactor parses a single variable, single constant number or (...some expr...).
func (p *Parser) ParseFactor() Expression {
	switch p.lookahead.TokenType {
	case lexer.LPAREN:
		p.match(lexer.LPAREN)

		expr := p.ParseExpression()

		if token, ok := p.match(lexer.RPAREN); !ok {
			p.addError(newParseError(token.Lexeme, []string{")"}, token.Position))
		}

		return expr

	case lexer.POD:
		return p.ParsePod()

	default:
		p.addError(newParseError(p.lookahead.Lexeme, []string{"(", "pod"},
			p.lookahead.Position))
		return nil
	}
}

func (p *Parser) ParsePod() Expression {
	// pod
	podToken, ok := p.match(lexer.POD)
	if !ok {
		p.addError(newParseError(podToken.Lexeme, []string{"pod"}, podToken.Position))
	}

	// (
	if token, ok := p.match(lexer.LPAREN); !ok {
		p.addError(newParseError(token.Lexeme, []string{"("}, token.Position))
	}

	pod := &PodExpression{Position: podToken.Position}

	for p.lookahead.TokenType != lexer.RPAREN {
		switch p.lookahead.TokenType {
		case lexer.MEMORY:
			p.match(lexer.MEMORY)
			if token, ok := p.match(lexer.COLON); !ok {
				p.addError(newParseError(token.Lexeme, []string{":"}, token.Position))
			}

			pod.Memory = p.ParseStringOrInteger()

		case lexer.CPU:
			p.match(lexer.CPU)
			if token, ok := p.match(lexer.COLON); !ok {
				p.addError(newParseError(token.Lexeme, []string{":"}, token.Position))
			}

			pod.CPU = p.ParseStringOrInteger()

		case lexer.GPU:
			p.match(lexer.GPU)
			if token, ok := p.match(lexer.COLON); !ok {
				p.addError(newParseError(token.Lexeme, []string{":"}, token.Position))
			}

			pod.GPU = p.ParseStringOrInteger()

		default:
			p.addError(newParseError(p.lookahead.Lexeme, []string{"cpu", "memory", "gpu", ")"},
				p.lookahead.Position))
			return pod
		}

		switch p.lookahead.TokenType {
		case lexer.RPAREN:
			p.match(lexer.RPAREN)
			return pod

		case lexer.COMMA:
			p.match(lexer.COMMA)
			continue

		default:
			p.addError(newParseError(p.lookahead.Lexeme, []string{",", ")"},
				p.lookahead.Position))
			return pod
		}
	}

	return pod
}

func (p *Parser) ParseStringOrInteger() Expression {
	switch p.lookahead.TokenType {
	case lexer.STRING:
		return p.ParseString()

	case lexer.INTEGER:
		return p.ParseInteger()

	default:
		p.addError(newParseError(p.lookahead.Lexeme, []string{"STRING", "INTEGER"},
			p.lookahead.Position))
		return nil
	}
}

func (p *Parser) ParseString() Expression {
	token, ok := p.match(lexer.STRING)
	if !ok {
		p.addError(newParseError(token.Lexeme, []string{"STRING"}, token.Position))
	}

	return &StringLiteral{Position: token.Position, Value: token.Lexeme}
}

func (p *Parser) ParseInteger() Expression {
	token, ok := p.match(lexer.INTEGER)
	if !ok {
		p.addError(newParseError(token.Lexeme, []string{"INTEGER"}, token.Position))
	}

	value, err := strconv.ParseInt(token.Lexeme, 10, 64)
	if err != nil {
		p.addError(ParseError{Message: fmt.Sprintf("%s is not number", token.Lexeme)})
	}

	return &IntLiteral{Position: token.Position, Value: value}
}

func (p *Parser) addError(e ParseError) {
	for _, err := range p.Errors {
		if err.Pos == e.Pos {
			return
		}
	}

	p.Errors = append(p.Errors, e)
}
