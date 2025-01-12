package common

import (
	"os"
	"runtime"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitRuntime(t *testing.T) {
	testCases := []struct {
		name    string
		withEcs bool
	}{
		{"with ECS", true},
		{"without ECS", false},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			givenMaxProcs := runtime.GOMAXPROCS(0)
			givenMemLimit := "500MiB"
			t.Setenv("GOMAXPROCS", strconv.Itoa(givenMaxProcs))
			t.Setenv("GOMEMLIMIT", givenMemLimit)
			if testCase.withEcs {
				t.Setenv("ECS_CONTAINER_METADATA_URI_V4", "xxx")
			} else {
				t.Setenv("ECS_CONTAINER_METADATA_URI_V4", "")
			}
			InitRuntime()
			assert.Equal(t, givenMaxProcs, runtime.GOMAXPROCS(0))
			assert.Equal(t, givenMemLimit, os.Getenv("GOMEMLIMIT"))
		})
	}
}
