package compiler

import "fmt"

type tokenType int

type token struct {
	ttype tokenType
	value string
	line  int
}

func newToken(ttype tokenType, value string, line int) token {
	return token{ttype: ttype, value: value, line: line}
}

func (t token) String() string {
	return fmt.Sprintf("{Token: %s Value: %s Line %d}",
		t.ttype.String(), t.value, t.line)
}

const (
	// A none type used for code flow
	// not a real token
	TOK_NONE tokenType = iota

	// OP Codes
	TOK_CREATE
	TOK_DELETE
	TOK_READ
	TOK_WRITE
	TOK_OPEN
	TOK_CLOSE
	TOK_TAKE

	// Type Keywords
	TOK_UNION
	TOK_STRUCT
	TOK_ENUM
	TOK_PRIM

	// Structure
	TOK_LEFT_BRACKET
	TOK_RIGHT_BRACKET
	TOK_LEFT_CURLY
	TOK_RIGHT_CURLY
	TOK_LEFT_PAREN
	TOK_RIGHT_PAREN
	TOK_COMMA
	TOK_COLON

	// Values
	TOK_IDENTIFIER
	TOK_STRING
	TOK_INTEGER
	TOK_FLOAT

	// Indicate end of file
	TOK_EOF
)

func (t tokenType) String() string {

	switch t {
	// A none type used for code flow
	// not a real token
	case TOK_NONE:
		return "TOK_NONE"

		// OP Codes
	case TOK_CREATE:
		return "TOK_CREATE"
	case TOK_DELETE:
		return "TOK_DELETE"
	case TOK_READ:
		return "TOK_READ"
	case TOK_WRITE:
		return "TOK_WRITE"
	case TOK_OPEN:
		return "TOK_OPEN"
	case TOK_CLOSE:
		return "TOK_CLOSE"
	case TOK_TAKE:
		return "TOK_TAKE"

	// Type Keywords
	case TOK_UNION:
		return "TOK_UNION"
	case TOK_STRUCT:
		return "TOK_STRUCT"
	case TOK_ENUM:
		return "TOK_ENUM"
	case TOK_PRIM:
		return "TOK_PRIM"

		// Structure
	case TOK_LEFT_BRACKET:
		return "TOK_LEFT_BRACKET"
	case TOK_RIGHT_BRACKET:
		return "TOK_RIGHT_BRACKET"
	case TOK_LEFT_CURLY:
		return "TOK_LEFT_CURLY"
	case TOK_RIGHT_CURLY:
		return "TOK_RIGHT_CURLY"
	case TOK_RIGHT_PAREN:
		return "TOK_RIGHT_PAREN"
	case TOK_LEFT_PAREN:
		return "TOK_LEFT_PAREN"
	case TOK_COMMA:
		return "TOK_COMMA"
	case TOK_COLON:
		return "TOK_COLON"

		// Values
	case TOK_IDENTIFIER:
		return "TOK_IDENTIFIER"
	case TOK_STRING:
		return "TOK_STRING"
	case TOK_INTEGER:
		return "TOK_INTEGER"
	case TOK_FLOAT:
		return "TOK_FLOAT"

	case TOK_EOF:
		return "TOK_EOF"
	}
	panic("Unknown tokenType")
}
