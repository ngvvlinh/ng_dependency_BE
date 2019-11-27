package types

type Carrier int

const (
	Default Carrier = 0
	Ahamove Carrier = 1
)

var Carrier_name = map[int]string{
	0: "default",
	1: "ahamove",
}

var Carrier_value = map[string]int{
	"default": 0,
	"ahamove": 1,
}

func CarrierToString(s Carrier) string {
	if s == 0 {
		return ""
	}
	return Carrier_name[int(s)]
}

func CarrierFromString(s string) Carrier {
	st := Carrier_value[s]
	return Carrier(st)
}

func (c Carrier) String() string {
	return Carrier_name[int(c)]
}
