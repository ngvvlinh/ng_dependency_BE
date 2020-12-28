package _all

import (
	"o.o/api/shopping/carrying"
	apishop "o.o/api/top/int/shop"
	"o.o/backend/pkg/common/apifw/cmapi"
)

func PbCarrier(m *carrying.ShopCarrier) *apishop.Carrier {
	return &apishop.Carrier{
		Id:        m.ID,
		ShopId:    m.ShopID,
		FullName:  m.FullName,
		Note:      m.Note,
		Status:    m.Status,
		CreatedAt: cmapi.PbTime(m.CreatedAt),
		UpdatedAt: cmapi.PbTime(m.UpdatedAt),
	}
}

func PbCarriers(ms []*carrying.ShopCarrier) []*apishop.Carrier {
	res := make([]*apishop.Carrier, len(ms))
	for i, m := range ms {
		res[i] = PbCarrier(m)
	}
	return res
}
