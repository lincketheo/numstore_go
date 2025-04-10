package numstore

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/lincketheo/numstore/internal/nserror"
)

type Database struct {
	Name string `json:"name"`
}

/////////////////////////////////// Public

func CreateDatabase(name string) error {
	const op = "Create Database"

	d := Database{
		Name: name,
	}

	if exists, err := dbExistsAndValid(d); err != nil {
		return nserror.ErrorStack(err)
	} else if exists {
		return fmt.Errorf("Database: %s already exists", name)
	}

	if err := createDbFolder(d); err != nil {
		return nserror.ErrorStack(err)
	}

	if err := createDbMetaFile(d); err != nil {
		return nserror.ErrorStack(err)
	}

	return nil
}

/////////////////////////////////// Private

func dbFolderName(d Database) string {
	return d.Name
}

func dbMetaFileName(d Database) string {
	return dbFolderName(d) + "/meta.json"
}

func dbExistsAndValid(d Database) (bool, error) {
	_, err := os.Stat(dbFolderName(d))
	if err != nil && os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, nserror.ErrorStack(err)
	}

	// TODO - check if valid
	return true, nil
}

func createDbFolder(v Database) error {
	os.Mkdir(dbFolderName(v), 0700)
	return nil
}

func createDbMetaFile(v Database) error {
	fname := dbMetaFileName(v)
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
