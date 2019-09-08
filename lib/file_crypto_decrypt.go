package lib

import (
	"bufio"
	"io"
)

func decryptFile(sourcePath string, targetPath string, pass string) error {
	inputFile, err := openFile(sourcePath)

	if err != nil {
		return err
	}

	outputFile, err := createFile(targetPath)

	if err != nil {
		return err
	}

	inputReader := bufio.NewReader(inputFile)
	outputWriter := bufio.NewWriter(outputFile)

	keySalt, fileLen, err := consumeHeader(inputReader)

	if err != nil {
		removeFile(targetPath)
		return err
	}

	key := deriveKey(pass, keySalt)
	blockCount := getBlockCount(fileLen)

	for i := 0; i < blockCount; i++ {
		err := decryptChunkToFile(inputReader, outputWriter, key, encryptedChunkSize)

		if err != nil {
			removeFile(targetPath)
			return err
		}
	}

	deltaChunkSize := getDeltaChunkSize(fileLen)

	if deltaChunkSize > 0 {
		err := decryptChunkToFile(inputReader, outputWriter, key, getEncryptedSize(deltaChunkSize))

		if err != nil {
			removeFile(targetPath)
			return err
		}
	}

	outputWriter.Flush()

	if err != nil {
		removeFile(targetPath)
		return err
	}

	return nil
}

func decryptChunkToFile(r *bufio.Reader, w *bufio.Writer, key []byte, chunkSize int) error {
	encryptedBytes := make([]byte, chunkSize)
	_, err := io.ReadFull(r, encryptedBytes)

	if err != nil {
		return err
	}

	nonceBytes := make([]byte, nonceSize)

	_, err = io.ReadFull(r, nonceBytes)

	if err != nil {
		return err
	}

	decryptedBytes, err := decrypt(encryptedBytes, key, nonceBytes)

	if err != nil {
		return err
	}

	_, err = w.Write(decryptedBytes)

	if err != nil {
		return err
	}

	return nil
}
