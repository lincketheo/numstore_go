package compiler

import (
	"fmt"

	"github.com/lincketheo/ndbgo/internal/dtypes"
	"github.com/lincketheo/ndbgo/internal/utils"
)

// ////////////////////////////////// Entity
type entity byte

const (
	E_DB entity = iota
	E_REL
	E_VAR
)

func ByteToEntity(b byte) (entity, error) {
	switch b {
	case byte(E_DB):
		return E_DB, nil
	case byte(E_REL):
		return E_REL, nil
	case byte(E_VAR):
		return E_VAR, nil
	}
	return 0, fmt.Errorf("Expected entity, got byte: %d", b)
}

/*
REL
STRLEN
f
o
o
*/
func (c *Program) pushEntityWithName(e entity, name string) error {
	c.pushByte(byte(e))
	if err := c.pushStr(name); err != nil {
		return err
	}
	return nil
}

func (c *Program) PopEntity() (entity, error) {
	b, err := c.popByte()
	if err != nil {
		return 0, err
	}
	return ByteToEntity(b)
}

//////////////////////////////////// Variable Configuration Options

type varConfig byte

const (
	DTYPE varConfig = iota
	SHAPE
)

func ByteToVarConfig(b byte) (varConfig, error) {
	switch b {
	case byte(DTYPE):
		return DTYPE, nil
	case byte(SHAPE):
		return SHAPE, nil
	}
	return 0, fmt.Errorf("Expected variable config, got byte: %d", b)
}

/*
DTYPE CODE (byte)
[DTYPE] (byte)
*/
func (c *Program) pushDtype(d dtypes.Dtype) {
	c.pushByte(byte(DTYPE))
	c.pushByte(byte(d))
}

/*
SHAPE CODE (byte)
SHAPE LEN (byte)
SHAPE0 (u32)
SHAPE1 (u32)
*/
func (c *Program) pushShape(s []uint32) error {
	if !utils.CanIntBeByte(len(s)) {
		return fmt.Errorf("Shape is too long: %v", s)
	}

	c.pushByte(byte(SHAPE))
	sb := utils.UInt32ArrBytes(s)
	c.pushByte(byte(len(s)))
	c.pushBytes(sb)

	return nil
}

// Assumes already consumed DTYPE CODE
// (how else are you gonna get here?)
func (c *Program) PopDtype() (dtypes.Dtype, error) {
	b, err := c.popByte() // DTYPE
	if err != nil {
		return 0, err
	}
	dtype, ok := dtypes.ByteToDtype(b)
	if !ok {
		return 0, fmt.Errorf("Failed to convert byte: %d to dtype", b)
	}
	return dtype, nil
}

// Assumes already consumed SHAPE CODE
// (how else are you gonna get here?)
func (c *Program) PopShape() ([]uint32, error) {
	b, err := c.popByte() // LEN
	if err != nil {
		return nil, err
	}

	retb, err := c.popBytes(int(b) * 4) // DATA
	if err != nil {
		return nil, err
	}

	return utils.BytesToUInt32Arr(retb), nil
}
