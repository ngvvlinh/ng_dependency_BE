package types

type Carrier int32

const (
	Default Carrier = 0
	Ahamove Carrier = 1
)

var Carrier_name = map[int32]string{
	0: "default",
	1: "ahamove",
}

var Carrier_value = map[string]int32{
	"default": 0,
	"ahamove": 1,
}

func CarrierToString(s Carrier) string {
	if s == 0 {
		return ""
	}
	return Carrier_name[int32(s)]
}

func CarrierFromString(s string) Carrier {
	st := Carrier_value[s]
	return Carrier(st)
}

func (c Carrier) String() string {
	return Carrier_name[int32(c)]
}
