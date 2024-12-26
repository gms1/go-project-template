package test

import (
	"os"
)

func CaptureOutput(f func() error) (string, string, error) {
	stdout := os.Stdout
	stderr := os.Stderr
	outFile, _ := os.CreateTemp("", "gotest")
	errFile, _ := os.CreateTemp("", "gotest")
	defer os.Remove(outFile.Name())
	defer os.Remove(errFile.Name())
	os.Stdout = outFile
	os.Stderr = errFile
	err := f()
	os.Stdout = stdout
	os.Stderr = stderr

	outFile.Close()
	errFile.Close()
	outBytes, _ := os.ReadFile(outFile.Name())
	errBytes, _ := os.ReadFile(errFile.Name())
	return string(outBytes), string(errBytes), err
}
