package ndb

import (
	"fmt"
	"os"

	"github.com/lincketheo/ndbgo/internal/dtypes"
	"github.com/lincketheo/ndbgo/internal/logging"
	"github.com/lincketheo/ndbgo/internal/utils"
)

// /////////////////////////////// PRIVATE
// -------------------------------- DB Utils
func dbFolderName(db string) (string, error) {
	return utils.CombinePath(db)
}

// -------------------------------- Rel Utils
func relFolderName(db, rel string) (string, error) {
	return utils.CombinePath(db, rel)
}

// -------------------------------- Variable Utils
func varFolderName(db, rel, vari string) (string, error) {
	return utils.CombinePath(db, rel, vari)
}

func varMetaName(db, rel, vari string) (string, error) {
	return utils.CombinePath(db, rel, vari, "meta")
}

func varMetaExists(db, rel, vari string) (bool, error) {
	if fname, err := varMetaName(db, rel, vari); err != nil {
		return false, err
	} else if exists, err := utils.FileExists(fname); err != nil {
		return false, err
	} else {
		return exists, err
	}
}

func varCreateMeta(
	db, rel, vari string,
	dtype dtypes.Dtype,
	shape []uint32,
) error {

	var fname string
	var err error

	// Check if it already exists
	if exists, err := varMetaExists(db, rel, vari); err != nil {
		return utils.ErrorContext(err)
	} else if exists {
		return fmt.Errorf("Meta for variable: %s already exists", vari)
	}

	fname, err = varMetaName(db, rel, vari)
	if err != nil {
		return err
	}

	// Create the file and write the header
	fp, err := os.Create(fname)
	if err != nil {
		return utils.ErrorContext(err)
	}
	defer fp.Close()

	if err = varWriteHeader(fp, dtype, shape); err != nil {
		return err
	}

	return err
}

func varWriteHeader(
	fp *os.File,
	dtype dtypes.Dtype,
	shape []uint32,
) error {
	if !utils.CanIntBeByte(len(shape)) {
		return fmt.Errorf("Shape is too long: %v", shape)
	}

	// DTYPE
	if _, err := fp.Write([]byte{byte(dtype)}); err != nil {
		return err
	}

	// SHAPE LEN
	if _, err := fp.Write([]byte{byte(len(shape))}); err != nil {
		return err
	}

	// SHAPE
	sb := utils.UInt32ArrBytes(shape)
	if _, err := fp.Write(sb); err != nil {
		return err
	}

	return nil
}

func dbExists(db string) (bool, error) {
	if name, err := dbFolderName(db); err != nil {
		return false, err
	} else if exists, err := utils.DirExists(name); err != nil {
		return false, err
	} else {
		return exists, nil
	}
}

func relExists(db, rel string) (bool, error) {
	if name, err := relFolderName(db, rel); err != nil {
		return false, err
	} else if exists, err := utils.DirExists(name); err != nil {
		return false, err
	} else {
		return exists, nil
	}
}

func varExists(db, rel, vari string) (bool, error) {
	if name, err := varFolderName(db, rel, vari); err != nil {
		return false, err
	} else if exists, err := utils.DirExists(name); err != nil {
		return false, err
	} else {
		return exists, nil
	}
}

// ///////////////////////////////// Connections
func (f *NDBimpl) disconnectDb() {
	f.dbConnected = false
	f.relConnected = false
	f.variConnected = false

	f.db = ""
	f.rel = ""
	f.vari = ""
}

func (f *NDBimpl) disconnectRel() {
	f.relConnected = false
	f.variConnected = false

	f.rel = ""
	f.vari = ""
}

func (f *NDBimpl) disconnectVar() {
	f.variConnected = false

	f.vari = ""
}

// ///////////////////////////////// Utils
func (f *NDBimpl) isDisconnected() bool {
	return !f.dbConnected && !f.variConnected && !f.relConnected
}

func (f *NDBimpl) logConnectionState() {
	logging.Info("Connection state:")
	if f.isDisconnected() {
		logging.Info("  Disconnected")
	}
	if f.dbConnected {
		logging.Info("  Database: %s\n", f.db)
	}
	if f.relConnected {
		logging.Info("  Relation: %s\n", f.rel)
	}
	if f.variConnected {
		logging.Info("  Variable: %s\n", f.vari)
	}
}
