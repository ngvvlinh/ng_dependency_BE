package admin

import (
	"context"

	"etop.vn/api/main/connectioning"
	"etop.vn/api/main/identity"
	"etop.vn/api/main/location"
	"etop.vn/api/main/moneytx"
	"etop.vn/api/main/shipmentpricing/pricelist"
	"etop.vn/api/main/shipmentpricing/shipmentprice"
	"etop.vn/api/main/shipmentpricing/shipmentservice"
	"etop.vn/api/top/int/admin"
	"etop.vn/api/top/int/etop"
	"etop.vn/api/top/int/types"
	pbcm "etop.vn/api/top/types/common"
	notimodel "etop.vn/backend/com/handler/notifier/model"
	creditmodelx "etop.vn/backend/com/main/credit/modelx"
	identitymodel "etop.vn/backend/com/main/identity/model"
	identitymodelx "etop.vn/backend/com/main/identity/modelx"
	shippingcarrier "etop.vn/backend/com/main/shipping/carrier"
	shippingmodelx "etop.vn/backend/com/main/shipping/modelx"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/apifw/cmapi"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/etop/api"
	"etop.vn/backend/pkg/etop/api/convertpb"
	"etop.vn/backend/pkg/etop/authorize/login"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/backend/pkg/etop/sqlstore"
	"etop.vn/capi"
	"etop.vn/capi/dot"
	"etop.vn/common/l"
)

var ll = l.New()

var (
	eventBus               capi.EventBus
	moneyTxQuery           moneytx.QueryBus
	moneyTxAggr            moneytx.CommandBus
	connectionAggr         connectioning.CommandBus
	connectionQuery        connectioning.QueryBus
	identityQuery          identity.QueryBus
	shipmentPriceAggr      shipmentprice.CommandBus
	shipmentPriceQuery     shipmentprice.QueryBus
	shipmentServiceAggr    shipmentservice.CommandBus
	shipmentServiceQuery   shipmentservice.QueryBus
	shipmentPriceListAggr  pricelist.CommandBus
	shipmentPriceListQuery pricelist.QueryBus
	locationAggr           location.CommandBus
	locationQuery          location.QueryBus
	shipmentManager        *shippingcarrier.ShipmentManager
)

func Init(
	eventB capi.EventBus,
	moneyTxQ moneytx.QueryBus,
	moneyTxA moneytx.CommandBus,
	connectionA connectioning.CommandBus,
	connectionQ connectioning.QueryBus,
	identityQ identity.QueryBus,
	shipmentpriceA shipmentprice.CommandBus,
	shipmentpriceQ shipmentprice.QueryBus,
	shipmentServiceA shipmentservice.CommandBus,
	shipmentServiceQ shipmentservice.QueryBus,
	shipmentPriceListA pricelist.CommandBus,
	shipmentPriceListQ pricelist.QueryBus,
	locationA location.CommandBus,
	locationQ location.QueryBus,
	shipmentM *shippingcarrier.ShipmentManager,
) {
	eventBus = eventB
	moneyTxQuery = moneyTxQ
	moneyTxAggr = moneyTxA
	connectionAggr = connectionA
	connectionQuery = connectionQ
	identityQuery = identityQ
	shipmentPriceAggr = shipmentpriceA
	shipmentPriceQuery = shipmentpriceQ
	shipmentServiceAggr = shipmentServiceA
	shipmentServiceQuery = shipmentServiceQ
	shipmentPriceListAggr = shipmentPriceListA
	shipmentPriceListQuery = shipmentPriceListQ
	locationAggr = locationA
	locationQuery = locationQ
	shipmentManager = shipmentM
}

