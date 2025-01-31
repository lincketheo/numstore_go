package compiler

import (
	"github.com/lincketheo/ndbgo/internal/dtypes"
	"github.com/lincketheo/ndbgo/internal/logging"
)

// ////////////////////////////////////// CREATE
/*
CREATE
DB
3
F
O
O
;
*/
func (c *Program) CreateDb(db string) error {
	logging.Trace("Pushing CREATE DB: %s\n", db)
	c.pushOpcode(OP_CREATE)

	if err := c.pushEntityWithName(E_DB, db); err != nil {
		return err
	}

	c.pushOpcode(OP_TERM)
	return nil
}

/*
CREATE
REL
3
F
O
O
;
*/
func (c *Program) CreateRel(rel string) error {
	logging.Trace("Pushing CREATE REL: %s\n", rel)

	c.pushOpcode(OP_CREATE)

	if err := c.pushEntityWithName(E_REL, rel); err != nil {
		return err
	}

	c.pushOpcode(OP_TERM)

	return nil
}

/*
CREATE
VAR
1
a
DTYPE
U32
SHAPE
1       - shape len
3
;
*/
func (c *Program) CreateVar(
	vari string,
	dtype dtypes.Dtype,
	shape []uint32,
) error {
	logging.Trace("Pushing CREATE VAR %s %v %v", vari, dtype, shape)

	c.pushOpcode(OP_CREATE)
	if err := c.pushEntityWithName(E_VAR, vari); err != nil {
		return err
	}

	c.pushDtype(dtype)
	if err := c.pushShape(shape); err != nil {
		return err
	}

	c.pushOpcode(OP_TERM)

	return nil
}

// ////////////////////////////////////// CONNECT
/*
CONNECT
DB
3
F
O
O
;
*/
func (c *Program) ConnectDB(db string) error {
	logging.Trace("Pushing CONNECT DB: %s\n", db)

	c.pushOpcode(OP_CONNECT)

	if err := c.pushEntityWithName(E_DB, db); err != nil {
		return err
	}

	c.pushOpcode(OP_TERM)
	return nil
}

/*
CONNECT
REL
3
F
O
O
;
*/
func (c *Program) ConnectRel(rel string) error {
	logging.Trace("Pushing CONNECT REL: %s\n", rel)

	c.pushOpcode(OP_CONNECT)

	if err := c.pushEntityWithName(E_REL, rel); err != nil {
		return err
	}

	c.pushOpcode(OP_TERM)

	return nil
}

/*
CONNECT
VAR
3
F
O
O
;
*/
func (c *Program) ConnectVar(vari string) error {
	logging.Trace("Pushing CONNECT VAR. %s", vari)

	c.pushOpcode(OP_CONNECT)
	if err := c.pushEntityWithName(E_VAR, vari); err != nil {
		return err
	}

	c.pushOpcode(OP_TERM)

	return nil
}
