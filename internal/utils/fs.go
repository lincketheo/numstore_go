package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func CombinePath(start string, next ...string) (string, error) {
	if start == "" {
		return "", fmt.Errorf("start path cannot be empty")
	}

	start = filepath.Clean(start)

	if !filepath.IsAbs(start) {
		abs, err := filepath.Abs(start)
		if err != nil {
			return "", err
		}
		start = abs
	}

	return filepath.Join(start, filepath.Join(next...)), nil
}

func DirExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return info.IsDir(), nil
}

func FileExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return !info.IsDir(), nil
}
