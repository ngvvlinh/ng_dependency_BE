package admin

import (
	"o.o/api/etelecom/mobile_network"
	"o.o/api/top/types/etc/status3"
	"o.o/capi/dot"
	"o.o/common/jsonx"
)

type CreateHotlineRequest struct {
	OwnerID      dot.ID                       `json:"owner_id"`
	Name         string                       `json:"name"`
	Hotline      string                       `json:"hotline"`
	Network      mobile_network.MobileNetwork `json:"network"`
	ConnectionID dot.ID                       `json:"connection_id"`
	Description  string                       `json:"description"`
	IsFreeCharge dot.NullBool                 `json:"is_free_charge"`
}

func (m *CreateHotlineRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateHotlineRequest struct {
	ID           dot.ID             `json:"id"`
	Name         string             `json:"name"`
	Description  string             `json:"description"`
	Status       status3.NullStatus `json:"status"`
	IsFreeCharge dot.NullBool       `json:"is_free_charge"`
}

func (m *UpdateHotlineRequest) String() string { return jsonx.MustMarshalToString(m) }
