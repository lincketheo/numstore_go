package usecases

import (
	"github.com/lincketheo/ndbgo/internal/dtypes"
)

// ///////////////////////////// INPUTS
type VarConfig struct {
	Dtype dtypes.Dtype
	Shape []uint32
}

// Continuous reader
type ReaderConfig struct {
	fmt [][]string
}

// Continuous writer
type WriterConfig struct {
	fmt [][]string
}

// One time write
type WriteArgs struct {
	writeto [][]string // "[a, b], c, d, [e, f]"
	data    []byte     // The contiguous data - C style for each
}

// Indexing types (a[1] or a[1:2])
type Index interface {
	isIndex()
}

type Slice struct {
	start int
	stop  int
	step  int
}

type Integer int

func (i Integer) isIndex() {}
func (s Slice) isIndex()   {}

// Read from a single variable
type ReadFromVar struct {
	readfrom string
	indexes  []Index
}

// Do a one time read
type ReadArgs struct {
	callback func(data []byte) error
	readfrom [][]ReadFromVar
}

/////////////////////////////// Interface

type NDB interface {
	// Create things
	CreateDB(db string) error
	CreateRel(rel string) error
	CreateVar(vari string, config VarConfig) error

	// Connect to things
	ConnectDB(db string) error
	ConnectRel(rel string) error

	// Disconnect from things
	DisconnectDB() error
	DisconnectRel() error

	// Delete things
	DeleteDB(db string) error
	DeleteRel(rel string) error
	DeleteVar(vari string) error

	// Manage readers
	AddReader(config ReaderConfig) (int, error)
	RemoveReader(id int) error

	// Manage writers
	AddWriter(config WriterConfig) (int, error)
	RemoveWriter(id int) error

	// One time write
	Write(args WriteArgs) error
	Read(args ReadArgs) error
}
