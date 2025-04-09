package numstore

import (
	"encoding/json"
	"fmt"
)

//
// TYPE DECLARATIONS
//

type NSTypeKind string

const (
	PrimKind        NSTypeKind = "primitive"
	StrictArrayKind NSTypeKind = "strict_array"
	VarArrayKind    NSTypeKind = "var_array"
	StructKind      NSTypeKind = "struct"
	EnumKind        NSTypeKind = "enum"
	UnionKind       NSTypeKind = "union"
)

type Type interface {
	Kind() NSTypeKind
}

type PrimitiveType struct {
	PT PrimType `json:"type"`
}

type StrictArrayType struct {
	Dims []uint32 `json:"dims"`
	Of   Type     `json:"of"`
}

type VarArrayType struct {
	Rank uint32 `json:"rank"`
	Of   Type   `json:"of"`
}

type StructType struct {
	Fields map[string]Type `json:"fields"`
}

type EnumType struct {
	Options []string `json:"options"`
}

type UnionType struct {
	Fields map[string]Type `json:"fields"`
}

//
// Kind() IMPLEMENTATIONS
//

func (PrimitiveType) Kind() NSTypeKind   { return PrimKind }
func (StrictArrayType) Kind() NSTypeKind { return StrictArrayKind }
func (VarArrayType) Kind() NSTypeKind    { return VarArrayKind }
func (StructType) Kind() NSTypeKind      { return StructKind }
func (EnumType) Kind() NSTypeKind        { return EnumKind }
func (UnionType) Kind() NSTypeKind       { return UnionKind }

//
// JSON MARSHAL/UNMARSHAL
//

func (p PrimitiveType) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Kind NSTypeKind    `json:"kind"`
		PT   PrimType `json:"type"`
	}{
		Kind: PrimKind,
		PT:   p.PT,
	})
}

