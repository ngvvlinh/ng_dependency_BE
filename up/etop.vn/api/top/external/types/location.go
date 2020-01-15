package types

import "etop.vn/common/jsonx"

type Ward struct {
	Name string `json:"name"`
}

func (m *Ward) String() string { return jsonx.MustMarshalToString(m) }

type District struct {
	Name  string `json:"name"`
	Wards []Ward `json:"wards"`
}

func (m *District) String() string { return jsonx.MustMarshalToString(m) }

type Province struct {
	Name      string     `json:"name"`
	Districts []District `json:"districts"`
}

func (m *Province) String() string { return jsonx.MustMarshalToString(m) }

type LocationResponse struct {
	Provinces []Province `json:"provinces"`
}

func (m *LocationResponse) String() string { return jsonx.MustMarshalToString(m) }

type LocationAddress struct {
	Province string `json:"province"`
	District string `json:"district"`
}

func (m *LocationAddress) String() string { return jsonx.MustMarshalToString(m) }

type Coordinates struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

func (m *Coordinates) String() string { return jsonx.MustMarshalToString(m) }
