package convert

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"o.o/api/shopping/addressing"
	"o.o/api/shopping/customering"
	"o.o/api/top/types/etc/customer_type"
	"o.o/api/top/types/etc/status3"
	addressconvert "o.o/backend/com/main/address/convert"
	"o.o/backend/com/shopping/customering/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/validate"
	"o.o/common/l"
)

// +gen:convert: o.o/backend/com/shopping/customering/model  -> o.o/api/shopping/customering, o.o/api/shopping/addressing, o.o/api/shopping/tradering
// +gen:convert: o.o/api/shopping/customering
// +gen:convert: o.o/api/shopping/addressing

var ll = l.New()

const (
	MaxCodeNorm = 999999
	codeRegex   = "^KH([0-9]{6})$"
	codePrefix  = "KH"
)

var reCode = regexp.MustCompile(codeRegex)

func ParseCodeNorm(code string) (_ int, ok bool) {
	parts := reCode.FindStringSubmatch(code)
	if len(parts) == 0 {
		return 0, false
	}
	number, err := strconv.Atoi(parts[1])
	if err != nil {
		ll.Panic("unexpected", l.Error(err))
	}
	return number, true
}

func GenerateCode(codeNorm int) string {
	return fmt.Sprintf("%v%06v", codePrefix, codeNorm)
}

func CreateShopCustomer(args *customering.CreateCustomerArgs) (out *model.ShopCustomer) {
	if args == nil {
		return nil
	}
	result := &model.ShopCustomer{
		ExternalID:   args.ExternalID,
		ExternalCode: args.ExternalCode,
		PartnerID:    args.PartnerID,
		ID:           cm.NewID(),
		ShopID:       args.ShopID,
		FullName:     args.FullName,
		Gender:       args.Gender,
		Type:         args.Type,
		Birthday:     args.Birthday,
		Note:         args.Note,
		Phone:        args.Phone,
		Email:        args.Email,
		CreatedBy:    args.CreatedBy,
		Status:       1,
	}
	result.PhoneNorm = validate.NormalizeSearchPhone(args.Phone)
	result.FullNameNorm = validate.NormalizeSearchCharacter(args.FullName)
	return result
}

func updateShopCustomer(args *customering.UpdateCustomerArgs, in *customering.ShopCustomer) (out *customering.ShopCustomer) {
	if in == nil {
		return nil
	}
	apply_customering_UpdateCustomerArgs_customering_ShopCustomer(args, in)
	in.UpdatedAt = time.Now()
	return in
}

func shopCustomerDB(args *customering.ShopCustomer, out *model.ShopCustomer) {
	convert_customering_ShopCustomer_customeringmodel_ShopCustomer(args, out)
	out.PhoneNorm = validate.NormalizeSearchPhone(args.Phone)
	out.FullNameNorm = validate.NormalizeSearchCharacter(args.FullName)
}

func shopCustomer(args *model.ShopCustomer, out *customering.ShopCustomer) {
	convert_customeringmodel_ShopCustomer_customering_ShopCustomer(args, out)
	if out.Type == customer_type.Independent {
		out.FullName = "Kh??ch l???"
	}
}

func ShopTraderAddress(args *model.ShopTraderAddress, out *addressing.ShopTraderAddress) {
	convert_customeringmodel_ShopTraderAddress_addressing_ShopTraderAddress(args, out)
	out.Coordinates = addressconvert.Coordinates(args.Coordinates)
}

func ShopTraderAddressDB(args *addressing.ShopTraderAddress, out *model.ShopTraderAddress) {
	convert_addressing_ShopTraderAddress_customeringmodel_ShopTraderAddress(args, out)
	out.Coordinates = addressconvert.CoordinatesDB(args.Coordinates)
	out.Status = status3.P
}

func CreateShopTraderAddress(args *addressing.CreateAddressArgs, out *addressing.ShopTraderAddress) {
	apply_addressing_CreateAddressArgs_addressing_ShopTraderAddress(args, out)
	out.ID = cm.NewID()
}

func UpdateShopTraderAddress(args *addressing.UpdateAddressArgs, out *addressing.ShopTraderAddress) {
	// TODO: coordinates
	coordinates := out.Coordinates
	apply_addressing_UpdateAddressArgs_addressing_ShopTraderAddress(args, out)
	if args.Coordinates == nil {
		out.Coordinates = coordinates
	}
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
