package common

import (
	"fmt"
	"os"
)

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func DirExists(path string) bool {
	stat, err := os.Stat(path)
	return err == nil && stat.IsDir()
}

func FileExists(path string) bool {
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
