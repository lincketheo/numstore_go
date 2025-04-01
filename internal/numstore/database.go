package numstore

import (
	"encoding/json"
	"os"

	"github.com/lincketheo/numstore/internal/nserror"
)

type CreateDatabaseArgs struct {
	Name string
}

func (c CreateDatabaseArgs) toDatabase() Database {
	return Database{
		Name: c.Name,
	}
}

type Database struct {
	Name string `json:"name"`
}

/////////////////////////////////// Public

func CreateDatabase(con Connection, args CreateDatabaseArgs) error {
	const op = "Create Database"

	d := args.toDatabase()

	if exists, err := dbExistsAndValid(con, d); err != nil {
		return err
	} else if exists {
		return nserror.DBAlreadyExists
	}

	if err := createDbFolder(con, d); err != nil {
		return err
	}

	if err := createDbMetaFile(con, d); err != nil {
		return err
	}

	return nil
}

/////////////////////////////////// Private

func dbFolderName(con Connection, d Database) string {
	return d.Name
}

func dbMetaFileName(con Connection, d Database) string {
	return dbFolderName(con, d) + "/meta.json"
}

func dbExistsAndValid(con Connection, d Database) (bool, error) {
	_, err := os.Stat(dbFolderName(con, d))
	if err != nil && os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	// TODO - check if valid
	return true, nil
}

func createDbFolder(con Connection, v Database) error {
	os.Mkdir(dbFolderName(con, v), 0700)
	return nil
}

func createDbMetaFile(con Connection, v Database) error {
	fname := dbMetaFileName(con, v)
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
