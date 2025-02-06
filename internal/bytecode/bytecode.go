package bytecode

import (
	"fmt"
)

type bytecode byte

const (
	BC_CREATE_DB      bytecode = iota // create(db)
	BC_CREATE_REL                     // create(rel)
	BC_CREATE_VAR                     // create(vari, dtype, shape)
	BC_CONNECT_DB                     // connect(db)
	BC_CONNECT_REL                    // connect(rel)
	BC_DISCONNECT_DB                  // disconnect(db)
	BC_DISCONNECT_REL                 // disconnect(db)
	BC_DELETE_DB                      // delete(db)
	BC_DELETE_REL                     // delete(rel)
	BC_DELETE_VAR                     // delete(vari)
	BC_ADD_READER                     // ...
	BC_REMOVE_READER                  // ...
	BC_ADD_WRITER                     // ...
	BC_REMOVE_WRITER                  // ...
	BC_WRITE                          // ...
	BC_READ                           // ...
	BC_ERROR
)

func byteToBytecode(b byte) (bytecode, bool) {
	if b <= byte(BC_READ) && b >= byte(BC_CREATE_DB) {
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
