package compiler

import (
	"bufio"
	"os"
	"strconv"

	"github.com/lincketheo/numstore/internal/logging"
	"github.com/lincketheo/numstore/internal/nserror"
	"github.com/lincketheo/numstore/internal/utils"
)

type tokenType int

type token struct {
	ttype tokenType
	value string
	line  int
}

func WriteTokensToFileClean(fname string, tokens []token) error {
	file, err := os.Create(fname)
	if err != nil {
		return nserror.ErrorStack(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			logging.Warn("Failed to close file: %s. %v\n", fname, err)
		}
	}()

	writer := bufio.NewWriter(file)

	for _, t := range tokens {
		if _, err := writer.WriteString(
			strconv.Itoa(t.line) +
				" " + t.ttype.String() +
				" " + t.value + "\n",
		); err != nil {
			return nserror.ErrorStack(err)
		}
	}

	if err := writer.Flush(); err != nil {
		logging.Warn("Failed to flush file: %s. %v\n", fname, err)
	}

	return nil
}

const (
	// A none type used for code flow
	// not a real token
	TOK_NONE tokenType = iota

	// OP Codes
	TOK_CREATE
	TOK_DELETE
	TOK_READ
	TOK_WRITE

	// Structure
	TOK_LEFT_BRACKET
	TOK_RIGHT_BRACKET
	TOK_LEFT_PAREN
	TOK_RIGHT_PAREN
	TOK_LEFT_CURLY
	TOK_RIGHT_CURLY
	TOK_COMMA
	TOK_COLON

	// Numbers
	TOK_INTEGER
	TOK_FLOAT

	// Variables
	TOK_IDENTIFIER
	TOK_STRING

	// DType
	TOK_DTYPE

	// To indicate done parsing
	TOK_EOF
)

func (t tokenType) String() string {

	switch t {
	case TOK_NONE:
		return "TOK_NONE"

		// OP Codes
	case TOK_CREATE:
		return "TOK_CREATE"
	case TOK_DELETE:
		return "TOK_DELETE"
	case TOK_READ:
		return "TOK_READ"
	case TOK_WRITE:
		return "TOK_WRITE"

		// Structure
	case TOK_LEFT_BRACKET:
		return "TOK_LEFT_BRACKET"
	case TOK_RIGHT_BRACKET:
		return "TOK_RIGHT_BRACKET"
	case TOK_LEFT_PAREN:
		return "TOK_LEFT_PAREN"
	case TOK_RIGHT_PAREN:
		return "TOK_RIGHT_PAREN"
	case TOK_LEFT_CURLY:
		return "TOK_LEFT_CURLY"
	case TOK_RIGHT_CURLY:
		return "TOK_RIGHT_CURLY"
	case TOK_COMMA:
		return "TOK_COMMA"
	case TOK_COLON:
		return "TOK_COLON"

	// Numbers
	case TOK_INTEGER:
		return "TOK_INTEGER"
	case TOK_FLOAT:
		return "TOK_FLOAT"

	// Variables
	case TOK_IDENTIFIER:
		return "TOK_IDENTIFIER"
	case TOK_STRING:
		return "TOK_STRING"

		// DType
	case TOK_DTYPE:
		return "TOK_DType"

	// To indicate done parsing
	case TOK_EOF:
		return "TOK_EOF"
	}
	panic("Unknown tokenType")
}

func assertTokens(toks []token, cur int) {
	utils.Assert(cur < len(toks))
	if cur == len(toks)-1 {
		utils.Assert(toks[cur].ttype == TOK_EOF)
	} else {
		utils.Assert(toks[cur].ttype != TOK_EOF)
	}
}
