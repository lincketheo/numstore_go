package compiler

import (
	"fmt"
	"github.com/lincketheo/numstore/internal/logging"
)

// ////////////////////////////// PARSER
type parser struct {
	tokens  []token
	cur     int
	isError bool
}

func (t parser) isEnd() bool {
	return t.cur >= len(t.tokens)
}

func (t *parser) parserError(exp ...tokenType) {
	t.isError = true
	logging.Error("Expected token types: %v", exp)
}

func (t *parser) nextToken() (token, bool) {
	if t.isEnd() {
		return token{}, false
	}

	cur := t.cur
	t.cur += 1
	return t.tokens[cur], true
}

func (t parser) peekToken() (token, bool) {
	if t.isEnd() {
		return token{}, false
	}

	return t.tokens[t.cur], true
}

func (p *parser) expect(ttypes ...tokenType) (token, bool) {
	tok, ok := p.peekToken()
	if !ok {
		return token{}, false
	}
	match := false
	for _, expected := range ttypes {
		if tok.ttype == expected {
			match = true
			break
		}
	}
	if !match {
		return token{}, false
	}
	return p.nextToken()
}

type tokenHandler func() bool

// ////////////////////////////// PARSE
func Parse(tokens []token) bool {
	if len(tokens) == 0 {
		return true
	}

	runner := parser{
		tokens: tokens,
		cur:    0,
	}

	// Create a map of token types to their corresponding handler functions.
	handlers := map[tokenType]tokenHandler{
		TOK_CREATE: runner.handleTokCreate,
		TOK_DELETE: runner.handleTokDelete,
		TOK_READ:   runner.handleTokRead,
		TOK_WRITE:  runner.handleTokWrite,
		TOK_TAKE:   runner.handleTokTake,
	}

	// Iterate over tokens.
	for t, ok := runner.nextToken(); ok; t, ok = runner.nextToken() {
		switch t.ttype {
		case TOK_EOF:
			return true
		default:
			// If a handler exists for the token, call it.
			if handler, exists := handlers[t.ttype]; exists {
				if !handler() {
					return false
				}
			} else {
				runner.parserError()
				return false
			}
		}
	}

	return true
}

// DONE
func (t *parser) handleTokCreate() bool {
	// create VAR TYPE

	if v, ok := t.expect(TOK_IDENTIFIER); !ok {
		// Parse VAR
		t.parserError(TOK_IDENTIFIER)
		return false

	} else if nstype, ok := t.parseType(); !ok {
		// PARSE TYPE
		return false

	} else {
		// Execute
		fmt.Printf("CREATING variable: %v with type %v\n", v.value, nstype)
		return true
	}
}

func (t *parser) handleTokDelete() bool {
	// delete VAR

	if v, ok := t.expect(TOK_IDENTIFIER); !ok {
		// Parse VAR
		t.parserError(TOK_IDENTIFIER)
		return false

	} else {
		// Execute
		fmt.Printf("DELETING: %v\n", v.value)
		return true
	}
}

func (t *parser) handleTokRead() bool {
	// read RFMT

	if rfmt, ok := t.parseRFMT(); !ok {
		// Parse RFMT
		return false

	} else {
		// Execute
		fmt.Printf("READING: RFMT: %v\n", rfmt)
		return true
	}
}

// DONE
func (t *parser) handleTokWrite() bool {
	// write WFMT

	if wfmt, ok := t.parseWFMT(); !ok {
		// Parse WFMT
		return false

	} else {
		// Execute
		fmt.Printf("Writing: WFMT: %v\n", wfmt)
		return true
	}
}

func (t *parser) handleTokTake() bool {
	// take VAR RFMT

	if v, ok := t.expect(TOK_IDENTIFIER); !ok {
		// Parse VAR
		return false

	} else if rfmt, ok := t.parseRFMT(); !ok {
		// Parse RFMT
		return false

	} else {
		// Execute
		fmt.Printf("TAKING: %v with RFMT: %v\n", v.value, rfmt)
		return true
	}
}
