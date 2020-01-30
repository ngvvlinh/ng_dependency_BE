package sqlstore

import (
	"encoding/base64"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	cc "etop.vn/backend/pkg/common/config"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/backend/pkg/common/sql/sq"
	"etop.vn/capi/dot"
)

var db *cmsql.Database

func init() {
	cfg := cc.DefaultPostgres()
	db = cmsql.MustConnect(cfg)
}

func TestGetPagingFieldMapping(t *testing.T) {
	type A struct {
		Name      string
		ProductID dot.ID    `paging:"id"`
		UpdatedAt time.Time `paging:"updated_at"`
	}
	t.Run("Test GetPagingFieldMapping 1", func(t *testing.T) {
		s := reflect.ValueOf(A{})
		mapping := getPagingFieldMapping(s.Type())
		assert.Equal(t, PagingFieldMapping{
			{
				Index: -1,
			},
			{
				Index: 1,
			},
			{
				Index: 2,
			},
		}, mapping)
	})
}

func cursorItems(args ...interface{}) (result []PagingCursorItem) {
	for i := 0; i < len(args)/2; i++ {
		result = append(result, PagingCursorItem{
			Field: args[i*2].(PagingField),
			Value: args[i*2+1],
		})
	}
	return result
}

func TestEncodeDecodeCursor(t *testing.T) {
	sampleTime := time.Date(2020, time.January, 2, 20, 10, 5, 123456789, time.UTC)
	encodedTime := sampleTime.Truncate(1000 * time.Nanosecond)

	tests := []struct {
		name         string
		items        []PagingCursorItem
		decodedItems []PagingCursorItem
		negativeSort bool
		expected     string
		err          string
	}{
		{
			name:         "ID",
			items:        cursorItems(PagingID, int64(9876)),
			decodedItems: cursorItems(PagingID, int64(9876)),
			negativeSort: false,
			expected:     "AAAAAAAAJpQB",
		},
		{
			name:         "ID and DescOrderBy",
			items:        cursorItems(PagingID, int64(9876)),
			decodedItems: cursorItems(PagingID, int64(9876)),
			negativeSort: true,
			expected:     "AAAAAAAAJpT_",
		},
		{
			name:         "ID, UpdatedAt",
			items:        cursorItems(PagingID, int64(1234), PagingUpdatedAt, sampleTime),
			decodedItems: cursorItems(PagingID, int64(1234), PagingUpdatedAt, encodedTime),
			negativeSort: false,
			expected:     "AAWbLcdr44ACAAAAAAAABNIB",
		},
		{
			name:         "ID, UpdatedAt with negative sort",
			items:        cursorItems(PagingID, int64(1234), PagingUpdatedAt, sampleTime),
			decodedItems: cursorItems(PagingID, int64(1234), PagingUpdatedAt, encodedTime),
			negativeSort: true,
			expected:     "AAWbLcdr44ACAAAAAAAABNL_",
		},
		{
			// still ok because decodeCursor does not validate PagingID
			name:         "UpdatedAt with negative sort",
			items:        cursorItems(PagingUpdatedAt, sampleTime),
			decodedItems: cursorItems(PagingUpdatedAt, encodedTime),
			negativeSort: true,
			expected:     "AAWbLcdr44D-",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := encodeCursor(tt.items, tt.negativeSort)
			require.Equal(t, tt.expected, output)

			decodedItems, negativeSort, err := decodeCursor(output)
			if tt.err != "" {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.decodedItems, decodedItems)
				assert.Equal(t, tt.negativeSort, negativeSort)
			}
		})
	}
	t.Run("Invalid field (panic)", func(t *testing.T) {
		cursorItems := cursorItems(PagingUnknown, int64(0))
		assert.Panics(t, func() { encodeCursor(cursorItems, false) })
	})
	t.Run("Invalid value (panic)", func(t *testing.T) {
		cursorItems := cursorItems(PagingUpdatedAt, "")
		assert.Panics(t, func() { encodeCursor(cursorItems, false) })
	})
}

