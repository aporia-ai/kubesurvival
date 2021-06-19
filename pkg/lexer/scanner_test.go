package lexer_test

import (
	"strings"
	"testing"

	"github.com/aporia-ai/kubesurvival/v2/pkg/lexer"
	"github.com/stretchr/testify/assert"
)

func TestScannerOneDigit(t *testing.T) {
	s := lexer.NewScanner(strings.NewReader("9"))
	assertToken(t, s, lexer.INTEGER, "9")
	assertToken(t, s, lexer.EOF, "EOF")
}

func TestScannerMultipleDigits(t *testing.T) {
	s := lexer.NewScanner(strings.NewReader("123"))
	assertToken(t, s, lexer.INTEGER, "123")
	assertToken(t, s, lexer.EOF, "EOF")
}

func TestScannerString(t *testing.T) {
	s := lexer.NewScanner(strings.NewReader("\"abc\""))
	assertToken(t, s, lexer.STRING, "abc")
	assertToken(t, s, lexer.EOF, "EOF")
}

func TestScannerBadString(t *testing.T) {
	s := lexer.NewScanner(strings.NewReader("\"abc"))
	assertToken(t, s, lexer.BADSTRING, "abc")
	assertToken(t, s, lexer.EOF, "EOF")
}

func TestScannerBadString2(t *testing.T) {
	s := lexer.NewScanner(strings.NewReader("\"abc\n"))
	assertToken(t, s, lexer.BADSTRING, "abc")
	assertToken(t, s, lexer.EOF, "EOF")
}

func TestScannerStringAndQuote(t *testing.T) {
	s := lexer.NewScanner(strings.NewReader("\"abc\"\""))
	assertToken(t, s, lexer.STRING, "abc")
	assertToken(t, s, lexer.BADSTRING, "")
	assertToken(t, s, lexer.EOF, "EOF")
}

func TestScanner2Strings(t *testing.T) {
	s := lexer.NewScanner(strings.NewReader("\"abc\" \"test\""))
	assertToken(t, s, lexer.STRING, "abc")
	assertToken(t, s, lexer.STRING, "test")
	assertToken(t, s, lexer.EOF, "EOF")
}

func TestScannerLiterals(t *testing.T) {
	s := lexer.NewScanner(strings.NewReader("\"hello1\" \n\n\n 1234 \t\n\t   \"hhh4h33\" 111 34"))
	assertToken(t, s, lexer.STRING, "hello1")
	assertToken(t, s, lexer.INTEGER, "1234")
	assertToken(t, s, lexer.STRING, "hhh4h33")
	assertToken(t, s, lexer.INTEGER, "111")
	assertToken(t, s, lexer.INTEGER, "34")
	assertToken(t, s, lexer.EOF, "EOF")
}

func TestScannerKeywords(t *testing.T) {
	s := lexer.NewScanner(strings.NewReader(`
		pod cpu   memory
		 gpu gpu pod da
	`))
	assertToken(t, s, lexer.POD, "pod")
	assertToken(t, s, lexer.CPU, "cpu")
	assertToken(t, s, lexer.MEMORY, "memory")
	assertToken(t, s, lexer.GPU, "gpu")
	assertToken(t, s, lexer.GPU, "gpu")
	assertToken(t, s, lexer.POD, "pod")
	assertToken(t, s, lexer.ILLEGAL, "da")
	assertToken(t, s, lexer.EOF, "EOF")
}

func TestScannerSymbols(t *testing.T) {
	s := lexer.NewScanner(strings.NewReader(`(),,    :`))
	assertToken(t, s, lexer.LPAREN, "(")
	assertToken(t, s, lexer.RPAREN, ")")
	assertToken(t, s, lexer.COMMA, ",")
	assertToken(t, s, lexer.COMMA, ",")
	assertToken(t, s, lexer.COLON, ":")
	assertToken(t, s, lexer.EOF, "EOF")
}

