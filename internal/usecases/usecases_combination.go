package usecases

import "github.com/lincketheo/ndbgo/internal/utils"

type NDBCombination []NDB

func (n NDBCombination) CreateDB(db string) error {
	for _, i := range n {
		if err := i.CreateDB(db); err != nil {
			return utils.ErrorContext(err)
		}
	}
	return nil
}

func (n NDBCombination) CreateRel(rel string) error {
	for _, i := range n {
		if err := i.CreateRel(rel); err != nil {
			return utils.ErrorContext(err)
		}
	}
	return nil
}

func (n NDBCombination) CreateVar(vari string, args VarConfig) error {
	for _, i := range n {
		if err := i.CreateVar(vari, args); err != nil {
			return utils.ErrorContext(err)
		}
	}
	return nil
}

func (n NDBCombination) ConnectDB(db string) error {
	for _, i := range n {
		if err := i.ConnectDB(db); err != nil {
			return utils.ErrorContext(err)
		}
	}
	return nil
}

func (n NDBCombination) ConnectRel(rel string) error {
	for _, i := range n {
		if err := i.ConnectRel(rel); err != nil {
			return utils.ErrorContext(err)
		}
	}
	return nil
}

func (n NDBCombination) ConnectVar(vari string) error {
	for _, i := range n {
		if err := i.ConnectVar(vari); err != nil {
			return utils.ErrorContext(err)
		}
	}
	return nil
}