func init() {
	bus.AddHandlers("api",
		miscService.VersionInfo,
		miscService.AdminLoginAsAccount,
		moneyTransactionService.GetMoneyTransaction,
		moneyTransactionService.GetMoneyTransactions,
		moneyTransactionService.ConfirmMoneyTransaction,
		moneyTransactionService.UpdateMoneyTransaction,
		moneyTransactionService.GetMoneyTransactionShippingExternal,
		moneyTransactionService.GetMoneyTransactionShippingExternals,
		moneyTransactionService.RemoveMoneyTransactionShippingExternalLines,
		moneyTransactionService.DeleteMoneyTransactionShippingExternal,
		moneyTransactionService.ConfirmMoneyTransactionShippingExternals,
		moneyTransactionService.UpdateMoneyTransactionShippingExternal,
		shopService.GetShop,
		shopService.GetShops,
		creditService.CreateCredit,
		creditService.GetCredit,
		creditService.GetCredits,
		creditService.UpdateCredit,
		creditService.ConfirmCredit,
		creditService.DeleteCredit,
		accountService.CreatePartner,
		accountService.GenerateAPIKey,
		fulfillmentService.UpdateFulfillment,
		moneyTransactionService.CreateMoneyTransactionShippingEtop,
		moneyTransactionService.GetMoneyTransactionShippingEtop,
		moneyTransactionService.GetMoneyTransactionShippingEtops,
		moneyTransactionService.UpdateMoneyTransactionShippingEtop,
		moneyTransactionService.ConfirmMoneyTransactionShippingEtop,
		moneyTransactionService.DeleteMoneyTransactionShippingEtop,
		notificationService.CreateNotifications,
	)
}

type MiscService struct{}
type AccountService struct{}
type OrderService struct{}
type FulfillmentService struct{}
type MoneyTransactionService struct{}
type ShopService struct{}
type CreditService struct{}
type NotificationService struct{}
type ConnectionService struct{}
type ShipmentPriceService struct{}
type LocationService struct{}

var miscService = &MiscService{}
var accountService = &AccountService{}
var orderService = &OrderService{}
var fulfillmentService = &FulfillmentService{}
var moneyTransactionService = &MoneyTransactionService{}
var shopService = &ShopService{}
var creditService = &CreditService{}
var notificationService = &NotificationService{}
var connectionService = &ConnectionService{}
var shipmentPriceService = &ShipmentPriceService{}
var locationService = &LocationService{}

func (s *MiscService) VersionInfo(ctx context.Context, q *VersionInfoEndpoint) error {
	q.Result = &pbcm.VersionInfoResponse{
		Service: "etop.Admin",
		Version: "0.1",
	}
	return nil
}

func (s *MiscService) AdminLoginAsAccount(ctx context.Context, q *AdminLoginAsAccountEndpoint) error {
	loginQuery := &login.LoginUserQuery{
		UserID:   q.Context.UserID,
		Password: q.Password,
	}
	if err := bus.Dispatch(ctx, loginQuery); err != nil {
		return cm.MapError(err).
			Mapf(cm.Unauthenticated, cm.Unauthenticated, "Admin password: %v", err).
			DefaultInternal()
	}

	switch cm.GetTag(q.AccountId) {
	case model.TagShop:
	default:
		return cm.Error(cm.InvalidArgument, "Must be shop account", nil)
	}

	resp, err := s.adminCreateLoginResponse(ctx, q.Context.UserID, q.UserId, q.AccountId)
	q.Result = resp
	return err
}

func (s *MiscService) adminCreateLoginResponse(ctx context.Context, adminID, userID, accountID dot.ID) (*etop.LoginResponse, error) {
	if adminID == 0 {
		return nil, cm.Error(cm.InvalidArgument, "Missing AdminID", nil)
	}

	resp, err := api.CreateLoginResponse(ctx, nil, "", userID, nil, accountID, 0, false, adminID)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *MoneyTransactionService) GetMoneyTransaction(ctx context.Context, q *GetMoneyTransactionEndpoint) error {
	query := &moneytx.GetMoneyTxShippingByIDQuery{
		MoneyTxShippingID: q.Id,
	}
	if err := moneyTxQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = convertpb.PbMoneyTxShipping(query.Result)
	return nil
}

func (s *MoneyTransactionService) GetMoneyTransactions(ctx context.Context, q *GetMoneyTransactionsEndpoint) error {
	paging := cmapi.CMPaging(q.Paging)
	query := &moneytx.ListMoneyTxShippingsQuery{
		MoneyTxShippingIDs: q.Ids,
		ShopID:             q.ShopId,
		Paging:             *paging,
		Filters:            cmapi.ToFilters(q.Filters),
		Result:             nil,
	}
	if err := moneyTxQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &types.MoneyTransactionsResponse{
		Paging:            cmapi.PbMetaPageInfo(query.Result.Paging),
		MoneyTransactions: convertpb.PbMoneyTxShippings(query.Result.MoneyTxShippings),
	}
	return nil
}

