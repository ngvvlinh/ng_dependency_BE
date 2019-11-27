package source

import (
	"etop.vn/common/jsonx"
)

type Source int

const (
	Source_unknown  Source = 0
	Source_self     Source = 1
	Source_import   Source = 2
	Source_api      Source = 3
	Source_etop_pos Source = 5
	Source_etop_pxs Source = 6
	Source_etop_cmx Source = 7
	Source_ts_app   Source = 8
)

var Source_name = map[int]string{
	0: "unknown",
	1: "self",
	2: "import",
	3: "api",
	5: "etop_pos",
	6: "etop_pxs",
	7: "etop_cmx",
	8: "ts_app",
}

var Source_value = map[string]int{
	"unknown":  0,
	"self":     1,
	"import":   2,
	"api":      3,
	"etop_pos": 5,
	"etop_pxs": 6,
	"etop_cmx": 7,
	"ts_app":   8,
}

func (x Source) Enum() *Source {
	p := new(Source)
	*p = x
	return p
}

func (x Source) String() string {
	return jsonx.EnumName(Source_name, int(x))
}

func (x *Source) UnmarshalJSON(data []byte) error {
	value, err := jsonx.UnmarshalJSONEnum(Source_value, data, "Source")
	if err != nil {
		return err
	}
	*x = Source(value)
	return nil
}
