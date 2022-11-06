package pack

import (
	"fmt"
	copy2 "github.com/otiai10/copy"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func MakeTargetDir(sourceDir string, targetDir string, ignoreDotfile bool, ignorePattern string) error {
	ignoreRegexp, err := regexp.Compile(ignorePattern)
	if err != nil {
		return fmt.Errorf("failed to compile ignore pattern: %w", err)
	}

	err = copy2.Copy(sourceDir, targetDir, copy2.Options{
		Skip: func(file string) (bool, error) {
			if ignoreDotfile && strings.HasPrefix(file, ".") {
				log.Printf("  ignored: %s\n", file)
				return true, nil
			}

			match := ignoreRegexp.MatchString(file)
			if match {
				log.Printf("  ignored: %s\n", file)
				return true, nil
			}
			return false, err
		},
	})

	if err != nil {
		return fmt.Errorf("failed to copy files: %w", err)
	}
	return nil
}

func PrintDirFiles(dir string) {
	log.Printf("%s\n", dir)
	_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		relPath, _ := filepath.Rel(dir, path)
		log.Printf("  /%s: %.2f kB\n", relPath, float32(info.Size())/1024.0)
		return nil
	})
}
