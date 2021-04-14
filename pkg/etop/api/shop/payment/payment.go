package payment

import (
	"context"
	"fmt"
	"strings"

	"o.o/api/external/payment"
	paymentmanager "o.o/api/external/payment/manager"
	"o.o/api/main/credit"
	api "o.o/api/top/int/shop"
	pbcm "o.o/api/top/types/common"
	"o.o/api/top/types/etc/payment_provider"
	"o.o/api/top/types/etc/payment_state"
	"o.o/api/top/types/etc/status4"
	"o.o/api/top/types/etc/subject_referral"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/etop/authorize/session"
)

type PaymentService struct {
	session.Session

	PaymentManager paymentmanager.CommandBus
	PaymentAggr    payment.CommandBus
	CreditQuery    credit.QueryBus
}

func (s *PaymentService) Clone() api.PaymentService { res := *s; return &res }

func (s *PaymentService) PaymentTradingOrder(ctx context.Context, q *api.PaymentTradingOrderRequest) (*api.PaymentTradingOrderResponse, error) {
	if q.OrderId == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing OrderID")
	}
	if q.ReturnUrl == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing ReturnURL")
	}

	argGenCode := &paymentmanager.GenerateCodeCommand{
		ReferralType: subject_referral.Order,
		ID:           q.OrderId.String(),
	}
	if err := s.PaymentManager.Dispatch(ctx, argGenCode); err != nil {
		return nil, err
	}
	args := &paymentmanager.BuildUrlConnectPaymentGatewayCommand{
		OrderID:           argGenCode.Result,
		Desc:              q.Desc,
		ReturnURL:         q.ReturnUrl,
		TransactionAmount: q.Amount,
		Provider:          q.PaymentProvider,
	}

	if err := s.PaymentManager.Dispatch(ctx, args); err != nil {
		return nil, err
	}
	result := &api.PaymentTradingOrderResponse{
		Url: args.Result,
	}
	return result, nil
}

func (s *PaymentService) PaymentCheckReturnData(ctx context.Context, q *api.PaymentCheckReturnDataRequest) (*pbcm.MessageResponse, error) {
	if q.Id == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Mã giao dịch không được để trống")
	}
	if q.Code == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Mã 'Code' không được để trống")
	}
	args := &paymentmanager.CheckReturnDataCommand{
		ID:                    q.Id,
		Code:                  q.Code,
		PaymentStatus:         q.PaymentStatus,
		Amount:                q.Amount,
		ExternalTransactionID: q.ExternalTransactionId,
		Provider:              payment_provider.PaymentProvider(q.PaymentProvider),
	}
	if err := s.PaymentManager.Dispatch(ctx, args); err != nil {
		return nil, err
	}
	result := &pbcm.MessageResponse{
		Code: "ok",
		Msg:  args.Result.Msg,
	}
	return result, nil
}

func (s *PaymentService) GetExternalPaymentUrl(ctx context.Context, req *api.GetExternalPaymenUrlRequest) (res *api.GetExternalPaymentUrlResponse, err error) {
	// TODO: Tuấn
	// Refactor & test this func, it's not working correctly. Move logic to provider
	var transactionAmount int
	switch req.Type {
	case subject_referral.Credit:
		getCreditQuery := &credit.GetCreditQuery{
			ID:     req.RefID,
			ShopID: s.SS.Shop().ID,
		}
		if err := s.CreditQuery.Dispatch(ctx, getCreditQuery); err != nil {
			return nil, err
		}
		transactionAmount = getCreditQuery.Result.Amount
	}

	creatPaymentCmd := &payment.CreatePaymentCommand{
		ShopID:          s.SS.Shop().ID,
		Amount:          transactionAmount,
		Status:          status4.Z,
		State:           payment_state.Created,
		PaymentProvider: req.PaymentProvider,
	}
	if err := s.PaymentAggr.Dispatch(ctx, creatPaymentCmd); err != nil {
		return nil, err
	}
	_payment := creatPaymentCmd.Result

	argGenCode := &paymentmanager.GenerateCodeCommand{
		ReferralType: req.Type,
		ID:           req.RefID.String(),
	}
	if err := s.PaymentManager.Dispatch(ctx, argGenCode); err != nil {
		return nil, err
	}

	buildUrlConnectPaymentCmd := &paymentmanager.BuildUrlConnectPaymentGatewayCommand{
		OrderID:           argGenCode.Result,
		Desc:              fmt.Sprintf("thanh toán %s %s", req.Type.String(), req.RefID.String()),
		ReturnURL:         fmt.Sprintf("%s?paymentID=%s", req.ReturnURL, _payment.ID.String()),
		CancelURL:         fmt.Sprintf("%s?paymentID=%s", req.CancelURL, _payment.ID.String()),
		TransactionAmount: transactionAmount,
		Provider:          req.PaymentProvider,
	}
	if err := s.PaymentManager.Dispatch(ctx, buildUrlConnectPaymentCmd); err != nil {
		return nil, err
	}

	// format of transaction URL: "https://sbx-ctt.payme.vn/payment/5855423270"
	externalTransUrl := buildUrlConnectPaymentCmd.Result

	// get transactionID from transactionURL
	externalTransID := externalTransUrl[strings.LastIndex(externalTransUrl, "/")+1:]
	updatePaymentInfoCmd := &payment.UpdateExternalPaymentInfoCommand{
		ID:              _payment.ID,
		ExternalTransID: externalTransID,
	}
	if err := s.PaymentAggr.Dispatch(ctx, updatePaymentInfoCmd); err != nil {
		return nil, err
	}

	return &api.GetExternalPaymentUrlResponse{
		PaymentUrl: buildUrlConnectPaymentCmd.Result,
	}, err
}

func (s *PaymentService) UpdatePaymentStatus(
	ctx context.Context, req *api.UpdatePaymentStatusRequest,
) (*pbcm.UpdatedResponse, error) {
	updatePaymentInfoCmd := &payment.UpdateExternalPaymentInfoCommand{
		ID:     req.ID,
		State:  req.State,
		Status: req.Status,
	}
	if err := s.PaymentAggr.Dispatch(ctx, updatePaymentInfoCmd); err != nil {
		return nil, err
	}
	return &pbcm.UpdatedResponse{
		Updated: 1,
	}, nil
}
