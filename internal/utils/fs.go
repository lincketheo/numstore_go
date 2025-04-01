package utils

import "os"

func FileExists(name string) (os.FileInfo, error) {
	if info, err := os.Stat(name); err != nil && os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return info, nil
	}
}
