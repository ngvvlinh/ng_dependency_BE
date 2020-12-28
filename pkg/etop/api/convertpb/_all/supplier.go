package _all

import (
	"o.o/api/shopping/suppliering"
	apishop "o.o/api/top/int/shop"
	"o.o/backend/pkg/common/apifw/cmapi"
)

func PbSupplier(m *suppliering.ShopSupplier) *apishop.Supplier {
	if m == nil {
		return nil
	}
	return &apishop.Supplier{
		Id:                m.ID,
		ShopId:            m.ShopID,
		FullName:          m.FullName,
		Note:              m.Note,
		Code:              m.Code,
		Phone:             m.Phone,
		Email:             m.Email,
		CompanyName:       m.CompanyName,
		TaxNumber:         m.TaxNumber,
		HeadquaterAddress: m.HeadquaterAddress,

		Status:    m.Status,
		CreatedAt: cmapi.PbTime(m.CreatedAt),
		UpdatedAt: cmapi.PbTime(m.UpdatedAt),
	}
}

func PbSuppliers(ms []*suppliering.ShopSupplier) []*apishop.Supplier {
	res := make([]*apishop.Supplier, len(ms))
	for i, m := range ms {
		res[i] = PbSupplier(m)
	}
	return res
}
