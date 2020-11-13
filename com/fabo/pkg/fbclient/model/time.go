package model

import "time"

type FacebookTime int64

func (t *FacebookTime) ToTime() time.Time {
	if t == nil {
		return time.Time{}
	}
	return time.Unix(int64(*t), 0)
}

func (t *FacebookTime) WebhookTimeToTime() time.Time {
	if t == nil {
		return time.Time{}
	}
	return time.Unix(int64(*t/1000), 0)
}
