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
	tests := []struct {
		name      string
		items     []PagingCursorItem
		descOrder bool
		expected  string
		err       string
	}{
		{
			name:      "ID",
			items:     cursorItems(PagingID, int64(9876)),
			descOrder: false,
			expected:  "AAAAAAAAJpQB",
		},
		{
			name:      "ID and DescOrderBy",
			items:     cursorItems(PagingID, int64(9876)),
			descOrder: true,
			expected:  "AAAAAAAAJpT_",
		},
		{
			name:      "ID, UpdatedAt",
			items:     cursorItems(PagingID, int64(1234), PagingUpdatedAt, time.Time{}),
			descOrder: false,
			expected:  "AAAAAAAAAAACAAAAAAAABNIB",
		},
		{
			name:      "ID, UpdatedAt and DescOrderBy",
			items:     cursorItems(PagingID, int64(1234), PagingUpdatedAt, time.Time{}),
			descOrder: true,
			expected:  "AAAAAAAAAAACAAAAAAAABNL_",
		},
		{
			name:      "UpdatedAt and DescOrderBy",
			items:     cursorItems(PagingUpdatedAt, time.Time{}),
			descOrder: true,
			expected:  "AAAAAAAAAAD-",
			err:       "paging is invalid",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := encodeCursor(tt.items, tt.descOrder)
			require.Equal(t, tt.expected, output)

			decodedItems, decodedDescOrderBy, err := decodeCursor(output)
			if tt.err != "" {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.items, decodedItems)
				assert.Equal(t, tt.descOrder, decodedDescOrderBy)
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
	tests := []struct {
		name         string
		paging       *Paging
		pagingCursor PagingCursor
		err          string
	}{
		{
			name: "After With Limit",
			paging: &Paging{
				After: "AAAAAAAAAAACAAAAAAAABNL_",
				Limit: 3,
			},
			pagingCursor: PagingCursor{
				Items: []PagingCursorItem{
					{
						Field: PagingID,
						Value: int64(1234),
					},
					{
						Field: PagingUpdatedAt,
						Value: time.Time{},
					},
				},
				DescOrderBy: true,
				Reverse:     false,
			},
		},
		{
			name: "Last",
			paging: &Paging{
				Last: 100,
			},
			pagingCursor: PagingCursor{
				Items: []PagingCursorItem{
					{
						Field: PagingID,
						Value: 0,
					},
				},
				DescOrderBy: true,
				Reverse:     true,
			},
		},
		{
			name: "Last with Sort(updated_At)",
			paging: &Paging{
				Last: 100,
				Sort: []string{"updated_at"},
			},
			pagingCursor: PagingCursor{
				Items: []PagingCursorItem{
					{
						Field: PagingUpdatedAt,
						Value: 0,
					},
					{
						Field: PagingID,
						Value: 0,
					},
				},
				DescOrderBy: true,
				Reverse:     true,
			},
		},
		{
			name: "Last with Sort(-updated_at)",
			paging: &Paging{
				Last: 100,
				Sort: []string{"-updated_at"},
			},
			pagingCursor: PagingCursor{
				Items: []PagingCursorItem{
					{
						Field: PagingUpdatedAt,
						Value: 0,
					},
					{
						Field: PagingID,
						Value: 0,
					},
				},
				Reverse: true,
			},
		},
		{
			name: "Before with Limit",
			paging: &Paging{
				Before: "AAAAAAAAAAACAAAAAAAABNL_",
				Limit:  3,
			},
			pagingCursor: PagingCursor{
				Items: []PagingCursorItem{
					{
						Field: PagingID,
						Value: int64(1234),
					},
					{
						Field: PagingUpdatedAt,
						Value: time.Time{},
					},
				},
				DescOrderBy: true,
				Reverse:     true,
			},
		},
		{
			name: "First",
			paging: &Paging{
				First: 100,
			},
			pagingCursor: PagingCursor{
				Items: []PagingCursorItem{
					{
						Field: PagingID,
						Value: 0,
					},
				},
				Reverse: false,
			},
		},
		{
			name: "First with Sort(-updated_at)",
			paging: &Paging{
				First: 100,
				Sort:  []string{"-updated_at"},
			},
			pagingCursor: PagingCursor{
				Items: []PagingCursorItem{
					{
						Field: PagingUpdatedAt,
						Value: 0,
					},
					{
						Field: PagingID,
						Value: 0,
					},
				},
				DescOrderBy: true,
				Reverse:     false,
			},
		},
		{
			name: "First with Sort(updated_at)",
			paging: &Paging{
				First: 100,
				Sort:  []string{"updated_at"},
			},
			pagingCursor: PagingCursor{
				Items: []PagingCursorItem{
					{
						Field: PagingUpdatedAt,
						Value: 0,
					},
					{
						Field: PagingID,
						Value: 0,
					},
				},
				Reverse: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.paging.decodeCursor()
			if tt.err != "" {
				assert.NoError(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.pagingCursor, result)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name   string
		paging Paging
		err    string
	}{
		{
			name: "Limit(1)",
			paging: Paging{
				Limit: 1,
			},
			err: "paging is invalid (0 < Limit < 1000)",
		},
		{
			name: "Limit(1001)",
			paging: Paging{
				Limit: 1001,
			},
			err: "paging is invalid (0 < Limit < 1000)",
		},
		{
			name: "Limit(0)",
			paging: Paging{
				Limit: 0,
			},
			err: "paging is invalid",
		},
		{
			name: "First(-1)",
			paging: Paging{
				Limit: -1,
			},
			err: "paging is invalid (First must have a value between 0 and 1000)",
		},
		{
			name: "First(1001)",
			paging: Paging{
				First: 1001,
			},
			err: "paging is invalid (First must have a value between 0 and 1000)",
		},
		{
			name: "First(101)",
			paging: Paging{
				First: 101,
			},
		},
		{
			name: "Last(-1)",
			paging: Paging{
				Last: -1,
			},
			err: "paging is invalid (Last must have a value between 0 and 1000)",
		},
		{
			name: "Last(1001)",
			paging: Paging{
				Last: 1001,
			},
			err: "paging is invalid (Last must have a value between 0 and 1000)",
		},
		{
			name: "First(100), Last(2)",
			paging: Paging{
				First: 100,
				Last:  2,
			},
			err: "paging is invalid",
		},
		{
			name:   "Empty paging",
			paging: Paging{},
			err:    "paging is invalid",
		},
		{
			name: "Before",
			paging: Paging{
				Before: "id~1",
			},
			err: "paging is invalid (Before must be used together with Limit)",
		},
		{
			name: "Before with Limit",
			paging: Paging{
				Before: "id~1",
				Limit:  1,
			},
		},
		{
			name: "Before with Sort(-updated_at)",
			paging: Paging{
				Before: "id~1",
				Sort:   []string{"-updated_at"},
			},
			err: "paging is invalid (Before must not be used together with Sort)",
		},
		{
			name: "Before, After and Limit(100)",
			paging: Paging{
				Before: "id~1",
				After:  "id~2",
				Limit:  100,
			},
			err: "paging is invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.err != "" {
				assert.Error(t, tt.paging.Validate())
			} else {
				assert.NoError(t, tt.paging.Validate())
			}
		})
	}
}

func TestValidateSort(t *testing.T) {
	sortWhitelist := map[string]string{
		"id":         "",
		"updated_at": "",
		"created_at": "",
	}

	tests := []struct {
		name   string
		paging *Paging
		err    string
	}{
		{
			name: "Before with Limit",
			paging: &Paging{
				Before: "id~1",
				Limit:  1,
				Sort:   []string{"id"},
			},
		},
		{
			name: "Before with Limit, Sort(created_at)",
			paging: &Paging{
				Before: "id~1",
				Limit:  1,
				Sort:   []string{"created_at"},
			},
			err: "Sort by created_at is not allowed",
		},
		{
			name: "Before with Limit, Sort(id, updated_at)",
			paging: &Paging{
				Before: "id~1",
				Limit:  1,
				Sort:   []string{"id", "updated_at"},
			},
			err: "paging is invalid (Sort support only one field at the same time)",
		},
		{
			name: "Before with Limit, Sort(name)",
			paging: &Paging{
				Before: "id~1",
				Limit:  1,
				Sort:   []string{"name"},
			},
			err: "Sort by created_at is not allowed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err != "" {
				assert.Error(t, validateSort(tt.paging, sortWhitelist))
			} else {
				assert.NoError(t, validateSort(tt.paging, sortWhitelist))
			}
		})
	}
}

