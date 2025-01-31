package interpreter

import (
	"errors"
	"fmt"

	"github.com/lincketheo/ndbgo/internal/compiler"
	"github.com/lincketheo/ndbgo/internal/logging"
	"github.com/lincketheo/ndbgo/internal/utils"
)

func RunNext(n *NDB, p *compiler.Program) error {
	// Pop the opcode
	b, err := p.PopOpcode()
	if err != nil {
		return err
	}

	// Execute the stuff
	switch b {

	case compiler.OP_CREATE:
		if err = interpretCreate(n, p); err != nil {
			return err
		}
		break

	case compiler.OP_CONNECT:
		if err = interpretConnect(n, p); err != nil {
			return err
		}
		break

	case compiler.OP_EOF:
		return nil // Nothing to do

	case compiler.OP_TERM:
		return fmt.Errorf("Unexpected op code OP_TERM at top level")

	default:
		utils.UNREACHABLE()
	}

	return nil
}

func RunAll(n *NDB, p *compiler.Program) error {

	for !p.Done() {
		if err := RunNext(n, p); err != nil {
			return err
		}
	}

	return nil
}

// /////////////////////////////////////// CREATE
func interpretCreate(n *NDB, p *compiler.Program) error {

	// Expect entity name
	rel, err := p.PopEntity()
	if err != nil {
		return err
	}

	switch rel {

	case compiler.E_DB:
		return interpretCreateDB(n, p)

	case compiler.E_REL:
		return interpretCreateREL(n, p)

	case compiler.E_VAR:
		return interpretCreateVAR(n, p)
	}

	utils.UNREACHABLE()
	return nil
}

func interpretCreateDB(n *NDB, p *compiler.Program) error {
	// Pop the database name off
	db, err := p.PopString()
	if err != nil {
		return err
	}

	// Expect terminal
	if err := p.PopOpcodeExpect(compiler.OP_TERM); err != nil {
		return err
	}

	// Create the database
	logging.Debug("Creating Database: %s\n", db)
	if err = (*n).CreateDB(db); err != nil {
		return err
	}

	return nil
}

func interpretCreateREL(n *NDB, p *compiler.Program) error {
	// Pop the relation name
	rel, err := p.PopString()
	if err != nil {
		return err
	}

	// Expect terminal
	if err := p.PopOpcodeExpect(compiler.OP_TERM); err != nil {
		return err
	}

	// Create the relation
	logging.Debug("Creating Relation: %s\n", rel)
	if err = (*n).CreateRel(rel); err != nil {
		return err
	}

	return nil
}

func interpretCreateVAR(n *NDB, p *compiler.Program) error {
	vari, err := p.PopString()
	if err != nil {
		return err
	}

	// Initialize
	args := CreateVarArgs{
		Vari:  vari,
		Dtype: 0,
		Shape: nil,
	}
	var dtypeSet = false
	var shapeSet = false

	for {
		if dtypeSet && shapeSet {
			break
		}

		if !dtypeSet && p.MatchByte(byte(compiler.DTYPE)) {
			dtype, err := p.PopDtype()
			if err != nil {
				return err
			}
			args.Dtype = dtype
			dtypeSet = true
		} else if !shapeSet && p.MatchByte(byte(compiler.SHAPE)) {
			shape, err := p.PopShape()
			if err != nil {
				return err
			}
			args.Shape = shape
			shapeSet = true
		} else {
			return errors.New(`Expected a valid configuration
        option for create variable argument`)
		}
	}

	// Expect terminal
	if err := p.PopOpcodeExpect(compiler.OP_TERM); err != nil {
		return err
	}

	// Create the variable
	logging.Debug("Creating Variable: %v\n", args)
	if err = (*n).CreateVar(args); err != nil {
		return err
	}

	return nil
}

///////////////////////////////////////// CONNECT

func interpretConnect(n *NDB, p *compiler.Program) error {

	// Expect entity name
	e, err := p.PopEntity()
	if err != nil {
		return err
	}

	switch e {

	case compiler.E_DB:
		return interpretConnectDB(n, p)

	case compiler.E_REL:
		return interpretConnectREL(n, p)

	case compiler.E_VAR:
		return interpretConnectVAR(n, p)
	}

	utils.UNREACHABLE()
	return nil
}

func interpretConnectDB(n *NDB, p *compiler.Program) error {
	// Pop the database name off
	db, err := p.PopString()
	if err != nil {
		return err
	}

	// Expect terminal
	if err := p.PopOpcodeExpect(compiler.OP_TERM); err != nil {
		return err
	}

	// Connect the database
	logging.Debug("Connecting to Database: %s\n", db)
	if err = (*n).ConnectDB(db); err != nil {
		return err
	}

	return nil
}

func interpretConnectREL(n *NDB, p *compiler.Program) error {
	// Pop the relation name
	rel, err := p.PopString()
	if err != nil {
		return err
	}

	// Expect terminal
	if err := p.PopOpcodeExpect(compiler.OP_TERM); err != nil {
		return err
	}

	// Connect the relation
	logging.Debug("Connecting to Relation: %s\n", rel)
	if err = (*n).ConnectRel(rel); err != nil {
		return err
	}

	return nil
}

func interpretConnectVAR(n *NDB, p *compiler.Program) error {
	vari, err := p.PopString()
	if err != nil {
		return err
	}

	// Expect terminal
	if err := p.PopOpcodeExpect(compiler.OP_TERM); err != nil {
		return err
	}

	// Connect the relation
	logging.Debug("Connecting to Variable: %s\n", vari)
	if err = (*n).ConnectVar(vari); err != nil {
		return err
	}

	return nil
}
