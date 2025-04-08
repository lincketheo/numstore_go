package compiler

import (
	"strconv"

	"github.com/lincketheo/numstore/internal/numstore"
)

var typeStartTokens = []tokenType{TOK_UNION, TOK_STRUCT, TOK_ENUM, TOK_LEFT_BRACKET, TOK_PRIM}

func (p *parser) parseType() (numstore.Type, bool) {
	tok, ok := p.peekToken()
	if !ok {
		p.parserError(typeStartTokens...)
		return nil, false
	}

	// Check which type is next
	switch tok.ttype {
	case TOK_UNION:
		fields, ok := p.parseStructOrUnion(TOK_UNION)
		return numstore.UnionType{Fields: fields}, ok

	case TOK_STRUCT:
		fields, ok := p.parseStructOrUnion(TOK_STRUCT)
		return numstore.StructType{Fields: fields}, ok

	case TOK_ENUM:
		return p.parseEnum()

	case TOK_LEFT_BRACKET:
		return p.parseArray()

	case TOK_PRIM:
		return p.parsePrimitive()

	default:
		p.parserError(typeStartTokens...)
		return nil, false
	}
}

func (p *parser) parseArray() (numstore.Type, bool) {
	if _, ok := p.expect(TOK_LEFT_BRACKET); !ok {
		p.parserError(TOK_LEFT_BRACKET)
		return nil, false
	}

	// Check if next value is number (strict) or not (variable len)
	if tok, ok := p.peekToken(); ok && tok.ttype == TOK_INTEGER {
		typed, ok := p.parseStrictArray()
		return numstore.Type(typed), ok
	}

	typed, ok := p.parseVarArray()
	return numstore.Type(typed), ok
}

func (p *parser) parseStrictArray() (numstore.StrictArrayType, bool) {
	// Reminder - previous function consumed first token, expect ']' first
	dims := []uint32{}

	for {
		// Parse the value inside the brackets (10 in [10])
		numTok, ok := p.expect(TOK_INTEGER)
		if !ok {
			break
		}

		// Convert to integer
		n, err := strconv.Atoi(numTok.value)
		if err != nil {
			return numstore.StrictArrayType{}, false
		}

		// Dims must be > 0
		if n <= 0 {
			return numstore.StrictArrayType{}, false
		}

		// Append to return
		dims = append(dims, uint32(n))

		// Check for right
		if _, ok := p.expect(TOK_RIGHT_BRACKET); !ok {
			p.parserError(TOK_RIGHT_BRACKET)
			return numstore.StrictArrayType{}, false
		}

		// Done when no more left bracket
		if _, ok := p.expect(TOK_LEFT_BRACKET); !ok {
			break
		}
	}

	// Finally, parse the type
	elem, ok := p.parseType()
	if !ok {
		return numstore.StrictArrayType{}, false
	}

	return numstore.StrictArrayType{
		Dims: dims,
		Of:   elem,
	}, true
}

func (p *parser) parsePrimitive() (numstore.PrimitiveType, bool) {
	// Parse the primitive
	tok, ok := p.expect(TOK_PRIM)
	if !ok {
		p.parserError(TOK_PRIM)
		return numstore.PrimitiveType{}, false
	}

	// Convert string to primitive. This shouldn't fail
	pt, ok := numstore.PrimitiveTypeFromString(tok.value)
	if !ok {
		p.parserError()
		return numstore.PrimitiveType{}, false
	}

	return numstore.PrimitiveType{PT: pt}, true
}

func (p *parser) parseVarArray() (numstore.VarArrayType, bool) {
	// Reminder - previous function consumed first token, expect ']' first
	// Meaning the loop runs at least once
	rank := uint32(0)

	for {
		// End with right bracket
		if _, ok := p.expect(TOK_RIGHT_BRACKET); !ok {
			return numstore.VarArrayType{}, false
		}

		rank++

		// Start with left bracket
		if _, ok := p.expect(TOK_LEFT_BRACKET); !ok {
			break
		}
	}

	// Finally parse the type
	elem, ok := p.parseType()
	if !ok {
		return numstore.VarArrayType{}, false
	}

	return numstore.VarArrayType{
		Rank: rank,
		Of:   elem,
	}, true
}

func (p *parser) parseStructOrUnion(prefixTok tokenType) (map[string]numstore.Type, bool) {
	// Expect 'struct' or 'union'
	if _, ok := p.expect(prefixTok); !ok {
		p.parserError(prefixTok)
		return nil, false
	}

	// Left curly next
	if _, ok := p.expect(TOK_LEFT_CURLY); !ok {
		p.parserError(TOK_LEFT_CURLY)
		return nil, false
	}

	fields := make(map[string]numstore.Type)

	for {
		// First, parse the name of the field
		if nameTok, ok := p.expect(TOK_IDENTIFIER); !ok {
			p.parserError(TOK_IDENTIFIER)
			return nil, false
		} else {

			// Then parse the type afterwards
			fieldType, ok := p.parseType()
			if !ok {
				p.parserError()
				return nil, false
			}

			// Check if this field already exists
			if _, exists := fields[nameTok.value]; exists {
				return nil, false
			}
			fields[nameTok.value] = fieldType
		}

		// Check for comma or right curly
		if _, ok := p.expect(TOK_COMMA); ok {
			continue
		} else if _, ok := p.expect(TOK_RIGHT_CURLY); ok {
			break
		} else {
			p.parserError(TOK_COMMA, TOK_RIGHT_CURLY)
			return nil, false
		}
	}
	return fields, true
}

func (p *parser) parseEnum() (numstore.EnumType, bool) {
	// Parse enum keyword
	if _, ok := p.expect(TOK_ENUM); !ok {
		p.parserError(TOK_ENUM)
		return numstore.EnumType{}, false
	}

	// Next, check for '{'
	if _, ok := p.expect(TOK_LEFT_CURLY); !ok {
		p.parserError(TOK_LEFT_CURLY)
		return numstore.EnumType{}, false
	}

	opts := []string{}

	for {
		// First, the name of the field
		if label, ok := p.expect(TOK_IDENTIFIER); !ok {
			return numstore.EnumType{}, false
		} else {
			opts = append(opts, label.value)
		}

		// Then check for comma or right curly
		if _, ok := p.expect(TOK_COMMA); ok {
			continue
		} else if _, ok := p.expect(TOK_RIGHT_CURLY); ok {
			break
		} else {
			p.parserError(TOK_COMMA, TOK_RIGHT_CURLY)
			return numstore.EnumType{}, false
		}
	}

	return numstore.EnumType{Options: opts}, true
}
