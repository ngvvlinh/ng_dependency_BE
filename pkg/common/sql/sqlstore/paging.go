package sqlstore

import (
	"encoding/binary"
	"errors"
	"io"
	"reflect"
	"time"

	cm "etop.vn/backend/pkg/common"
)

// +enum
// +enum:zero=null
type PagingField int

const (
	// +enum=unknown
	PagingUnknown PagingField = 0

	// +enum=id
	PagingID PagingField = 1

	// +enum=updated_at
	PagingUpdatedAt PagingField = 2

	NumPagingField = int(PagingUpdatedAt) + 1
)

var pagingFieldDescs = map[PagingField]*PagingFieldDesc{
	PagingID: {
		FromField: func(field reflect.Value) interface{} { return field.Int() },
		Decode:    func(r io.Reader) (interface{}, error) { return readInt64(r) },
		Encode:    func(w io.Writer, v interface{}) error { return writeInt64(w, v.(int64)) },
	},
	PagingUpdatedAt: {
		FromField: func(field reflect.Value) interface{} { return field.Interface().(time.Time) },
		Decode:    func(r io.Reader) (interface{}, error) { v, err := readInt64(r); return cm.FromMicros(v), err },
		Encode:    func(w io.Writer, v interface{}) error { return writeInt64(w, cm.Micros(v.(time.Time))) },
	},
}

type PagingFieldDesc struct {
	FromField func(field reflect.Value) interface{}
	Decode    func(r io.Reader) (interface{}, error)
	Encode    func(w io.Writer, v interface{}) error
}

func readInt64(r io.Reader) (int64, error) {
	var p [8]byte
	n, err := r.Read(p[:])
	if err != nil {
		return 0, err
	}
	if n < 8 {
		return 0, errors.New("invalid")
	}
	return int64(binary.LittleEndian.Uint64(p[:])), nil
}

func writeInt64(w io.Writer, v int64) error {
	var p [8]byte
	binary.LittleEndian.PutUint64(p[:], uint64(v))
	_, err := w.Write(p[:])
	return err
}

type PagingCursor struct {
	Items       []PagingCursorItem
	Reverse     bool // before/after, first/last
	DescOrderBy bool
}

type PagingCursorItem struct {
	Field PagingField
	Value interface{}
}

type PagingFieldMappingItem struct {
	Index int
}

// PagingFieldMapping provides a map between Paging enum and field index in a
// struct
type PagingFieldMapping [NumPagingField]PagingFieldMappingItem

func (m PagingFieldMapping) GetField(st reflect.Value, pagingField PagingField) reflect.Value {
	return st.Field(m[pagingField].Index)
}

func (m PagingFieldMapping) GetFieldType(st reflect.Type, pagingField PagingField) reflect.StructField {
	return st.Field(m[pagingField].Index)
}

func (m PagingFieldMapping) GetFieldValue(st reflect.Value, pagingField PagingField) interface{} {
	field := st.Field(m[pagingField].Index)
	return pagingFieldDescs[pagingField].FromField(field)
}
