package compiler

import (
	"github.com/lincketheo/ndbgo/internal/dtypes"
	"github.com/lincketheo/ndbgo/internal/logging"
)

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
