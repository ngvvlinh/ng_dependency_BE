package summary

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDate(t *testing.T) {
	d := startOfDay(0)
	assert.Equal(t, d.Hour(), 0)
	assert.Equal(t, d.UTC().Hour(), 17)
}
