package util

import (
	"unsafe"

	"o.o/capi/dot"
)

func Int64ToIDs(ids []int64) []dot.ID {
	return *(*[]dot.ID)(unsafe.Pointer(&ids))
}

func IDsToInt64(ids []dot.ID) []int64 {
	return *(*[]int64)(unsafe.Pointer(&ids))
}

func CoalesceString(ss ...string) string {
	for _, s := range ss {
		if s != "" {
			return s
		}
	}
	return ""
}

func CoalesceInt32(ints ...int32) int32 {
	for _, i := range ints {
		if i != 0 {
			return i
		}
	}
	return 0
}

func CoalesceInt(ints ...int) int {
	for _, i := range ints {
		if i != 0 {
			return i
		}
	}
	return 0
}

func ListStringsContain(list []string, item string) bool {
	for _, x := range list {
		if x == item {
			return true
		}
	}
	return false
}

func ListStringsRemove(list []string, item string) ([]string, bool) {
	for i, x := range list {
		if x == item {
			return append(list[:i], list[i+1:]...), true
		}
	}
	return list, false
}

func ListStringsRemoveAll(list []string, item string) ([]string, int) {
	found := false
	for _, x := range list {
		if x == item {
			found = true
			break
		}
	}
	if !found {
		return list, 0
	}
	result := make([]string, 0, len(list))
	for _, x := range list {
		if x != item {
			result = append(result, x)
		}
	}
	return result, len(list) - len(result)
}
