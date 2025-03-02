package test

import (
	"os"

	"github.com/prashantv/gostub"
)

// CaptureOutput captures the output of stdout and stderr for calling the given function
//
//	in case you want to get the output of a logger, please make sure that the logging handler is created inside
//	the given function, otherwise it will not use the stubs
func CaptureOutput(f func() error) (string, string, error) {
	outFile, _ := os.CreateTemp("", "gotest")
	errFile, _ := os.CreateTemp("", "gotest")
	defer os.Remove(outFile.Name())
	defer os.Remove(errFile.Name())

	stubs := gostub.New()
	defer stubs.Reset()
	stubs.Stub(&os.Stdout, outFile)
	stubs.Stub(&os.Stderr, errFile)

	err := f()

	outFile.Close()
	errFile.Close()
	outBytes, _ := os.ReadFile(outFile.Name())
	errBytes, _ := os.ReadFile(errFile.Name())
	return string(outBytes), string(errBytes), err
}
