package compiler

import (
	"strconv"

	"github.com/lincketheo/numstore/internal/core"
	"github.com/lincketheo/numstore/internal/logging"
	"github.com/lincketheo/numstore/internal/utils"
)

type parser struct {
	isError bool
	ret     byteStack
	data    []token
	cur     int
}

func Parse(data []token) error {
	p := parserCreate(data)

	for !p.isEnd() {
		p.parseNext()
	}

	return nil
}

func (p *parser) parseNext() {
	t, _ := p.nextToken()
	switch t.ttype {
	case TOK_DELETE:
		{
			p.parseDelete()
		}
	case TOK_CREATE:
		{
			p.parseCreate()
		}
	case TOK_READ:
		{
			p.parseRead()
		}
	case TOK_WRITE:
		{
			p.parseWrite()
		}
	default:
		panic("Invalid token")
	}
}

func (p *parser) parseDelete() {
	p.ret.pushByteCode(BC_DELETE)
	// TODO
}

func (p *parser) parseCreate() bool {
	p.ret.pushByteCode(BC_CREATE)

	vname, ok := p.parseIdentifier()
	if !ok {
		return false
	}

	dtype, ok := p.parseDtype()
	if !ok {
		return false
	}

	shape, ok := p.parseShape()
	if !ok {
		return false
	}

	logging.Debug("Create: %s %v %d", vname, dtype, shape)
	return true
}

func (p *parser) parseIdentifier() (string, bool) {
	t, end := p.nextToken()
	if end {
		p.compileError("Unexpected end of input after create operation code")
		return "", false
	}
	if t.ttype != TOK_IDENTIFIER {
		p.compileError("Expected IDENTIFIER after create operation code")
		return "", false
	}
	return t.value, true
}

func (p *parser) parseDtype() (core.Dtype, bool) {
	t, end := p.nextToken()
	if end {
		p.compileError("Unexpected end of input after create IDENTIFIER")
		return 0, false
	}
	if t.ttype != TOK_DTYPE {
		p.compileError("Expected DTYPE after create IDENTIFIER")
		return 0, false
	}

	dtype, ok := core.DtypeFromString(t.value)
	utils.Assert(ok) // Already checked this in scanner - duplicate checks
	return dtype, true
}

func (p *parser) parseShape() ([]uint32, bool) {
	var shape []uint32
	for {
		t := p.peekToken()
		if t.ttype == TOK_INTEGER {
			i, err := strconv.ParseInt(t.value, 10, 32)

			// Parse check
			if err != nil {
				p.compileError("Expected shape dimension > 0. Got: %s", t.value)
				return nil, false
			}

			// Range check
			if i < 0 {
				p.compileError("Expected shape dimension > 0. Got: %s", t.value)
				return nil, false
			}

			shape = append(shape, uint32(i))
			_, end := p.nextToken()
			if end {
				p.compileError("Unexpected end of input while parsing shape dimensions")
				return nil, false
			}
		} else {
			break
		}
	}
	return shape, true
}

func (p *parser) parseRead() {
	p.ret.pushByteCode(BC_READ)
	// TODO
}

func (p *parser) parseWrite() {
	p.ret.pushByteCode(BC_WRITE)
	// TODO
}

func (p parser) isEnd() bool {
	assertTokens(p.data, p.cur)
	return p.data[p.cur].ttype == TOK_EOF
}

func (p *parser) peekToken() token {
	return p.data[p.cur]
}

func (p *parser) nextToken() (token, bool) {
  if p.isEnd() {
    return token{}, false
  }
	ret := p.data[p.cur]
	p.cur++
	return ret, true
}

func (p *parser) compileError(msg string, args ...any) {
	p.isError = true
	p.ret.pushByteCode(BC_ERROR)
	logging.Error(msg, args...)
}

func parserCreate(_data []token) parser {
	return parser{
		data: _data,
		cur:  0,
	}
}
