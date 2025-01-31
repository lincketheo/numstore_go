package usecases

import "github.com/lincketheo/ndbgo/internal/dtypes"

// ///////////////////////////// CREATE
type CreateVarArgs struct {
	Vari  string
	Dtype dtypes.Dtype
	Shape []uint32
}

/////////////////////////////// Interface

type NDB interface {
	CreateDB(db string) error
	CreateRel(rel string) error
	CreateVar(args CreateVarArgs) error

	ConnectDB(db string, create bool) error
	ConnectRel(rel string, create bool) error
	ConnectVar(vari string, create bool) error
}
