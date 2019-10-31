package admin

import (
	"context"

	"etop.vn/capi"

	notimodel "etop.vn/backend/com/handler/notifier/model"
	"etop.vn/backend/com/main/moneytx/modelx"
	shippingmodelx "etop.vn/backend/com/main/shipping/modelx"
	pbcm "etop.vn/backend/pb/common"
	pbetop "etop.vn/backend/pb/etop"
	pbadmin "etop.vn/backend/pb/etop/admin"
	pborder "etop.vn/backend/pb/etop/order"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/etop/api"
	"etop.vn/backend/pkg/etop/authorize/login"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/backend/pkg/etop/sqlstore"
	wrapadmin "etop.vn/backend/wrapper/etop/admin"
	"etop.vn/common/l"
)

var ll = l.New()
var s = &Service{}

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
		s.VersionInfo,
		s.LoginAsAccount,
		s.GetMoneyTransaction,
		s.GetMoneyTransactions,
		s.ConfirmMoneyTransaction,
		s.UpdateMoneyTransaction,
		s.GetMoneyTransactionShippingExternal,
		s.GetMoneyTransactionShippingExternals,
		s.RemoveMoneyTransactionShippingExternalLines,
		s.DeleteMoneyTransactionShippingExternal,
		s.ConfirmMoneyTransactionShippingExternal,
		s.ConfirmMoneyTransactionShippingExternals,
		s.UpdateMoneyTransactionShippingExternal,
		s.GetShop,
		s.GetShops,
		s.CreateCredit,
		s.GetCredit,
		s.GetCredits,
		s.UpdateCredit,
		s.ConfirmCredit,
		s.DeleteCredit,
		s.CreatePartner,
		s.GenerateAPIKey,
		s.UpdateFulfillment,
		s.CreateMoneyTransactionShippingEtop,
		s.GetMoneyTransactionShippingEtop,
		s.GetMoneyTransactionShippingEtops,
		s.UpdateMoneyTransactionShippingEtop,
		s.ConfirmMoneyTransactionShippingEtop,
		s.DeleteMoneyTransactionShippingEtop,
		s.CreateNotifications,
	)
}

type Service struct{}

func (s *Service) VersionInfo(ctx context.Context, q *wrapadmin.VersionInfoEndpoint) error {
	q.Result = &pbcm.VersionInfoResponse{
		Service: "etop.Admin",
		Version: "0.1",
	}
	return nil
}

