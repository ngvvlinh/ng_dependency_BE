package convert

import (
	"etop.vn/api/shopping/customering"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/shopping/customering/model"
)

func CreateShopCustomer(in *customering.CreateCustomerArgs) (out *customering.ShopCustomer) {
	if in == nil {
		return nil
	}
	return &customering.ShopCustomer{
		ID:       cm.NewID(),
		ShopID:   in.ShopID,
		Code:     in.Code,
		FullName: in.FullName,
		Gender:   in.Gender,
		Type:     in.Type,
		Birthday: in.Birthday,
		Note:     in.Note,
		Phone:    in.Phone,
		Email:    in.Email,
		Status:   1,
	}
}

func UpdateShopCustomer(in *customering.ShopCustomer, update *customering.UpdateCustomerArgs) (out *customering.ShopCustomer) {
	if in == nil {
		return nil
	}
	return &customering.ShopCustomer{
		ID:        in.ID,
		ShopID:    in.ShopID,
		Code:      update.Code.Apply(in.Code),
		FullName:  update.FullName.Apply(in.FullName),
		Gender:    update.Gender.Apply(in.Gender),
		Type:      update.Type.Apply(in.Type),
		Birthday:  update.Birthday.Apply(in.Birthday),
		Note:      update.Note.Apply(in.Note),
		Phone:     update.Phone.Apply(in.Phone),
		Email:     update.Email.Apply(in.Email),
		Status:    in.Status,
		CreatedAt: in.CreatedAt,
	}
}

func ShopCustomer(in *model.ShopCustomer) (out *customering.ShopCustomer) {
	if in == nil {
		return nil
	}
	return &customering.ShopCustomer{
		ID:        in.ID,
		ShopID:    in.ShopID,
		Code:      in.Code,
		FullName:  in.FullName,
		Gender:    in.Gender,
		Type:      in.Type,
		Birthday:  in.Birthday,
		Note:      in.Note,
		Phone:     in.Phone,
		Email:     in.Email,
		Status:    in.Status,
		CreatedAt: in.CreatedAt,
		UpdatedAt: in.UpdatedAt,
	}
}

func ShopCustomerDB(in *customering.ShopCustomer) (out *model.ShopCustomer) {
	if in == nil {
		return nil
	}
	return &model.ShopCustomer{
		ID:        in.ID,
		ShopID:    in.ShopID,
		Code:      in.Code,
		FullName:  in.FullName,
		Gender:    in.Gender,
		Type:      in.Type,
		Birthday:  in.Birthday,
		Note:      in.Note,
		Phone:     in.Phone,
		Email:     in.Email,
		Status:    in.Status,
		CreatedAt: in.CreatedAt,
		UpdatedAt: in.UpdatedAt,
	}
}

func ShopCustomers(ins []*model.ShopCustomer) (outs []*customering.ShopCustomer) {
	outs = make([]*customering.ShopCustomer, len(ins))
	for i, in := range ins {
		outs[i] = ShopCustomer(in)
	}
	return outs
}
