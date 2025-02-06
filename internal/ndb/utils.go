package ndb

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/lincketheo/ndbgo/internal/dtypes"
	"github.com/lincketheo/ndbgo/internal/logging"
	"github.com/lincketheo/ndbgo/internal/utils"
)

const noRel = "norel"

// /////////////////////////////// ERROR Wrappers
func expectDbExistance(db string, shouldExist bool) error {
	if exists, err := dbExists(db); err != nil {
		return utils.ErrorMoref(err,
			"Failed to check if Database: %s exists", db)
	} else if exists && !shouldExist {
		return fmt.Errorf("Database: %s should not exist\n", db)
	} else if !exists && shouldExist {
		return fmt.Errorf("Database: %s should exist\n", db)
	}

	return nil
}

func (n NDBimpl) expectDbConnection(shouldBeConnected bool) error {
	if !n.dbConnected && shouldBeConnected {
		return fmt.Errorf("Must be connected to a database")
	} else if n.dbConnected && !shouldBeConnected {
		return fmt.Errorf("Must not be connected to a database")
	}

	return nil
}

func createDBFolder(db string) error {
	dname := dbFolderName(db)
	if err := os.Mkdir(dname, 0755); err != nil {
		return utils.ErrorMoref(err,
			"Failed to create directory: %s for database: %s", dname, db)
	}

	return nil
}

func createDBNoRelFolder(db string) error {
	dname := noRelFolderName(db, noRel)
	if err := os.Mkdir(dname, 0755); err != nil {
		return utils.ErrorMoref(err,
			"Failed to create directory: %s for database: %s", dname, db)
	}
	return nil
}

func (n NDBimpl) expectRelExistance(rel string, shouldExist bool) error {
	utils.ASSERT(n.dbConnected)

	if exists, err := n.relExists(rel); err != nil {
		return utils.ErrorMoref(err,
			"Failed to check if Relation: %s exists", rel)
	} else if exists && !shouldExist {
		return fmt.Errorf("Relation: %s should not exist\n", rel)
	} else if !exists && shouldExist {
		return fmt.Errorf("Relation: %s should exist\n", rel)
	}

	return nil
}

func (n NDBimpl) expectRelConnection(shouldBeConnected bool) error {
	if !n.relConnected && shouldBeConnected {
		return fmt.Errorf("Must be connected to a relation")
	} else if n.dbConnected && !shouldBeConnected {
		return fmt.Errorf("Must not be connected to a relation")
	}

	return nil
}

func (n NDBimpl) createRelFolder(rel string) error {
	utils.ASSERT(n.dbConnected)
	dname := n.relFolderName(rel)

	if err := os.Mkdir(dname, 0755); err != nil {
		return utils.ErrorMoref(err,
			"Failed to create directory: %s for relation: %s", dname, rel)
	}
	return nil
}

func (n NDBimpl) expectVarExistance(
	rel string,
	vari string,
	shouldExist bool) error {

	utils.ASSERT(n.dbConnected)

	if exists, err := n.varExists(rel, vari); err != nil {
		return utils.ErrorMoref(err,
			"Failed to check if Variable: %s exists", rel)
	} else if exists && !shouldExist {
		return fmt.Errorf("Relation: %s should not exist\n", rel)
	} else if !exists && shouldExist {
		return fmt.Errorf("Relation: %s should exist\n", rel)
	}

	return nil
}

func (n NDBimpl) createVariFolder(rel, vari string) error {
	utils.ASSERT(n.dbConnected)
	dname := n.varFolderName(rel, vari)

	if err := os.Mkdir(dname, 0755); err != nil {
		return utils.ErrorMoref(err,
			"Failed to create directory: %s for variable: %s", dname, rel)
	}
	return nil
}

// /////////////////////////////// PRIVATE
func dbFolderName(db string) string {
	return filepath.Join(".", db)
}

func noRelFolderName(db, rel string) string {
	return filepath.Join(dbFolderName(db), rel)
}

func (n NDBimpl) relFolderName(rel string) string {
	utils.ASSERT(n.dbConnected)
	return filepath.Join(dbFolderName(n.db), rel)
}

func (n NDBimpl) varFolderName(rel, vari string) string {
	utils.ASSERT(n.dbConnected)
	return filepath.Join(n.relFolderName(rel), vari)
}

func (n NDBimpl) varMetaFileName(rel, vari string) string {
	utils.ASSERT(n.dbConnected)
	return filepath.Join(n.varFolderName(rel, vari), "meta")
}

func (n NDBimpl) varMetaFileExists(rel, vari string) (bool, error) {
	utils.ASSERT(n.dbConnected)
	fname := n.varMetaFileName(rel, vari)

	if exists, err := utils.FileExists(fname); err != nil {
		return false, utils.ErrorContext(err)
	} else {
		return exists, utils.ErrorContext(err)
	}
}

func (n NDBimpl) varCreateMeta(
	rel,
	vari string,
	dtype dtypes.Dtype,
	shape []uint32,
) error {
	var fname string = n.varMetaFileName(rel, vari)

	// Check if shape is valid
	if !utils.CanIntBeByte(len(shape)) {
		return fmt.Errorf("Shape is too long: %v", shape)
	}

	// Create the file and write the header
	fp, err := os.Create(fname)
	if err != nil {
		return utils.ErrorContext(err)
	}
	defer func() {
		if err := fp.Close(); err != nil {
			logging.Warn("Failed to close file: %v. Cause: %v\n", fp, err)
		}
	}()

	// Write DTYPE
	if _, err := fp.Write([]byte{byte(dtype)}); err != nil {
		return err
	}

	// Write SHAPE LEN
	if _, err := fp.Write([]byte{byte(len(shape))}); err != nil {
		return err
	}

	// Write SHAPE
	sb := utils.UInt32ArrBytes(shape)
	if _, err := fp.Write(sb); err != nil {
		return err
	}

	return err
}

func dbExists(db string) (bool, error) {
	name := dbFolderName(db)

	if exists, err := utils.DirExists(name); err != nil {
		return false, err
	} else {
		return exists, nil
	}
}

func (n NDBimpl) relExists(rel string) (bool, error) {
	name := n.relFolderName(rel)

	if exists, err := utils.DirExists(name); err != nil {
		return false, err
	} else {
		return exists, nil
	}
}

func (n NDBimpl) varExists(rel, vari string) (bool, error) {
	name := n.varFolderName(rel, vari)

	if exists, err := utils.DirExists(name); err != nil {
		return false, err
	} else {
		return exists, nil
	}
}

// ///////////////////////////////// Connections
func (f *NDBimpl) disconnectDb() {
	f.dbConnected = false
	f.relConnected = false

	f.db = ""
	f.rel = ""
}

func (f *NDBimpl) disconnectRel() {
	f.relConnected = false

	f.rel = ""
}

// ///////////////////////////////// Utils
func (f *NDBimpl) isDisconnected() bool {
	return !f.dbConnected && !f.relConnected
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
}

func parseRelVarStr(rvstr string) (string, string, bool) {
	parts := strings.Split(rvstr, ":")
	switch len(parts) {
	case 0:
		panic("Unreachable")
	case 1:
		return "", "", false
	case 2:
		return parts[0], parts[1], true
	default:
		return "", "", false
	}
}
