package usecases

import "github.com/lincketheo/ndbgo/internal/utils"

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
