package convert

import (
	"time"

	"etop.vn/api/main/etop"
	"etop.vn/api/shopping/addressing"
	"etop.vn/api/shopping/customering"
	orderconvert "etop.vn/backend/com/main/ordering/convert"
	"etop.vn/backend/com/shopping/customering/model"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/validate"
)

// +gen:convert: etop.vn/backend/com/shopping/customering/model -> etop.vn/api/shopping/customering
// +gen:convert: etop.vn/api/shopping/customering

func CreateShopCustomer(args *customering.CreateCustomerArgs) (out *customering.ShopCustomer) {
	if args == nil {
		return nil
	}
	return &customering.ShopCustomer{
		ID:       cm.NewID(),
		ShopID:   args.ShopID,
		Code:     args.Code,
		FullName: args.FullName,
		Gender:   args.Gender,
		Type:     args.Type,
		Birthday: args.Birthday,
		Note:     args.Note,
		Phone:    args.Phone,
		Email:    args.Email,
		Status:   1,
	}
}

func UpdateShopCustomer(in *customering.ShopCustomer, args *customering.UpdateCustomerArgs) (out *customering.ShopCustomer) {
	if in == nil {
		return nil
	}
	return &customering.ShopCustomer{
		ID:        in.ID,
		ShopID:    in.ShopID,
		Code:      args.Code.Apply(in.Code),
		FullName:  args.FullName.Apply(in.FullName),
		Gender:    args.Gender.Apply(in.Gender),
		Type:      args.Type.Apply(in.Type),
		Birthday:  args.Birthday.Apply(in.Birthday),
		Note:      args.Note.Apply(in.Note),
		Phone:     args.Phone.Apply(in.Phone),
		Email:     args.Email.Apply(in.Email),
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
		ID:           in.ID,
		ShopID:       in.ShopID,
		Code:         in.Code,
		FullName:     in.FullName,
		Gender:       in.Gender,
		Type:         in.Type,
		Birthday:     in.Birthday,
		Note:         in.Note,
		Phone:        in.Phone,
		Email:        in.Email,
		PhoneNorm:    validate.NormalizeSearchPhone(in.Phone),
		Status:       in.Status,
		CreatedAt:    in.CreatedAt,
		UpdatedAt:    in.UpdatedAt,
		FullNameNorm: validate.NormalizeSearch(in.FullName),
	}
}

func ShopCustomers(ins []*model.ShopCustomer) (outs []*customering.ShopCustomer) {
	outs = make([]*customering.ShopCustomer, len(ins))
	for i, in := range ins {
		outs[i] = ShopCustomer(in)
	}
	return outs
}

func ShopTraderAddress(in *model.ShopTraderAddress) (out *addressing.ShopTraderAddress) {
	if in == nil {
		return nil
	}
	return &addressing.ShopTraderAddress{
		ID:           in.ID,
		ShopID:       in.ShopID,
		TraderID:     in.TraderID,
		FullName:     in.FullName,
		Phone:        in.Phone,
		Email:        in.Email,
		Company:      in.Company,
		Address1:     in.Address1,
		Address2:     in.Address2,
		DistrictCode: in.DistrictCode,
		WardCode:     in.WardCode,
		Coordinates:  orderconvert.Coordinates(in.Coordinates),
		CreatedAt:    in.CreatedAt,
		UpdatedAt:    in.UpdatedAt,
		IsDefault:    in.IsDefault,
	}
}

func Addresses(ins []*model.ShopTraderAddress) (outs []*addressing.ShopTraderAddress) {
	outs = make([]*addressing.ShopTraderAddress, len(ins))
	for i, in := range ins {
		outs[i] = ShopTraderAddress(in)
	}
	return outs
}

