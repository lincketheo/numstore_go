package usecases

import (
	"github.com/lincketheo/ndbgo/internal/dtypes"
	"github.com/lincketheo/ndbgo/internal/logging"
)

// ///////////////////////////// CREATE
type VarConfig struct {
	Dtype dtypes.Dtype
	Shape []uint32
}

/////////////////////////////// Interface

type NDB interface {
	CreateDB(db string) error
	CreateRel(rel string) error
	CreateVar(vari string, config VarConfig) error

	ConnectDB(db string) error
	ConnectRel(rel string) error
	ConnectVar(vari string) error
}

// ///////////////////////////// Simple logging implementation

type NDBlogger struct{}

func (n NDBlogger) CreateDB(db string) error {
	logging.Debug("Executing CREATE DB: %s\n", db)
	return nil
}

func (n NDBlogger) CreateRel(rel string) error {
	logging.Debug("Executing CREATE REL: %s\n", rel)
	return nil
}

func (n NDBlogger) CreateVar(vari string, args VarConfig) error {
	logging.Debug(`Executing CREATE Var: %s\n
      Shape: %v Dtype: %v`, vari, args.Shape, args.Dtype)
	return nil
}

func (n NDBlogger) ConnectDB(db string) error {
	logging.Debug("Executing CONNECT db: %s\n", db)
	return nil
}

func (n NDBlogger) ConnectRel(rel string) error {
	logging.Debug("Executing CONNECT rel: %s\n", rel)
	return nil
}

func (n NDBlogger) ConnectVar(vari string) error {
	logging.Debug("Executing CONNECT var: %s\n", vari)
	return nil
}
