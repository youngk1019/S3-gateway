package edu_backend

import (
	"context"
	"github.com/go-playground/assert/v2"
	"testing"
)

func Test_EduBackend(t *testing.T) {
	set, err := GetDataSet(context.Background(), "f3842b6948c3406385683d8236b83112")
	assert.Equal(t, err, nil)
	assert.Equal(t, len(set) != 0, true)
}
