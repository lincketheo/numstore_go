package bytecode

import (
	"fmt"

	"github.com/lincketheo/ndbgo/internal/dtypes"
)

// ////////////////////////////////// Entity
type entity byte
type varConfig byte

const (
	E_DB entity = iota
	E_REL
	E_VAR
)

const (
	DTYPECONFIG varConfig = iota
	SHAPECONFIG
)

func byteToEntity(b byte) (entity, error) {
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

func byteToVarConfig(b byte) (varConfig, error) {
	switch b {
	case byte(DTYPECONFIG):
		return DTYPECONFIG, nil
	case byte(SHAPECONFIG):
		return SHAPECONFIG, nil
	}
	return 0, fmt.Errorf("Expected variable config, got byte: %d", b)
}

// ////////////////////////////////////// POP

func (c *ByteStack) popEntity() (entity, error) {
	if b, err := c.popByte(); err != nil {
		return 0, err
	} else {
		return byteToEntity(b)
	}
}

// ////////////////////////////////////// PUSH

/*
ENTITY Type Code
NAME (string)
*/
func (c *ByteStack) pushEntityWithName(e entity, name string) error {
	c.pushByte(byte(e))
	if err := c.pushStr(name); err != nil {
		return err
	}
	return nil
}

/*
DTYPE Config Code
DTYPE
*/
func (c *ByteStack) pushVarDtypeConfig(d dtypes.Dtype) {
	c.pushByte(byte(DTYPECONFIG))
	c.pushDtype(d)
}

/*
SHAPE Config Code
SHAPE ([]uint32)
*/
func (c *ByteStack) pushVarShapeConfig(shape []uint32) error {
	c.pushByte(byte(SHAPECONFIG))
	if err := c.pushUint32Arr(shape); err != nil {
		return err
	}
	return nil
}

//////////////////////////////////// POP
