package main

import (
	"encoding/json"
	"fmt"
	"log"
	"unicode"
	"unicode/utf8"
)

// TokenType represents the type of token.
type TokenType int

const (
	// Special tokens
	ILLEGAL TokenType = iota
	EOF

	// Literals
	IDENT // json key or value
	STRING
	NUMBER

	// Delimiters
	COMMA
	COLON
	LBRACE // {
	RBRACE // }
	LBRACKET
	RBRACKET
)

// Token represents a token returned by the lexer.
type Token struct {
	Type   TokenType // Type of the token
	Value  string    // Value of the token
	Line   int       // Line number in the input (for debugging)
	Column int       // Column number in the input (for debugging)
}

// Lexer represents a JSON lexer.
type Lexer struct {
	input  string // Input string to tokenize
	pos    int    // Current position in the input
	line   int    // Current line in the input (for debugging)
	column int    // Current column in the input (for debugging)
}

// NewLexer creates a new Lexer instance with the given input.
func NewLexer(input string) *Lexer {
	return &Lexer{
		input: input,
		line:  1,
	}
}

// nextToken returns the next token from the input.
func (l *Lexer) nextToken() Token {
	l.skipWhitespace()

	if l.pos >= len(l.input) {
		return Token{Type: EOF, Line: l.line, Column: l.column}
	}

	ch := l.input[l.pos]

	switch {
	case ch == '{':
		l.consume()
		return Token{Type: LBRACE, Value: "{", Line: l.line, Column: l.column}
	case ch == '}':
		l.consume()
		return Token{Type: RBRACE, Value: "}", Line: l.line, Column: l.column}
	case ch == '[':
		l.consume()
		return Token{Type: LBRACKET, Value: "[", Line: l.line, Column: l.column}
	case ch == ']':
		l.consume()
		return Token{Type: RBRACKET, Value: "]", Line: l.line, Column: l.column}
	case ch == ':':
		l.consume()
		return Token{Type: COLON, Value: ":", Line: l.line, Column: l.column}
	case ch == ',':
		l.consume()
		return Token{Type: COMMA, Value: ",", Line: l.line, Column: l.column}
	case ch == '"':
		return l.readString()
	case unicode.IsDigit(rune(ch)) || ch == '-':
		return l.readNumber()
	default:
		if isValidIdentifierStart(ch) {
			return l.readIdentifier()
		}
		return Token{Type: ILLEGAL, Value: string(ch), Line: l.line, Column: l.column}
	}
}

// readString reads a string token from the input.
func (l *Lexer) readString() Token {
	start := l.pos
	l.consume() // consume opening quote
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == '"' && l.input[l.pos-1] != '\\' {
			l.consume() // consume closing quote
			return Token{Type: STRING, Value: l.input[start:l.pos], Line: l.line, Column: l.column}
		}
		l.consume()
	}
	return Token{Type: ILLEGAL, Value: l.input[start:l.pos], Line: l.line, Column: l.column}
}

// readNumber reads a number token from the input.
func (l *Lexer) readNumber() Token {
	start := l.pos
	for l.pos < len(l.input) && (unicode.IsDigit(rune(l.input[l.pos])) || l.input[l.pos] == '.' || l.input[l.pos] == '-' || l.input[l.pos] == 'e' || l.input[l.pos] == 'E') {
		l.consume()
	}
	return Token{Type: NUMBER, Value: l.input[start:l.pos], Line: l.line, Column: l.column}
}

// readIdentifier reads an identifier token from the input.
func (l *Lexer) readIdentifier() Token {
	start := l.pos
	for l.pos < len(l.input) && isValidIdentifierPart(l.input[l.pos]) {
		l.consume()
	}
	return Token{Type: IDENT, Value: l.input[start:l.pos], Line: l.line, Column: l.column}
}

// skipWhitespace skips whitespace characters in the input.
func (l *Lexer) skipWhitespace() {
	for l.pos < len(l.input) && (l.input[l.pos] == ' ' || l.input[l.pos] == '\t' || l.input[l.pos] == '\n' || l.input[l.pos] == '\r') {
		if l.input[l.pos] == '\n' {
			l.line++
			l.column = 0
		}
		l.consume()
	}
}

// consume consumes the current character in the input.
func (l *Lexer) consume() {
	_, size := utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += size
	l.column++
}

// isValidIdentifierStart checks if the given character can be the start of an identifier.
func isValidIdentifierStart(ch byte) bool {
	return ch == '_' || unicode.IsLetter(rune(ch))
}

// isValidIdentifierPart checks if the given character can be part of an identifier.
func isValidIdentifierPart(ch byte) bool {
	return ch == '_' || unicode.IsLetter(rune(ch)) || unicode.IsDigit(rune(ch))
}

func main() {

	input := `{"list": [[[[[[[[[[[[[[["siddhi", {"name" : "siddhant"},            2, [["siddhant"]]]]]]]]]]]]]]]]]}` // valid json
	// input1 := `{[[[[[[[[[[[[[[["siddhi", {"name" : "siddhant"},            2, [["siddhant"]]]]]]]]]]]]]]]]]}` // invalid as key-value pair is required
	// input2 := `[[[[[[[[[[[[[[[[[[["siddhi"]]]]]]]]]]]]]]]]]]]` // valid json as json can be a list
	// input3 := `{"name" : "siddhi", "number" : 1234567890, "marks" : [98, 95, 25],}` // invalid as comma before closing ("}") brace
	// input4 := `[[[[[["siddhi",     [[8,       9]]]]]]]]` // valid json
	lexer := NewLexer(input)

	for {
		token := lexer.nextToken()
		// fmt.Printf("Token: %v\n", token)

		if token.Type == EOF {
			break
		}
	}

	// Parsing JSON using the standard library for comparison
	var jsonData interface{}
	err := json.Unmarshal([]byte(input), &jsonData)
	if err != nil {
		fmt.Println(false)
		log.Fatal(err)
	}
	fmt.Println(true)
	fmt.Println("Parsed JSON:", jsonData)
}