func (p *PrimitiveType) UnmarshalJSON(data []byte) error {
	aux := struct {
		PT PrimType `json:"type"`
	}{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	p.PT = aux.PT
	return nil
}

func (t StrictArrayType) MarshalJSON() ([]byte, error) {
	ofBytes, err := json.Marshal(t.Of)
	if err != nil {
		return nil, err
	}
	return json.Marshal(struct {
		Kind NSTypeKind      `json:"kind"`
		Rank uint32          `json:"rank"`
		Dims []uint32        `json:"dims"`
		Of   json.RawMessage `json:"of"`
	}{
		Kind: StrictArrayKind,
		Dims: t.Dims,
		Of:   ofBytes,
	})
}

func (t *StrictArrayType) UnmarshalJSON(data []byte) error {
	aux := struct {
		Rank uint32          `json:"rank"`
		Dims []uint32        `json:"dims"`
		Of   json.RawMessage `json:"of"`
	}{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	t.Dims = aux.Dims
	inner, err := unmarshalType(aux.Of)
	if err != nil {
		return err
	}
	t.Of = inner
	return nil
}

func (t VarArrayType) MarshalJSON() ([]byte, error) {
	ofBytes, err := json.Marshal(t.Of)
	if err != nil {
		return nil, err
	}
	return json.Marshal(struct {
		Kind NSTypeKind      `json:"kind"`
		Rank uint32          `json:"rank"`
		Of   json.RawMessage `json:"of"`
	}{
		Kind: VarArrayKind,
		Rank: t.Rank,
		Of:   ofBytes,
	})
}

func (t *VarArrayType) UnmarshalJSON(data []byte) error {
	aux := struct {
		Rank uint32          `json:"rank"`
		Of   json.RawMessage `json:"of"`
	}{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	t.Rank = aux.Rank
	inner, err := unmarshalType(aux.Of)
	if err != nil {
		return err
	}
	t.Of = inner
	return nil
}

func (s StructType) MarshalJSON() ([]byte, error) {
	fields := make(map[string]json.RawMessage, len(s.Fields))
	for name, t := range s.Fields {
		b, err := json.Marshal(t)
		if err != nil {
			return nil, err
		}
		fields[name] = b
	}
	return json.Marshal(struct {
		Kind   NSTypeKind                 `json:"kind"`
		Fields map[string]json.RawMessage `json:"fields"`
	}{
		Kind:   StructKind,
		Fields: fields,
	})
}

func (s *StructType) UnmarshalJSON(data []byte) error {
	aux := struct {
		Fields map[string]json.RawMessage `json:"fields"`
	}{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	s.Fields = make(map[string]Type, len(aux.Fields))
	for name, raw := range aux.Fields {
		inner, err := unmarshalType(raw)
		if err != nil {
			return err
		}
		s.Fields[name] = inner
	}
	return nil
}

func (e EnumType) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Kind    NSTypeKind `json:"kind"`
		Options []string   `json:"options"`
	}{
		Kind:    EnumKind,
		Options: e.Options,
	})
}

func (e *EnumType) UnmarshalJSON(data []byte) error {
	aux := struct {
		Options []string `json:"options"`
	}{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	e.Options = aux.Options
	return nil
}

func (u UnionType) MarshalJSON() ([]byte, error) {
	fields := make(map[string]json.RawMessage, len(u.Fields))
	for name, t := range u.Fields {
		b, err := json.Marshal(t)
		if err != nil {
			return nil, err
		}
		fields[name] = b
	}
	return json.Marshal(struct {
		Kind   NSTypeKind                 `json:"kind"`
		Fields map[string]json.RawMessage `json:"fields"`
	}{
		Kind:   UnionKind,
		Fields: fields,
	})
}

func (u *UnionType) UnmarshalJSON(data []byte) error {
	aux := struct {
		Fields map[string]json.RawMessage `json:"fields"`
	}{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	u.Fields = make(map[string]Type, len(aux.Fields))
	for name, raw := range aux.Fields {
		inner, err := unmarshalType(raw)
		if err != nil {
			return err
		}
		u.Fields[name] = inner
	}
	return nil
}

// Helper to dispatch based on "kind"
func unmarshalType(raw json.RawMessage) (Type, error) {
	var kind struct {
		Kind NSTypeKind `json:"kind"`
	}
	if err := json.Unmarshal(raw, &kind); err != nil {
		return nil, err
	}
	var t Type
	switch kind.Kind {
	case PrimKind:
		t = new(PrimitiveType)
	case StrictArrayKind:
		t = new(StrictArrayType)
	case VarArrayKind:
		t = new(VarArrayType)
	case StructKind:
		t = new(StructType)
	case EnumKind:
		t = new(EnumType)
	case UnionKind:
		t = new(UnionType)
	default:
		return nil, fmt.Errorf("unknown kind %q", kind.Kind)
	}
	if err := json.Unmarshal(raw, t); err != nil {
		return nil, err
	}
	return t, nil
}

//
// PRIMITIVE‑TYPE
//

type PrimType uint8

const (
	I8 PrimType = iota
	I16
	I32
	I64
	I128

	U8
	U16
	U32
	U64
	U128

	F16
	F32
	F64
	F128

	CF32
	CF64
	CF128
	CF256

	Char
	Bool
)

func PrimTypeSizeof(pt PrimType) uint32 {
	switch pt {
	case I8, U8, Char, Bool:
		return 1

	case I16, U16, F16:
		return 2

	case I32, U32, F32:
		return 4
	case CF32:
		return 8 // 2 * 4

	case I64, U64, F64:
		return 8
	case CF64:
		return 16 // 2 * 8

	case I128, U128, F128:
		return 16
	case CF128:
		return 32 // 2 * 16

	case CF256:
		return 64 // 2 * 32

	default:
		panic("Unreachable")
	}
}

var (
	_ptToString = map[PrimType]string{
		I8:    "i8",
		I16:   "i16",
		I32:   "i32",
		I64:   "i64",
		I128:  "i128",
		U8:    "u8",
		U16:   "u16",
		U32:   "u32",
		U64:   "u64",
		U128:  "u128",
		F16:   "f16",
		F32:   "f32",
		F64:   "f64",
		F128:  "f128",
		CF32:  "cf32",
		CF64:  "cf64",
		CF128: "cf128",
		CF256: "cf256",
		Char:  "char",
		Bool:  "bool",
	}
	_stringToPT = func() map[string]PrimType {
		m := make(map[string]PrimType, len(_ptToString))
		for k, v := range _ptToString {
			m[v] = k
		}
		return m
	}()
)

func (d PrimType) String() string {
	if s, ok := _ptToString[d]; ok {
		return s
	}
	panic("invalid PrimType")
}

func PrimitiveTypeFromString(s string) (PrimType, bool) {
	pt, ok := _stringToPT[s]
	return pt, ok
}

func (d PrimType) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *PrimType) UnmarshalJSON(data []byte) error {
	var name string
	if err := json.Unmarshal(data, &name); err != nil {
		return err
	}
	if pt, ok := PrimitiveTypeFromString(name); ok {
		*d = pt
		return nil
	}
	return fmt.Errorf("invalid PrimType %q", name)
}
