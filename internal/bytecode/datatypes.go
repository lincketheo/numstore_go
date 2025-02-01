package bytecode

import (
	"fmt"

	"github.com/lincketheo/ndbgo/internal/dtypes"
	"github.com/lincketheo/ndbgo/internal/utils"
)

type dataTypePrefix byte

const (
	STRING dataTypePrefix = iota
	UINT32ARR
	DTYPE
)

// ////////////////////////////////////// POP

func (c *ByteStack) popString() (string, error) {
	if err := c.popByteExpect(byte(STRING)); err != nil {
		return "", utils.ErrorContext(err)
	} else if byteLen, err := c.popByte(); err != nil {
		return "", utils.ErrorContext(err)
	} else if ret, err := c.popBytes(int(byteLen)); err != nil {
		return "", utils.ErrorContext(err)
	} else {
		return string(ret), nil
	}
}

func (c *ByteStack) popUint32Arr() ([]uint32, error) {
	if err := c.popByteExpect(byte(UINT32ARR)); err != nil {
		return nil, utils.ErrorContext(err)
	} else if b, err := c.popByte(); err != nil { // LEN
		return nil, utils.ErrorContext(err)
	} else if retb, err := c.popBytes(int(b) * 4); err != nil { // DATA
		return nil, utils.ErrorContext(err)
	} else {
		return utils.BytesToUInt32Arr(retb), nil
	}
}

func (c *ByteStack) popDtype() (dtypes.Dtype, error) {
	if err := c.popByteExpect(byte(DTYPE)); err != nil {
		return 0, utils.ErrorContext(err)
	} else if retb, err := c.popByte(); err != nil {
		return 0, utils.ErrorContext(err)
	} else {
		return dtypes.ByteToDtype(retb)
	}
}

// ////////////////////////////////////// PUSH
func (c *ByteStack) pushStr(data string) error {
	prefix := len(data)

	if !utils.CanIntBeByte(prefix) {
		return fmt.Errorf("String of length: %d is too long\n", prefix)
	}

	c.pushByte(byte(STRING))
	c.pushByte(byte(prefix))
	c.pushBytes([]byte(data))

	return nil
}

func (c *ByteStack) pushUint32Arr(s []uint32) error {
	prefix := len(s)

	if !utils.CanIntBeByte(prefix) {
		return fmt.Errorf("Shape is too long: %v", s)
	}

	c.pushByte(byte(UINT32ARR))
	sb := utils.UInt32ArrBytes(s)
	c.pushByte(byte(prefix))
	c.pushBytes(sb)

	return nil
}

func (c *ByteStack) pushDtype(d dtypes.Dtype) {
	c.pushByte(byte(DTYPE))
	c.pushByte(byte(d))
}
