package bytecode

import (
	"encoding/binary"

	"github.com/lincketheo/numstore/internal/numstore"
)

type ByteCode struct {
	Bytes []byte
	pos   int 
}

const (
	OP_CREATE byte = iota
	OP_DELETE
	OP_WRITE
	OP_READ
	OP_TAKE

	OP_PRIM
	OP_STRICT_ARR
	OP_VAR_ARR
	OP_STRUCT
	OP_ENUM
	OP_UNION
)

// //////////////////////////////////// Main Methods
func (b *ByteCode) HandleCreate(vname string, tp numstore.Type) {
	b.pushByte(OP_CREATE)
	b.pushString(vname)
	b.pushType(tp)
}

func (b *ByteCode) PopHandleCreate() (vname string, tp numstore.Type, ok bool) {
	op, ok := b.popByte()
	if !ok || op != OP_CREATE {
		return "", nil, false
	}
	vname, ok = b.popString()
	if !ok {
		return "", nil, false
	}
	tp, ok = b.popType()
	return vname, tp, ok
}

func (b *ByteCode) HandleDelete(vname string) {
	b.pushByte(OP_DELETE)
	b.pushString(vname)
}

func (b *ByteCode) PopHandleDelete() (vname string, ok bool) {
	op, ok := b.popByte()
	if !ok || op != OP_DELETE {
		return "", false
	}
	return b.popString()
}

func (b *ByteCode) HandleRead(fmt numstore.ReadFormat) {
	b.pushByte(OP_READ)
	b.pushReadFormat(fmt)
}

func (b *ByteCode) PopHandleRead() (rf numstore.ReadFormat, ok bool) {
	op, ok := b.popByte()
	if !ok || op != OP_READ {
		return numstore.ReadFormat{}, false
	}
	return b.popReadFormat()
}

func (b *ByteCode) HandleWrite(fmt numstore.WriteFormat) {
	b.pushByte(OP_WRITE)
	b.pushWriteFormat(fmt)
}

func (b *ByteCode) PopHandleWrite() (wf numstore.WriteFormat, ok bool) {
	op, ok := b.popByte()
	if !ok || op != OP_WRITE {
		return numstore.WriteFormat{}, false
	}
	return b.popWriteFormat()
}

func (b *ByteCode) HandleTake(fmt numstore.ReadFormat) {
	b.pushByte(OP_TAKE)
	b.pushReadFormat(fmt)
}

func (b *ByteCode) PopHandleTake() (rf numstore.ReadFormat, ok bool) {
	op, ok := b.popByte()
	if !ok || op != OP_TAKE {
		return numstore.ReadFormat{}, false
	}
	return b.popReadFormat()
}

// /////////////////////////////// Utils
func (b *ByteCode) pushByte(bt byte) {
	b.Bytes = append(b.Bytes, bt)
}

func (b *ByteCode) popByte() (byte, bool) {
	if b.pos >= len(b.Bytes) {
		return 0, false
	}
	bt := b.Bytes[b.pos]
	b.pos++
	return bt, true
}

func (b *ByteCode) pushUint32(v uint32) {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, v)
	b.Bytes = append(b.Bytes, buf...)
}

func (b *ByteCode) popUint32() (uint32, bool) {
	if b.pos+4 > len(b.Bytes) {
		return 0, false
	}
	v := binary.LittleEndian.Uint32(b.Bytes[b.pos : b.pos+4])
	b.pos += 4
	return v, true
}

func (b *ByteCode) pushString(str string) {
	strBytes := []byte(str)
	length := uint32(len(strBytes))
	b.pushUint32(length)
	b.Bytes = append(b.Bytes, strBytes...)
}

func (b *ByteCode) popString() (string, bool) {
	length, ok := b.popUint32()
	if !ok {
		return "", false
	}
	if b.pos+int(length) > len(b.Bytes) {
		return "", false
	}
	str := string(b.Bytes[b.pos : b.pos+int(length)])
	b.pos += int(length)
	return str, true
}

func (b *ByteCode) pushDimRange(dr numstore.DimRange) {
	b.pushUint32(uint32(dr.Start))
	b.pushUint32(uint32(dr.Stop))
	b.pushUint32(uint32(dr.Step))
	var inf byte
	if dr.IsInf {
		inf = 1
	}
	b.pushByte(inf)
	var single byte
	if dr.IsSingle {
		single = 1
	}
	b.pushByte(single)
}

