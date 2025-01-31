package bytecode

import (
	"github.com/lincketheo/ndbgo/internal/logging"
	"github.com/lincketheo/ndbgo/internal/usecases"
)

func (b *ByteStack) CreateDB(db string) error {
	logging.Debug("Pushing CREATE DB: %s onto the stack\n", db)
	b.pushOpcode(OP_CREATE)
	if err := b.pushEntityWithName(E_DB, db); err != nil {
		return err
	}
	b.pushOpcode(OP_TERM)
	return nil
}

func (b *ByteStack) CreateRel(rel string) error {
	logging.Debug("Pushing CREATE REL: %s onto the stack\n", rel)
	b.pushOpcode(OP_CREATE)
	if err := b.pushEntityWithName(E_REL, rel); err != nil {
		return err
	}
	b.pushOpcode(OP_TERM)
	return nil
}

func (b *ByteStack) CreateVar(args usecases.CreateVarArgs) error {
	logging.Debug(`Pushing CREATE Var: %s onto the stack\n
      Shape: %v Dtype: %v`, args.Vari, args.Shape, args.Dtype)

	// CREATE
	b.pushOpcode(OP_CREATE)

	// NAME
	if err := b.pushEntityWithName(E_VAR, args.Vari); err != nil {
		return err
	}

	// DTYPE CONFIG
	b.pushVarDtypeConfig(args.Dtype)

	// SHAPE CONFIG
	if err := b.pushVarShapeConfig(args.Shape); err != nil {
		return err
	}

	// TERM
	b.pushOpcode(OP_TERM)
	return nil
}

func (b *ByteStack) ConnectDB(db string, create bool) error {
	logging.Debug("Pushing CONNECT db: %s onto the stack\n", db)

	// CONNECT
	b.pushOpcode(OP_CONNECT)

	// NAME
	if err := b.pushEntityWithName(E_DB, db); err != nil {
		return err
	}

	// TERM
	b.pushOpcode(OP_TERM)
	return nil
}

func (b *ByteStack) ConnectRel(rel string, create bool) error {
	logging.Debug("Pushing CONNECT rel: %s onto the stack\n", rel)

	// CONNECT
	b.pushOpcode(OP_CONNECT)

	// NAME
	if err := b.pushEntityWithName(E_REL, rel); err != nil {
		return err
	}

	// TERM
	b.pushOpcode(OP_TERM)

	return nil
}

func (b *ByteStack) ConnectVar(vari string, create bool) error {
	logging.Debug("Pushing CONNECT var: %s onto the stack\n", vari)

	b.pushOpcode(OP_CONNECT)
	if err := b.pushEntityWithName(E_VAR, vari); err != nil {
		return err
	}

	b.pushOpcode(OP_TERM)

	return nil
}
