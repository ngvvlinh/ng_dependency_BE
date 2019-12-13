package convert

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"etop.vn/api/shopping/addressing"
	"etop.vn/api/shopping/customering"
	"etop.vn/api/top/types/etc/status3"
	orderconvert "etop.vn/backend/com/main/ordering/convert"
	"etop.vn/backend/com/shopping/customering/model"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/validate"
	"etop.vn/common/l"
)

// +gen:convert: etop.vn/backend/com/shopping/customering/model -> etop.vn/api/shopping/customering,etop.vn/api/shopping/addressing,etop.vn/api/shopping/tradering
// +gen:convert: etop.vn/api/shopping/customering
// +gen:convert: etop.vn/api/shopping/addressing

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
		ID:       cm.NewID(),
		ShopID:   args.ShopID,
		FullName: args.FullName,
		Gender:   args.Gender,
		Type:     args.Type,
		Birthday: args.Birthday,
		Note:     args.Note,
		Phone:    args.Phone,
		Email:    args.Email,
		Status:   1,
	}
	result.PhoneNorm = validate.NormalizeSearchPhone(args.Phone)
	result.FullNameNorm = validate.NormalizeSearch(args.FullName)
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
	out.FullNameNorm = validate.NormalizeSearch(args.FullName)
}

func shopCustomer(args *model.ShopCustomer, out *customering.ShopCustomer) {
	convert_customeringmodel_ShopCustomer_customering_ShopCustomer(args, out)
	if out.Type == customering.CustomerTypeIndependent {
		out.FullName = "Khách lẻ"
	}
}

func ShopTraderAddress(args *model.ShopTraderAddress, out *addressing.ShopTraderAddress) {
	convert_customeringmodel_ShopTraderAddress_addressing_ShopTraderAddress(args, out)
	out.Coordinates = orderconvert.Coordinates(args.Coordinates)
}

func ShopTraderAddressDB(args *addressing.ShopTraderAddress, out *model.ShopTraderAddress) {
	convert_addressing_ShopTraderAddress_customeringmodel_ShopTraderAddress(args, out)
	out.Coordinates = orderconvert.CoordinatesDB(args.Coordinates)
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
