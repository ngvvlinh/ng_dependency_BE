package meta

import (
	uuid "github.com/satori/go.uuid"
	cmutil "o.o/capi/util"
	"o.o/common/xerrors"
)

type EventMeta struct {
	ID string
}

func NewEvent() EventMeta {
	return EventMeta{
		ID: uuid.NewV4().String(),
	}
}

type Empty struct{}

type UUID struct{ Data []byte }

type UpdatedResponse struct {
	Updated int
}

func NewUUID() UUID {
	u := uuid.NewV4()
	return UUID{u[:]}
}

type PageInfo struct {
	Offset int
	Limit  int
	Sort   []string

	Before string
	After  string

	Next string
	Prev string
}

type Paging struct {
	Offset int
	Limit  int
	Sort   []string

	Before string
	After  string
}

type Filter struct {
	// Comma separated properties: "name,s_name"
	Name string
	// Can be = ≠ (!=) < ≤ (<=) > ≥ (>=) ⊃ (c) ∈ (in) ∩ (n)
	//
	// - Text or set: ⊃ ∩
	// - Exactly: = ≠ ∈
	// - Numeric: = ≠ ∈ < ≤ > ≥
	Op string
	// Must always be string
	Value string
}

type Filters []Filter

func (f Filters) ListFilters() []Filter {
	return f
}

type Error struct {
	Code string
	Msg  string
	Meta map[string]string
}

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

type UpdateListInt struct {
	Op      UpdateOp
	Changes []int
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
		return u.Changes, 0, nil

	case OpDeleteAll:
		return []string{}, 0, nil

	default:
		return nil, 0, xerrors.Errorf(xerrors.InvalidArgument, nil, "invalid update operator")
	}
}
