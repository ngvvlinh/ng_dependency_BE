package shop

import (
	"context"

	"etop.vn/backend/com/main/shipping/carrier"
)

func (s *ConnectionService) RegisterShopConnection(ctx context.Context, q *RegisterShopConnectionEndpoint) error {
	cmd := &carrier.ShopConnectionSignUpArgs{
		ConnectionID: q.ConnectionID,
		ShopID:       q.Context.Shop.ID,
		Name:         q.Name,
		Email:        q.Email,
		Password:     q.Password,
		Phone:        q.Phone,
		Province:     q.Province,
		District:     q.District,
		Address:      q.Address,
	}
	shopConnection, err := shipmentManager.ShopConnectionSignUp(ctx, cmd)
	if err != nil {
		return err
	}
	q.Result = PbShopConnection(shopConnection)
	return nil
}

func (s *ConnectionService) LoginShopConnection(ctx context.Context, q *LoginShopConnectionEndpoint) error {
	cmd := &carrier.ShopConnectionSignInArgs{
		ConnectionID: q.ConnectionID,
		ShopID:       q.Context.Shop.ID,
		Email:        q.Email,
		Password:     q.Password,
	}
	shopConnection, err := shipmentManager.ShopConnectionSignIn(ctx, cmd)
	if err != nil {
		return err
	}
	q.Result = PbShopConnection(shopConnection)
	return nil
}
