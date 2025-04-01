package numstore

import (
	"encoding/json"

	"github.com/lincketheo/numstore/internal/nserror"
)

type Dtype uint16

const (
	U32 Dtype = iota
	U16
)

func (d Dtype) String() string {
	switch d {
	case U32:
		return "U32"
	case U16:
		return "U16"
	}
	panic("Unreachable")
}

func (d Dtype) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func dtypeSizeof(dtype Dtype) uint32 {
	switch dtype {
	case U32:
		return 4
	case U16:
		return 2
	default:
		panic("Unreachable")
	}
}

func (d *Dtype) UnmarshalJSON(data []byte) error {
	var name string
	if err := json.Unmarshal(data, &name); err != nil {
		return err
	}

	switch name {
	case U32.String():
		*d = U32
	case U16.String():
		*d = U16
	default:
		return nserror.JSONInvalidDtype
	}

	return nil
}
