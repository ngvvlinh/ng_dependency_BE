package meta

import (
	"time"

	uuid "github.com/satori/go.uuid"

	metav1 "etop.vn/api/meta/v1"
)

type Empty = metav1.Empty
type UUID = metav1.UUID
type Timestamp = metav1.Timestamp

func NewUUID() UUID {
	u := uuid.NewV4()
	return UUID{u[:]}
}

func PbTime(t time.Time) *Timestamp {
	return metav1.PbTime(t)
}

type PageInfo struct {
	Offset int32
	Limit  int32
	Sort   []string

	// TODO: next, prev
}

func FromPaging(paging Paging) PageInfo {
	return PageInfo(paging)
}

type Paging = metav1.Paging
type Filter = metav1.Filter
type Filters []Filter

type Error = metav1.Error
