package cm_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	. "o.o/backend/pkg/common"
)

func TestID(t *testing.T) {
	t.Run("Luna New Year", func(t *testing.T) {
		ta := time.Date(2018, 02, 15, 17, 0, 0, 0, time.UTC)
		id := NewIDWithTime(ta)
		bit := (id >> 24) & 1
		t.Logf("%b (%b)", id, bit)
		assert.Equal(t, int64(1012345678921203712), (id>>24)<<24)
	})
	t.Run("Biggest available id", func(t *testing.T) {
		ta := time.Date(2173, 03, 19, 23, 58, 13, 50e6, time.UTC)
		id := NewIDWithTime(ta)
		bit := (id >> 24) & 1
		t.Logf("%b (%b)", id, bit)
		assert.Equal(t, int64(9223372036821221376), (id>>24)<<24)
	})
	t.Run("No tag", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			id := NewID()
			assert.True(t, id > 1e18, "ID must have 19 characters")
			assert.True(t, uint64(id) < uint64(1e19), "ID must have 19 characters")

			bit := (id >> 24) & 1
			t.Logf("%b (%b)", id, bit)
			assert.True(t, bit == 0, "Bit must be 0")
			assert.Equal(t, int64(0), GetTag(id), "Tag must be 0")
		}
	})
	t.Run("With tag", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			tag := 31 + byte(i)*2 // Tag must be odd number
			id := NewIDWithTag(tag)
			assert.True(t, id > 1e18, "ID must have 19 characters")
			assert.True(t, uint64(id) < uint64(1e19), "ID must have 19 characters")

			bit := (id >> 24) & 1
			t.Logf("%b (%b)", id, bit)
			assert.True(t, bit == 1, "Bit must be 1")
			assert.Equal(t, int64(tag), GetTag(id), "Tag must be extracted")
		}
	})
}

func TestNewBase54ID(t *testing.T) {
	id := NewBase54ID()
	t.Logf("%v", id)
	assert.Equal(t, Base54IDLength, len(id))
}