func (s *MoneyTransactionService) UpdateMoneyTransaction(ctx context.Context, q *UpdateMoneyTransactionEndpoint) error {
	cmd := &moneytx.UpdateMoneyTxShippingInfoCommand{
		MoneyTxShippingID: q.Id,
		Note:              q.Note,
		InvoiceNumber:     q.InvoiceNumber,
		BankAccount:       convertpb.Convert_api_BankAccount_To_core_BankAccount(q.BankAccount),
	}
	if err := moneyTxAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = convertpb.PbMoneyTxShipping(cmd.Result)
	return nil
}

func (s *MoneyTransactionService) ConfirmMoneyTransaction(ctx context.Context, q *ConfirmMoneyTransactionEndpoint) error {
	cmd := &moneytx.ConfirmMoneyTxShippingCommand{
		MoneyTxShippingID: q.MoneyTransactionId,
		ShopID:            q.ShopId,
		TotalCOD:          q.TotalCod,
		TotalAmount:       q.TotalAmount,
		TotalOrders:       q.TotalOrders,
	}
	if err := moneyTxAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.UpdatedResponse{
		Updated: 1,
	}
	return nil
}

func (s *MoneyTransactionService) GetMoneyTransactionShippingExternal(ctx context.Context, q *GetMoneyTransactionShippingExternalEndpoint) error {
	query := &moneytx.GetMoneyTxShippingExternalQuery{
		ID: q.Id,
	}
	if err := moneyTxQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = convertpb.PbMoneyTxShippingExternalFtLine(query.Result)
	return nil
}

func (s *MoneyTransactionService) GetMoneyTransactionShippingExternals(ctx context.Context, q *GetMoneyTransactionShippingExternalsEndpoint) error {
	paging := cmapi.CMPaging(q.Paging)
	query := &moneytx.ListMoneyTxShippingExternalsQuery{
		MoneyTxShippingExternalIDs: q.Ids,
		Paging:                     *paging,
		Filters:                    cmapi.ToFilters(q.Filters),
	}
	if err := moneyTxQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &types.MoneyTransactionShippingExternalsResponse{
		Paging:            cmapi.PbMetaPageInfo(query.Result.Paging),
		MoneyTransactions: convertpb.PbMoneyTxShippingExternalsFtLine(query.Result.MoneyTxShippingExternals),
	}
	return nil
}

func (s *MoneyTransactionService) RemoveMoneyTransactionShippingExternalLines(ctx context.Context, q *RemoveMoneyTransactionShippingExternalLinesEndpoint) error {
	cmd := &moneytx.RemoveMoneyTxShippingExternalLinesCommand{
		MoneyTxShippingExternalID: q.MoneyTransactionShippingExternalId,
		LineIDs:                   q.LineIds,
	}
	if err := moneyTxAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = convertpb.PbMoneyTxShippingExternalFtLine(cmd.Result)
	return nil
}

func (s *MoneyTransactionService) DeleteMoneyTransactionShippingExternal(ctx context.Context, q *DeleteMoneyTransactionShippingExternalEndpoint) error {
	cmd := &moneytx.DeleteMoneyTxShippingExternalCommand{
		ID: q.Id,
	}
	if err := moneyTxAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.RemovedResponse{
		Removed: cmd.Result,
	}
	return nil
}

func (s *MoneyTransactionService) ConfirmMoneyTransactionShippingExternals(ctx context.Context, q *ConfirmMoneyTransactionShippingExternalsEndpoint) error {
	cmd := &moneytx.ConfirmMoneyTxShippingExternalsCommand{
		IDs: q.Ids,
	}
	if err := moneyTxAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}

	q.Result = &pbcm.UpdatedResponse{
		Updated: cmd.Result,
	}
	return nil
}

func (s *MoneyTransactionService) UpdateMoneyTransactionShippingExternal(ctx context.Context, q *UpdateMoneyTransactionShippingExternalEndpoint) error {
	cmd := &moneytx.UpdateMoneyTxShippingExternalInfoCommand{
		MoneyTxShippingExternalID: q.Id,
		BankAccount:               convertpb.BankAccountToCoreBankAccount(q.BankAccount),
		Note:                      q.Note,
		InvoiceNumber:             q.InvoiceNumber,
	}
	if err := moneyTxAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = convertpb.PbMoneyTxShippingExternalFtLine(cmd.Result)
	return nil
}

