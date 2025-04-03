package compiler

import (
	"github.com/lincketheo/numstore/internal/logging"
	"github.com/lincketheo/numstore/internal/utils"
)

type byteStack struct {
	data    []byte
	ip      int
	isError bool
}

func createByteStack() byteStack {
	data := make([]byte, 0, 20)
	return byteStack{
		data:    data,
		ip:      0,
		isError: false,
	}
}

// ////////////////////////////////////// Utilities
func (c byteStack) head() []byte {
	return c.data[c.ip:]
}

func (c byteStack) hasNLeft(n int) bool {
	return len(c.data)-c.ip >= n
}

func (c byteStack) isEnd() bool {
	return c.ip == len(c.data)
}

func (c *byteStack) compileError(fmt string, args ...any) {
	logging.Error(fmt, args...)
	c.pushByteCode(BC_ERROR)
}

// ////////////////////////////////////// ADVANCE
func (c *byteStack) nextByte() byte {
	utils.Assert(!c.isEnd())
	ret := c.data[c.ip]
	c.ip++
	return ret
}

func (c *byteStack) nextBytes(n int) ([]byte, bool) {
	ret := c.head()[0:n]

	if len(ret) != n {
		return nil, false
	}

	c.ip += n
	return ret, true
}

func (c *byteStack) nextByteCode() (bytecode, bool) {
	return byteToBytecode(c.nextByte())
}

func (c *byteStack) nextString() (string, bool) {
	if ret, ok := c.nextBytes(int(c.nextByte())); !ok {
		return "", ok
	} else {
		return string(ret), true
	}
}

func (c *byteStack) nextUint32Arr() ([]uint32, bool) {
	if ret, ok := c.nextBytes(int(c.nextByte()) * 4); !ok {
		return nil, ok
	} else {
		return utils.BytesToUInt32Arr(ret), true
	}
}

// ////////////////////////////////////// PUSH
func (c *byteStack) pushByte(b byte) {
	c.data = append(c.data, b)
}

func (c *byteStack) pushBytes(data []byte) {
	c.data = append(c.data, data...)
}

func (c *byteStack) pushByteCode(b bytecode) {
	c.pushByte(byte(b))
}

func (c *byteStack) pushStr(data string) {
	prefix := len(data)

	// String is too long
	if !utils.CanIntBeByte(prefix) {
		c.compileError("String of length: %d is too long\n", prefix)
		return
	}

	c.pushByte(byte(prefix))
	c.pushBytes([]byte(data))
}

func (c *byteStack) pushUint32Arr(s []uint32) {
	prefix := len(s)

	if !utils.CanIntBeByte(prefix) {
		c.compileError("Shape is too long: %v", s)
		return
	}

	sb := utils.UInt32ArrBytes(s)
	c.pushByte(byte(prefix))
	c.pushBytes(sb)
}

