package lib

import (
	"bufio"
	"encoding/base64"
	"strconv"
)

func writeHeader(w *bufio.Writer, keySalt []byte, fileSize int) error {
	if err := writeFilePart(w, keySalt); err != nil {
		return err
	}

	if err := writeFilePart(w, []byte(strconv.Itoa(fileSize))); err != nil {
		return err
	}

	return nil
}

func consumeHeader(r *bufio.Reader) ([]byte, int, error) {
	keySalt, err := consumeFilePart(r)

	if err != nil {
		return nil, 0, err
	}

	fileSizeBytes, err := consumeFilePart(r)

	if err != nil {
		return nil, 0, err
	}

	fileSize, err := strconv.Atoi(string(fileSizeBytes))

	if err != nil {
		return nil, 0, err
	}

	return keySalt, fileSize, nil
}

const filePartSep = '\n'

func writeFilePart(w *bufio.Writer, part []byte) error {
	_, err := w.WriteString(base64.StdEncoding.EncodeToString(part))

	if err != nil {
		return err
	}

	err = w.WriteByte(filePartSep)

	if err != nil {
		return err
	}

	return nil
}

func consumeFilePart(r *bufio.Reader) ([]byte, error) {
	encodedPart, err := r.ReadBytes(filePartSep)

	if err != nil {
		return nil, err
	}

	decodedPart, err := base64.StdEncoding.DecodeString(string(encodedPart))

	if err != nil {
		return nil, err
	}

	return decodedPart, nil
}
