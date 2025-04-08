package compiler

import "github.com/lincketheo/numstore/internal/numstore"

func (p *parser) parseWriteFormat() (numstore.WriteFormat, error) {

	tok, ok := p.peekToken()
	if !ok {
		return numstore.WriteFormat{},
			expectedStringButEarlyTermination("WFMT")
	}

	// Options:
	// 1. [a, b, ....]
	// 2. (a, b, c...]
	// 3. a

	var list [][]string
	var err error

	switch tok.ttype {

	// 1.
	case TOK_LEFT_BRACKET:
		{
			list, err = p.parseWriteFormatList()
		}

		// 2.
	case TOK_LEFT_PAREN:
		{
			var ret []string
			ret, err = p.parseWriteFormatTuple()
			list = [][]string{ret}
		}

		// 3.
	case TOK_IDENTIFIER:
		{
			list = [][]string{{tok.value}}
      p.nextToken()
		}

		// Invalid
	default:
		{
			return numstore.WriteFormat{}, invalidTokenExpectedAny(tok.ttype,
				TOK_LEFT_BRACKET, TOK_LEFT_PAREN, TOK_IDENTIFIER)
		}
	}

	if err != nil {
		return numstore.WriteFormat{}, err
	}
	return numstore.WriteFormat{Variables: list}, nil
}

func (p *parser) parseWriteFormatList() ([][]string, error) {
	// Expect '['
	_, err := p.expect(TOK_LEFT_BRACKET)
	if err != nil {
		return nil, err
	}

	var list [][]string

	for {
		tok, ok := p.peekToken()
		if !ok {
			return nil, expectedAnyTokenButEarlyTermination(
				TOK_LEFT_PAREN, TOK_IDENTIFIER)
		}

		var next []string

		// Options:
		// 1. (a, b, c)
		// 2. a
		switch tok.ttype {
		case TOK_LEFT_PAREN:
			{
				if next, err = p.parseWriteFormatTuple(); err != nil {
					return nil, err
				}
			}
		case TOK_IDENTIFIER:
			{
				next = []string{tok.value}
				p.nextToken()
			}
		default:
			{
				return nil, invalidTokenAfterTokenExpected(
					tok.ttype, TOK_LEFT_PAREN, TOK_IDENTIFIER)
			}
		}

		// Append to list
		list = append(list, next)

		// Now consume either a comma (and loop) or the closing bracket
		tok, _ = p.peekToken()
		switch tok.ttype {
		case TOK_COMMA:
			{
				p.nextToken()
				continue
			}
		case TOK_RIGHT_BRACKET:
			{
				p.nextToken()
				return list, nil
			}
		default:
			{
				return nil, invalidTokenAfterTokenExpected(tok.ttype,
					TOK_COMMA, TOK_RIGHT_BRACKET)
			}
		}
	}
}

func (p *parser) parseWriteFormatTuple() ([]string, error) {
	// Expect '('
	if _, err := p.expect(TOK_LEFT_PAREN); err != nil {
		return nil, err
	}

	var ret []string
	for {
		tok, ok := p.nextToken()
		if !ok {
			return nil, expectedAnyTokenButEarlyTermination(TOK_IDENTIFIER)
		}

		switch tok.ttype {
		case TOK_IDENTIFIER:
			{
				ret = append(ret, tok.value)
			}
		default:
			{
				return nil, expectedTokenButGot(TOK_IDENTIFIER, tok.ttype)
			}
		}

		// Next must be comma or ')'
		peek, _ := p.peekToken()

		switch peek.ttype {
		case TOK_COMMA:
			{
				p.nextToken()
				continue
			}
		case TOK_RIGHT_PAREN:
			{
				p.nextToken()
				return ret, nil
			}
		default:
			{
				return nil, invalidTokenAfterTokenExpected(peek.ttype,
					TOK_COMMA, TOK_RIGHT_PAREN)
			}
		}
	}
}
