package ndb

import (
	"fmt"
	"os"

	"github.com/lincketheo/ndbgo/internal/logging"
	"github.com/lincketheo/ndbgo/internal/usecases"
	"github.com/lincketheo/ndbgo/internal/utils"
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
	// Check if db already exists

	// Check if relation already exists
	if exists, err := dbExists(db); err != nil {
		return utils.ErrorContext(err)
	} else if exists {
		return fmt.Errorf("Database: %s already exists", db)

		// Create database folder
	} else if fname, err := dbFolderName(db); err != nil {
		return utils.ErrorContext(err)
	} else if err = os.Mkdir(fname, 0755); err != nil {
		return utils.ErrorContext(err)

		// Log and return
	} else {
		logging.Info("Created db in: %s\n", fname)
		return nil
	}
}

func (f NDBimpl) CreateRel(rel string) error {
	if !f.dbConnected {
		return fmt.Errorf(`Trying to create a Relation,
      but you are not connected to any database`)
	}

	// Check if relation already exists
	if exists, err := relExists(f.db, rel); err != nil {
		return utils.ErrorContext(err)
	} else if exists {
		return fmt.Errorf("Relation: %s already exists", rel)

		// Create relation folder
	} else if fname, err := relFolderName(f.db, rel); err != nil {
		return utils.ErrorContext(err)
	} else if err = os.Mkdir(fname, 0755); err != nil {
		return utils.ErrorContext(err)

		// Log and return
	} else {
		logging.Info("Created rel in: %s\n", fname)
		return nil
	}
}

func (f NDBimpl) CreateVar(
	vari string,
	config usecases.VarConfig,
) error {
	if !f.dbConnected {
		return fmt.Errorf(`Trying to create a Relation,
      but you are not connected to any database`)
	} else if !f.relConnected {
		return fmt.Errorf(`Trying to create a Variable,
      but you are not connected to any relation`)
	}

	// Check if variable already exists
	if exists, err := varExists(f.db, f.rel, vari); err != nil {
		return utils.ErrorContext(err)
	} else if exists {
		return fmt.Errorf("Variable: %s already exists", vari)

		// Create the variable folder
	} else if dir, err := varFolderName(f.db, f.rel, vari); err != nil {
		return utils.ErrorContext(err)
	} else if err = os.Mkdir(dir, 0755); err != nil {
		return utils.ErrorContext(err)

		// Create Meta File
	} else if meta, err := varMetaName(f.db, f.rel, vari); err != nil {
		return utils.ErrorContext(err)
	} else {
		// Create meta file
		fp, err := os.Create(meta)
		if err != nil {
			return utils.ErrorContext(err)
		}
		defer fp.Close()

		// Write the header
		if err = varWriteHeader(fp, config.Dtype, config.Shape); err != nil {
			return utils.ErrorContext(err)
		}

		// Log and return
		logging.Info("Created db in: %s\n", dir)
		logging.Info("Created db meta in: %s\n", meta)
		return nil

	}
}

func (f *NDBimpl) ConnectDB(db string) error {
	if exists, err := dbExists(db); err != nil {
		return utils.ErrorContext(err)
	} else if !exists {
		return fmt.Errorf(`Trying to connect to database: %s,
      but it doesn't exist`, db)
	}

	f.disconnectDb()
	f.db = db
	f.dbConnected = true
	f.logConnectionState()

	return nil
}

func (f *NDBimpl) ConnectRel(rel string) error {
	if !f.dbConnected {
		return fmt.Errorf(`Trying to connect to a Relation,
      but you are not connected to any database`)
	}

	if exists, err := relExists(f.db, rel); err != nil {
		return utils.ErrorContext(err)
	} else if !exists {
		return fmt.Errorf(`Trying to connect to relation: %s,
      but it doesn't exist`, rel)
	}

	f.disconnectRel()
	f.rel = rel
	f.relConnected = true
	f.logConnectionState()

	return nil
}

func (f *NDBimpl) ConnectVar(vari string) error {
	if !f.dbConnected {
		return fmt.Errorf(`Trying to connect to a Variable,
      but you are not connected to any database`)
	} else if !f.relConnected {
		return fmt.Errorf(`Trying to connect to a Variable,
      but you are not connected to any relation`)
	}

	if exists, err := varExists(f.db, f.rel, vari); err != nil {
		return utils.ErrorContext(err)
	} else if !exists {
		return fmt.Errorf(`Trying to connect to Variable: %s,
      but it doesn't exist`, vari)
	}

	f.disconnectVar()
	f.vari = vari
	f.variConnected = true
	f.logConnectionState()

	return nil
}
