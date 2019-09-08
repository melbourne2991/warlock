package lib

import (
	"os"
)

func openFile(filePath string) (*os.File, error) {
	file, err := os.Open(filePath)

	if err != nil {
		return nil, err
	}

	return file, nil
}

func createFile(filePath string) (*os.File, error) {
	file, err := os.Create(filePath)

	if err != nil {
		return nil, err
	}

	return file, nil
}

func removeFile(filePath string) error {
	err := os.Remove(filePath)

	if err != nil {
		return err
	}

	return nil
}

func getDeltaChunkSize(fileLength int) int {
	return (fileLength % plainTextChunkSize)
}

func getBlockCount(fileLength int) int {
	return (fileLength - getDeltaChunkSize(fileLength)) / plainTextChunkSize
}

func getEncryptedSize(inputByteLength int) int {
	return inputByteLength + gcmTagSize
}