func TestDecodeCursorPaging(t *testing.T) {
	sampleTime := time.Date(2020, time.January, 2, 20, 10, 5, 123456789, time.UTC)
	encodedTime := sampleTime.Truncate(1000 * time.Nanosecond)

	tests := []struct {
		name         string
		paging       *Paging
		pagingCursor PagingCursor
		descOrderBy  bool
		err          string
	}{
		{
			name: "First",
			paging: &Paging{
				After: ".",
				Limit: 100,
			},
			pagingCursor: PagingCursor{
				Items:        cursorItems(PagingID, nil),
				Reverse:      false,
				NegativeSort: false,
			},
			descOrderBy: false,
		},
		{
			name: "Last",
			paging: &Paging{
				Before: ".",
				Limit:  100,
			},
			pagingCursor: PagingCursor{
				Items:        cursorItems(PagingID, nil),
				Reverse:      true,
				NegativeSort: false,
			},
			descOrderBy: true,
		},
		{
			name: "First with sort (updated_at)",
			paging: &Paging{
				After: ".",
				Limit: 100,
				Sort:  []string{"updated_at"},
			},
			pagingCursor: PagingCursor{
				Items:        cursorItems(PagingUpdatedAt, nil, PagingID, nil),
				Reverse:      false,
				NegativeSort: false,
			},
			descOrderBy: false,
		},
		{
			name: "First with negative sort (-updated_at)",
			paging: &Paging{
				After: ".",
				Limit: 100,
				Sort:  []string{"-updated_at"},
			},
			pagingCursor: PagingCursor{
				Items:        cursorItems(PagingUpdatedAt, nil, PagingID, nil),
				Reverse:      false,
				NegativeSort: true,
			},
			descOrderBy: true,
		},
		{
			name: "Last with sort (updated_at)",
			paging: &Paging{
				Before: ".",
				Limit:  100,
				Sort:   []string{"updated_at"},
			},
			pagingCursor: PagingCursor{
				Items:        cursorItems(PagingUpdatedAt, nil, PagingID, nil),
				Reverse:      true,
				NegativeSort: false,
			},
			descOrderBy: true,
		},
		{
			name: "Last with negative sort (-updated_at)",
			paging: &Paging{
				Before: ".",
				Limit:  100,
				Sort:   []string{"-updated_at"},
			},
			pagingCursor: PagingCursor{
				Items:        cursorItems(PagingUpdatedAt, nil, PagingID, nil),
				Reverse:      true,
				NegativeSort: true,
			},
			descOrderBy: false,
		},
		{
			name: "After with sort (updated_at_",
			paging: &Paging{
				After: encodeCursor(cursorItems(PagingUpdatedAt, sampleTime, PagingID, int64(1234)), false),
				Limit: 3,
			},
			pagingCursor: PagingCursor{
				Items:        cursorItems(PagingUpdatedAt, encodedTime, PagingID, int64(1234)),
				Reverse:      false,
				NegativeSort: false,
			},
			descOrderBy: false,
		},
		{
			name: "After with negative sort",
			paging: &Paging{
				After: encodeCursor(cursorItems(PagingUpdatedAt, sampleTime, PagingID, int64(1234)), true),
				Limit: 3,
			},
			pagingCursor: PagingCursor{
				Items:        cursorItems(PagingUpdatedAt, encodedTime, PagingID, int64(1234)),
				Reverse:      false,
				NegativeSort: true,
			},
			descOrderBy: true,
		},
		{
			name: "Before with negative sort",
			paging: &Paging{
				Before: encodeCursor(cursorItems(PagingUpdatedAt, sampleTime, PagingID, int64(1234)), true),
				Limit:  3,
			},
			pagingCursor: PagingCursor{
				Items:        cursorItems(PagingUpdatedAt, encodedTime, PagingID, int64(1234)),
				Reverse:      true,
				NegativeSort: true,
			},
			descOrderBy: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.paging.decodeCursor()
			if tt.err != "" {
				assert.EqualError(t, err, tt.err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.pagingCursor, result)
			assert.Equal(t, tt.descOrderBy, tt.pagingCursor.DescOrderBy())
		})
	}
}

