package utils

import (
	"fmt"
	"io"
	"os"
)

func FileExists(name string) (os.FileInfo, error) {
	if info, err := os.Stat(name); err != nil && os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return info, nil
	}
}

func ReadFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening file %q: %w", path, err)
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("reading file %q: %w", path, err)
	}
	return data, nil
}
