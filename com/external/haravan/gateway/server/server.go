package server

import (
	"strconv"

	"etop.vn/api/external/haravan/gateway"
	haravanidentity "etop.vn/api/external/haravan/identity"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/apifw/httpx"
	"etop.vn/common/l"
)

var ll = l.New()

type Server struct {
	haravan           gateway.CommandBus
	haravanIdentityQS haravanidentity.QueryBus
}

func New(haravancb gateway.CommandBus, haravanIdentityQueryService haravanidentity.QueryBus) *Server {
	return &Server{
		haravan:           haravancb,
		haravanIdentityQS: haravanIdentityQueryService,
	}
}

func SetResult(c *httpx.Context, data interface{}, err error) {
	var res *CommonResponse
	if err != nil {
		res = ErrorResponse(err.Error())
	} else {
		res = SuccessResponse(data)
	}
	c.SetResult(res)
}

func (cb *Server) GetShippingRates(c *httpx.Context) error {
	ctx := c.Req.Context()
	account, err := cb.getExternalAccountHaravan(c)
	if err != nil {
		SetResult(c, nil, err)
		return nil
	}

	var cmd gateway.GetShippingRateCommand
	if _err := c.DecodeJson(&cmd); _err != nil {
		_err = cm.Errorf(cm.InvalidArgument, err, "Haravan: Can not decode JSON data in GetShippingRate")
		SetResult(c, nil, _err)
		return nil
	}

	cmd.EtopShopID = account.ShopID
	if _err := cb.haravan.Dispatch(ctx, &cmd); _err != nil {
		SetResult(c, nil, _err)
		return nil
	}
	SetResult(c, cmd.Result, err)
	return nil
}

func (cb *Server) CreateOrder(c *httpx.Context) error {
	ctx := c.Req.Context()
	account, err := cb.getExternalAccountHaravan(c)
	if err != nil {
		SetResult(c, nil, err)
		return nil
	}

	var cmd gateway.CreateOrderCommand
	if _err := c.DecodeJson(&cmd); _err != nil {
		_err = cm.Errorf(cm.InvalidArgument, err, "Haravan: Can not decode JSON data in CreateOrder")
		SetResult(c, nil, _err)
		return nil
	}

	cmd.EtopShopID = account.ShopID
	if _err := cb.haravan.Dispatch(ctx, &cmd); _err != nil {
		SetResult(c, nil, _err)
		return nil
	}
	SetResult(c, cmd.Result, nil)
	return nil
}

func (cb *Server) GetOrder(c *httpx.Context) error {
	ctx := c.Req.Context()
	account, err := cb.getExternalAccountHaravan(c)
	if err != nil {
		SetResult(c, nil, err)
		return nil
	}

	var cmd gateway.GetOrderCommand
	if _err := c.DecodeJson(&cmd); _err != nil {
		_err = cm.Errorf(cm.InvalidArgument, err, "Haravan: Can not decode JSON data in GetOrder")
		SetResult(c, nil, _err)
		return nil
	}

	cmd.EtopShopID = account.ShopID
	if _err := cb.haravan.Dispatch(ctx, &cmd); _err != nil {
		SetResult(c, nil, _err)
		return nil
	}
	SetResult(c, cmd.Result, nil)
	return nil
}

func (cb *Server) CancelOrder(c *httpx.Context) error {
	ctx := c.Req.Context()
	account, err := cb.getExternalAccountHaravan(c)
	if err != nil {
		SetResult(c, nil, err)
		return nil
	}

	var cmd gateway.CancelOrderCommand
	if _err := c.DecodeJson(&cmd); _err != nil {
		_err = cm.Errorf(cm.InvalidArgument, err, "Haravan: Can not decode JSON data in CancelOrder")
		SetResult(c, nil, _err)
		return nil
	}

	cmd.EtopShopID = account.ShopID
	if _err := cb.haravan.Dispatch(ctx, &cmd); _err != nil {
		SetResult(c, nil, _err)
		return nil
	}
	SetResult(c, cmd.Result, nil)
	return nil
}

func (cb *Server) getExternalAccountHaravan(c *httpx.Context) (*haravanidentity.ExternalAccountHaravan, error) {
	id := c.Params.ByName("shopid")
	shopID, err := strconv.Atoi(id)

	if err != nil || shopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing Shop ID")
	}

	query := &haravanidentity.GetExternalAccountHaravanByXShopIDQuery{
		ExternalShopID: shopID,
	}
	if err = cb.haravanIdentityQS.Dispatch(c.Req.Context(), query); err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "Shop does not valid")
	}
	return query.Result, nil
}

// follow the response format from Haravan
type CommonResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func SuccessResponse(data interface{}) *CommonResponse {
	return &CommonResponse{
		Error:   false,
		Message: "",
		Data:    data,
	}
}

func ErrorResponse(msg string) *CommonResponse {
	return &CommonResponse{
		Error:   true,
		Message: msg,
		Data:    nil,
	}
}
