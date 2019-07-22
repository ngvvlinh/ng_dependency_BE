package admin

import (
	"context"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/etop/api"
	"etop.vn/backend/pkg/etop/authorize/login"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/backend/pkg/etop/sqlstore"
	"etop.vn/backend/pkg/services/moneytx/modelx"
	shippingmodelx "etop.vn/backend/pkg/services/shipping/modelx"
	"etop.vn/common/bus"
	"etop.vn/common/l"

	pbcm "etop.vn/backend/pb/common"
	pbetop "etop.vn/backend/pb/etop"
	pbadmin "etop.vn/backend/pb/etop/admin"
	pborder "etop.vn/backend/pb/etop/order"
	notimodel "etop.vn/backend/pkg/notifier/model"
	wrapadmin "etop.vn/backend/wrapper/etop/admin"
)

var ll = l.New()

func init() {
	bus.AddHandlers("api",
		VersionInfo,
		LoginAsAccount,
		GetMoneyTransaction,
		GetMoneyTransactions,
		ConfirmMoneyTransaction,
		UpdateMoneyTransaction,
		GetMoneyTransactionShippingExternal,
		GetMoneyTransactionShippingExternals,
		RemoveMoneyTransactionShippingExternalLines,
		DeleteMoneyTransactionShippingExternal,
		ConfirmMoneyTransactionShippingExternal,
		ConfirmMoneyTransactionShippingExternals,
		UpdateMoneyTransactionShippingExternal,
		GetShop,
		GetShops,
		CreateCredit,
		GetCredit,
		GetCredits,
		UpdateCredit,
		ConfirmCredit,
		DeleteCredit,
		CreatePartner,
		GenerateAPIKey,
		UpdateFulfillment,
		CreateMoneyTransactionShippingEtop,
		GetMoneyTransactionShippingEtop,
		GetMoneyTransactionShippingEtops,
		UpdateMoneyTransactionShippingEtop,
		ConfirmMoneyTransactionShippingEtop,
		DeleteMoneyTransactionShippingEtop,
		CreateNotifications,
	)
}

func VersionInfo(ctx context.Context, q *wrapadmin.VersionInfoEndpoint) error {
	q.Result = &pbcm.VersionInfoResponse{
		Service: "etop.Admin",
		Version: "0.1",
	}
	return nil
}

func LoginAsAccount(ctx context.Context, q *wrapadmin.AdminLoginAsAccountEndpoint) error {
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
	case model.TagShop, model.TagSupplier:
	default:
		return cm.Error(cm.InvalidArgument, "Must be shop or supplier account", nil)
	}

	resp, err := adminCreateLoginResponse(ctx, q.Context.UserID, q.UserId, q.AccountId)
	q.Result = resp
	return err
}

