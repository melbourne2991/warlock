package lib

import (
	"fmt"
	"io"
	"log"
	"os"
	"testing"

	"github.com/melbourne2991/warlock/testutils"
)

const tmpPath = "../tmp"

const smallInput = `I am test data we all love test data, I am NOT encrypted omg!
This is another line.
And another line
test data is awesome 😊`

const pass = `fakepasswordlol`

func TestFileEncryptDecrypt(t *testing.T) {
	newFileCryptoTestCase(t, "simple", smallInput).run()
	newFileCryptoTestCase(t, "large", generateLargeStr()).run()
}

type testInputBuilder func() func(writer *io.Writer) bool

type fileCryptoTestCase struct {
	testing        *testing.T
	testName       string
	inputFilePath  string
	outputFilePath string
	inputStr       string
}

func newFileCryptoTestCase(testing *testing.T, testName string, inputStr string) *fileCryptoTestCase {
	inputFilePath := testutils.TmpPath(testName + "_in.txt")
	outputFilePath := testutils.TmpPath(testName + "_out.txt")

	testing.Log("Test input file path:", inputFilePath)
	testing.Log("Test output file path:", outputFilePath)

	return &fileCryptoTestCase{
		testing,
		testName,
		inputFilePath,
		outputFilePath,
		inputStr,
	}
}

func (testcase *fileCryptoTestCase) run() {
	testcase.testSetup()
	defer testcase.testTeardown()

	err := encryptFile(testcase.inputFilePath, testcase.outputFilePath, pass)

	if err != nil {
		testcase.testing.Fatal(err)
	}

	testcase.verifyFileEncrypted()

	err = decryptFile(testcase.outputFilePath, testcase.inputFilePath, pass)

	if err != nil {
		testcase.testing.Fatal(err)
	}

	testcase.verifyFileDecrypted()
}

func (testcase *fileCryptoTestCase) verifyFileEncrypted() {
	stat, err := os.Stat(testcase.outputFilePath)

	if err != nil {
		if os.IsNotExist(err) {
			testcase.testing.Errorf("Output file does not exist")
		} else {
			testcase.testing.Error(err)
		}

		return
	}

	if stat.Size() == 0 {
		testcase.testing.Errorf("Output file is empty")
	}

	if _, err := os.Stat(testcase.inputFilePath); os.IsNotExist(err) {
		testcase.testing.Errorf("Should be non destructive op but input file was removed")
	}
}

func (testcase *fileCryptoTestCase) verifyFileDecrypted() {
	file, err := os.Open(testcase.inputFilePath)
	fi, err := file.Stat()

	if err != nil {
		testcase.testing.Fatal(err)
	}

	contents := make([]byte, fi.Size())

	bytesRead, err := file.Read(contents)

	if err != nil {
		testcase.testing.Fatal(err)
	}

	contentsStr := string(contents[:bytesRead])

	if contentsStr != testcase.inputStr {
		testcase.testing.Errorf("Decrypted file contents:\n%s\n\n does not match original input:\n%s\n\n", contentsStr, testcase.inputStr)
	}
}

func (testcase *fileCryptoTestCase) testSetup() {
	fmt.Println("Creating test file")

	f, err := os.Create(testcase.inputFilePath)
	defer f.Close()

	if err != nil {
		testcase.testing.Fatal(err)
	}

	_, err = f.WriteString(testcase.inputStr)

	if err != nil {
		testcase.testing.Fatal(err)
	} else {
		testcase.testing.Log("File created")
	}
}

func (testcase *fileCryptoTestCase) testTeardown() {
	remove(testcase.inputFilePath)  // remove test input
	remove(testcase.outputFilePath) // remove test output
}

func makeDir(dir string) {
	err := os.Mkdir(dir, os.FileMode(int(0755)))
	if err != nil {
		log.Fatal(err)
	}
}

func remove(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		fmt.Println(fmt.Sprintf("Warning - tried to remove non-existing file/directory: \"%s\"", dir))
		return
	}

	err := os.Remove(dir)

	if err != nil {
		log.Fatal(err)
	}
}

func generateLargeStr() string {
	numOfLines := 10000
	var str string

	for i := 0; i < numOfLines; i++ {
		str += "YOU DO THE HOKIE POKIE AND YOU TURN AROUND, THAT'S WHAT IT'S ALL ABOUT........\n"
	}

	return str
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func setup() {
	err := os.Mkdir(testutils.TmpPath(""), 0777)

	if err != nil {
		log.Fatal(err)
	}
}

func shutdown() {
	err := os.RemoveAll(testutils.TmpPath(""))

	if err != nil {
		log.Fatal(err)
	}
}
