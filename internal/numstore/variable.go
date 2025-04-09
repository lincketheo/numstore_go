package numstore

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/lincketheo/numstore/internal/nserror"
	"github.com/lincketheo/numstore/internal/utils"
)

type VariableMeta struct {
	Name  string        `json:"name"`
	Dtype PrimType `json:"dtype"`
	Shape []uint32      `json:"shape"`
}

type WritingVariable struct {
	vfd        *os.File
	tfd        *os.File
	dataBuffer []byte
	timeBuffer []byte
}

/////////////////////////////////// Public

func CreateVariable(dbname string, v VariableMeta) error {
	// Check if variable exists
	if exists, err := varExistsAndValid(dbname, v.Name); err != nil {
		return nserror.ErrorStack(err)
	} else if exists {
		return fmt.Errorf("Variable: %s already exists in db: %s", v.Name, dbname)
	}

	if err := createVarFolder(dbname, v.Name); err != nil {
		return nserror.ErrorStack(err)
	}

	if err := createVarMetaFile(dbname, v); err != nil {
		return nserror.ErrorStack(err)
	}

	return nil
}

func LoadVariableMeta(dbname, vname string) (VariableMeta, error) {
	fname := varMetaFileName(dbname, vname)
	data, err := os.ReadFile(fname)
	if err != nil {
		return VariableMeta{}, nserror.ErrorStack(err)
	}

	var v VariableMeta
	if err := json.Unmarshal(data, &v); err != nil {
		return VariableMeta{}, nserror.ErrorStack(err)
	}

	return v, nil
}

func OpenVariable(dbname string, v VariableMeta) (WritingVariable, error) {
	dfd, err := os.OpenFile(varDataFileName(dbname, v.Name), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return WritingVariable{}, nserror.ErrorStack(err)
	}

	tfd, err := os.OpenFile(varTimeFileName(dbname, v.Name), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		dfd.Close()
		return WritingVariable{}, nserror.ErrorStack(err)
	}

	dataSize := PrimTypeSizeof(v.Dtype)
	if len(v.Shape) > 0 {
		dataSize *= utils.ReduceMultU32(v.Shape)
	}
	timeSize := PrimTypeSizeof(U64)

	data := make([]byte, dataSize)
	time := make([]byte, timeSize)

	ret := WritingVariable{
		vfd:        dfd,
		tfd:        tfd,
		dataBuffer: data,
		timeBuffer: time,
	}

	return ret, nil
}

func (v WritingVariable) WriteNext(r io.Reader, t uint64) error {
	if n, err := r.Read(v.dataBuffer); err != nil {
		return nserror.ErrorStack(err)
	} else if n != len(v.dataBuffer) {
		panic("Invalid read")
	}

	if n, err := v.vfd.Write(v.dataBuffer); err != nil {
		return nserror.ErrorStack(err)
	} else if n != len(v.dataBuffer) {
		panic("Invalid write")
	}

	binary.BigEndian.PutUint64(v.timeBuffer, t)
	if n, err := v.tfd.Write(v.timeBuffer); err != nil {
		return nserror.ErrorStack(err)
	} else if n != len(v.timeBuffer) {
		panic("Invalid write")
	}

	return nil
}

func (v WritingVariable) Close() error {
	var err error = nil

	err = v.tfd.Close()
	err = v.vfd.Close()

	return nserror.ErrorStack(err)
}

/////////////////////////////////// Private

func varFolderName(db, vname string) string {
	return db + "/" + vname
}

func varMetaFileName(db, vname string) string {
	return varFolderName(db, vname) + "/meta.json"
}

func varDataFileName(db, vname string) string {
	return varFolderName(db, vname) + "/data.bin"
}

func varTimeFileName(db, vname string) string {
	return varFolderName(db, vname) + "/time.bin"
}

func varExistsAndValid(db, v string) (bool, error) {
	if stat, err := os.Stat(varFolderName(db, v)); err != nil && os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, nserror.ErrorStack(err)
	} else {
		// Todo - check validity of stat
		return stat != nil, nil
	}
}

func createVarFolder(dbname, vname string) error {
	if err := os.Mkdir(varFolderName(dbname, vname), 0700); err != nil {
		return nserror.ErrorStack(err)
	}
	return nil
}

func createVarMetaFile(dbname string, v VariableMeta) error {
	fname := varMetaFileName(dbname, v.Name)
	file, err := os.OpenFile(fname, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return nserror.ErrorStack(err)
	}
	defer file.Close()

	meta, err := json.Marshal(v)
	if err != nil {
		return nserror.ErrorStack(err)
	}

	_, err = file.Write(meta)
	if err != nil {
		return nserror.ErrorStack(err)
	}

	return nil
}