func adminCreateLoginResponse(ctx context.Context, adminID, userID, accountID int64) (*pbetop.LoginResponse, error) {
	if adminID == 0 {
		return nil, cm.Error(cm.InvalidArgument, "Missing AdminID", nil)
	}

	resp, err := api.CreateLoginResponse(ctx, nil, "", userID, nil, accountID, 0, false, adminID)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func GetMoneyTransaction(ctx context.Context, q *wrapadmin.GetMoneyTransactionEndpoint) error {
	query := &modelx.GetMoneyTransaction{
		ID: q.Id,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = pborder.PbMoneyTransactionExtended(query.Result)
	return nil
}

func GetMoneyTransactions(ctx context.Context, q *wrapadmin.GetMoneyTransactionsEndpoint) error {
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

func UpdateMoneyTransaction(ctx context.Context, q *wrapadmin.UpdateMoneyTransactionEndpoint) error {
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

func ConfirmMoneyTransaction(ctx context.Context, q *wrapadmin.ConfirmMoneyTransactionEndpoint) error {
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

func GetMoneyTransactionShippingExternal(ctx context.Context, q *wrapadmin.GetMoneyTransactionShippingExternalEndpoint) error {
	query := &modelx.GetMoneyTransactionShippingExternal{
		ID: q.Id,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = pborder.PbMoneyTransactionShippingExternalExtended(query.Result)
	return nil
}

func GetMoneyTransactionShippingExternals(ctx context.Context, q *wrapadmin.GetMoneyTransactionShippingExternalsEndpoint) error {
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

func RemoveMoneyTransactionShippingExternalLines(ctx context.Context, q *wrapadmin.RemoveMoneyTransactionShippingExternalLinesEndpoint) error {
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

func DeleteMoneyTransactionShippingExternal(ctx context.Context, q *wrapadmin.DeleteMoneyTransactionShippingExternalEndpoint) error {
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

func ConfirmMoneyTransactionShippingExternal(ctx context.Context, q *wrapadmin.ConfirmMoneyTransactionShippingExternalEndpoint) error {
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

func ConfirmMoneyTransactionShippingExternals(ctx context.Context, q *wrapadmin.ConfirmMoneyTransactionShippingExternalsEndpoint) error {
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

func UpdateMoneyTransactionShippingExternal(ctx context.Context, q *wrapadmin.UpdateMoneyTransactionShippingExternalEndpoint) error {
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

func GetShop(ctx context.Context, q *wrapadmin.GetShopEndpoint) error {
	query := &model.GetShopExtendedQuery{
		ShopID: q.Id,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = pbetop.PbShopExtended(query.Result)
	return nil
}

func GetShops(ctx context.Context, q *wrapadmin.GetShopsEndpoint) error {
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

func CreateCredit(ctx context.Context, q *wrapadmin.CreateCreditEndpoint) error {
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

func GetCredit(ctx context.Context, q *wrapadmin.GetCreditEndpoint) error {
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

func GetCredits(ctx context.Context, q *wrapadmin.GetCreditsEndpoint) error {
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

func UpdateCredit(ctx context.Context, q *wrapadmin.UpdateCreditEndpoint) error {
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

func ConfirmCredit(ctx context.Context, q *wrapadmin.ConfirmCreditEndpoint) error {
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

func DeleteCredit(ctx context.Context, q *wrapadmin.DeleteCreditEndpoint) error {
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

func CreatePartner(ctx context.Context, q *wrapadmin.CreatePartnerEndpoint) error {
	cmd := &model.CreatePartnerCommand{
		Partner: q.ToModel(),
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = pbetop.PbPartner(cmd.Result.Partner)
	return nil
}

func UpdateFulfillment(ctx context.Context, q *wrapadmin.UpdateFulfillmentEndpoint) error {
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

func GenerateAPIKey(ctx context.Context, q *wrapadmin.GenerateAPIKeyEndpoint) error {
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

func GetMoneyTransactionShippingEtop(ctx context.Context, q *wrapadmin.GetMoneyTransactionShippingEtopEndpoint) error {
	query := &modelx.GetMoneyTransactionShippingEtop{
		ID: q.Id,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = pborder.PbMoneyTransactionShippingEtopExtended(query.Result)
	return nil
}

func GetMoneyTransactionShippingEtops(ctx context.Context, q *wrapadmin.GetMoneyTransactionShippingEtopsEndpoint) error {
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

func CreateMoneyTransactionShippingEtop(ctx context.Context, q *wrapadmin.CreateMoneyTransactionShippingEtopEndpoint) error {
	cmd := &modelx.CreateMoneyTransactionShippingEtop{
		MoneyTransactionShippingIDs: q.Ids,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = pborder.PbMoneyTransactionShippingEtopExtended(cmd.Result)
	return nil
}

func UpdateMoneyTransactionShippingEtop(ctx context.Context, q *wrapadmin.UpdateMoneyTransactionShippingEtopEndpoint) error {
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

func DeleteMoneyTransactionShippingEtop(ctx context.Context, q *wrapadmin.DeleteMoneyTransactionShippingEtopEndpoint) error {
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

func ConfirmMoneyTransactionShippingEtop(ctx context.Context, q *wrapadmin.ConfirmMoneyTransactionShippingEtopEndpoint) error {
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

func CreateNotifications(ctx context.Context, q *wrapadmin.CreateNotificationsEndpoint) error {
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
