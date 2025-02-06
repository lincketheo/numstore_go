package ndb

import (
	"errors"
	"os"

	"github.com/lincketheo/ndbgo/internal/usecases"
	"github.com/lincketheo/ndbgo/internal/utils"
)

type NDBimpl struct {
	db  string
	rel string

	dbConnected  bool
	relConnected bool
}

// Constructor
func CreateNDBimpl() NDBimpl {
	return NDBimpl{
		db:  "",
		rel: "",

		dbConnected:  false,
		relConnected: false,
	}
}

// TODO - check validity of names

func (f NDBimpl) CreateDB(db string) error {
	if err := expectDbExistance(db, false); err != nil {
		return err
	}
	if err := createDBFolder(db); err != nil {
		return err
	}
	if err := createDBNoRelFolder(db); err != nil {
		return err
	}
	return nil
}

func (f NDBimpl) CreateRel(rel string) error {
	if err := f.expectDbConnection(true); err != nil {
		return err
	}
	if err := f.expectRelExistance(rel, true); err != nil {
		return err
	}
	if err := f.createRelFolder(rel); err != nil {
		return err
	}
	return nil
}

func (f NDBimpl) CreateVar(variStr string, config usecases.VarConfig) error {
	if err := f.expectDbConnection(true); err != nil {
		return err
	}

	var rel string = f.rel
	var vari string = variStr
	var err error

	// If not connected to rel
	if !f.relConnected {

		// Check if the format is rel:vari
		if cvari, crel, ok := parseRelVarStr(vari); ok {
			vari = cvari
			rel = crel

			// Connect to relationship temporarily
			if err = f.ConnectRel(rel); err != nil {
				return utils.ErrorMoref(err,
					`Failed to connect to relation: %s
          while trying to create var: %s`, rel, vari)
			}

			// Disconnect from relationship
			defer func() {
				if cerr := f.DisconnectRel(); cerr != nil {
					err = cerr
				}
			}()

			// Otherwise use no rel
		} else {
			rel = noRel
		}
	}

	if err := f.expectVarExistance(rel, vari, true); err != nil {
		return err
	}
	if err := f.createVariFolder(rel, vari); err != nil {
		return err
	}
	if err := f.varCreateMeta(rel, vari, config.Dtype, config.Shape); err != nil {
		return err
	}

	return nil
}

func (f *NDBimpl) ConnectDB(db string) error {
	if err := expectDbExistance(db, true); err != nil {
		return err
	}

	f.disconnectDb()
	f.db = db
	f.dbConnected = true
	f.logConnectionState()

	return nil
}

func (f *NDBimpl) ConnectRel(rel string) error {
	if err := f.expectDbConnection(true); err != nil {
		return err
	}
	if err := f.expectRelExistance(rel, true); err != nil {
		return err
	}

	f.disconnectRel()
	f.rel = rel
	f.relConnected = true
	f.logConnectionState()

	return nil
}

func (f *NDBimpl) DisconnectDB() error {
	if err := f.expectDbConnection(true); err != nil {
		return err
	}
	f.disconnectDb()
	return nil
}

func (f *NDBimpl) DisconnectRel() error {
	if err := f.expectRelConnection(false); err != nil {
		return err
	}
	f.disconnectRel()
	return nil
}

func (f *NDBimpl) DeleteDB(db string) error {
	if err := f.expectDbConnection(false); err != nil {
		return err
	}
	if err := expectDbExistance(db, true); err != nil {
		return err
	}
	return os.RemoveAll(dbFolderName(db))
}

func (f *NDBimpl) DeleteRel(rel string) error {
	if err := f.expectRelConnection(false); err != nil {
		return err
	}
	if err := f.expectDbConnection(true); err != nil {
		return err
	}
	if err := f.expectRelExistance(rel, true); err != nil {
		return err
	}
	return os.RemoveAll(f.relFolderName(rel))
}

func (f *NDBimpl) DeleteVar(variStr string) error {
	if err := f.expectDbConnection(true); err != nil {
		return err
	}

	var rel string = f.rel
	var vari string = variStr
	var err error

	// If not connected to rel
	if !f.relConnected {

		// Check if the format is rel:vari
		if cvari, crel, ok := parseRelVarStr(vari); ok {
			vari = cvari
			rel = crel

			// Connect to relationship temporarily
			if err = f.ConnectRel(rel); err != nil {
				return utils.ErrorMoref(err,
					`Failed to connect to relation: %s
          while trying to create var: %s`, rel, vari)
			}

			// Disconnect from relationship
			defer func() {
				if cerr := f.DisconnectRel(); cerr != nil {
					err = cerr
				}
			}()

			// Otherwise use no rel
		} else {
			rel = noRel
		}
	}

	if err := f.expectVarExistance(rel, vari, true); err != nil {
		return err
	}

	return os.RemoveAll(f.varFolderName(rel, vari))
}

func (n *NDBimpl) AddReader(config usecases.ReaderConfig) (int, error) {
	return 0, errors.New("Unimplemented")
}

func (n *NDBimpl) RemoveReader(id int) error {
	return errors.New("Unimplemented")
}

func (n *NDBimpl) AddWriter(config usecases.WriterConfig) (int, error) {
	return 0, errors.New("Unimplemented")
}

func (n *NDBimpl) RemoveWriter(id int) error {
	return errors.New("Unimplemented")
}

func (n *NDBimpl) Write(args usecases.WriteArgs) error {
	return errors.New("Unimplemented")
}

func (n *NDBimpl) Read(args usecases.ReadArgs) error {
	return errors.New("Unimplemented")
}
