package common

import (
	"fmt"
	"os"
)

//func IsPathExists(path string) bool {
//	_, err := os.Stat(path)
//	return err == nil
//}

func IsDirExists(path string) bool {
	stat, err := os.Stat(path)
	return err == nil && stat.IsDir()
}

func IsFileExists(path string) bool {
	stat, err := os.Stat(path)
	return err == nil && !stat.IsDir()
}

func MakeTempDir() (string, error) {
	temp, err := os.MkdirTemp("", "esp-fs-tool-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp dir: %w", err)
	}
	return temp, nil
}