func TestScannerOperators(t *testing.T) {
	s := lexer.NewScanner(strings.NewReader(`+ ++ * | |`))
	assertToken(t, s, lexer.ADD, "+")
	assertToken(t, s, lexer.ADD, "+")
	assertToken(t, s, lexer.ADD, "+")
	assertToken(t, s, lexer.MUL, "*")
	assertToken(t, s, lexer.ILLEGAL, "|")
	assertToken(t, s, lexer.ILLEGAL, "|")
	assertToken(t, s, lexer.EOF, "EOF")
}

func TestScannerComments(t *testing.T) {
	s := lexer.NewScanner(strings.NewReader(
		`# hello
		5 + 4 # good 3
		# # comment inside comment 3+2
		pod(6 * 7) # )
		`))

	// 5 + 4
	assertToken(t, s, lexer.INTEGER, "5")
	assertToken(t, s, lexer.ADD, "+")
	assertToken(t, s, lexer.INTEGER, "4")

	// pod(6 * 7)
	assertToken(t, s, lexer.POD, "pod")
	assertToken(t, s, lexer.LPAREN, "(")
	assertToken(t, s, lexer.INTEGER, "6")
	assertToken(t, s, lexer.MUL, "*")
	assertToken(t, s, lexer.INTEGER, "7")
	assertToken(t, s, lexer.RPAREN, ")")

	assertToken(t, s, lexer.EOF, "EOF")
}

func TestScannerPosition(t *testing.T) {
	input := "\"hello\"\n  \"world\"\n\"test\" 314 # comment \n \tpod\n * +"

	s := lexer.NewScanner(strings.NewReader(input))

	hello := assertToken(t, s, lexer.STRING, "hello")
	assert.Truef(t, hello.Position.Line == 0 && hello.Position.Column == 1,
		"Invalid `hello` position - Ln %d, Col %d", hello.Position.Line, hello.Position.Column)

	world := assertToken(t, s, lexer.STRING, "world")
	assert.Truef(t, world.Position.Line == 1 && world.Position.Column == 3,
		"Invalid `world` position - Ln %d, Col %d", world.Position.Line, world.Position.Column)

	test := assertToken(t, s, lexer.STRING, "test")
	assert.Truef(t, test.Position.Line == 2 && test.Position.Column == 1,
		"Invalid `test` position - Ln %d, Col %d", test.Position.Line, test.Position.Column)

	pi := assertToken(t, s, lexer.INTEGER, "314")
	assert.Truef(t, pi.Position.Line == 2 && pi.Position.Column == 7,
		"Invalid `314` position - Ln %d, Col %d", pi.Position.Line, pi.Position.Column)

	pod := assertToken(t, s, lexer.POD, "pod")
	assert.Truef(t, pod.Position.Line == 3 && pod.Position.Column == 2,
		"Invalid `pod` position - Ln %d, Col %d", pod.Position.Line, pod.Position.Column)

	mul := assertToken(t, s, lexer.MUL, "*")
	assert.Truef(t, mul.Position.Line == 4 && mul.Position.Column == 1,
		"Invalid `*` position - Ln %d, Col %d", mul.Position.Line, mul.Position.Column)

	add := assertToken(t, s, lexer.ADD, "+")
	assert.Truef(t, add.Position.Line == 4 && add.Position.Column == 3,
		"Invalid `+` position - Ln %d, Col %d", add.Position.Line, add.Position.Column)
}

func assertToken(t *testing.T, s *lexer.Scanner, tokenType lexer.TokenType, lexeme string) lexer.Token {
	token := s.Scan()
	if token.TokenType != tokenType {
		t.Errorf("Unexpected token %v (lexeme = %v)", token.TokenType, token.Lexeme)
		return token
	} else if token.Lexeme != lexeme {
		t.Errorf("Token %v has unexpected lexeme: %v", token.TokenType, token.Lexeme)
		return token
	}

	t.Logf("Read token %v (lexeme %s)", token.TokenType, token.Lexeme)
	return token
}
