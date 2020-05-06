package model

import "time"

type FacebookTime int

func (t *FacebookTime) ToTime() time.Time {
	if t == nil {
		return time.Time{}
	}
	return time.Unix(int64(*t), 0)
}
