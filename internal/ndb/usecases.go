package ndb

import (
	"fmt"
	"os"

	"github.com/lincketheo/ndbgo/internal/interpreter"
)

type NDBimpl struct {
	db   string
	rel  string
	vari string

	dbConnected   bool
	relConnected  bool
	variConnected bool
}

func CreateNDBimpl() NDBimpl {
	return NDBimpl{
		db:   "",
		rel:  "",
		vari: "",

		dbConnected:   false,
		relConnected:  false,
		variConnected: false,
	}
}

func (f NDBimpl) CreateDB(db string) error {
	if fname, err := dbFolderName(db); err != nil {
		return err
	} else if err = os.Mkdir(fname, 0755); err != nil {
		return err
	}

	return nil
}

func (f NDBimpl) CreateRel(rel string) error {
	if !f.dbConnected {
		return fmt.Errorf(`Trying to create a Relation,
      but you are not connected to any database`)
	}

	// Assume db folder exists

	if fname, err := relFolderName(f.db, rel); err != nil {
		return err
	} else if err = os.Mkdir(fname, 0755); err != nil {
		return err
	}

	return nil
}

func (f NDBimpl) CreateVar(args interpreter.CreateVarArgs) error {
	if !f.dbConnected {
		return fmt.Errorf(`Trying to create a Relation,
      but you are not connected to any database`)
	} else if !f.relConnected {
		return fmt.Errorf(`Trying to create a Variable,
      but you are not connected to any relation`)
	}

	// Assume db and rel folders exists

	if err := varCreateDir(
		f.db,
		f.rel,
		args.Vari,
	); err != nil {
		return err
	} else if err = varCreateMeta(
		f.db,
		f.rel,
		args.Vari,
		args.Dtype,
		args.Shape,
	); err != nil {
		return err
	}
	return nil
}

func (f *NDBimpl) disconnectDb() {
	f.dbConnected = false
	f.relConnected = false
	f.variConnected = false

	f.db = ""
	f.rel = ""
	f.vari = ""
}

func (f *NDBimpl) ConnectDB(db string) error {
	if exists, err := DbExists(db); err != nil {
		return err
	} else if !exists {
		return fmt.Errorf(`Trying to connect to database: %s,
      but it doesn't exist`, db)
	}

	f.disconnectDb()
	f.db = db
	f.dbConnected = true

	return nil
}

func (f *NDBimpl) disconnectRel() {
	f.relConnected = false
	f.variConnected = false

	f.rel = ""
	f.vari = ""
}

func (f *NDBimpl) ConnectRel(rel string) error {
	if !f.dbConnected {
		return fmt.Errorf(`Trying to connect to a Relation,
      but you are not connected to any database`)
	}

	if exists, err := RelExists(f.db, rel); err != nil {
		return err
	} else if !exists {
		return fmt.Errorf(`Trying to connect to relation: %s,
      but it doesn't exist`, rel)
	}

	f.disconnectRel()
	f.rel = rel
	f.relConnected = true

	return nil
}

func (f *NDBimpl) disconnectVar() {
	f.variConnected = false

	f.vari = ""
}

func (f *NDBimpl) ConnectVar(vari string) error {
	if !f.dbConnected {
		return fmt.Errorf(`Trying to connect to a Variable,
      but you are not connected to any database`)
	} else if !f.relConnected {
		return fmt.Errorf(`Trying to connect to a Variable,
      but you are not connected to any relation`)
	}

	if exists, err := VarExists(f.db, f.rel, vari); err != nil {
		return err
	} else if !exists {
		return fmt.Errorf(`Trying to connect to Variable: %s,
      but it doesn't exist`, vari)
	}

	f.disconnectVar()
	f.vari = vari
	f.variConnected = true

	return nil
}
