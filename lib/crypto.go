package lib

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"golang.org/x/crypto/pbkdf2"
	"io"
	"log"
)

func encrypt(chunk []byte, key []byte, nonce []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)

	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)

	if err != nil {
		log.Fatal(err)
	}

	ciphertext := gcm.Seal(nil, nonce, chunk, nil)

	return ciphertext, nil
}

func decrypt(chunk []byte, key []byte, nonce []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)

	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)

	if err != nil {
		return nil, err
	}

	plaintext, err := gcm.Open(nil, nonce, chunk, nil)

	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func makeNonce() []byte {
	// Never use more than 2^32 random nonces with a given key because of the risk of a repeat.
	var nonce = make([]byte, nonceSize)

	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	return nonce
}

func deriveKey(pass string, salt []byte) []byte {
	key := pbkdf2.Key([]byte(pass), salt, iterationCount, 32, sha1.New)

	return key
}

func makeKeySalt() []byte {
	salt := make([]byte, pwSaltBytes)

	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		log.Fatal(err)
	}

	return salt
}
