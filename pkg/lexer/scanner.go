package lexer

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"strings"
)

// Scanner represents a lexical scanner.
type Scanner struct {
	Reader      *bufio.Reader
	position    Position
	eof         bool
	bufferIndex int
	bufferSize  int
	buffer      [1024]struct {
		ch       rune
		position Position
	}
	DisablePositions bool // for testing.
}

// NewScanner returns a new instance of Scanner.
func NewScanner(reader io.Reader) *Scanner {
	return &Scanner{
		Reader: bufio.NewReader(reader),
	}
}

// Scan returns the next token and literal value.
func (s *Scanner) Scan() Token {
	// Read the next rune.
	ch, pos := s.read()

	for {
		// Skip comments and whitespaces.
		if ch == '#' {
			if err := s.skipUntilEndComment(); err != nil {
				return Token{TokenType: ILLEGAL, Lexeme: "", Position: pos}
			}
		} else if isWhitespace(ch) {
			s.scanWhitespace()
		} else {
			break
		}

		ch, pos = s.read()
	}

	// If we see a letter then consume as an keyword.
	if isLetter(ch) {
		s.Unscan()
		return s.scanKeyword()
	} else if isDigit(ch) {
		s.Unscan()
		return s.scanInteger()
	}

	// Otherwise read the individual character.
	switch ch {
	case eof:
		return Token{TokenType: EOF, Lexeme: "EOF", Position: pos}

	case '"':
		return s.scanString()

	case '(':
		return Token{TokenType: LPAREN, Lexeme: string(ch), Position: pos}

	case ')':
		return Token{TokenType: RPAREN, Lexeme: string(ch), Position: pos}

	case ',':
		return Token{TokenType: COMMA, Lexeme: string(ch), Position: pos}

	case ':':
		return Token{TokenType: COLON, Lexeme: string(ch), Position: pos}

	case '+':
		return Token{TokenType: ADD, Lexeme: string(ch), Position: pos}

	case '*':
		return Token{TokenType: MUL, Lexeme: string(ch), Position: pos}
	}

	return Token{TokenType: ILLEGAL, Lexeme: string(ch), Position: pos}
}

// scanWhitespace consumes the current rune and all contiguous whitespace.
func (s *Scanner) scanWhitespace() {
	// Read every subsequent whitespace character into the buffer.
	// Non-whitespace characters and EOF will cause the loop to exit.
	for {
		if ch, _ := s.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			s.Unscan()
			break
		}
	}
}

// scanKeyword consumes the current rune and all contiguous identifier runes.
func (s *Scanner) scanKeyword() Token {
	ch, pos := s.read()

	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(ch)

	// Read every subsequent ident character into the buffer.
	// Non-ident characters and EOF will cause the loop to exit.
	for {
		if ch, _ = s.read(); ch == eof {
			break
		} else if !isLetter(ch) && !isDigit(ch) && ch != '_' {
			s.Unscan()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	// If the string matches a keyword then return that keyword.
	switch buf.String() {
	case "pod":
		return Token{TokenType: POD, Lexeme: buf.String(), Position: pos}
	case "cpu":
		return Token{TokenType: CPU, Lexeme: buf.String(), Position: pos}
	case "memory":
		return Token{TokenType: MEMORY, Lexeme: buf.String(), Position: pos}
	case "gpu":
		return Token{TokenType: GPU, Lexeme: buf.String(), Position: pos}
	}

	return Token{TokenType: ILLEGAL, Lexeme: buf.String(), Position: pos}
}

// scanInteger consumes a contiguous series of digits.
func (s *Scanner) scanInteger() Token {
	var buf bytes.Buffer
	ch, pos := s.read()

	for {
		if !isDigit(ch) {
			s.Unscan()
			break
		}
		_, _ = buf.WriteRune(ch)
		ch, _ = s.read()
	}

	return Token{TokenType: INTEGER, Lexeme: buf.String(), Position: pos}
}

// scanString consumes a contiguous string of non-quote characters.
// Quote characters can be consumed if they're first escaped with a backslash.
func (s *Scanner) scanString() Token {
	var buf strings.Builder
	ch, pos := s.read()

	for {
		if ch == '"' {
			return Token{TokenType: STRING, Lexeme: buf.String(), Position: pos}
		} else if ch == eof || ch == '\n' {
			return Token{TokenType: BADSTRING, Lexeme: buf.String(), Position: pos}
		} else {
			_, _ = buf.WriteRune(ch)
		}

		ch, _ = s.read()
	}
}

// skipUntilEndComment skips characters until it reaches the end of the line.
func (s *Scanner) skipUntilEndComment() error {
	for {
		if ch, _ := s.read(); ch == '\n' {
			return nil
		} else if ch == eof {
			return io.EOF
		}
	}
}

// read reads the next rune from the bufferred reader.
// Returns the rune(0) if an error occurs (or io.EOF is returned).
func (s *Scanner) read() (rune, Position) {
	// If we have unread characters then read them off the buffer first.
	if s.bufferSize > 0 {
		s.bufferSize--
		return s.curr()
	}

	// Read next rune from underlying reader.
	// Any error (including io.EOF) should return as EOF.
	ch, _, err := s.Reader.ReadRune()
	if err != nil {
		ch = eof
	} else if ch == '\r' {
		if ch, _, err := s.Reader.ReadRune(); err != nil {
			// nop
		} else if ch != '\n' {
			_ = s.Reader.UnreadRune()
		}
		ch = '\n'
	}

	// Save character and position to the buffer.
	s.bufferIndex = (s.bufferIndex + 1) % len(s.buffer)
	buffer := &s.buffer[s.bufferIndex]
	buffer.ch, buffer.position = ch, s.position

	// Update position.
	// Only count EOF once.
	if ch == '\n' {
		s.position.Line++
		s.position.Column = 0
	} else if !s.eof {
		s.position.Column++
	}

	// Mark the reader as EOF.
	// This is used so we don't double count EOF characters.
	if ch == eof {
		s.eof = true
	}

	return s.curr()
}

// curr returns the last read character and position.
func (s *Scanner) curr() (ch rune, pos Position) {
	bufferIndex := (s.bufferIndex - s.bufferSize + len(s.buffer)) % len(s.buffer)
	buffer := &s.buffer[bufferIndex]

	if s.DisablePositions {
		return buffer.ch, Position{}
	}

	return buffer.ch, buffer.position
}

// Unscan pushes the previously token back onto the buffer.
func (s *Scanner) Unscan() {
	s.bufferSize++
}

var errBadString = errors.New("bad string")
var errBadEscape = errors.New("bad escape")
