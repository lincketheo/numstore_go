package compiler

import (
	"fmt"
	"github.com/lincketheo/ndbgo/internal/utils"
)

type opcode byte

const (
	OP_CREATE opcode = iota
	OP_CONNECT
	OP_TERM
	OP_EOF
)

func ByteToOpcode(b byte) (opcode, bool) {
	switch b {
	case byte(OP_CREATE):
		return OP_CREATE, true
	case byte(OP_TERM):
		return OP_TERM, true
	case byte(OP_CONNECT):
		return OP_CONNECT, true
	}
	return 0, false
}

type Program struct {
	byteCodes []byte
	ip        int
}

func (p Program) Bytes() []byte {
	return p.byteCodes
}


