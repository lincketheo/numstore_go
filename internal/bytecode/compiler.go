package bytecode

import (
	"fmt"

	"github.com/lincketheo/ndbgo/internal/utils"
)

type parser struct {
	program  byteStack
	tokens   []token
	current  token
	previous token
	isError  bool
}

func parserCreate(_tokens []token) parser {
	utils.ASSERT(len(_tokens) != 0)

	return parser{
		program:  createByteStack(),
		tokens:   _tokens,
		current:  _tokens[0],
		previous: _tokens[0],
		isError:  false,
	}
}

func Compile(data string) {
	s := scannerCreate(data)
	line := -1

	for {
		tok := s.scanNextToken()
		if tok.ttype == TOK_NONE {
			continue
		}

		if tok.line != line {
			fmt.Printf("%4d ", tok.line)
			line = tok.line
		} else {
			fmt.Printf("   | ")
		}

		fmt.Printf("%s '%s'\n", tok.ttype, tok.value)

		if tok.ttype == TOK_EOF {
			break
		}
	}
}
