package tests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi/dot"
)

var _db *cmsql.Database

func connectDB() *cmsql.Database {
	if _db == nil {
		_db = cmsql.MustConnect(cc.DefaultPostgres())
	}

	return _db
}

func TestSQL(t *testing.T) {
	db := connectDB()

	_db.MustExec("drop table if exists a")
	_db.MustExec(`
create table a (
    id int8,
    name text,
    created_at timestamptz,
    updated_at timestamptz
)`)

	// time different between postgres and local machine can be high
	timeDelta := 5 * time.Minute

	start := time.Now()
	t.Run("insert", func(t *testing.T) {
		a := &A{
			ID:   100,
			Name: "one hundred",
		}
		n, err := db.Insert(a)
		require.NoError(t, err)
		require.Equal(t, 1, n)
	})

	var lastCreatedAt time.Time
	t.Run("get back", func(t *testing.T) {
		var a A
		ok, err := db.Where("id = 100").Get(&a)
		require.NoError(t, err)
		require.True(t, ok)

		assert.Equal(t, dot.ID(100), a.ID)
		assert.Equal(t, "one hundred", a.Name)
		assert.True(t, a.CreatedAt.Equal(a.UpdatedAt))
		assert.WithinDuration(t, start, a.CreatedAt, timeDelta)

		lastCreatedAt = a.CreatedAt
	})
	t.Run("update", func(t *testing.T) {
		a := &A{
			Name: "one hundred (100)",
		}
		n, err := db.Where("id = 100").Update(a)
		require.NoError(t, err)
		require.Equal(t, 1, n)
	})
	t.Run("get back after update", func(t *testing.T) {
		var a A
		ok, err := db.Where("id = 100").Get(&a)
		require.NoError(t, err)
		require.True(t, ok)

		assert.Equal(t, dot.ID(100), a.ID)
		assert.Equal(t, "one hundred (100)", a.Name)
		assert.True(t, a.CreatedAt.Equal(lastCreatedAt), "created_at should not change")
		assert.True(t, a.UpdatedAt.After(lastCreatedAt), "updated_at should change")
		assert.WithinDuration(t, start, a.CreatedAt, timeDelta)
	})
}