func TestValidate(t *testing.T) {
	sampleCursor := encodeCursor(cursorItems(PagingID, int64(1000)), false)
	tests := []struct {
		name   string
		paging Paging
		err    string
	}{
		{
			name: "First(101)",
			paging: Paging{
				After: ".",
				Limit: 101,
			},
		},
		{
			name: "Before with Limit",
			paging: Paging{
				Before: sampleCursor,
				Limit:  1,
			},
		},
		{
			name: "Before with sort",
			paging: Paging{
				Before: encodeCursor(cursorItems(PagingUpdatedAt, time.Now(), PagingID, int64(1000)), false),
				Sort:   []string{"updated_at"},
				Limit:  1,
			},
		},
		{
			name: "Before with negative sort",
			paging: Paging{
				Before: encodeCursor(cursorItems(PagingUpdatedAt, time.Now(), PagingID, int64(1000)), true),
				Sort:   []string{"-updated_at"},
				Limit:  1,
			},
		},
		{
			name: "Before with sort (invalid)",
			paging: Paging{
				Before: sampleCursor,
				Sort:   []string{"updated_at"},
				Limit:  1,
			},
			err: "paging is invalid (sort does not match)",
		},
		{
			name: "Missing after and before",
			paging: Paging{
				Limit: 1,
			},
			err: "paging is invalid",
		},
		{
			name: "Both after and before",
			paging: Paging{
				Before: ".",
				After:  ".",
			},
			err: "paging is invalid",
		},
		{
			name: "Limit(0)",
			paging: Paging{
				After: sampleCursor,
				Limit: 0,
			},
			err: "paging is invalid (limit is required)",
		},
		{
			name: "Limit(1001)",
			paging: Paging{
				After: sampleCursor,
				Limit: 1001,
			},
			err: "paging is invalid (limit outside of range)",
		},
		{
			name: "Limit(-1)",
			paging: Paging{
				After: sampleCursor,
				Limit: -1,
			},
			err: "paging is invalid (limit outside of range)",
		},
		{
			name: "Last(-1)",
			paging: Paging{
				Before: ".",
				Limit:  1001,
			},
			err: "paging is invalid (limit outside of range)",
		},
		{
			name: "Last(1001)",
			paging: Paging{
				Before: ".",
				Limit:  1001,
			},
			err: "paging is invalid (limit outside of range)",
		},
		{
			name:   "empty",
			paging: Paging{},
			err:    "paging is invalid",
		},
		{
			name: "Before",
			paging: Paging{
				Before: sampleCursor,
			},
			err: "paging is invalid (limit is required)",
		},
		{
			name: "Before with Sort(-updated_at)",
			paging: Paging{
				Before: sampleCursor,
				Sort:   []string{"-updated_at"},
			},
			err: "paging is invalid (limit is required)",
		},
		{
			name: "Before, After and Limit(100)",
			paging: Paging{
				Before: sampleCursor,
				After:  sampleCursor,
				Limit:  100,
			},
			err: "paging is invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.paging.validateCursorPaging()
			if tt.err != "" {
				assert.EqualError(t, err, tt.err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestValidateSort(t *testing.T) {
	sortWhitelist := map[string]string{
		"id":         "",
		"updated_at": "",
		"created_at": "",
	}

	sampleCursor := encodeCursor(cursorItems(PagingID, int64(1000)), false)
	tests := []struct {
		name   string
		paging *Paging
		err    string
	}{
		{
			name: "Before with Limit",
			paging: &Paging{
				Before: sampleCursor,
				Limit:  1,
				Sort:   []string{"id"},
			},
		},
		{
			name: "Before with Limit, Sort(created_at)",
			paging: &Paging{
				Before: sampleCursor,
				Limit:  1,
				Sort:   []string{"created_at"},
			},
			err: "Sort by created_at is not allowed",
		},
		{
			name: "Before with Limit, Sort(name)",
			paging: &Paging{
				Before: sampleCursor,
				Limit:  1,
				Sort:   []string{"name"},
			},
			err: "Sort by name is not allowed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err != "" {
				assert.Error(t, validateCursorPagingSort(tt.paging.Sort, sortWhitelist))
			} else {
				assert.NoError(t, validateCursorPagingSort(tt.paging.Sort, sortWhitelist))
			}
		})
	}
}

func TestReverseAny(t *testing.T) {
	t.Run("Reverse Any", func(t *testing.T) {
		a := []string{"a", "b", "c"}
		reverseSlice(a)
		assert.Equal(t, []string{"c", "b", "a"}, a)
	})
}

func TestLimitSort(t *testing.T) {
	sortWhitelist := map[string]string{
		"id":         "id",
		"updated_at": "updated_at",
		"name":       "",
	}
	sampleTime := time.Date(2020, time.January, 2, 20, 10, 5, 123456789, time.UTC)
	encodedTime := sampleTime.Truncate(1000 * time.Nanosecond)

	tests := []struct {
		name   string
		paging *Paging
		cond   sq.WriterTo
		sql    string
		args   []interface{}
		err    string
	}{
		{
			name: "First",
			paging: &Paging{
				After: ".",
				Limit: 2,
			},
			sql:  `SELECT * FROM foo ORDER BY "id" LIMIT 2`,
			args: []interface{}{},
			cond: nil,
		},
		{
			name: "Last",
			paging: &Paging{
				Before: ".",
				Limit:  2,
			},
			sql:  `SELECT * FROM foo ORDER BY "id" DESC LIMIT 2`,
			args: []interface{}{},
			cond: nil,
		},
		{
			name: "First with negative sort (updated_at)",
			paging: &Paging{
				After: ".",
				Limit: 2,
				Sort:  []string{"-updated_at"},
			},
			sql:  `SELECT * FROM foo ORDER BY "updated_at" DESC,"id" DESC LIMIT 2`,
			args: []interface{}{},
			cond: nil,
		},
		{
			name: "Last with negative sort (updated_at)",
			paging: &Paging{
				Before: ".",
				Limit:  2,
				Sort:   []string{"-updated_at"},
			},
			sql:  `SELECT * FROM foo ORDER BY "updated_at","id" LIMIT 2`,
			args: []interface{}{},
			cond: nil,
		},
		{
			name: "After (id)",
			paging: &Paging{
				After: encodeCursor(cursorItems(PagingID, int64(9876)), false),
				Limit: 1,
			},
			cond: &CursorPagingCondition{
				prefix:    "",
				operation: ">",
				cols:      []string{"id"},
				args:      []interface{}{int64(9876)},
			},
			sql:  `SELECT * FROM foo WHERE ("id" > $1) ORDER BY "id" LIMIT 1`,
			args: []interface{}{int64(9876)},
		},
		{
			name: "Before (id)",
			paging: &Paging{
				Before: encodeCursor(cursorItems(PagingID, int64(9876)), false),
				Limit:  1,
			},
			cond: &CursorPagingCondition{
				prefix:    "",
				operation: "<",
				cols:      []string{"id"},
				args:      []interface{}{int64(9876)},
			},
			sql:  `SELECT * FROM foo WHERE ("id" < $1) ORDER BY "id" DESC LIMIT 1`,
			args: []interface{}{int64(9876)},
		},
		{
			name: "After with negative sort (-updated_at)",
			paging: &Paging{
				After: encodeCursor(cursorItems(PagingUpdatedAt, sampleTime, PagingID, int64(1234)), true),
				Limit: 10,
			},
			cond: &CursorPagingCondition{
				operation: ">",
				prefix:    "",
				cols:      []string{"updated_at", "id"},
				args:      []interface{}{encodedTime, int64(1234)},
			},
			sql:  `SELECT * FROM foo WHERE (("updated_at","id") < ($1,$2)) ORDER BY "updated_at" DESC,"id" DESC LIMIT 10`,
			args: []interface{}{encodedTime, int64(1234)},
		},
		{
			name: "Before with negative sort (-updated_at)",
			paging: &Paging{
				Before: encodeCursor(cursorItems(PagingUpdatedAt, sampleTime, PagingID, int64(1234)), true),
				Limit:  10,
			},
			cond: &CursorPagingCondition{
				args:      []interface{}{encodedTime, int64(1234)},
				cols:      []string{"updated_at", "id"},
				operation: "<",
				prefix:    "",
			},
			sql:  `SELECT * FROM foo WHERE (("updated_at","id") > ($1,$2)) ORDER BY "updated_at","id" LIMIT 10`,
			args: []interface{}{encodedTime, int64(1234)},
		},
		{
			name: "Before (invalid)",
			paging: &Paging{
				Before: base64.StdEncoding.EncodeToString([]byte("011")),
				Limit:  10,
			},
			cond: nil,
			err:  "paging is invalid (invalid cursor) original=invalid cursor",
		},
		{
			name: "After (invalid)",
			paging: &Paging{
				After: "AAAAAAAAJpQZ",
				Limit: 1,
			},
			err: "paging is invalid (invalid cursor) original=invalid cursor",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := db.NewQuery().SQL(`SELECT * FROM foo`)
			query, err := LimitSort(query, tt.paging, sortWhitelist)
			if tt.err != "" {
				assert.EqualError(t, err, tt.err)
				return
			}

			require.NoError(t, err)
			sql, args, err := query.Build()
			require.NoError(t, err)
			assert.Equal(t, tt.sql, sql)
			assert.Equal(t, tt.args, args)
		})
	}
}
