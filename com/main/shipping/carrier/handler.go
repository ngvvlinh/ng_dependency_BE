package carrier

import (
	"context"
	"fmt"

	"o.o/api/main/connectioning"
	"o.o/api/main/location"
	"o.o/api/top/int/types"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/etop/model"
)

/*
	- Sign In to connection using driver.
	- Create order Update Shop Connection.
	- Sign in email + password or send OTP (loginShopConnectionWithOTP)
*/
func (m *ShipmentManager) ShopConnectionSignIn(ctx context.Context, args *ShopConnectionSignInArgs) (*types.LoginShopConnectionResponse, error) {
	signInCmd := &ConnectionSignInArgs{
		ConnectionID: args.ConnectionID,
		Identifier:   args.Identifier,
		Password:     args.Password,
	}
	account, err := m.SignIn(ctx, signInCmd)
	if err != nil {
		return nil, err
	}

	// chỉ tạo hoặc cập nhật thông tin shopConnection khi login bằng email + password
	if !account.IsRequiredOTP {
		cmd2 := &connectioning.CreateOrUpdateShopConnectionCommand{
			ShopID:       args.ShopID,
			ConnectionID: args.ConnectionID,
			Token:        account.Token,
			ExternalData: &connectioning.ShopConnectionExternalData{
				UserID:     account.UserID,
				Identifier: args.Identifier,
				ShopID:     account.ShopID,
			},
		}
		if err := m.connectionAggr.Dispatch(ctx, cmd2); err != nil {
			return nil, err
		}
		return &types.LoginShopConnectionResponse{
			Code: "OK",
		}, nil
	}
	return &types.LoginShopConnectionResponse{
		Code: "OTP_SENT",
		Msg:  fmt.Sprintf("đã gửi OTP đến sđt %s", args.Identifier),
	}, nil
}

func (m *ShipmentManager) ShopConnectionSignInWithOTP(ctx context.Context, args *ShopConnectionSignInWithOTPArgs) (*connectioning.ShopConnection, error) {
	if args.OTP == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "OTP không được rỗng.")
	}

	signInCmd := &ConnectionSignInArgs{
		ConnectionID: args.ConnectionID,
		Identifier:   args.Identifier,
		OTP:          args.OTP,
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
			UserID:     account.UserID,
			Identifier: args.Identifier,
			ShopID:     account.ShopID,
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
		Identifier:   args.Identifier,
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
			Identifier:   args.Identifier,
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
			UserID:     newAccount.UserID,
			Identifier: args.Identifier,
		},
	}
	if err := m.connectionAggr.Dispatch(ctx, cmd2); err != nil {
		return nil, err
	}
	return cmd2.Result, nil
}

func (m *ShipmentManager) PrepareDataGetShippingServices(ctx context.Context, q *types.GetShippingServicesRequest) (*GetShippingServicesArgs, error) {
	fromQuery := &location.FindOrGetLocationQuery{
		Province:     q.FromProvince,
		District:     q.FromDistrict,
		Ward:         q.FromWard,
		ProvinceCode: q.FromProvinceCode,
		DistrictCode: q.FromDistrictCode,
		WardCode:     q.FromWardCode,
	}
	toQuery := &location.FindOrGetLocationQuery{
		Province:     q.ToProvince,
		District:     q.ToDistrict,
		Ward:         q.ToWard,
		ProvinceCode: q.ToProvinceCode,
		DistrictCode: q.ToDistrictCode,
		WardCode:     q.ToWardCode,
	}
	if err := m.locationQS.Dispatch(ctx, fromQuery); err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "Địa chỉ gửi không hợp lệ: %v", err)
	}
	if err := m.locationQS.Dispatch(ctx, toQuery); err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "Địa chỉ nhận không hợp lệ: %v", err)
	}

	fromDistrict, fromProvince, fromWard := fromQuery.Result.District, fromQuery.Result.Province, fromQuery.Result.Ward
	toDistrict, toProvince, toWard := toQuery.Result.District, toQuery.Result.Province, toQuery.Result.Ward
	if fromDistrict == nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Địa chỉ gửi không hợp lệ")
	}
	if toDistrict == nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Địa chỉ nhận không hợp lệ")
	}

	weight := q.GrossWeight
	length := q.Length
	width := q.Width
	height := q.Height

	chargeableWeight, err := ValidateFfmWeight(weight, length, width, height, q.ChargeableWeight)
	if err != nil {
		return nil, err
	}

	res := &GetShippingServicesArgs{
		ConnectionIDs:    q.ConnectionIDs,
		FromDistrictCode: fromDistrict.Code,
		FromProvinceCode: fromProvince.Code,
		ToDistrictCode:   toDistrict.Code,
		ToProvinceCode:   toProvince.Code,
		ChargeableWeight: chargeableWeight,
		Length:           length,
		Width:            width,
		Height:           height,
		IncludeInsurance: q.IncludeInsurance.Apply(false),
		InsuranceValue:   q.InsuranceValue,
		BasketValue:      q.BasketValue,
		CODAmount:        q.TotalCodAmount,
		Coupon:           q.Coupon,
	}
	if fromWard != nil {
		res.FromWardCode = fromWard.Code
	}
	if toWard != nil {
		res.ToWardCode = toWard.Code
	}
	return res, nil
}

func ValidateFfmWeight(grossWeight, length, width, height, chargeableWeight int) (int, error) {
	calculatedChargeableWeight := model.CalcChargeableWeight(grossWeight, length, width, height)
	if chargeableWeight == 0 {
		chargeableWeight = calculatedChargeableWeight
	} else if chargeableWeight < calculatedChargeableWeight {
		return 0, cm.Errorf(cm.InvalidArgument, nil, "Khối lượng tính phí không hợp lệ.").
			WithMetap("chargeable_weight", chargeableWeight).
			WithMetap("gross_weight", grossWeight).
			WithMetap("volumetric_weight (= length*width*height / 5)", length*width*height/5).
			WithMetap("expected chargeable_weight (= MAX(gross_weight, volumetric_weight))", calculatedChargeableWeight)
	}
	return chargeableWeight, nil
}
