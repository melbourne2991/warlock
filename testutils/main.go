package testutils

import (
	"log"
	"os"
	"path"
)

// ProjectDir returns path to current project
func ProjectDir() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

// TmpPath gives project tmp dir path
func TmpPath(filePath string) string {
	return path.Join(toAbsPath("tmp"), filePath)
}

func toAbsPath(relativePath string) string {
	return path.Join(ProjectDir(), relativePath)
}

// MkdirInTmp creates a directory in tmp folder
func MkdirInTmp(filePath string) {
	err := os.Mkdir(TmpPath(filePath), 0777)

	if err != nil {
		log.Fatal(err)
	}
}
