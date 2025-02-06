package utils

import (
	"os"
)

func DirExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, ErrorContext(err)
	}
	return info.IsDir(), nil
}

func FileExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, ErrorContext(err)
	}
	return !info.IsDir(), nil
}
