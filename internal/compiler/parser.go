package compiler

import (
	"fmt"
)

//////////////////////////////// PARSER
type parser struct {
	tokens  []token
	cur     int
	isError bool
}

func (t parser) isEnd() bool {
	return t.cur >= len(t.tokens)
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

func (p *parser) expect(t tokenType) (token, error) {
	tok, ok := p.peekToken()
	if !ok {
		return token{}, fmt.Errorf("Expected token: %v, "+
			"got end of token stream", t)
	}

	if tok.ttype != t {
		return token{}, fmt.Errorf("Expected token: %v, "+
			"got token: %v", t, tok.ttype)
	}

	// advance
	_, ok = p.nextToken()

	return tok, nil
}

//////////////////////////////// PARSE
func RunTokens(tokens []token) error {
	if len(tokens) == 0 {
		return nil
	}

	runner := parser{
		tokens: tokens,
		cur:    0,
	}

	for t, ok := runner.nextToken(); ok; t, ok = runner.nextToken() {
		switch t.ttype {
		case TOK_CREATE:
			{
				if err := runner.handleTokCreate(); err != nil {
					return err
				}
			}
		case TOK_DELETE:
			{
				if err := runner.handleTokDelete(); err != nil {
					return err
				}
			}
		case TOK_READ:
			{
				if err := runner.handleTokRead(); err != nil {
					return err
				}
			}
		case TOK_WRITE:
			{
				if err := runner.handleTokWrite(); err != nil {
					return err
				}
			}
		case TOK_OPEN:
			{
				if err := runner.handleTokOpen(); err != nil {
					return err
				}
			}
		case TOK_CLOSE:
			{
				if err := runner.handleTokClose(); err != nil {
					return err
				}
			}
		case TOK_TAKE:
			{
				if err := runner.handleTokTake(); err != nil {
					return err
				}
			}
		case TOK_EOF:
			{
				return nil
			}
		default:
			return fmt.Errorf("Invalid token: %v", t)
		}

		if _, err := runner.expect(TOK_SEMICOLON); err != nil {
			return err
		}
	}

	panic("Unreachable")
}

// DONE
func (t *parser) handleTokCreate() error {
	v, ok := t.nextToken()
	if !ok {
		return expectedAfter(TOK_IDENTIFIER, TOK_CREATE)
	}

	if v.ttype != TOK_IDENTIFIER {
		return invalidAfterExpected(v.ttype, TOK_CREATE, TOK_IDENTIFIER)
	}

	if nstype, err := t.parseType(); err != nil {
		return err
	} else {
		fmt.Printf("CREATING: %v\n", nstype)
		return nil
	}
}

func (t *parser) handleTokDelete() error {
	v, ok := t.nextToken()
	if !ok {
		return expectedAfter(TOK_IDENTIFIER, TOK_DELETE)
	}

	if v.ttype != TOK_IDENTIFIER {
		return invalidAfterExpected(v.ttype, TOK_DELETE, TOK_IDENTIFIER)
	}

	fmt.Printf("DELETING: %v\n", v.value)
	return nil
}

func (t *parser) handleTokRead() error {
	v, ok := t.nextToken()
	if !ok {
		return expectedAfter(TOK_IDENTIFIER, TOK_READ)
	}

	if v.ttype != TOK_IDENTIFIER {
		return invalidAfterExpected(v.ttype, TOK_READ, TOK_IDENTIFIER)
	}

	fmt.Printf("READING: %v\n", v.value)
	return nil
}

// DONE
func (t *parser) handleTokWrite() error {
	if wfmt, err := t.parseWriteFormat(); err != nil {
		return err
	} else {
		fmt.Printf("Writing: %v\n", wfmt)
		return nil
	}
}

func (t *parser) handleTokOpen() error {
	v, ok := t.nextToken()
	if !ok {
		return expectedAfter(TOK_IDENTIFIER, TOK_OPEN)
	}

	if v.ttype != TOK_IDENTIFIER {
		return invalidAfterExpected(v.ttype, TOK_OPEN, TOK_IDENTIFIER)
	}

	fmt.Printf("OPENING: %v\n", v.value)
	return nil
}

func (t *parser) handleTokClose() error {
	v, ok := t.nextToken()
	if !ok {
		return expectedAfter(TOK_IDENTIFIER, TOK_CLOSE)
	}

	if v.ttype != TOK_IDENTIFIER {
		return invalidAfterExpected(v.ttype, TOK_CLOSE, TOK_IDENTIFIER)
	}

	fmt.Printf("CLOSEING: %v\n", v.value)
	return nil
}

func (t *parser) handleTokTake() error {
	v, ok := t.nextToken()
	if !ok {
		return expectedAfter(TOK_IDENTIFIER, TOK_TAKE)
	}

	if v.ttype != TOK_IDENTIFIER {
		return invalidAfterExpected(v.ttype, TOK_TAKE, TOK_IDENTIFIER)
	}

	fmt.Printf("TAKEING: %v\n", v.value)
	return nil
}
