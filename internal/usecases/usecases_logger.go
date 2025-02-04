package usecases

import "github.com/lincketheo/ndbgo/internal/logging"

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
