package bytecode

import (
	"fmt"
)

type ByteStack struct {
	data []byte
	ip   int
}

// ////////////////////////////////////// Utilities
func createByteStack() ByteStack {
	data := make([]byte, 0, 20)
	return ByteStack{data, 0}
}

func (c ByteStack) head() []byte {
	return c.data[c.ip:]
}

func (c ByteStack) hasNLeft(n int) bool {
	return len(c.data)-c.ip >= n
}

func (c ByteStack) empty() bool {
	return c.ip == len(c.data)
}

// ////////////////////////////////////// POP
func (c *ByteStack) popByte() (byte, error) {
	if c.empty() {
		return 0, fmt.Errorf("Poping byte but there are no bytes left")
	}
	ret := c.data[c.ip]
	c.ip++
	return ret, nil
}

func (c *ByteStack) popByteExpect(b byte) error {
	if ret, err := c.popByte(); err != nil {
		return err
	} else if ret != b {
		return fmt.Errorf(`Poping byte expected byte:
      %v but popped byte was: %v`, b, ret)
	} else {
		return nil
	}
}

func (c *ByteStack) popBytes(n int) ([]byte, error) {
	ret := c.head()[0:n]
	if len(ret) != n {
		return nil, fmt.Errorf(`Poping %d bytes left
      but had %d leftover bytes`, n, len(ret))
	}

	c.ip += n
	return ret, nil
}

// ////////////////////////////////////// PEEK
func (c *ByteStack) peekByte() (byte, bool) {
	if c.empty() {
		return 0, false
	}
	return c.data[c.ip], true
}

func (c *ByteStack) peekByteExpect(o byte) error {
	if ret, ok := c.peekByte(); !ok {
		return fmt.Errorf(`Peeking byte expected byte:
      %v but stack is empty`, o)
	} else if ret != o {
		return fmt.Errorf(`Peeking byte expected byte:
      %v but peeked byte was: %v`, o, ret)
	} else {
		return nil
	}
}

func (c *ByteStack) peekByteCheck(o byte) bool {
	if ret, ok := c.peekByte(); !ok {
		return false
	} else {
		return ret == o
	}
}

func (c *ByteStack) peekBytes(n int) ([]byte, bool) {
	if c.hasNLeft(n) {
		return c.data[c.ip : c.ip+n], true
	}
	return nil, false
}

// ////////////////////////////////////// PUSH
func (c *ByteStack) pushByte(b byte) {
	c.data = append(c.data, b)
}

func (c *ByteStack) pushBytes(data []byte) {
	c.data = append(c.data, data...)
}
