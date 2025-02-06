package bytecode

type tokenType int

type token struct {
	ttype tokenType
	value string
	line  int
}

const (
	TOK_NONE tokenType = iota

	// OP Codes
	TOK_CREATE
	TOK_DELETE
	TOK_READ
	TOK_WRITE

	// Type Identifiers
	TOK_DATABASE
	TOK_RELATION
	TOK_VARIABLE
	//TOK_I

	// Options for read / write
	TOK_VALUES
	TOK_OPEN
	TOK_SORTBY
	TOK_ASC
	TOK_DESC

	// Structure
	TOK_LEFT_BRACKET
	TOK_RIGHT_BRACKET
	TOK_LEFT_PAREN
	TOK_RIGHT_PAREN
	TOK_COMMA

	// Indexes and slices
	TOK_COLON
	TOK_NUMBER

	// Expressions (for the future)
	//TOK_PLUS // For complex
	//TOK_MINUS
	//TOK_MULT
	//TOK_DIV

	// Other
	TOK_IDENTIFIER
	TOK_EOF
)

func (t tokenType) String() string {

	switch t {
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

	// Type Identifiers
	case TOK_DATABASE:
		return "TOK_DATABASE"
	case TOK_RELATION:
		return "TOK_RELATION"
	case TOK_VARIABLE:
		return "TOK_VARIABLE"
	//case TOK_I:
	//	return "TOK_I"

	// Options for read / write
	case TOK_VALUES:
		return "TOK_VALUES"
	case TOK_OPEN:
		return "TOK_OPEN"
	case TOK_SORTBY:
		return "TOK_SORTBY"
	case TOK_ASC:
		return "TOK_ASC"
	case TOK_DESC:
		return "TOK_DESC"

	// Structure
	case TOK_LEFT_BRACKET:
		return "TOK_LEFT_BRACKET"
	case TOK_RIGHT_BRACKET:
		return "TOK_RIGHT_BRACKET"
	case TOK_LEFT_PAREN:
		return "TOK_LEFT_PAREN"
	case TOK_RIGHT_PAREN:
		return "TOK_RIGHT_PAREN"
	case TOK_COMMA:
		return "TOK_COMMA"

	// Indexes and slices
	case TOK_COLON:
		return "TOK_COLON"
	case TOK_NUMBER:
		return "TOK_NUMBER"

	// Expressions
	//case TOK_PLUS:
	//		return "TOK_PLUS"
	//case TOK_MINUS:
	//		return "TOK_MINUS"
	//case TOK_MULT:
	//		return "TOK_MULT"
	//case TOK_DIV:
	//		return "TOK_DIV"

	// Other
	case TOK_IDENTIFIER:
		return "TOK_IDENTIFIER"
	case TOK_EOF:
		return "TOK_EOF"
	}
	panic("Unknown tokenType")
}