func (b *ByteCode) popDimRange() (numstore.DimRange, bool) {
	start, ok := b.popUint32()
	if !ok {
		return numstore.DimRange{}, false
	}
	stop, ok := b.popUint32()
	if !ok {
		return numstore.DimRange{}, false
	}
	step, ok := b.popUint32()
	if !ok {
		return numstore.DimRange{}, false
	}
	inf, ok := b.popByte()
	if !ok {
		return numstore.DimRange{}, false
	}
	single, ok := b.popByte()
	if !ok {
		return numstore.DimRange{}, false
	}
	return numstore.DimRange{
		Start:    int(start),
		Stop:     int(stop),
		Step:     int(step),
		IsInf:    inf != 0,
		IsSingle: single != 0,
	}, true
}

func (b *ByteCode) pushReadVariable(rv numstore.ReadVariable) {
	b.pushString(rv.Vname)
	b.pushUint32(uint32(len(rv.Range)))
	for _, dr := range rv.Range {
		b.pushDimRange(dr)
	}
}

func (b *ByteCode) popReadVariable() (numstore.ReadVariable, bool) {
	vname, ok := b.popString()
	if !ok {
		return numstore.ReadVariable{}, false
	}
	count, ok := b.popUint32()
	if !ok {
		return numstore.ReadVariable{}, false
	}
	drs := make([]numstore.DimRange, 0, count)
	for i := 0; i < int(count); i++ {
		dr, ok := b.popDimRange()
		if !ok {
			return numstore.ReadVariable{}, false
		}
		drs = append(drs, dr)
	}
	return numstore.ReadVariable{
		Vname: vname,
		Range: drs,
	}, true
}

func (b *ByteCode) pushReadFormat(rf numstore.ReadFormat) {
	b.pushUint32(uint32(rf.ToRead))
	b.pushUint32(uint32(len(rf.Variables)))
	for _, group := range rf.Variables {
		b.pushUint32(uint32(len(group)))
		for _, rv := range group {
			b.pushReadVariable(rv)
		}
	}
}

func (b *ByteCode) popReadFormat() (numstore.ReadFormat, bool) {
	toRead, ok := b.popUint32()
	if !ok {
		return numstore.ReadFormat{}, false
	}
	groupCount, ok := b.popUint32()
	if !ok {
		return numstore.ReadFormat{}, false
	}
	variables := make([][]numstore.ReadVariable, 0, groupCount)
	for i := 0; i < int(groupCount); i++ {
		varCount, ok := b.popUint32()
		if !ok {
			return numstore.ReadFormat{}, false
		}
		group := make([]numstore.ReadVariable, 0, varCount)
		for j := 0; j < int(varCount); j++ {
			rv, ok := b.popReadVariable()
			if !ok {
				return numstore.ReadFormat{}, false
			}
			group = append(group, rv)
		}
		variables = append(variables, group)
	}
	return numstore.ReadFormat{
		ToRead:    int(toRead),
		Variables: variables,
	}, true
}

func (b *ByteCode) pushWriteFormat(wf numstore.WriteFormat) {
	b.pushUint32(uint32(wf.ToWrite))
	b.pushUint32(uint32(len(wf.Variables)))
	for _, group := range wf.Variables {
		b.pushUint32(uint32(len(group)))
		for _, str := range group {
			b.pushString(str)
		}
	}
}

func (b *ByteCode) popWriteFormat() (numstore.WriteFormat, bool) {
	toWrite, ok := b.popUint32()
	if !ok {
		return numstore.WriteFormat{}, false
	}
	groupCount, ok := b.popUint32()
	if !ok {
		return numstore.WriteFormat{}, false
	}
	variables := make([][]string, 0, groupCount)
	for i := 0; i < int(groupCount); i++ {
		strCount, ok := b.popUint32()
		if !ok {
			return numstore.WriteFormat{}, false
		}
		group := make([]string, 0, strCount)
		for j := 0; j < int(strCount); j++ {
			s, ok := b.popString()
			if !ok {
				return numstore.WriteFormat{}, false
			}
			group = append(group, s)
		}
		variables = append(variables, group)
	}
	return numstore.WriteFormat{
		ToWrite:   int(toWrite),
		Variables: variables,
	}, true
}

