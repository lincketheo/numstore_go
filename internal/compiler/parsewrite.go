package compiler

import "github.com/lincketheo/numstore/internal/numstore"

var writeTokStarts = []tokenType{
	TOK_LEFT_BRACKET,
	TOK_LEFT_PAREN,
	TOK_IDENTIFIER,
}

func (p *parser) parseWFMT() (numstore.WriteFormat, bool) {
  // TODO ToWrite

	tok, ok := p.peekToken()
	if !ok {
		p.parserError(writeTokStarts...)
		return numstore.WriteFormat{}, false
	}

	// Options:
	// 1. [a, b, ....]
	// 2. (a, b, c...]
	// 3. a

	var variables [][]string

	switch tok.ttype {

	// 1.
	case TOK_LEFT_BRACKET:
		{
			variables, ok = p.parseWFMTList()
		}

		// 2.
	case TOK_LEFT_PAREN:
		{
			if ret, ok := p.parseWFMTTuple(); !ok {
				return numstore.WriteFormat{}, false
			} else {
				variables = [][]string{ret}
			}
		}

		// 3.
	case TOK_IDENTIFIER:
		{
			variables = [][]string{{tok.value}}
			p.nextToken()
		}

		// Invalid
	default:
		{
			p.parserError(writeTokStarts...)
			return numstore.WriteFormat{}, false
		}
	}

	return numstore.WriteFormat{
    Variables: variables,
  }, true
}

func (p *parser) parseWFMTList() ([][]string, bool) {
	// Expect '['
	if _, ok := p.expect(TOK_LEFT_BRACKET); !ok {
		p.parserError(TOK_LEFT_BRACKET)
		return nil, false
	}

	var list [][]string

	for {
		tok, ok := p.peekToken()
		if !ok {
			p.parserError(TOK_LEFT_PAREN, TOK_IDENTIFIER)
			return nil, false
		}

		var next []string

		switch tok.ttype {

		// 1. (a, b, c)
		case TOK_LEFT_PAREN:
			{
				if next, ok = p.parseWFMTTuple(); !ok {
					return nil, false
				}
			}

		// 2. a
		case TOK_IDENTIFIER:
			{
				next = []string{tok.value}
				p.nextToken()
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

func (p *parser) parseWFMTTuple() ([]string, bool) {
	// Expect '('
	if _, ok := p.expect(TOK_LEFT_PAREN); !ok {
		p.parserError(TOK_LEFT_PAREN)
		return nil, false
	}

	var ret []string

	for {
		tok, ok := p.nextToken()
		if !ok {
			p.parserError(TOK_IDENTIFIER)
			return nil, false
		}

		switch tok.ttype {
		case TOK_IDENTIFIER:
			{
				ret = append(ret, tok.value)
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
