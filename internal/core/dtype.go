package core

import (
	"encoding/json"
	"fmt"
)

type Dtype uint16

const (
	U16 Dtype = iota
	U32
	U64
)

func (d Dtype) String() string {
	switch d {
	case U16:
		return "U16"
	case U32:
		return "U32"
	case U64:
		return "U64"
	}
	panic("Unreachable")
}

func DtypeFromString(s string) (Dtype, bool) {
	switch s {
	case U16.String():
		{
			return U16, true
		}
	case U32.String():
		{
			return U32, true
		}
	case U64.String():
		{
			return U64, true
		}
	}
  return 0, false
}

func (d Dtype) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func DtypeSizeof(dtype Dtype) uint32 {
	switch dtype {
	case U16:
		return 2
	case U32:
		return 4
	case U64:
		return 8
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
	case U16.String():
		*d = U16
	case U32.String():
		*d = U32
	case U64.String():
		*d = U64
	default:
		return fmt.Errorf("Invalid JSON Dtype")
	}

	return nil
}
