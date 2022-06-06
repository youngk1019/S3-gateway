package list_objects

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func Test_GenListObjectsResult(t *testing.T) {
	ret, err := GenListObjectsResult([]string{"BasicDataset"})
	assert.Equal(t, err, nil)
	println(string(ret))
}