func ShopTraderAddressDB(in *addressing.ShopTraderAddress) (out *model.ShopTraderAddress) {
	if in == nil {
		return nil
	}
	return &model.ShopTraderAddress{
		ID:           in.ID,
		ShopID:       in.ShopID,
		TraderID:     in.TraderID,
		FullName:     in.FullName,
		Phone:        in.Phone,
		Email:        in.Email,
		Company:      in.Company,
		Address1:     in.Address1,
		Address2:     in.Address2,
		DistrictCode: in.DistrictCode,
		WardCode:     in.WardCode,
		Coordinates:  orderconvert.CoordinatesDB(in.Coordinates),
		CreatedAt:    in.CreatedAt,
		UpdatedAt:    in.UpdatedAt,
		IsDefault:    in.IsDefault,

		// currently not use TODO
		Status: etop.S3Positive,
	}
}

func CreateShopTraderAddress(in *addressing.CreateAddressArgs) (out *addressing.ShopTraderAddress) {
	if in == nil {
		return nil
	}
	return &addressing.ShopTraderAddress{
		ID:           cm.NewID(),
		ShopID:       in.ShopID,
		TraderID:     in.TraderID,
		FullName:     in.FullName,
		Phone:        in.Phone,
		Email:        in.Email,
		Company:      in.Company,
		Address1:     in.Address1,
		Address2:     in.Address2,
		DistrictCode: in.DistrictCode,
		WardCode:     in.WardCode,
		Coordinates:  in.Coordinates,
		IsDefault:    in.IsDefault,
	}
}

func UpdateShopTraderAddress(in *addressing.ShopTraderAddress, update *addressing.UpdateAddressArgs) (out *addressing.ShopTraderAddress) {
	if in == nil {
		return nil
	}
	out = &addressing.ShopTraderAddress{
		ID:           in.ID,
		ShopID:       in.ShopID,
		TraderID:     in.TraderID,
		FullName:     update.FullName.Apply(in.FullName),
		Phone:        update.Phone.Apply(in.Phone),
		Email:        update.Email.Apply(in.Email),
		Company:      update.Company.Apply(in.Company),
		Address1:     update.Address1.Apply(in.Address1),
		Address2:     update.Address2.Apply(in.Address2),
		DistrictCode: update.DistrictCode.Apply(in.DistrictCode),
		WardCode:     update.WardCode.Apply(in.WardCode),
		Coordinates:  in.Coordinates,
		IsDefault:    in.IsDefault,
	}
	if update.Coordinates != nil {
		out.Coordinates = update.Coordinates
	}
	return out
}

func ShopCustomerGroup(in *model.ShopCustomerGroup, out *customering.ShopCustomerGroup) {
	convert_customeringmodel_ShopCustomerGroup_customering_ShopCustomerGroup(in, out)
}

func ShopCustomerGroupDB(in *customering.ShopCustomerGroup, out *model.ShopCustomerGroup) {
	convert_customering_ShopCustomerGroup_customeringmodel_ShopCustomerGroup(in, out)
}

func ShopCustomerGroups(ins []*model.ShopCustomerGroup) (outs []*customering.ShopCustomerGroup) {
	return Convert_customeringmodel_ShopCustomerGroups_customering_ShopCustomerGroups(ins)
}

func ShopCustomerGroupCustomer(in *model.ShopCustomerGroupCustomer, out *customering.ShopCustomerGroupCustomer) {
	convert_customeringmodel_ShopCustomerGroupCustomer_customering_ShopCustomerGroupCustomer(in, out)
}

func ShopCustomerGroupCustomerDB(in *customering.ShopCustomerGroupCustomer, out *model.ShopCustomerGroupCustomer) {
	convert_customering_ShopCustomerGroupCustomer_customeringmodel_ShopCustomerGroupCustomer(in, out)
}

func UpdateCustomerGroup(in *model.ShopCustomerGroup, update *customering.UpdateCustomerGroupArgs) (out *model.ShopCustomerGroup) {
	if in == nil {
		return nil
	}
	return &model.ShopCustomerGroup{
		ID:        update.ID,
		Name:      update.Name,
		CreatedAt: in.CreatedAt,
		UpdatedAt: time.Now(),
	}
}
