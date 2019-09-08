package lib

import (
	"bufio"
	"os"
)

func encryptFile(sourcePath string, targetPath string, pass string) error {
	keySalt := makeKeySalt()
	key := deriveKey(pass, keySalt)

	inputFile, err := os.Open(sourcePath)

	if err != nil {
		return err
	}

	fi, err := inputFile.Stat()
	fileLen := int(fi.Size())

	if err != nil {
		return err
	}

	outputFile, err := os.Create(targetPath)

	if err != nil {
		removeFile(targetPath)
		return err
	}

	inputReader := bufio.NewReader(inputFile)
	outputWriter := bufio.NewWriter(outputFile)

	writeHeader(outputWriter, keySalt, fileLen)

	blockCount := getBlockCount(fileLen)

	for i := 0; i < blockCount; i++ {
		err := encryptChunkToFile(inputReader, outputWriter, key, plainTextChunkSize)

		if err != nil {
			removeFile(targetPath)
			return err
		}
	}

	deltaChunkSize := getDeltaChunkSize(fileLen)

	if deltaChunkSize > 0 {
		err := encryptChunkToFile(inputReader, outputWriter, key, deltaChunkSize)

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

func encryptChunkToFile(r *bufio.Reader, w *bufio.Writer, key []byte, chunkSize int) error {
	nonce := makeNonce()
	plaintextBuffer := make([]byte, chunkSize)
	_, err := r.Read(plaintextBuffer)

	if err != nil {
		return err
	}

	encryptedBytes, err := encrypt(plaintextBuffer, key, nonce)

	if err != nil {
		return err
	}

	_, err = w.Write(encryptedBytes)

	if err != nil {
		return err
	}

	_, err = w.Write(nonce)

	if err != nil {
		return err
	}

	return nil
}
