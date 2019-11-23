package admin

import (
	"context"

	pbcm "etop.vn/api/pb/common"
	pbetop "etop.vn/api/pb/etop"
	pbadmin "etop.vn/api/pb/etop/admin"
	pborder "etop.vn/api/pb/etop/order"
	notimodel "etop.vn/backend/com/handler/notifier/model"
	"etop.vn/backend/com/main/moneytx/modelx"
	shippingmodelx "etop.vn/backend/com/main/shipping/modelx"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmapi"
	"etop.vn/backend/pkg/etop/api"
	"etop.vn/backend/pkg/etop/api/convertpb"
	"etop.vn/backend/pkg/etop/authorize/login"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/backend/pkg/etop/sqlstore"
	"etop.vn/capi"
	"etop.vn/common/l"
)

var ll = l.New()

var (
	eventBus capi.EventBus
)

func Init(
	eventB capi.EventBus,
) {
	eventBus = eventB
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
		moneyTransactionService.ConfirmMoneyTransactionShippingExternal,
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

var miscService = &MiscService{}
var accountService = &AccountService{}
var orderService = &OrderService{}
var fulfillmentService = &FulfillmentService{}
var moneyTransactionService = &MoneyTransactionService{}
var shopService = &ShopService{}
var creditService = &CreditService{}
var notificationService = &NotificationService{}

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

func (s *MiscService) adminCreateLoginResponse(ctx context.Context, adminID, userID, accountID int64) (*pbetop.LoginResponse, error) {
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
	query := &modelx.GetMoneyTransaction{
		ID: q.Id,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = convertpb.PbMoneyTransactionExtended(query.Result)
	return nil
}

func (s *MoneyTransactionService) GetMoneyTransactions(ctx context.Context, q *GetMoneyTransactionsEndpoint) error {
	paging := cmapi.CMPaging(q.Paging)
	query := &modelx.GetMoneyTransactions{
		IDs:                                q.Ids,
		ShopID:                             q.ShopId,
		Paging:                             paging,
		MoneyTransactionShippingExternalID: q.MoneyTransactionShippingExternalId,
		Filters:                            cmapi.ToFilters(q.Filters),
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &pborder.MoneyTransactionsResponse{
		MoneyTransactions: convertpb.PbMoneyTransactionExtendeds(query.Result.MoneyTransactions),
		Paging:            cmapi.PbPageInfo(paging, int32(query.Result.Total)),
	}
	return nil
}

func (s *MoneyTransactionService) UpdateMoneyTransaction(ctx context.Context, q *UpdateMoneyTransactionEndpoint) error {
	cmd := &modelx.UpdateMoneyTransaction{
		ID:            q.Id,
		Note:          q.Note,
		InvoiceNumber: q.InvoiceNumber,
		BankAccount:   convertpb.BankAccountToModel(q.BankAccount),
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = convertpb.PbMoneyTransactionExtended(cmd.Result)
	return nil
}

func (s *MoneyTransactionService) ConfirmMoneyTransaction(ctx context.Context, q *ConfirmMoneyTransactionEndpoint) error {
	cmd := &modelx.ConfirmMoneyTransaction{
		MoneyTransactionID: q.MoneyTransactionId,
		ShopID:             q.ShopId,
		TotalCOD:           int(q.TotalCod),
		TotalAmount:        int(q.TotalAmount),
		TotalOrders:        int(q.TotalOrders),
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.UpdatedResponse{
		Updated: int32(cmd.Result.Updated),
	}
	return nil
}

func (s *MoneyTransactionService) GetMoneyTransactionShippingExternal(ctx context.Context, q *GetMoneyTransactionShippingExternalEndpoint) error {
	query := &modelx.GetMoneyTransactionShippingExternal{
		ID: q.Id,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = convertpb.PbMoneyTransactionShippingExternalExtended(query.Result)
	return nil
}

func (s *MoneyTransactionService) GetMoneyTransactionShippingExternals(ctx context.Context, q *GetMoneyTransactionShippingExternalsEndpoint) error {
	paging := cmapi.CMPaging(q.Paging)
	query := &modelx.GetMoneyTransactionShippingExternals{
		IDs:     q.Ids,
		Paging:  paging,
		Filters: cmapi.ToFilters(q.Filters),
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &pborder.MoneyTransactionShippingExternalsResponse{
		MoneyTransactions: convertpb.PbMoneyTransactionShippingExternalExtendeds(query.Result.MoneyTransactionShippingExternals),
		Paging:            cmapi.PbPageInfo(paging, int32(query.Result.Total)),
	}
	return nil
}

func (s *MoneyTransactionService) RemoveMoneyTransactionShippingExternalLines(ctx context.Context, q *RemoveMoneyTransactionShippingExternalLinesEndpoint) error {
	cmd := &modelx.RemoveMoneyTransactionShippingExternalLines{
		MoneyTransactionShippingExternalID: q.MoneyTransactionShippingExternalId,
		LineIDs:                            q.LineIds,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = convertpb.PbMoneyTransactionShippingExternalExtended(cmd.Result)
	return nil
}

func (s *MoneyTransactionService) DeleteMoneyTransactionShippingExternal(ctx context.Context, q *DeleteMoneyTransactionShippingExternalEndpoint) error {
	cmd := &modelx.DeleteMoneyTransactionShippingExternal{
		ID: q.Id,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.RemovedResponse{
		Removed: int32(cmd.Result.Deleted),
	}
	return nil
}

func (s *MoneyTransactionService) ConfirmMoneyTransactionShippingExternal(ctx context.Context, q *ConfirmMoneyTransactionShippingExternalEndpoint) error {
	cmd := &modelx.ConfirmMoneyTransactionShippingExternal{
		ID: q.Id,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}

	q.Result = &pbcm.UpdatedResponse{
		Updated: int32(cmd.Result.Updated),
	}
	return nil
}

func (s *MoneyTransactionService) ConfirmMoneyTransactionShippingExternals(ctx context.Context, q *ConfirmMoneyTransactionShippingExternalsEndpoint) error {
	cmd := &modelx.ConfirmMoneyTransactionShippingExternals{
		IDs: q.Ids,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}

	q.Result = &pbcm.UpdatedResponse{
		Updated: int32(cmd.Result.Updated),
	}
	return nil
}

func (s *MoneyTransactionService) UpdateMoneyTransactionShippingExternal(ctx context.Context, q *UpdateMoneyTransactionShippingExternalEndpoint) error {
	cmd := &modelx.UpdateMoneyTransactionShippingExternal{
		ID:            q.Id,
		Note:          q.Note,
		InvoiceNumber: q.InvoiceNumber,
		BankAccount:   convertpb.BankAccountToModel(q.BankAccount),
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = convertpb.PbMoneyTransactionShippingExternalExtended(cmd.Result)
	return nil
}

func (s *ShopService) GetShop(ctx context.Context, q *GetShopEndpoint) error {
	query := &model.GetShopExtendedQuery{
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
	query := &model.GetAllShopExtendedsQuery{
		Paging: paging,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &pbadmin.GetShopsResponse{
		Paging: cmapi.PbPageInfo(paging, int32(query.Result.Total)),
		Shops:  convertpb.PbShopExtendeds(query.Result.Shops),
	}
	return nil
}

func (s *CreditService) CreateCredit(ctx context.Context, q *CreateCreditEndpoint) error {
	cmd := &model.CreateCreditCommand{
		Amount: int(q.Amount),
		ShopID: q.ShopId,
		Type:   model.AccountType(q.Type.ToModel()),
		PaidAt: cmapi.PbTimeToModel(q.PaidAt),
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = convertpb.PbCreditExtended(cmd.Result)
	return nil
}

func (s *CreditService) GetCredit(ctx context.Context, q *GetCreditEndpoint) error {
	query := &model.GetCreditQuery{
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
	query := &model.GetCreditsQuery{
		ShopID: q.ShopId,
		Paging: paging,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &pbetop.CreditsResponse{
		Credits: convertpb.PbCreditExtendeds(query.Result.Credits),
		Paging:  cmapi.PbPageInfo(paging, int32(query.Result.Total)),
	}
	return nil
}

func (s *CreditService) UpdateCredit(ctx context.Context, q *UpdateCreditEndpoint) error {
	cmd := &model.UpdateCreditCommand{
		ID:     q.Id,
		ShopID: q.ShopId,
		Amount: int(q.Amount),
		PaidAt: cmapi.PbTimeToModel(q.PaidAt),
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = convertpb.PbCreditExtended(cmd.Result)
	return nil
}

func (s *CreditService) ConfirmCredit(ctx context.Context, q *ConfirmCreditEndpoint) error {
	cmd := &model.ConfirmCreditCommand{
		ID:     q.Id,
		ShopID: q.ShopId,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.UpdatedResponse{
		Updated: int32(cmd.Result.Updated),
	}
	return nil
}

func (s *CreditService) DeleteCredit(ctx context.Context, q *DeleteCreditEndpoint) error {
	cmd := &model.DeleteCreditCommand{
		ID: q.Id,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.RemovedResponse{
		Removed: int32(cmd.Result.Deleted),
	}
	return nil
}

func (s *AccountService) CreatePartner(ctx context.Context, q *CreatePartnerEndpoint) error {
	cmd := &model.CreatePartnerCommand{
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
		TotalCODAmount:           cm.PInt32(q.TotalCodAmount),
		IsPartialDelivery:        q.IsPartialDelivery,
		AdminNote:                q.AdminNote,
		ActualCompensationAmount: int(q.ActualCompensationAmount),
		ShippingState:            convertpb.ShippingStateToModel(q.ShippingState),
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.UpdatedResponse{
		Updated: int32(cmd.Result.Updated),
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

	aa := &model.AccountAuth{
		AccountID:   q.AccountId,
		Status:      1,
		Roles:       nil,
		Permissions: nil,
	}
	err = sqlstore.AccountAuth(ctx).Create(aa)
	q.Result = &pbadmin.GenerateAPIKeyResponse{
		AccountId: q.AccountId,
		ApiKey:    aa.AuthKey,
	}
	return err
}

func (s *MoneyTransactionService) GetMoneyTransactionShippingEtop(ctx context.Context, q *GetMoneyTransactionShippingEtopEndpoint) error {
	query := &modelx.GetMoneyTransactionShippingEtop{
		ID: q.Id,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = convertpb.PbMoneyTransactionShippingEtopExtended(query.Result)
	return nil
}

func (s *MoneyTransactionService) GetMoneyTransactionShippingEtops(ctx context.Context, q *GetMoneyTransactionShippingEtopsEndpoint) error {
	paging := cmapi.CMPaging(q.Paging)
	query := &modelx.GetMoneyTransactionShippingEtops{
		IDs:     q.Ids,
		Status:  convertpb.Status3ToModel(q.Status),
		Paging:  paging,
		Filters: cmapi.ToFilters(q.Filters),
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &pborder.MoneyTransactionShippingEtopsResponse{
		Paging:                        cmapi.PbPageInfo(paging, int32(query.Result.Total)),
		MoneyTransactionShippingEtops: convertpb.PbMoneyTransactionShippingEtopExtendeds(query.Result.MoneyTransactionShippingEtops),
	}
	return nil
}

func (s *MoneyTransactionService) CreateMoneyTransactionShippingEtop(ctx context.Context, q *CreateMoneyTransactionShippingEtopEndpoint) error {
	cmd := &modelx.CreateMoneyTransactionShippingEtop{
		MoneyTransactionShippingIDs: q.Ids,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = convertpb.PbMoneyTransactionShippingEtopExtended(cmd.Result)
	return nil
}

func (s *MoneyTransactionService) UpdateMoneyTransactionShippingEtop(ctx context.Context, q *UpdateMoneyTransactionShippingEtopEndpoint) error {
	cmd := &modelx.UpdateMoneyTransactionShippingEtop{
		ID:            q.Id,
		Adds:          q.Adds,
		Deletes:       q.Deletes,
		ReplaceAll:    q.ReplaceAll,
		Note:          q.Note,
		InvoiceNumber: q.InvoiceNumber,
		BankAccount:   convertpb.BankAccountToModel(q.BankAccount),
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = convertpb.PbMoneyTransactionShippingEtopExtended(cmd.Result)
	return nil
}

func (s *MoneyTransactionService) DeleteMoneyTransactionShippingEtop(ctx context.Context, q *DeleteMoneyTransactionShippingEtopEndpoint) error {
	cmd := &modelx.DeleteMoneyTransactionShippingEtop{
		ID: q.Id,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.DeletedResponse{
		Deleted: int32(cmd.Result.Deleted),
	}
	return nil
}

func (s *MoneyTransactionService) ConfirmMoneyTransactionShippingEtop(ctx context.Context, q *ConfirmMoneyTransactionShippingEtopEndpoint) error {
	cmd := &modelx.ConfirmMoneyTransactionShippingEtop{
		ID:          q.Id,
		TotalCOD:    int(q.TotalCod),
		TotalAmount: int(q.TotalAmount),
		TotalOrders: int(q.TotalOrders),
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.UpdatedResponse{
		Updated: int32(cmd.Result.Updated),
	}
	return nil
}

func (s *NotificationService) CreateNotifications(ctx context.Context, q *CreateNotificationsEndpoint) error {
	cmd := &notimodel.CreateNotificationsArgs{
		AccountIDs:       q.AccountIds,
		Title:            q.Title,
		Message:          q.Message,
		EntityID:         q.EntityId,
		Entity:           notimodel.NotiEntity(q.Entity.ToModel()),
		SendAll:          q.SendAll,
		SendNotification: true,
	}
	created, errored, err := sqlstore.CreateNotifications(ctx, cmd)
	if err != nil {
		return err
	}
	q.Result = &pbadmin.CreateNotificationsResponse{
		Created: int32(created),
		Errored: int32(errored),
	}
	return nil
}
