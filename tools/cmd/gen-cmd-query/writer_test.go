package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchVersion(t *testing.T) {
	assert.Equal(t, "main/v1", reVx.FindString("example.com/main/v1"))
	assert.Equal(t, "main/v1beta", reVx.FindString("example.com/main/v1beta"))
	assert.Equal(t, "main/v1/foo", reVx.FindString("example.com/main/v1/foo"))
	assert.Equal(t, "", reVx.FindString("example.com/main/v1beta/foo/bar"))
}
