package compiler

import "github.com/lincketheo/numstore/internal/numstore"

var readTokStarts = []tokenType{
	TOK_LEFT_BRACKET,
	TOK_LEFT_PAREN,
	TOK_IDENTIFIER,
}

func (p *parser) parseRFMT() (numstore.ReadFormat, bool) {

	tok, ok := p.peekToken()
	if !ok {
		p.parserError(readTokStarts...)
		return numstore.ReadFormat{}, false
	}

	var variables [][]numstore.ReadVariable

	switch tok.ttype {

	// 1. [a, (b, c), ....]
	case TOK_LEFT_BRACKET:
		{
			if variables, ok = p.parseRFMTList(); !ok {
				return numstore.ReadFormat{}, false
			}
		}

	// 2. (a, b, c...)
	case TOK_LEFT_PAREN:
		{
			if ret, ok := p.parseRFMTTuple(); !ok {
				return numstore.ReadFormat{}, false
			} else {
				variables = [][]numstore.ReadVariable{ret}
			}
		}

	// 3. a
	case TOK_IDENTIFIER:
		{
			if r, ok := p.parseRFMTIdent(); !ok {
				return numstore.ReadFormat{}, false
			} else {
				variables = [][]numstore.ReadVariable{{r}}
			}
		}

	// Invalid
	default:
		{
			p.parserError(readTokStarts...)
			return numstore.ReadFormat{}, false
		}
	}

	return numstore.ReadFormat{
		Variables: variables,
	}, true
}

func (p *parser) parseRFMTList() ([][]numstore.ReadVariable, bool) {
	// Expect '['
	if _, ok := p.expect(TOK_LEFT_BRACKET); !ok {
		p.parserError(TOK_LEFT_BRACKET)
		return nil, false
	}

	var list [][]numstore.ReadVariable

	for {
		tok, ok := p.peekToken()
		if !ok {
			p.parserError(TOK_LEFT_PAREN, TOK_IDENTIFIER)
			return nil, false
		}

		var next []numstore.ReadVariable

		switch tok.ttype {

		// 1. (a, b, c)
		case TOK_LEFT_PAREN:
			{
				if next, ok = p.parseRFMTTuple(); !ok {
					return nil, false
				}
			}

		// 2. a
		case TOK_IDENTIFIER:
			{
				if r, ok := p.parseRFMTIdent(); !ok {
					return nil, false
				} else {
					// Update
					next = []numstore.ReadVariable{r}
				}
			}

		// Invalid
		default:
			{
				p.parserError(TOK_LEFT_PAREN, TOK_IDENTIFIER)
				return nil, false
			}
		}

		// Append to list of variables
		list = append(list, next)

		// Now consume either a comma (and loop) or the closing bracket
		if _, ok := p.expect(TOK_COMMA); ok {
			continue
		} else if _, ok := p.expect(TOK_RIGHT_BRACKET); ok {
			return list, true
		} else {
			p.parserError(TOK_COMMA, TOK_RIGHT_BRACKET)
			return nil, false
		}
	}
}

func (p *parser) parseRFMTTuple() ([]numstore.ReadVariable, bool) {
	// Expect '('
	if _, ok := p.expect(TOK_LEFT_PAREN); !ok {
		p.parserError(TOK_LEFT_PAREN)
		return nil, false
	}

	var ret []numstore.ReadVariable

	for {
		tok, ok := p.peekToken()
		if !ok {
			p.parserError(TOK_IDENTIFIER)
			return nil, false
		}

		switch tok.ttype {

		case TOK_IDENTIFIER:
			{
				// Parse Identity and range
				if r, ok := p.parseRFMTIdent(); !ok {
					return nil, false
				} else {
					ret = append(ret, r)
				}
			}

		default:
			{
				p.parserError(TOK_IDENTIFIER)
				return nil, false
			}
		}

		if _, ok := p.expect(TOK_COMMA); ok {
			continue
		} else if _, ok := p.expect(TOK_RIGHT_PAREN); ok {
			return ret, true
		} else {
			p.parserError(TOK_COMMA, TOK_RIGHT_PAREN)
			return nil, false
		}
	}
}

func (p *parser) parseRFMTIdent() (numstore.ReadVariable, bool) {
	var tok token
	var ok bool
	if tok, ok = p.expect(TOK_IDENTIFIER); !ok {
		return numstore.ReadVariable{}, false
	}

	r := []numstore.DimRange{}
	if tok, ok = p.peekToken(); ok && tok.ttype == TOK_LEFT_BRACKET {
		if r, ok = p.parseRange(); !ok {
			return numstore.ReadVariable{}, false
		}
	}

	return numstore.ReadVariable{
		Vname: tok.value,
		Range: r,
	}, true
}
