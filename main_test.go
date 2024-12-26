package main

import (
	"testing"

	"github.com/gms1/go-project-template/test"
	"github.com/stretchr/testify/assert"
)

func TestMainFunction(t *testing.T) {
	stdout, _, _ := test.CaptureOutput(func() error {
		main()
		return nil
	})
	assert.Contains(t, stdout, "Usage:")
}
