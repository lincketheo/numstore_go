package bytecode

import (
	"fmt"

	"github.com/lincketheo/ndbgo/internal/dtypes"
	"github.com/lincketheo/ndbgo/internal/usecases"
	"github.com/lincketheo/ndbgo/internal/utils"
)

// ////////////////////////////////// Entity
type entity byte

const (
	E_DB entity = iota
	E_REL
	E_VAR
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

type varConfigCode byte

const (
	DTYPECONFIG varConfigCode = iota
	SHAPECONFIG
)

func byteToVarConfig(b byte) (varConfigCode, error) {
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
		return 0, utils.ErrorContext(err)
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
		return utils.ErrorContext(err)
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
		return utils.ErrorContext(err)
	}
	return nil
}

func (c *ByteStack) pushVarConfig(config usecases.VarConfig) error {
	if err := c.pushVarShapeConfig(config.Shape); err != nil {
		return utils.ErrorContext(err)
	} else {
		c.pushVarDtypeConfig(config.Dtype)
		return nil
	}
}

// ////////////////////////////////// POP
func (c *ByteStack) popVarDtypeConfig() (dtypes.Dtype, error) {
	if err := c.popByteExpect(byte(DTYPECONFIG)); err != nil {
		return 0, utils.ErrorContext(err)
	}
	return c.popDtype()
}

func (c *ByteStack) popVarShapeConfig() ([]uint32, error) {
	if err := c.popByteExpect(byte(SHAPECONFIG)); err != nil {
		return nil, utils.ErrorContext(err)
	}

	return c.popUint32Arr()
}

func (c *ByteStack) popVarConfig() (usecases.VarConfig, error) {
	ret := usecases.VarConfig{Dtype: 0, Shape: nil}

	if dtype, err := c.popVarDtypeConfig(); err != nil {
		return ret, utils.ErrorContext(err)
	} else if shape, err := c.popVarShapeConfig(); err != nil {
		return ret, utils.ErrorContext(err)
	} else {
		return usecases.VarConfig{
			Dtype: dtype,
			Shape: shape,
		}, nil
	}
}