func (b *ByteCode) pushType(tp numstore.Type) {
	switch tp.Kind() {
	case numstore.PrimKind:
		prim, ok := tp.(*numstore.PrimitiveType)
		if !ok {
			panic("pushType: expected *PrimitiveType")
		}
		b.pushByte(OP_PRIM)
		b.pushByte(byte(prim.PT))
	case numstore.StrictArrayKind:
		arr, ok := tp.(*numstore.StrictArrayType)
		if !ok {
			panic("pushType: expected *StrictArrayType")
		}
		b.pushByte(OP_STRICT_ARR)
		b.pushUint32(uint32(len(arr.Dims)))
		for _, dim := range arr.Dims {
			b.pushUint32(dim)
		}
		b.pushType(arr.Of)
	case numstore.VarArrayKind:
		varr, ok := tp.(*numstore.VarArrayType)
		if !ok {
			panic("pushType: expected *VarArrayType")
		}
		b.pushByte(OP_VAR_ARR)
		b.pushUint32(varr.Rank)
		b.pushType(varr.Of)
	case numstore.StructKind:
		st, ok := tp.(*numstore.StructType)
		if !ok {
			panic("pushType: expected *StructType")
		}
		b.pushByte(OP_STRUCT)
		b.pushUint32(uint32(len(st.Fields)))
		for name, fieldType := range st.Fields {
			b.pushString(name)
			b.pushType(fieldType)
		}
	case numstore.EnumKind:
		enm, ok := tp.(*numstore.EnumType)
		if !ok {
			panic("pushType: expected *EnumType")
		}
		b.pushByte(OP_ENUM)
		b.pushUint32(uint32(len(enm.Options)))
		for _, option := range enm.Options {
			b.pushString(option)
		}
	case numstore.UnionKind:
		uni, ok := tp.(*numstore.UnionType)
		if !ok {
			panic("pushType: expected *UnionType")
		}
		b.pushByte(OP_UNION)
		b.pushUint32(uint32(len(uni.Fields)))
		for name, fieldType := range uni.Fields {
			b.pushString(name)
			b.pushType(fieldType)
		}
	default:
		panic("pushType: unknown type kind")
	}
}

func (b *ByteCode) popType() (numstore.Type, bool) {
	op, ok := b.popByte()
	if !ok {
		return nil, false
	}
	switch op {
	case OP_PRIM:
		primCode, ok := b.popByte()
		if !ok {
			return nil, false
		}
		return &numstore.PrimitiveType{PT: numstore.PrimType(primCode)}, true
	case OP_STRICT_ARR:
		dimCount, ok := b.popUint32()
		if !ok {
			return nil, false
		}
		dims := make([]uint32, dimCount)
		for i := 0; i < int(dimCount); i++ {
			d, ok := b.popUint32()
			if !ok {
				return nil, false
			}
			dims[i] = d
		}
		of, ok := b.popType()
		if !ok {
			return nil, false
		}
		return &numstore.StrictArrayType{Dims: dims, Of: of}, true
	case OP_VAR_ARR:
		rank, ok := b.popUint32()
		if !ok {
			return nil, false
		}
		of, ok := b.popType()
		if !ok {
			return nil, false
		}
		return &numstore.VarArrayType{Rank: rank, Of: of}, true
	case OP_STRUCT:
		fieldCount, ok := b.popUint32()
		if !ok {
			return nil, false
		}
		fields := make(map[string]numstore.Type)
		for i := 0; i < int(fieldCount); i++ {
			name, ok := b.popString()
			if !ok {
				return nil, false
			}
			ft, ok := b.popType()
			if !ok {
				return nil, false
			}
			fields[name] = ft
		}
		return &numstore.StructType{Fields: fields}, true
	case OP_ENUM:
		optionCount, ok := b.popUint32()
		if !ok {
			return nil, false
		}
		options := make([]string, optionCount)
		for i := 0; i < int(optionCount); i++ {
			opt, ok := b.popString()
			if !ok {
				return nil, false
			}
			options[i] = opt
		}
		return &numstore.EnumType{Options: options}, true
	case OP_UNION:
		fieldCount, ok := b.popUint32()
		if !ok {
			return nil, false
		}
		fields := make(map[string]numstore.Type)
		for i := 0; i < int(fieldCount); i++ {
			name, ok := b.popString()
			if !ok {
				return nil, false
			}
			ft, ok := b.popType()
			if !ok {
				return nil, false
			}
			fields[name] = ft
		}
		return &numstore.UnionType{Fields: fields}, true
	default:
		return nil, false
	}
}
