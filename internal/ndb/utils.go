package ndb

import (
	"fmt"
	"os"

	"github.com/lincketheo/ndbgo/internal/dtypes"
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

func varCreateDir(db, rel, vari string) error {
	if fname, err := varFolderName(db, rel, vari); err != nil {
		return err
	} else if err = os.Mkdir(fname, 0755); err != nil {
		return err
	} else {
		return nil
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
		return nil
	} else if exists {
		return fmt.Errorf("Meta for variable: %s already exists", vari)
	}

	// Create the file and write the header
	fp, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := fp.Close(); cerr != nil {
			err = cerr
		}
	}()

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

// /////////////////////////////// PUBLIC
func DbExists(db string) (bool, error) {
	if name, err := dbFolderName(db); err != nil {
		return false, err
	} else if exists, err := utils.DirExists(name); err != nil {
		return false, err
	} else {
		return exists, nil
	}
}

func RelExists(db, rel string) (bool, error) {
	if name, err := relFolderName(db, rel); err != nil {
		return false, err
	} else if exists, err := utils.DirExists(name); err != nil {
		return false, err
	} else {
		return exists, nil
	}
}

func VarExists(db, rel, vari string) (bool, error) {
	if name, err := varFolderName(db, rel, vari); err != nil {
		return false, err
	} else if exists, err := utils.DirExists(name); err != nil {
		return false, err
	} else {
		return exists, nil
	}
}
