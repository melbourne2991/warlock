package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/melbourne2991/warlock/testutils"
)

func TestDisplaysHelp(t *testing.T) {
	cmd := warlockCommand("")
	out, err := cmd.Output()

	if err != nil {
		t.Error(err)
		return
	}

	if !strings.Contains(string(out), "Usage") {
		t.Fatalf("String should have included \"Usage\"")
	}
}

func warlockCommand(cmd string) *exec.Cmd {
	return exec.Command("./warlock", "--dir", testutils.TmpPath(".warlock"), cmd)
}

func execCmdByLine(cmd *exec.Cmd, readCallback func(str string, lineNumber int), errCallback func(err string, lineNumber int)) error {
	cmd.Dir = testutils.ProjectDir()

	rc, err := cmd.StdoutPipe()

	if err != nil {
		return err
	}

	rcerr, err := cmd.StderrPipe()

	if err != nil {
		return err
	}

	err = cmd.Start()

	if err != nil {
		return err
	}

	readFrom(rc, readCallback)
	readFrom(rcerr, errCallback)

	err = cmd.Wait()

	if err != nil {
		return err
	}

	return nil
}

func readFrom(rc io.ReadCloser, callback func(str string, lineNumber int)) {
	reader := bufio.NewReader(rc)

	lineNumber := 0

	for {
		str, err := reader.ReadString('\n')

		if err != nil {
			if err != io.EOF {
				log.Fatalln(err)
			}

			break
		}

		callback(str, lineNumber)
		lineNumber++
	}
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

	err = os.Mkdir(testutils.TmpPath(".warlock"), 0777)

	if err != nil {
		log.Fatal(err)
	}
}

func shutdown() {
	err := os.RemoveAll(testutils.TmpPath(""))

	if err != nil {
		log.Fatal(err)
	}

	err = os.RemoveAll(testutils.TmpPath(".warlock"))

	if err != nil {
		log.Fatal(err)
	}
}
