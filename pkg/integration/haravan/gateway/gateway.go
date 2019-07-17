package gateway

import (
	"strconv"

	haravanidentity "etop.vn/api/external/haravan/identity"

	"etop.vn/api/external/haravan/gateway"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/httpx"
	"etop.vn/common/l"
)

var ll = l.New()

type Callback struct {
	haravan           gateway.CommandBus
	haravanIdentityQS haravanidentity.QueryBus
}

func New(haravancb gateway.CommandBus, haravanIdentityQueryService haravanidentity.QueryBus) *Callback {
	return &Callback{
		haravan:           haravancb,
		haravanIdentityQS: haravanIdentityQueryService,
	}
}

func (cb *Callback) GetShippingRates(c *httpx.Context) error {
	ctx := c.Req.Context()
	account, err := cb.GetXAccountHaravan(c)
	if err != nil {
		return err
	}

	var cmd gateway.GetShippingRateCommand
	if err := c.DecodeJson(&cmd); err != nil {
		return cm.Errorf(cm.InvalidArgument, err, "Haravan: Can not decode JSON data in GetShippingRate")
	}

	cmd.EtopShopID = account.ShopID
	if err := cb.haravan.Dispatch(ctx, &cmd); err != nil {
		return err
	}
	res := cmd.Result
	c.SetResult(res)

	return nil
}

func (cb *Callback) CreateOrder(c *httpx.Context) error {
	ctx := c.Req.Context()
	account, err := cb.GetXAccountHaravan(c)
	if err != nil {
		return err
	}

	var cmd gateway.CreateOrderCommand
	if err := c.DecodeJson(&cmd); err != nil {
		return cm.Errorf(cm.InvalidArgument, err, "Haravan: Can not decode JSON data in CreateOrder")
	}

	cmd.EtopShopID = account.ShopID
	if err := cb.haravan.Dispatch(ctx, &cmd); err != nil {
		return err
	}
	res := cmd.Result
	c.SetResult(res)
	return nil
}

func (cb *Callback) GetOrder(c *httpx.Context) error {
	ctx := c.Req.Context()
	account, err := cb.GetXAccountHaravan(c)
	if err != nil {
		return err
	}

	var cmd gateway.GetOrderCommand
	if err := c.DecodeJson(&cmd); err != nil {
		return cm.Errorf(cm.InvalidArgument, err, "Haravan: Can not decode JSON data in GetOrder")
	}

	cmd.EtopShopID = account.ShopID
	if err := cb.haravan.Dispatch(ctx, &cmd); err != nil {
		return err
	}
	res := cmd.Result
	c.SetResult(res)
	return nil
}

func (cb *Callback) CancelOrder(c *httpx.Context) error {
	ctx := c.Req.Context()
	account, err := cb.GetXAccountHaravan(c)
	if err != nil {
		return err
	}

	var cmd gateway.CancelOrderCommand
	if err := c.DecodeJson(&cmd); err != nil {
		return cm.Errorf(cm.InvalidArgument, err, "Haravan: Can not decode JSON data in CancelOrder")
	}

	cmd.EtopShopID = account.ShopID
	if err := cb.haravan.Dispatch(ctx, &cmd); err != nil {
		return err
	}
	res := cmd.Result
	c.SetResult(res)
	return nil
}

func (cb *Callback) GetXAccountHaravan(c *httpx.Context) (*haravanidentity.ExternalAccountHaravan, error) {
	id := c.Params.ByName("shopid")
	shopID, err := strconv.Atoi(id)

	if err != nil || shopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing Shop ID")
	}

	query := &haravanidentity.GetExternalAccountHaravanByXShopIDQuery{
		ExternalShopID: shopID,
	}
	if err := cb.haravanIdentityQS.Dispatch(c.Req.Context(), query); err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "Shop does not valid")
	}
	return query.Result, nil
}
