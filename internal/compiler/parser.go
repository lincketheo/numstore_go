package compiler

import (
	"github.com/lincketheo/numstore/internal/bytecode"
)

// ////////////////////////////// PARSER
type parser struct {
	tokens  []token
	cur     int
	isError bool
	bc      bytecode.ByteCode
}


// ////////////////////////////// PARSE

type tokenHandler func() bool

func Parse(tokens []token) ([]byte, bool) {
	if len(tokens) == 0 {
		return []byte{}, true
	}

	p := parser{
		tokens: tokens,
		cur:    0,
	}

	// Create a map of token types to their corresponding handler functions.
	handlers := map[tokenType]tokenHandler{
		TOK_CREATE: p.handleTokCreate,
		TOK_DELETE: p.handleTokDelete,
		TOK_READ:   p.handleTokRead,
		TOK_WRITE:  p.handleTokWrite,
		TOK_TAKE:   p.handleTokTake,
	}

	// Iterate over tokens.
	for t, ok := p.nextToken(); ok; t, ok = p.nextToken() {
		switch t.ttype {

		case TOK_EOF:
			return p.bc.Bytes, true

		default:
			{
				if handler, exists := handlers[t.ttype]; exists {
					if !handler() {
						p.windNext()
					}

				} else {
					p.parserError()
					p.windNext()
				}
			}
		}
	}

	// TODO - get rid of this duplicate code
	return p.bc.Bytes, true
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
		t.bc.HandleCreate(v.value, nstype)
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
		t.bc.HandleDelete(v.value)
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
		t.bc.HandleRead(rfmt)
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
		t.bc.HandleWrite(wfmt)
		return true
	}
}

func (t *parser) handleTokTake() bool {
	// take VAR RFMT

	if rfmt, ok := t.parseRFMT(); !ok {
		// Parse RFMT
		return false

	} else {
		// Execute
		t.bc.HandleTake(rfmt)
		return true
	}
}