func TestReverseAny(t *testing.T) {
	t.Run("Reverse Any", func(t *testing.T) {
		a := []string{"a", "b", "c"}
		reverseAny(a)
		assert.Equal(t, []string{"c", "b", "a"}, a)
	})
}

func TestLimitSort(t *testing.T) {
	sortWhitelist := map[string]string{
		"id":         "id",
		"updated_at": "updated_at",
		"name":       "",
	}

	tests := []struct {
		name   string
		paging *Paging
		limit  int
		err    string
	}{
		{
			name: "First",
			paging: &Paging{
				First: 2,
			},
			limit: 2,
		},
		{
			name: "Before with Limit",
			paging: &Paging{
				Before: base64.StdEncoding.EncodeToString([]byte("011")),
				Limit:  10,
			},
			err: "paging is invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := applyCursorPaging(db.NewQuery(), tt.paging, sortWhitelist, "")
			if tt.err != "" {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, 2, tt.limit)
			}
		})
	}
}

func TestBuild(t *testing.T) {
	sortWhitelist := map[string]string{
		"id":         "id",
		"updated_at": "updated_at",
		"name":       "",
	}

	tests := []struct {
		name                  string
		paging                *Paging
		cursorPagingCondition CursorPagingCondition
		err                   string
	}{
		{
			name: "Before with Limit",
			paging: &Paging{
				Before: "AAAAAAAAAAACAAAAAAAABNIB",
				Limit:  10,
			},
			cursorPagingCondition: CursorPagingCondition{
				args: []interface{}{
					time.Time{},
					int64(1234),
				},
				cols: []string{
					"updated_at",
					"id",
				},
				operation: "<",
				prefix:    "",
			},
		},
		{
			name: "After(sort=id) with Limit",
			paging: &Paging{
				After: "AAAAAAAAJpQB",
				Limit: 1,
			},
			cursorPagingCondition: CursorPagingCondition{
				prefix:    "",
				operation: ">",
				cols: []string{
					"id",
				},
				args: []interface{}{
					int64(9876),
				},
			},
		},
		{
			name: "After(contain updated_at, id) with Limit",
			paging: &Paging{
				After: "AAAAAAAAAAACAAAAAAAABNIB",
				Limit: 10,
			},
			cursorPagingCondition: CursorPagingCondition{
				operation: ">",
				prefix:    "",
				cols: []string{
					"updated_at",
					"id",
				},
				args: []interface{}{
					time.Time{},
					int64(1234),
				},
			},
		},
		{
			name: "After(invalid) with Limit",
			paging: &Paging{
				After: "AAAAAAAAJpQZ",
				Limit: 1,
			},
			err: "paging is invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := build(tt.paging, sortWhitelist, "")
			if tt.err != "" {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.cursorPagingCondition, result)
			}
		})
	}

	t.Run("Paging nil", func(t *testing.T) {
		result, err := build(nil, sortWhitelist, "")
		assert.NoError(t, err)
		assert.Nil(t, result)
	})
}
