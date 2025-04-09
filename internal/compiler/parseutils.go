package compiler

import "github.com/lincketheo/numstore/internal/logging"

func (t *parser) windNext() {
	for {
		if _, ok := t.expectNot(
			TOK_CREATE,
			TOK_DELETE,
			TOK_READ,
			TOK_WRITE,
			TOK_TAKE); !ok {
			return
		}
	}
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

func (p *parser) expectNot(ttypes ...tokenType) (token, bool) {
	tok, ok := p.peekToken()
	if !ok {
		return token{}, false
	}
	match := false
	for _, expected := range ttypes {
		if tok.ttype != expected {
			match = true
			break
		}
	}
	if !match {
		return token{}, false
	}
	return p.nextToken()
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
