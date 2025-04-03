package compiler

import (
	"fmt"
)

type bytecode byte

const (
	BC_WRITE bytecode = iota
	BC_DELETE
	BC_READ
	BC_CREATE
  BC_ERROR
)

func byteToBytecode(b byte) (bytecode, bool) {
	if b <= byte(BC_CREATE) && b >= byte(BC_WRITE) {
		return bytecode(b), false
	}
	return 0, true
}

//////////////////////////// MAIN case Methods

// ////////////////////////// ERROR Wrappers
func (p *byteStack) nextStringOrError() (string, error) {
	db, ok := p.nextString()
	if !ok {
		return "", fmt.Errorf("Failed to get db name\n")
	}
	return db, nil
}
