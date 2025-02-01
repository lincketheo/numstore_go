package bytecode

import (
	"fmt"

	"github.com/lincketheo/ndbgo/internal/usecases"
	"github.com/lincketheo/ndbgo/internal/utils"
)

type Interpreter interface {
	RunNext(n *usecases.NDB, p *ByteStack)
}

func RunNext(n *usecases.NDB, p *ByteStack) error {
	if b, err := p.popOpcode(); err != nil {
		return utils.ErrorContext(err)
	} else {
		switch b {

		case OP_CREATE:
			return utils.ErrorContext(interpretCreate(n, p))

		case OP_CONNECT:
			return utils.ErrorContext(interpretConnect(n, p, false))

		case OP_CONNECT_CREATE:
			return utils.ErrorContext(interpretConnect(n, p, true))

		case OP_TERM:
			return fmt.Errorf("Unexpected op code OP_TERM at top level")
		}
	}

	panic("Unreachable")
}

func RunAll(n *usecases.NDB, p *ByteStack) error {
	for !p.empty() {
		if err := RunNext(n, p); err != nil {
			return utils.ErrorContext(err)
		}
	}
	return nil
}

// /////////////////////////////////////// CREATE
func interpretCreate(n *usecases.NDB, p *ByteStack) error {
	if rel, err := p.popEntity(); err != nil {
		return utils.ErrorContext(err)
	} else {
		switch rel {

		case E_DB:
			return utils.ErrorContext(interpretCreateDB(n, p))

		case E_REL:
			return utils.ErrorContext(interpretCreateREL(n, p))

		case E_VAR:
			return utils.ErrorContext(interpretCreateVAR(n, p))
		}
	}

	panic("Unreachable")
}

func interpretCreateDB(n *usecases.NDB, p *ByteStack) error {
	if db, err := p.popString(); err != nil {
		return utils.ErrorContext(err)
	} else if err = p.popOpcodeExpect(OP_TERM); err != nil {
		return utils.ErrorContext(err)
	} else {
		return utils.ErrorContext((*n).CreateDB(db))
	}
}

func interpretCreateREL(n *usecases.NDB, p *ByteStack) error {
	if rel, err := p.popString(); err != nil {
		return utils.ErrorContext(err)
	} else if err := p.popOpcodeExpect(OP_TERM); err != nil {
		return utils.ErrorContext(err)
	} else {
		return utils.ErrorContext((*n).CreateRel(rel))
	}
}

func interpretCreateVAR(n *usecases.NDB, p *ByteStack) error {
	if vari, err := p.popString(); err != nil {
		return utils.ErrorContext(err)
	} else if config, err := p.popVarConfig(); err != nil {
		return utils.ErrorContext(err)
	} else if err = p.popOpcodeExpect(OP_TERM); err != nil {
		return utils.ErrorContext(err)
	} else {
		return utils.ErrorContext((*n).CreateVar(vari, config))
	}
}

///////////////////////////////////////// CONNECT

func interpretConnect(
	n *usecases.NDB,
	p *ByteStack,
	create bool,
) error {
	if e, err := p.popEntity(); err != nil {
		return utils.ErrorContext(err)
	} else {
		switch e {

		case E_DB:
			return utils.ErrorContext(interpretConnectDB(n, p, create))

		case E_REL:
			return utils.ErrorContext(interpretConnectREL(n, p, create))

		case E_VAR:
			return utils.ErrorContext(interpretConnectVAR(n, p, create))
		}
	}

	panic("Unreachable")
}

func interpretConnectDB(
	n *usecases.NDB,
	p *ByteStack,
	create bool,
) error {
	if db, err := p.popString(); err != nil {
		return utils.ErrorContext(err)
	} else if err = p.popOpcodeExpect(OP_TERM); err != nil {
		return utils.ErrorContext(err)
	} else {
		return utils.ErrorContext((*n).ConnectDB(db, create))
	}
}

func interpretConnectREL(
	n *usecases.NDB,
	p *ByteStack,
	create bool,
) error {
	if rel, err := p.popString(); err != nil {
		return utils.ErrorContext(err)
	} else if err = p.popOpcodeExpect(OP_TERM); err != nil {
		return utils.ErrorContext(err)
	} else {
		return utils.ErrorContext((*n).ConnectRel(rel, create))
	}
}

func interpretConnectVAR(
	n *usecases.NDB,
	p *ByteStack,
	create bool,
) error {
	if vari, err := p.popString(); err != nil {
		return utils.ErrorContext(err)
	} else if err = p.popOpcodeExpect(OP_TERM); err != nil {
		return utils.ErrorContext(err)
	} else {
		return utils.ErrorContext((*n).ConnectVar(vari, create))
	}
}
