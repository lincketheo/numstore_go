package compiler

import (
	"fmt"
	"strconv"

	"github.com/lincketheo/numstore/internal/numstore"
)

func (p *parser) parseType() (numstore.Type, error) {
	tok, ok := p.peekToken()
	if !ok {
		return nil, expectedStrButEOF("type")
	}

	switch tok.ttype {
	case TOK_UNION:
		{
			if fields, err := p.parseStructOrUnion(TOK_UNION); err != nil {
				return nil, err
			} else {
				return numstore.UnionType{Fields: fields}, nil
			}
		}

	case TOK_STRUCT:
		{
			if fields, err := p.parseStructOrUnion(TOK_STRUCT); err != nil {
				return nil, err
			} else {
				return numstore.StructType{Fields: fields}, nil
			}
		}

	case TOK_ENUM:
		return p.parseEnum()

	case TOK_LEFT_BRACKET:
		return p.parseArray()

	case TOK_PRIM:
		return p.parsePrimitive()

	default:
		return nil, fmt.Errorf("expected type, got %v", tok)
	}
}

func (p *parser) parseArray() (numstore.Type, error) {
	if _, err := p.expect(TOK_LEFT_BRACKET); err != nil {
		return nil, err
	}

	// Check if next value is number (strict) or not (variable len)
	if tok, _ := p.peekToken(); tok.ttype == TOK_INTEGER {
		return p.parseStrictArray()
	}

	return p.parseVarArray()
}

func (p *parser) parseStrictArray() (numstore.StrictArrayType, error) {
	// Reminder - previous function consumed first token, expect ']' first

	dims := []uint32{}

	for {
		// Parse Number
		if num, err := p.expect(TOK_INTEGER); err != nil {
			break

			// Failed to parse
		} else if n, err := strconv.Atoi(num.value); err != nil {
			return numstore.StrictArrayType{}, err

			// Number is < 0
		} else if n <= 0 {
			return numstore.StrictArrayType{},
				fmt.Errorf("Invalid array shape: %d, expecting shape > 0", n)

			// Success
		} else {
			dims = append(dims, uint32(n))
		}

		// Expect "]"
		if _, err := p.expect(TOK_RIGHT_BRACKET); err != nil {
			return numstore.StrictArrayType{}, err
		}

		// Maybe Expect "["
		if _, err := p.expect(TOK_LEFT_BRACKET); err != nil {
			break
		}
	}

	elem, err := p.parseType()
	if err != nil {
		return numstore.StrictArrayType{}, err
	}
	return numstore.StrictArrayType{
		Rank: uint32(len(dims)),
		Dims: dims,
		Of:   elem,
	}, nil
}

func (p *parser) parsePrimitive() (numstore.PrimitiveType, error) {
	tok, err := p.expect(TOK_PRIM)
	if err != nil {
		return numstore.PrimitiveType{}, err
	}

	pt, ok := numstore.PrimitiveTypeFromString(tok.value)
	if !ok {
		return numstore.PrimitiveType{},
			fmt.Errorf("Invalid primitive %q", tok.value)
	}

	return numstore.PrimitiveType{PT: pt}, nil
}

func (p *parser) parseVarArray() (numstore.VarArrayType, error) {
	// Reminder - previous function consumed first token, expect ']' first

	rank := uint32(0)
	for {
		// Expect "]"
		if _, err := p.expect(TOK_RIGHT_BRACKET); err != nil {
			return numstore.VarArrayType{}, fmt.Errorf("expected ] for var array")
		}
		rank++

		// Maybe Expect "["
		if _, err := p.expect(TOK_LEFT_BRACKET); err != nil {
			break
		}
	}

	elem, err := p.parseType()
	if err != nil {
		return numstore.VarArrayType{}, err
	}

	return numstore.VarArrayType{
		Rank: rank,
		Of:   elem,
	}, nil
}

func (p *parser) parseStructOrUnion(prefixTok tokenType) (map[string]numstore.Type, error) {
	// Expect 'struct'
	_, err := p.expect(prefixTok)
	if err != nil {
		return nil, err
	}

	// Expect '{'
	_, err = p.expect(TOK_LEFT_CURLY)
	if err != nil {
		return nil, err
	}

	// Parse Fields
	fields := make(map[string]numstore.Type)
	for {
		// Name
		nameTok, err := p.expect(TOK_IDENTIFIER)
		if err != nil {
			return nil, err
		}

		// Type
		fieldType, err := p.parseType()
		if err != nil {
			return nil, err
		}

		// Check conflicts
		if _, ok := fields[nameTok.value]; ok {
			return nil, fmt.Errorf("Invalid struct type,"+
				"got conflicting fields: %v\n", nameTok.value)
		}
		fields[nameTok.value] = fieldType

		// Expect ',' or '}'
		if tok, _ := p.peekToken(); tok.ttype == TOK_COMMA {
			p.nextToken()
			continue
		} else if tok.ttype == TOK_RIGHT_CURLY {
			p.nextToken()
			break
		} else {
			return nil, expectedAnyButGot(tok.ttype,
				TOK_COMMA, TOK_RIGHT_CURLY)
		}
	}

	return fields, nil
}

func (p *parser) parseEnum() (numstore.EnumType, error) {
	_, err := p.expect(TOK_ENUM)
	if err != nil {
		return numstore.EnumType{}, err
	}

	// Expect '{'
	_, err = p.expect(TOK_LEFT_CURLY)
	if err != nil {
		return numstore.EnumType{}, err
	}

	opts := []string{}

	for {
		// Expect name
		if label, err := p.expect(TOK_IDENTIFIER); err != nil {
			return numstore.EnumType{}, err
		} else {
			opts = append(opts, label.value)
		}

		// Expect ',' or '}'
		if tok, _ := p.peekToken(); tok.ttype == TOK_COMMA {
			p.nextToken()
			continue
		} else if tok.ttype == TOK_RIGHT_CURLY {
			p.nextToken()
			break
		}
	}

	return numstore.EnumType{Options: opts}, nil
}
