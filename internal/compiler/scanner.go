package compiler

import (
	"github.com/lincketheo/numstore/internal/logging"
	"github.com/lincketheo/numstore/internal/numstore"
	"github.com/lincketheo/numstore/internal/utils"
)

type scanner struct {
	data    string
	start   int
	current int
	line    int
	col     int
	isError bool
}

func Scan(data string) []token {
	s := scannerCreate(data)
	ret := make([]token, 0, 20)

	for {
		tok := s.scanNextToken()
		if tok.ttype == TOK_NONE {
			continue
		}
		ret = append(ret, tok)
		if tok.ttype == TOK_EOF {
			break
		}
	}

	return ret
}

func scannerCreate(_data string) scanner {
	return scanner{
		data:    _data,
		start:   0,
		current: 0,
		line:    1,
		isError: false,
	}
}

func (s *scanner) skipWhitespace() {
	for {
		c := s.peekChar()
		switch c {
		case ' ', '\t', '\r':
			{
				_ = s.nextChar()
				break
			}
		case '\n':
			{
				s.line++
				s.col = 0
				_ = s.nextChar()
				break
			}
		default:
			return
		}
	}
}

func (s *scanner) scanNextToken() token {
	col, line := s.col, s.line

	return token{
		ttype: s.scanNextTokenType(),
		value: s.data[s.start:s.current],
		col:   col,
		line:  line,
	}
}

func (s *scanner) scanNextTokenType() tokenType {
	s.skipWhitespace()
	s.start = s.current

	if s.isEnd() {
		return TOK_EOF
	}

	c := s.nextChar()

	switch c {
	case '[':
		return TOK_LEFT_BRACKET
	case ']':
		return TOK_RIGHT_BRACKET
	case '{':
		return TOK_LEFT_CURLY
	case '}':
		return TOK_RIGHT_CURLY
	case '(':
		return TOK_LEFT_PAREN
	case ')':
		return TOK_RIGHT_PAREN
	case ',':
		return TOK_COMMA
	case ':':
		return TOK_COLON
	case '"':
		{
			s.parseString()
			return TOK_STRING
		}
	default:
		{
			if utils.IsDigit(c) || c == '-' || c == '+' {
				return s.parseNumber()
			} else if utils.IsAlpha(c) {
				return s.parseIdent()
			}

			s.compileError("Unexpected char: %c\n", c)
			return TOK_NONE
		}
	}
}

func (s *scanner) parseString() {
	for !s.isEnd() && s.peekChar() != '"' {
		s.nextChar()
	}

	if s.isEnd() {
		s.compileError("Unterminated string")
		return
	}

	s.nextChar()
}

func (s *scanner) parseNumber() tokenType {
	for utils.IsDigit(s.peekChar()) {
		_ = s.nextChar()
	}
	if s.peekChar() == '.' && utils.IsDigit(s.peek2Char()) {
		_ = s.nextChar()
		for utils.IsDigit(s.peekChar()) {
			_ = s.nextChar()
		}
		return TOK_FLOAT
	}
	return TOK_INTEGER
}

func (s *scanner) checkKeyword() tokenType {
	lexeme := s.data[s.start:s.current]

	if _, isPrim := numstore.PrimitiveTypeFromString(lexeme); isPrim {
		return TOK_PRIM
	}

	switch lexeme {
	case "create":
		return TOK_CREATE
	case "delete":
		return TOK_DELETE
	case "read":
		return TOK_READ
	case "write":
		return TOK_WRITE
	case "open":
		return TOK_OPEN
	case "close":
		return TOK_CLOSE
	case "take":
		return TOK_TAKE
	case "union":
		return TOK_UNION
	case "struct":
		return TOK_STRUCT
	default:
		return TOK_NONE
	}
}

func (s *scanner) parseIdent() tokenType {
	for utils.IsAlphaNum(s.peekChar()) {
		s.nextChar()
	}

	if ret := s.checkKeyword(); ret == TOK_NONE {
		return TOK_IDENTIFIER
	} else {
		return ret
	}
}

func (s *scanner) compileError(msg string, args ...any) {
	s.isError = true
	logging.Error(msg, args...)
}

func (s scanner) isEnd() bool {
	return s.current == len(s.data)
}

func (s scanner) nextIsEnd() bool {
	return s.current+1 == len(s.data)
}

func (s scanner) peekChar() byte {
	if s.isEnd() {
		return 0
	}
	return s.data[s.current]
}

func (s scanner) peek2Char() byte {
	if s.isEnd() || s.nextIsEnd() {
		return 0
	}
	return s.data[s.current+1]
}

func (s *scanner) nextChar() byte {
	utils.Assert(!s.isEnd())
	i := s.current
	s.current++
	s.col++
	return s.data[i]
}
