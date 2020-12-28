package _min

import (
	"o.o/api/shopping/customering"
	exttypes "o.o/api/top/external/types"
	"o.o/capi/dot"
)

func PbShopCustomer(customer *customering.ShopCustomer) *exttypes.Customer {
	if customer == nil {
		return nil
	}
	if customer.Deleted {
		return &exttypes.Customer{
			Id:      customer.ID,
			Deleted: true,
		}
	}
	return &exttypes.Customer{
		Id:           customer.ID,
		ShopId:       customer.ShopID,
		ExternalId:   dot.String(customer.ExternalID),
		ExternalCode: dot.String(customer.ExternalCode),
		FullName:     dot.String(customer.FullName),
		Code:         dot.String(customer.Code),
		Note:         dot.String(customer.Note),
		Phone:        dot.String(customer.Phone),
		Email:        dot.String(customer.Email),
		Gender:       customer.Gender.Wrap(),
		Type:         customer.Type.Wrap(),
		Birthday:     dot.String(customer.Birthday),
		CreatedAt:    dot.Time(customer.CreatedAt),
		UpdatedAt:    dot.Time(customer.UpdatedAt),
		Status:       customer.Status.Wrap(),
		Deleted:      customer.Deleted,
	}
}
