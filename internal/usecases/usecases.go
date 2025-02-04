package usecases

import (
	"github.com/lincketheo/ndbgo/internal/dtypes"
	"github.com/lincketheo/ndbgo/internal/utils"
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

  Read(a *NDBApp)
}

func CreateDBRel(
	db,
	rel string,
	n *NDB,
) error {
	if err := (*n).ConnectDB(db); err != nil {
		return utils.ErrorContext(err)
	} else {
		return utils.ErrorContext((*n).CreateRel(rel))
	}
}

func CreateDBRelVar(
	db,
	rel,
	vari string,
	config VarConfig,
	n *NDB,
) error {
	if err := (*n).ConnectDB(db); err != nil {
		return utils.ErrorContext(err)
	} else if err = (*n).ConnectRel(rel); err != nil {
		return utils.ErrorContext(err)
	} else {
		return utils.ErrorContext((*n).CreateVar(vari, config))
	}
}
