package lib

import (
	"crypto/md5"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/mholt/archiver"
)

const compressionType = "tar.gz"
const archivePrefix = "arc"
const encryptedPrefix = "enc"

const nonceSize = 12
const gcmTagSize = 16
const plainTextChunkSize = 1024
const pwSaltBytes = 32
const iterationCount = 150000

var encryptedChunkSize = getEncryptedSize(plainTextChunkSize)

// ErrPathAlreadyLocked when the path is already locked
var ErrPathAlreadyLocked = errors.New("path is already locked")

// ErrEncryptedFileMissing when the encrypted file provided does not exist
var ErrEncryptedFileMissing = errors.New("encrypted file is missing")

// ErrPathDoesNotExist when the source path does not exist
var ErrPathDoesNotExist = errors.New("directory or file does not exist")

// NewPathLocker instantiates a new path locker
func NewPathLocker(storePath string) *PathLocker {
	return &PathLocker{storePath}
}

// PathLocker stores locked sources
type PathLocker struct {
	storePath string
}

// LockPath compresses and encrypts specified path to the store
func (p PathLocker) LockPath(dirPath string, pass string) error {
	if err := p.validateLocking(dirPath); err != nil {
		return err
	}

	fmt.Println("Archiving...")
	if err := p.archiveSource(dirPath); err != nil {
		return err
	}

	fmt.Println("Encrypting...")
	if err := encryptFile(
		p.getArchivedPath(dirPath),
		p.getEncryptedPath(dirPath),
		pass,
	); err != nil {
		return err
	}

	fmt.Println("Removing temporary archive...")
	// Remove temporary archive
	if err := os.RemoveAll(p.getArchivedPath(dirPath)); err != nil {
		return err
	}

	fmt.Println("Removing original files...")
	// Remove everything at the original path
	if err := os.RemoveAll(dirPath); err != nil {
		return err
	}

	return nil
}

// UnlockPath decrypts and decompresses the archive to the source
func (p PathLocker) UnlockPath(dirPath string, pass string) error {
	fmt.Println("Decrypting...")
	if err := decryptFile(
		p.getEncryptedPath(dirPath),
		p.getArchivedPath(dirPath),
		pass,
	); err != nil {
		return err
	}

	fmt.Println("Unarchiving...")
	if err := p.unarchiveSource(dirPath); err != nil {
		return err
	}

	fmt.Println("Removing temporary archive...")
	// Remove temporary archive
	if err := os.RemoveAll(p.getArchivedPath(dirPath)); err != nil {
		return err
	}

	fmt.Println("Removing encrypted files...")
	// Remove encrypted files
	if err := os.RemoveAll(p.getEncryptedPath(dirPath)); err != nil {
		return err
	}

	return nil
}

func (p PathLocker) validateLocking(dirPath string) error {
	if !pathExists(dirPath) {
		return ErrPathDoesNotExist
	}

	if pathExists(p.getEncryptedPath(dirPath)) {
		return ErrPathAlreadyLocked
	}

	return nil
}

func (p PathLocker) validateUnlocking(dirPath string) error {
	if !pathExists(p.getEncryptedPath(dirPath)) {
		return ErrEncryptedFileMissing
	}

	return nil
}

func pathExists(dirPath string) bool {
	_, err := os.Stat(dirPath)

	if err != nil {
		if os.IsNotExist(err) {
			return false
		}

		log.Fatal(err)
	}

	return true
}

func (p PathLocker) archiveSource(sourcePath string) error {
	return archiver.Archive([]string{sourcePath}, p.getArchivedPath(sourcePath))
}

func (p PathLocker) unarchiveSource(sourcePath string) error {
	return archiver.Unarchive(p.getArchivedPath(sourcePath), filepath.Dir(sourcePath))
}

func (p PathLocker) getArchivedPath(sourcePath string) string {
	fileName := hashNameWithPrefix(archivePrefix, sourcePath)
	fileNameWithExt := fmt.Sprintf("%s.%s", fileName, compressionType)

	return path.Join(p.storePath, fileNameWithExt)
}

func (p PathLocker) getEncryptedPath(sourcePath string) string {
	return path.Join(p.storePath, hashNameWithPrefix(encryptedPrefix, sourcePath))
}

func hashNameWithPrefix(prefix string, val string) string {
	return fmt.Sprintf("%s_%x", prefix, md5.Sum([]byte(val)))
}
