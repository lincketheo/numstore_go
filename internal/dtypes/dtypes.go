package dtypes

import "fmt"

type Dtype byte

const (
	U8 Dtype = iota
	U16
	U32
	U64
	U128

	I8
	I16
	I32
	I64
	I128

	F32
	F64
	F128

	CF32
	CF64
	CF128

	CU16
	CU32
	CU64
	CU128

	CI16
	CI32
	CI64
	CI128
)

func ByteToDtype(b byte) (Dtype, error) {
	if !(b <= byte(CF128) && b >= byte(U8)) {
		return 0, fmt.Errorf("Expected dtype, got byte: %v", b)
	}
	return Dtype(b), nil
}

func StrtoDtype(dtstr string) (Dtype, error) {
	if len(dtstr) == 0 {
		return 0, fmt.Errorf("Invalid dtype: %s", dtstr)
	}

	switch dtstr {
	case "U8":
		return U8, nil
	case "U16":
		return U16, nil
	case "U32":
		return U32, nil
	case "U64":
		return U64, nil
	case "U128":
		return U128, nil
	case "I8":
		return I8, nil
	case "I16":
		return I16, nil
	case "I32":
		return I32, nil
	case "I64":
		return I64, nil
	case "I128":
		return I128, nil
	case "F32":
		return F32, nil
	case "F64":
		return F64, nil
	case "F128":
		return F128, nil
	case "CF32":
		return CF32, nil
	case "CF64":
		return CF64, nil
	case "CF128":
		return CF128, nil
	case "CU16":
		return CU16, nil
	case "CU32":
		return CU32, nil
	case "CU64":
		return CU64, nil
	case "CU128":
		return CU128, nil
	case "CI16":
		return CI16, nil
	case "CI32":
		return CI32, nil
	case "CI64":
		return CI64, nil
	case "CI128":
		return CI128, nil
	default:
		return 0, fmt.Errorf("Invalid dtype: %s", dtstr)
	}
}

func (d Dtype) String() string {
	switch d {
	case U8:
		return "U8"
	case U16:
		return "U16"
	case U32:
		return "U32"
	case U64:
		return "U64"
	case U128:
		return "U128"
	case I8:
		return "I8"
	case I16:
		return "I16"
	case I32:
		return "I32"
	case I64:
		return "I64"
	case I128:
		return "I128"
	case F32:
		return "F32"
	case F64:
		return "F64"
	case F128:
		return "F128"
	case CF32:
		return "CF32"
	case CF64:
		return "CF64"
	case CF128:
		return "CF128"
	case CU16:
		return "CU16"
	case CU32:
		return "CU32"
	case CU64:
		return "CU64"
	case CU128:
		return "CU128"
	case CI16:
		return "CI16"
	case CI32:
		return "CI32"
	case CI64:
		return "CI64"
	case CI128:
		return "CI128"
	default:
		panic("Invalid dtype encountered. This is a bug")
	}
}
