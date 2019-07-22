package meta

import (
	"database/sql"
	"time"

	metav1 "etop.vn/api/meta/v1"
	cmutil "etop.vn/common/util"
	"etop.vn/common/xerrors"

	uuid "github.com/satori/go.uuid"
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

type UpdateOp string

const (
	OpAdd        UpdateOp = "add"
	OpRemove     UpdateOp = "remove"
	OpReplace    UpdateOp = "replace"
	OpReplaceAll UpdateOp = "replace_all"
	OpDeleteAll  UpdateOp = "delete_all"
)

type UpdateSet struct {
	Op      UpdateOp
	Changes []string
}

type UpdateListInt64 struct {
	Op      UpdateOp
	Changes []int64
}

type UpdateListInt32 struct {
	Op      UpdateOp
	Changes []int32
}

type NullBool sql.NullBool
type NullString sql.NullString
type NullInt64 sql.NullInt64
type NullInt32 struct {
	Int32 int32
	Valid bool
}

func (u *UpdateSet) Update(set []string) (result []string, count int, err error) {
	switch u.Op {
	case OpAdd:
		for _, item := range u.Changes {
			if cmutil.ListStringsContain(set, item) {
				continue
			}
			set = append(set, item)
			count++
		}
		return set, count, nil

	case OpRemove:
		found := true
	found:
		for _, item := range u.Changes {
			for _, x := range set {
				if x == item {
					found = true
					break found
				}
			}
		}
		if !found {
			return set, 0, nil
		}
		result := make([]string, 0, len(set))
		for _, item := range set {
			if !cmutil.ListStringsContain(u.Changes, item) {
				result = append(result, item)
			}
		}
		return result, len(set) - len(result), nil

	case OpReplace:
		if len(u.Changes) != 2 {
			return nil, 0, xerrors.Errorf(xerrors.InvalidArgument, nil, "replace operator expects 2 params")
		}
		result, n := cmutil.ListStringsRemoveAll(set, u.Changes[0])
		result = append(result, u.Changes[1])
		return result, n, nil

	case OpReplaceAll:
		return set, 0, nil

	case OpDeleteAll:
		return []string{}, 0, nil

	default:
		return nil, 0, xerrors.Errorf(xerrors.InvalidArgument, nil, "invalid update operator")
	}
}
