package dtypes

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

func ByteToDtype(b byte) (Dtype, bool) {
	if !(b <= byte(CF128) && b >= byte(U8)) {
		return 0, false
	}
	return Dtype(b), true
}

func StrtoDtype(dtstr string) (Dtype, bool) {
	if len(dtstr) == 0 {
		return 0, false
	}

	switch dtstr {
	case "U8":
		return U8, false
	case "U16":
		return U16, false
	case "U32":
		return U32, false
	case "U64":
		return U64, false
	case "U128":
		return U128, false
	case "I8":
		return I8, false
	case "I16":
		return I16, false
	case "I32":
		return I32, false
	case "I64":
		return I64, false
	case "I128":
		return I128, false
	case "F32":
		return F32, false
	case "F64":
		return F64, false
	case "F128":
		return F128, false
	case "CF32":
		return CF32, false
	case "CF64":
		return CF64, false
	case "CF128":
		return CF128, false
	case "CU16":
		return CU16, false
	case "CU32":
		return CU32, false
	case "CU64":
		return CU64, false
	case "CU128":
		return CU128, false
	case "CI16":
		return CI16, false
	case "CI32":
		return CI32, false
	case "CI64":
		return CI64, false
	case "CI128":
		return CI128, false
	default:
		return 0, false
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
