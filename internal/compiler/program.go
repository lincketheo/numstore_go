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

func CreateProgram() Program {
	byteCodes := make([]byte, 0, 20)
	return Program{byteCodes, 0}
}

// ////////////////////////////////////// Utilities
func (c Program) head() []byte {
	return c.byteCodes[c.ip:]
}

func (c Program) len() int {
	return len(c.byteCodes)
}

func (c Program) hasNLeft(n int) bool {
	return len(c.byteCodes)-c.ip >= n
}

func (c Program) Done() bool {
	return c.ip == len(c.byteCodes)
}

// ////////////////////////////////////// POP
func (c *Program) popByte() (byte, error) {
	if c.Done() {
		return 0, fmt.Errorf("Poping byte no bytes left")
	}
	ret := c.byteCodes[c.ip]
	c.ip++
	return ret, nil
}

func (c *Program) popBytes(n int) ([]byte, error) {
	ret := c.head()[0:n]
	if len(ret) != n {
		return nil, fmt.Errorf(`Poping %d bytes left
      but had %d leftover bytes`, n, len(ret))
	}

	c.ip += n
	return ret, nil
}

func (c *Program) PopOpcode() (opcode, error) {
	retb, err := c.popByte()

	if err != nil {
		return 0, err
	}

	ret, ok := ByteToOpcode(retb)

	if !ok {
		return 0, fmt.Errorf("Poping opcode")
	}

	return ret, nil
}

func (c *Program) PopOpcodeExpect(o opcode) error {
	ret, err := c.PopOpcode()
	if err != nil {
		return err
	}
	if ret != o {
		return fmt.Errorf("Expecting opcode: %d but got code: %d", o, ret)
	}
	return nil
}

func (c *Program) PopString() (string, error) {
	byteLen, err := c.popByte()
	if err != nil {
		return "", err
	}

	ret, err := c.popBytes(int(byteLen))
	if err != nil {
		return "", err
	}

	return string(ret), nil
}

func (c *Program) peekByte() byte {
	if c.Done() {
		return c.byteCodes[len(c.byteCodes)-1]
	}
	return c.byteCodes[c.ip]
}

func (c *Program) PeekByteExpect(o byte) bool {
	ret := c.peekByte()
	return ret == o
}

func (c *Program) MatchByte(o byte) bool {
	if c.PeekByteExpect(o) {
		_, err := c.popByte()
		utils.ASSERT(err == nil)
		return true
	}
	return false
}

// ////////////////////////////////////// PUSH
func (c *Program) pushByte(b byte) {
	c.byteCodes = append(c.byteCodes, b)
}

func (c *Program) pushBytes(data []byte) {
	c.byteCodes = append(c.byteCodes, data...)
}

func (c *Program) pushOpcode(code opcode) {
	c.pushByte(byte(code))
}

func (c *Program) pushStr(data string) error {
	lengthByte := len(data)

	if !utils.CanIntBeByte(lengthByte) {
		return fmt.Errorf("String of length: %d is too long\n", lengthByte)
	}

	c.pushByte(byte(lengthByte))
	c.pushBytes([]byte(data))

	return nil
}