func (s *Service) LoginAsAccount(ctx context.Context, q *wrapadmin.AdminLoginAsAccountEndpoint) error {
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

func (s *Service) adminCreateLoginResponse(ctx context.Context, adminID, userID, accountID int64) (*pbetop.LoginResponse, error) {
	if adminID == 0 {
		return nil, cm.Error(cm.InvalidArgument, "Missing AdminID", nil)
	}

	resp, err := api.CreateLoginResponse(ctx, nil, "", userID, nil, accountID, 0, false, adminID)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *Service) GetMoneyTransaction(ctx context.Context, q *wrapadmin.GetMoneyTransactionEndpoint) error {
	query := &modelx.GetMoneyTransaction{
		ID: q.Id,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = pborder.PbMoneyTransactionExtended(query.Result)
	return nil
}

func (s *Service) GetMoneyTransactions(ctx context.Context, q *wrapadmin.GetMoneyTransactionsEndpoint) error {
	paging := q.Paging.CMPaging()
	query := &modelx.GetMoneyTransactions{
		IDs:                                q.Ids,
		ShopID:                             q.ShopId,
		Paging:                             paging,
		MoneyTransactionShippingExternalID: q.MoneyTransactionShippingExternalId,
		Filters:                            pbcm.ToFilters(q.Filters),
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &pborder.MoneyTransactionsResponse{
		MoneyTransactions: pborder.PbMoneyTransactionExtendeds(query.Result.MoneyTransactions),
		Paging:            pbcm.PbPageInfo(paging, int32(query.Result.Total)),
	}
	return nil
}

func (s *Service) UpdateMoneyTransaction(ctx context.Context, q *wrapadmin.UpdateMoneyTransactionEndpoint) error {
	cmd := &modelx.UpdateMoneyTransaction{
		ID:            q.Id,
		Note:          q.Note,
		InvoiceNumber: q.InvoiceNumber,
		BankAccount:   q.BankAccount.ToModel(),
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = pborder.PbMoneyTransactionExtended(cmd.Result)
	return nil
}

func (s *Service) ConfirmMoneyTransaction(ctx context.Context, q *wrapadmin.ConfirmMoneyTransactionEndpoint) error {
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

func (s *Service) GetMoneyTransactionShippingExternal(ctx context.Context, q *wrapadmin.GetMoneyTransactionShippingExternalEndpoint) error {
	query := &modelx.GetMoneyTransactionShippingExternal{
		ID: q.Id,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = pborder.PbMoneyTransactionShippingExternalExtended(query.Result)
	return nil
}

func (s *Service) GetMoneyTransactionShippingExternals(ctx context.Context, q *wrapadmin.GetMoneyTransactionShippingExternalsEndpoint) error {
	paging := q.Paging.CMPaging()
	query := &modelx.GetMoneyTransactionShippingExternals{
		IDs:     q.Ids,
		Paging:  paging,
		Filters: pbcm.ToFilters(q.Filters),
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &pborder.MoneyTransactionShippingExternalsResponse{
		MoneyTransactions: pborder.PbMoneyTransactionShippingExternalExtendeds(query.Result.MoneyTransactionShippingExternals),
		Paging:            pbcm.PbPageInfo(paging, int32(query.Result.Total)),
	}
	return nil
}

func (s *Service) RemoveMoneyTransactionShippingExternalLines(ctx context.Context, q *wrapadmin.RemoveMoneyTransactionShippingExternalLinesEndpoint) error {
	cmd := &modelx.RemoveMoneyTransactionShippingExternalLines{
		MoneyTransactionShippingExternalID: q.MoneyTransactionShippingExternalId,
		LineIDs:                            q.LineIds,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = pborder.PbMoneyTransactionShippingExternalExtended(cmd.Result)
	return nil
}

func (s *Service) DeleteMoneyTransactionShippingExternal(ctx context.Context, q *wrapadmin.DeleteMoneyTransactionShippingExternalEndpoint) error {
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

func (s *Service) ConfirmMoneyTransactionShippingExternal(ctx context.Context, q *wrapadmin.ConfirmMoneyTransactionShippingExternalEndpoint) error {
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

func (s *Service) ConfirmMoneyTransactionShippingExternals(ctx context.Context, q *wrapadmin.ConfirmMoneyTransactionShippingExternalsEndpoint) error {
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

func (s *Service) UpdateMoneyTransactionShippingExternal(ctx context.Context, q *wrapadmin.UpdateMoneyTransactionShippingExternalEndpoint) error {
	cmd := &modelx.UpdateMoneyTransactionShippingExternal{
		ID:            q.Id,
		Note:          q.Note,
		InvoiceNumber: q.InvoiceNumber,
		BankAccount:   q.BankAccount.ToModel(),
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = pborder.PbMoneyTransactionShippingExternalExtended(cmd.Result)
	return nil
}

func (s *Service) GetShop(ctx context.Context, q *wrapadmin.GetShopEndpoint) error {
	query := &model.GetShopExtendedQuery{
		ShopID: q.Id,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = pbetop.PbShopExtended(query.Result)
	return nil
}

func (s *Service) GetShops(ctx context.Context, q *wrapadmin.GetShopsEndpoint) error {
	paging := q.Paging.CMPaging()
	query := &model.GetAllShopExtendedsQuery{
		Paging: paging,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &pbadmin.GetShopsResponse{
		Paging: pbcm.PbPageInfo(paging, int32(query.Result.Total)),
		Shops:  pbetop.PbShopExtendeds(query.Result.Shops),
	}
	return nil
}

func (s *Service) CreateCredit(ctx context.Context, q *wrapadmin.CreateCreditEndpoint) error {
	cmd := &model.CreateCreditCommand{
		Amount: int(q.Amount),
		ShopID: q.ShopId,
		Type:   model.AccountType(q.Type.ToModel()),
		PaidAt: pbcm.PbTimeToModel(q.PaidAt),
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = pbetop.PbCreditExtended(cmd.Result)
	return nil
}

func (s *Service) GetCredit(ctx context.Context, q *wrapadmin.GetCreditEndpoint) error {
	query := &model.GetCreditQuery{
		ID:     q.Id,
		ShopID: q.ShopId,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = pbetop.PbCreditExtended(query.Result)
	return nil
}

func (s *Service) GetCredits(ctx context.Context, q *wrapadmin.GetCreditsEndpoint) error {
	paging := q.Paging.CMPaging()
	query := &model.GetCreditsQuery{
		ShopID: q.ShopId,
		Paging: paging,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &pbetop.CreditsResponse{
		Credits: pbetop.PbCreditExtendeds(query.Result.Credits),
		Paging:  pbcm.PbPageInfo(paging, int32(query.Result.Total)),
	}
	return nil
}

func (s *Service) UpdateCredit(ctx context.Context, q *wrapadmin.UpdateCreditEndpoint) error {
	cmd := &model.UpdateCreditCommand{
		ID:     q.Id,
		ShopID: q.ShopId,
		Amount: int(q.Amount),
		PaidAt: pbcm.PbTimeToModel(q.PaidAt),
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = pbetop.PbCreditExtended(cmd.Result)
	return nil
}

func (s *Service) ConfirmCredit(ctx context.Context, q *wrapadmin.ConfirmCreditEndpoint) error {
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

func (s *Service) DeleteCredit(ctx context.Context, q *wrapadmin.DeleteCreditEndpoint) error {
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

func (s *Service) CreatePartner(ctx context.Context, q *wrapadmin.CreatePartnerEndpoint) error {
	cmd := &model.CreatePartnerCommand{
		Partner: q.ToModel(),
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = pbetop.PbPartner(cmd.Result.Partner)
	return nil
}

func (s *Service) UpdateFulfillment(ctx context.Context, q *wrapadmin.UpdateFulfillmentEndpoint) error {
	cmd := &shippingmodelx.AdminUpdateFulfillmentCommand{
		FulfillmentID:            q.Id,
		FullName:                 q.FullName,
		Phone:                    q.Phone,
		TotalCODAmount:           cm.PInt32(q.TotalCodAmount),
		IsPartialDelivery:        q.IsPartialDelivery,
		AdminNote:                q.AdminNote,
		ActualCompensationAmount: int(q.ActualCompensationAmount),
		ShippingState:            q.ShippingState.ToModel(),
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.UpdatedResponse{
		Updated: int32(cmd.Result.Updated),
	}
	return nil
}

func (s *Service) GenerateAPIKey(ctx context.Context, q *wrapadmin.GenerateAPIKeyEndpoint) error {
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

func (s *Service) GetMoneyTransactionShippingEtop(ctx context.Context, q *wrapadmin.GetMoneyTransactionShippingEtopEndpoint) error {
	query := &modelx.GetMoneyTransactionShippingEtop{
		ID: q.Id,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = pborder.PbMoneyTransactionShippingEtopExtended(query.Result)
	return nil
}

func (s *Service) GetMoneyTransactionShippingEtops(ctx context.Context, q *wrapadmin.GetMoneyTransactionShippingEtopsEndpoint) error {
	paging := q.Paging.CMPaging()
	query := &modelx.GetMoneyTransactionShippingEtops{
		IDs:     q.Ids,
		Status:  q.Status.ToModel(),
		Paging:  paging,
		Filters: pbcm.ToFilters(q.Filters),
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &pborder.MoneyTransactionShippingEtopsResponse{
		Paging:                        pbcm.PbPageInfo(paging, int32(query.Result.Total)),
		MoneyTransactionShippingEtops: pborder.PbMoneyTransactionShippingEtopExtendeds(query.Result.MoneyTransactionShippingEtops),
	}
	return nil
}

func (s *Service) CreateMoneyTransactionShippingEtop(ctx context.Context, q *wrapadmin.CreateMoneyTransactionShippingEtopEndpoint) error {
	cmd := &modelx.CreateMoneyTransactionShippingEtop{
		MoneyTransactionShippingIDs: q.Ids,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = pborder.PbMoneyTransactionShippingEtopExtended(cmd.Result)
	return nil
}

func (s *Service) UpdateMoneyTransactionShippingEtop(ctx context.Context, q *wrapadmin.UpdateMoneyTransactionShippingEtopEndpoint) error {
	cmd := &modelx.UpdateMoneyTransactionShippingEtop{
		ID:            q.Id,
		Adds:          q.Adds,
		Deletes:       q.Deletes,
		ReplaceAll:    q.ReplaceAll,
		Note:          q.Note,
		InvoiceNumber: q.InvoiceNumber,
		BankAccount:   q.BankAccount.ToModel(),
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = pborder.PbMoneyTransactionShippingEtopExtended(cmd.Result)
	return nil
}

func (s *Service) DeleteMoneyTransactionShippingEtop(ctx context.Context, q *wrapadmin.DeleteMoneyTransactionShippingEtopEndpoint) error {
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

func (s *Service) ConfirmMoneyTransactionShippingEtop(ctx context.Context, q *wrapadmin.ConfirmMoneyTransactionShippingEtopEndpoint) error {
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

func (s *Service) CreateNotifications(ctx context.Context, q *wrapadmin.CreateNotificationsEndpoint) error {
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
