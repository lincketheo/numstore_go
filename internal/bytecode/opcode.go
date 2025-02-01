package bytecode

import (
	"fmt"

	"github.com/lincketheo/ndbgo/internal/utils"
)

type opcode byte

const (
	OP_CREATE opcode = iota
	OP_CONNECT
	OP_TERM
)

func byteToOpcode(b byte) (opcode, bool) {
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

// ////////////////////////////////////// POP

func (c *ByteStack) popOpcode() (opcode, error) {
	if retb, err := c.popByte(); err != nil {
		return 0, utils.ErrorContext(err)
	} else if ret, ok := byteToOpcode(retb); !ok {
		return 0, fmt.Errorf("Popped byte: %d was not an expected OPCode",
			retb)
	} else {
		return ret, nil
	}
}

func (c *ByteStack) popOpcodeExpect(o opcode) error {
	if ret, err := c.popOpcode(); err != nil {
		return utils.ErrorContext(err)
	} else if ret != o {
		return fmt.Errorf("Expecting opcode: %d but got code: %d", o, ret)
	}
	return nil
}

func (c *ByteStack) popOpcodeIfMatches(o opcode) bool {
	if c.peekByteCheck(byte(o)) {
		_, err := c.popByte()
		utils.ASSERT(err == nil)
		return true
	}
	return false
}

// ////////////////////////////////////// PUSH

func (c *ByteStack) pushOpcode(code opcode) {
	c.pushByte(byte(code))
}