func (s *ShopService) GetShop(ctx context.Context, q *GetShopEndpoint) error {
	query := &identitymodelx.GetShopExtendedQuery{
		ShopID: q.Id,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = convertpb.PbShopExtended(query.Result)
	return nil
}

func (s *ShopService) GetShops(ctx context.Context, q *GetShopsEndpoint) error {
	paging := cmapi.CMPaging(q.Paging)
	query := &identity.ListShopExtendedsQuery{
		Paging:  *paging,
		Filters: cmapi.ToFilters(q.Filters),
	}
	if err := identityQuery.Dispatch(ctx, query); err != nil {
		return err
	}

	q.Result = &admin.GetShopsResponse{
		Paging: cmapi.PbPageInfo(paging),
		Shops:  convertpb.Convert_core_ShopExtendeds_To_api_ShopExtendeds(query.Result.Shops),
	}
	return nil
}

func (s *ShopService) GetShopsByIDs(ctx context.Context, q *GetShopsByIDsEndpoint) error {
	query := &identity.ListShopsByIDsQuery{
		IDs: q.Ids,
	}
	if err := identityQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &admin.GetShopsResponse{
		Shops: convertpb.Convert_core_Shops_To_api_Shops(query.Result),
	}
	return nil
}

func (s *CreditService) CreateCredit(ctx context.Context, q *CreateCreditEndpoint) error {
	cmd := &creditmodelx.CreateCreditCommand{
		Amount: q.Amount,
		ShopID: q.ShopId,
		PaidAt: cmapi.PbTimeToModel(q.PaidAt),
		Type:   q.Type,
	}

	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = convertpb.PbCreditExtended(cmd.Result)
	return nil
}

func (s *CreditService) GetCredit(ctx context.Context, q *GetCreditEndpoint) error {
	query := &creditmodelx.GetCreditQuery{
		ID:     q.Id,
		ShopID: q.ShopId,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = convertpb.PbCreditExtended(query.Result)
	return nil
}

func (s *CreditService) GetCredits(ctx context.Context, q *GetCreditsEndpoint) error {
	paging := cmapi.CMPaging(q.Paging)
	query := &creditmodelx.GetCreditsQuery{
		ShopID: q.ShopId,
		Paging: paging,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &etop.CreditsResponse{
		Credits: convertpb.PbCreditExtendeds(query.Result.Credits),
		Paging:  cmapi.PbPageInfo(paging),
	}
	return nil
}

func (s *CreditService) UpdateCredit(ctx context.Context, q *UpdateCreditEndpoint) error {
	cmd := &creditmodelx.UpdateCreditCommand{
		ID:     q.Id,
		ShopID: q.ShopId,
		Amount: q.Amount,
		PaidAt: cmapi.PbTimeToModel(q.PaidAt),
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = convertpb.PbCreditExtended(cmd.Result)
	return nil
}

func (s *CreditService) ConfirmCredit(ctx context.Context, q *ConfirmCreditEndpoint) error {
	cmd := &creditmodelx.ConfirmCreditCommand{
		ID:     q.Id,
		ShopID: q.ShopId,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.UpdatedResponse{
		Updated: cmd.Result.Updated,
	}
	return nil
}

func (s *CreditService) DeleteCredit(ctx context.Context, q *DeleteCreditEndpoint) error {
	cmd := &creditmodelx.DeleteCreditCommand{
		ID: q.Id,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.RemovedResponse{
		Removed: cmd.Result.Deleted,
	}
	return nil
}

func (s *AccountService) CreatePartner(ctx context.Context, q *CreatePartnerEndpoint) error {
	cmd := &identitymodelx.CreatePartnerCommand{
		Partner: convertpb.CreatePartnerRequestToModel(q.CreatePartnerRequest),
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = convertpb.PbPartner(cmd.Result.Partner)
	return nil
}

func (s *FulfillmentService) UpdateFulfillment(ctx context.Context, q *UpdateFulfillmentEndpoint) error {
	cmd := &shippingmodelx.AdminUpdateFulfillmentCommand{
		FulfillmentID:            q.Id,
		FullName:                 q.FullName,
		Phone:                    q.Phone,
		TotalCODAmount:           q.TotalCodAmount,
		IsPartialDelivery:        q.IsPartialDelivery,
		AdminNote:                q.AdminNote,
		ActualCompensationAmount: q.ActualCompensationAmount,
		ShippingState:            q.ShippingState,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.UpdatedResponse{
		Updated: cmd.Result.Updated,
	}
	return nil
}

func (s *AccountService) GenerateAPIKey(ctx context.Context, q *GenerateAPIKeyEndpoint) error {
	_, err := sqlstore.AccountAuth(ctx).AccountID(q.AccountId).Get()
	if cm.ErrorCode(err) != cm.NotFound {
		return cm.MapError(err).
			Map(cm.OK, cm.AlreadyExists, "account already has an api_key").
			Throw()
	}

	aa := &identitymodel.AccountAuth{
		AccountID:   q.AccountId,
		Status:      1,
		Roles:       nil,
		Permissions: nil,
	}
	err = sqlstore.AccountAuth(ctx).Create(aa)
	q.Result = &admin.GenerateAPIKeyResponse{
		AccountId: q.AccountId,
		ApiKey:    aa.AuthKey,
	}
	return err
}

func (s *MoneyTransactionService) GetMoneyTransactionShippingEtop(ctx context.Context, q *GetMoneyTransactionShippingEtopEndpoint) error {
	query := &moneytx.GetMoneyTxShippingEtopQuery{
		ID: q.Id,
	}
	if err := moneyTxQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = convertpb.PbMoneyTxShippingEtop(query.Result)
	return nil
}

func (s *MoneyTransactionService) GetMoneyTransactionShippingEtops(ctx context.Context, q *GetMoneyTransactionShippingEtopsEndpoint) error {
	paging := cmapi.CMPaging(q.Paging)
	query := &moneytx.ListMoneyTxShippingEtopsQuery{
		MoneyTxShippingEtopIDs: q.Ids,
		Status:                 q.Status,
		Paging:                 *paging,
		Filter:                 cmapi.ToFilters(q.Filters),
	}
	if err := moneyTxQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &types.MoneyTransactionShippingEtopsResponse{
		Paging:                        cmapi.PbMetaPageInfo(query.Result.Paging),
		MoneyTransactionShippingEtops: convertpb.PbMoneyTxShippingEtops(query.Result.MoneyTxShippingEtops),
	}
	return nil
}

func (s *MoneyTransactionService) CreateMoneyTransactionShippingEtop(ctx context.Context, q *CreateMoneyTransactionShippingEtopEndpoint) error {
	cmd := &moneytx.CreateMoneyTxShippingEtopCommand{
		MoneyTxShippingIDs: q.Ids,
	}
	if err := moneyTxAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = convertpb.PbMoneyTxShippingEtop(cmd.Result)
	return nil
}

func (s *MoneyTransactionService) UpdateMoneyTransactionShippingEtop(ctx context.Context, q *UpdateMoneyTransactionShippingEtopEndpoint) error {
	cmd := &moneytx.UpdateMoneyTxShippingEtopCommand{
		MoneyTxShippingEtopID: q.Id,
		BankAccount:           convertpb.Convert_api_BankAccount_To_core_BankAccount(q.BankAccount),
		Note:                  q.Note,
		InvoiceNumber:         q.InvoiceNumber,
		Adds:                  q.Adds,
		Deletes:               q.Deletes,
		ReplaceAll:            q.ReplaceAll,
	}
	if err := moneyTxAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = convertpb.PbMoneyTxShippingEtop(cmd.Result)
	return nil
}

func (s *MoneyTransactionService) DeleteMoneyTransactionShippingEtop(ctx context.Context, q *DeleteMoneyTransactionShippingEtopEndpoint) error {
	cmd := &moneytx.DeleteMoneyTxShippingEtopCommand{
		ID: q.Id,
	}
	if err := moneyTxAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.DeletedResponse{Deleted: 1}
	return nil
}

func (s *MoneyTransactionService) ConfirmMoneyTransactionShippingEtop(ctx context.Context, q *ConfirmMoneyTransactionShippingEtopEndpoint) error {
	cmd := &moneytx.ConfirmMoneyTxShippingEtopCommand{
		MoneyTxShippingEtopID: q.Id,
		TotalCOD:              q.TotalCod,
		TotalAmount:           q.TotalAmount,
		TotalOrders:           q.TotalOrders,
	}
	if err := moneyTxAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.UpdatedResponse{Updated: 1}
	return nil
}

func (s *NotificationService) CreateNotifications(ctx context.Context, q *CreateNotificationsEndpoint) error {
	cmd := &notimodel.CreateNotificationsArgs{
		AccountIDs:       q.AccountIds,
		Title:            q.Title,
		Message:          q.Message,
		EntityID:         q.EntityId,
		Entity:           q.Entity,
		SendAll:          q.SendAll,
		SendNotification: true,
	}
	created, errored, err := sqlstore.CreateNotifications(ctx, cmd)
	if err != nil {
		return err
	}
	q.Result = &admin.CreateNotificationsResponse{
		Created: created,
		Errored: errored,
	}
	return nil
}
