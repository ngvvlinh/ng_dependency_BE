package cmutil

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
