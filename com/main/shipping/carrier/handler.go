package carrier

import (
	"context"

	"etop.vn/api/main/connectioning"
	"etop.vn/api/main/location"
	"etop.vn/api/top/int/types"
	carriertypes "etop.vn/backend/com/main/shipping/carrier/types"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/capi/dot"
)

/*
	- Sign In to connection using driver.
	- Create order Update Shop Connection.
*/
func (m *ShipmentManager) ShopConnectionSignIn(ctx context.Context, args *ShopConnectionSignInArgs) (*connectioning.ShopConnection, error) {
	signInCmd := &ConnectionSignInArgs{
		ConnectionID: args.ConnectionID,
		Email:        args.Email,
		Password:     args.Password,
	}
	account, err := m.SignIn(ctx, signInCmd)
	if err != nil {
		return nil, err
	}

	cmd2 := &connectioning.CreateOrUpdateShopConnectionCommand{
		ShopID:       args.ShopID,
		ConnectionID: args.ConnectionID,
		Token:        account.Token,
		ExternalData: &connectioning.ShopConnectionExternalData{
			UserID: account.UserID,
			Email:  args.Email,
		},
	}
	if err := m.connectionAggr.Dispatch(ctx, cmd2); err != nil {
		return nil, err
	}
	return cmd2.Result, nil
}

/*
	- Sign Up to connection using driver.
	- Create order Update Shop Connection.
*/
func (m *ShipmentManager) ShopConnectionSignUp(ctx context.Context, args *ShopConnectionSignUpArgs) (*connectioning.ShopConnection, error) {
	signUpCmd := &ConnectionSignUpArgs{
		ConnectionID: args.ConnectionID,
		Name:         args.Name,
		Email:        args.Email,
		Password:     args.Password,
		Phone:        args.Phone,
		Province:     args.Province,
		District:     args.District,
		Address:      args.Address,
	}
	newAccount, err := m.SignUp(ctx, signUpCmd)
	if err != nil {
		return nil, err
	}
	if newAccount.Token == "" {
		// Sign In
		signInCmd := &ConnectionSignInArgs{
			ConnectionID: args.ConnectionID,
			Email:        args.Email,
			Password:     args.Password,
		}
		newAccount, err = m.SignIn(ctx, signInCmd)
		if err != nil {
			return nil, err
		}
	}

	cmd2 := &connectioning.CreateOrUpdateShopConnectionCommand{
		ShopID:       args.ShopID,
		ConnectionID: args.ConnectionID,
		Token:        newAccount.Token,
		ExternalData: &connectioning.ShopConnectionExternalData{
			UserID: newAccount.UserID,
			Email:  args.Email,
		},
	}
	if err := m.connectionAggr.Dispatch(ctx, cmd2); err != nil {
		return nil, err
	}
	return cmd2.Result, nil
}

func (m *ShipmentManager) PrepareDataGetShippingService(ctx context.Context, accountID dot.ID, q *types.GetShippingServicesRequest) (*carriertypes.GetShippingServicesArgs, error) {
	fromQuery := &location.FindOrGetLocationQuery{
		Province:     q.FromProvince,
		District:     q.FromDistrict,
		ProvinceCode: q.FromProvinceCode,
		DistrictCode: q.FromDistrictCode,
	}
	toQuery := &location.FindOrGetLocationQuery{
		Province:     q.ToProvince,
		District:     q.ToDistrict,
		ProvinceCode: q.ToProvinceCode,
		DistrictCode: q.ToDistrictCode,
	}
	if err := m.LocationQS.Dispatch(ctx, fromQuery); err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "Địa chỉ gửi không hợp lệ: %v", err)
	}
	if err := m.LocationQS.Dispatch(ctx, toQuery); err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "Địa chỉ nhận không hợp lệ: %v", err)
	}

	fromDistrict, _ := fromQuery.Result.District, fromQuery.Result.Province
	topDistrict, _ := toQuery.Result.District, toQuery.Result.Province
	if fromDistrict == nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Địa chỉ gửi không hợp lệ")
	}
	if topDistrict == nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Địa chỉ nhận không hợp lệ")
	}

	weight := q.GrossWeight
	length := q.Length
	width := q.Width
	height := q.Height
	chargeableWeight := q.ChargeableWeight
	calculatedChargeableWeight := model.CalcChargeableWeight(weight, length, width, height)
	if chargeableWeight == 0 {
		chargeableWeight = calculatedChargeableWeight
	} else if chargeableWeight < calculatedChargeableWeight {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Khối lượng tính phí không hợp lệ.").
			WithMetap("chargeable_weight", chargeableWeight).
			WithMetap("gross_weight", q.GrossWeight).
			WithMetap("volumetric_weight (= length*width*height / 5)", length*width*height/5).
			WithMetap("expected chargeable_weight (= MAX(gross_weight, volumetric_weight))", calculatedChargeableWeight)
	}

	return &carriertypes.GetShippingServicesArgs{
		AccountID:        accountID,
		FromDistrictCode: fromDistrict.Code,
		ToDistrictCode:   topDistrict.Code,
		ChargeableWeight: chargeableWeight,
		Length:           length,
		Width:            width,
		Height:           height,
		IncludeInsurance: q.IncludeInsurance.Apply(false),
		BasketValue:      q.BasketValue,
		CODAmount:        q.TotalCodAmount,
	}, nil
}
