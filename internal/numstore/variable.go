package numstore

import (
	"encoding/json"
	"io"
	"os"

	"github.com/lincketheo/numstore/internal/nserror"
	"github.com/lincketheo/numstore/internal/utils"
)

type CreateVariableArgs struct {
	Name  string
	Dtype Dtype
	Shape []uint32
}

func (c CreateVariableArgs) toVariable() Variable {
	return Variable{
		Name:  c.Name,
		Dtype: c.Dtype,
		Shape: c.Shape,
	}
}

type Variable struct {
	Name  string   `json:"name"`
	Dtype Dtype    `json:"Dtype"`
	Shape []uint32 `json:"shape"`
}

type WritingVariable struct {
	vfd        *os.File
	tfd        *os.File
	dataBuffer []byte
	timeBuffer []byte
}

/////////////////////////////////// Public

func CreateVariable(con Connection, args CreateVariableArgs) error {
	if !con.IsConnectedToDB() {
		return nserror.NotConnectedToDB
	}

	v := args.toVariable()

	// Check if variable exists
	if exists, err := varExistsAndValid(con, v.Name); err != nil {
		return err
	} else if exists {
		return nserror.VarAlreadyExists
	}

	if err := createVarFolder(con, v.Name); err != nil {
		return err
	}

	if err := createVarMetaFile(con, v); err != nil {
		return err
	}

	return nil
}

func LoadVariable(con Connection, name string) (Variable, error) {
	if !con.IsConnectedToDB() {
		return Variable{}, nserror.NotConnectedToDB
	}
	return loadVarMetaFile(con, name)
}

func (v Variable) Open(con Connection) (WritingVariable, error) {
	dfd, err := os.OpenFile(varDataFileName(con, v.Name), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return WritingVariable{}, err
	}

	tfd, err := os.OpenFile(varTimeFileName(con, v.Name), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		dfd.Close()
		return WritingVariable{}, err
	}

	dataSize := utils.ReduceMultU32(v.Shape) * dtypeSizeof(v.Dtype)
	timeSize := dtypeSizeof(U32)

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

func (v WritingVariable) WriteNext(r io.Reader) error {
	if n, err := r.Read(v.dataBuffer); err != nil {
		return err
	} else if n != len(v.dataBuffer) {
		panic("Invalid read")
	}

	if n, err := v.vfd.Write(v.dataBuffer); err != nil {
		return err
	} else if n != len(v.dataBuffer) {
		panic("Invalid write")
	}

	return nil
}

func (v WritingVariable) Close() (error, error) {
	return v.tfd.Close(), v.vfd.Close()
}

/////////////////////////////////// Private

func varFolderName(con Connection, vname string) string {
	utils.Assert(con.IsConnectedToDB())
	return con.DbName + "/" + vname
}

func varMetaFileName(con Connection, vname string) string {
	return varFolderName(con, vname) + "/meta.json"
}

func varDataFileName(con Connection, vname string) string {
	return varFolderName(con, vname) + "/data.bin"
}

func varTimeFileName(con Connection, vname string) string {
	return varFolderName(con, vname) + "/time.bin"
}

func varExistsAndValid(con Connection, v string) (bool, error) {
	if stat, err := os.Stat(varFolderName(con, v)); err != nil && os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		// Todo - check validity of stat
		return stat != nil, nil
	}
}

func createVarFolder(con Connection, vname string) error {
	if err := os.Mkdir(varFolderName(con, vname), 0700); err != nil {
		return err
	}
	return nil
}

func createVarMetaFile(con Connection, v Variable) error {
	fname := varMetaFileName(con, v.Name)
	file, err := os.OpenFile(fname, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	meta, err := json.Marshal(v)
	if err != nil {
		return err
	}

	_, err = file.Write(meta)
	if err != nil {
		return err
	}

	return nil
}

func loadVarMetaFile(con Connection, vname string) (Variable, error) {
	fname := varMetaFileName(con, vname)
	data, err := os.ReadFile(fname)
	if err != nil {
		return Variable{}, err
	}

	var v Variable
	if err := json.Unmarshal(data, &v); err != nil {
		return Variable{}, err
	}

	return v, nil
}
