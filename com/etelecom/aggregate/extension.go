package aggregate

import (
	"context"
	"time"

	"o.o/api/etelecom"
	"o.o/api/main/connectioning"
	"o.o/api/main/identity"
	"o.o/api/main/invoicing"
	"o.o/api/subscripting/subscription"
	subscriptingtypes "o.o/api/subscripting/types"
	"o.o/api/top/types/etc/connection_type"
	"o.o/api/top/types/etc/payment_method"
	"o.o/api/top/types/etc/service_classify"
	"o.o/api/top/types/etc/status3"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi/dot"
)

func (a *EtelecomAggregate) CreateExtension(ctx context.Context, args *etelecom.CreateExtensionArgs) (*etelecom.Extension, error) {
	event := &etelecom.ExtensionCreatingEvent{
		OwnerID:   args.OwnerID,
		AccountID: args.AccountID,
		UserID:    args.UserID,
	}
	if err := a.eventBus.Publish(ctx, event); err != nil {
		return nil, err
	}
	return a.createExtension(ctx, args)
}

func (a *EtelecomAggregate) createExtension(ctx context.Context, args *etelecom.CreateExtensionArgs) (_ *etelecom.Extension, _err error) {
	if err := args.Validate(); err != nil {
		return nil, err
	}
	if args.UserID != 0 {
		queryUser := &identity.GetAccountUserQuery{
			UserID:    args.UserID,
			AccountID: args.AccountID,
		}
		if err := a.identityQuery.Dispatch(ctx, queryUser); err != nil {
			return nil, cm.Errorf(cm.ErrorCode(err), err, "Không tìm thấy nhân viên")
		}
	}
	tenant, err := a.getTenant(ctx, args)
	if err != nil {
		return nil, err
	}

	var ext *etelecom.Extension
	if args.UserID != 0 {
		ext, err = a.extensionStore(ctx).UserID(args.UserID).AccountID(args.AccountID).GetExtension()
		switch cm.ErrorCode(err) {
		case cm.NoError:
			if ext.ExternalData != nil && ext.ExternalData.ID != "" {
				return nil, cm.Errorf(cm.FailedPrecondition, nil, "Extension đã được tạo cho người dùng này.")
			}
		case cm.NotFound:
		default:
			return nil, err
		}
	}

	defer func() {
		if _err != nil && ext != nil {
			// delete extension
			a.extensionStore(ctx).ID(ext.ID).SoftDelete()
		}
		return
	}()

	if err = a.checkDuplicateExtensionNumber(ctx, tenant.ID, args.ExtensionNumber); err != nil {
		return nil, err
	}
	if ext == nil {
		var extension etelecom.Extension
		if err = scheme.Convert(args, &extension); err != nil {
			return nil, err
		}
		extension.ID = cm.NewID()
		extension.TenantID = tenant.ID
		ext, err = a.extensionStore(ctx).CreateExtension(&extension)
		if err != nil {
			return nil, err
		}
	}

	externalExtensionResp, err := a.telecomManager.CreateExtension(ctx, ext)
	if err != nil {
		return nil, err
	}
	updateExt := &etelecom.UpdateExternalExtensionInfoArgs{
		ID:                externalExtensionResp.ExtensionID,
		HotlineID:         externalExtensionResp.HotlineID,
		ExternalID:        externalExtensionResp.ExternalID,
		ExtensionNumber:   externalExtensionResp.ExtensionNumber,
		ExtensionPassword: externalExtensionResp.ExtensionPassword,
		TenantDomain:      tenant.Domain,
	}
	if err = a.UpdateExternalExtensionInfo(ctx, updateExt); err != nil {
		return nil, err
	}

	if args.UserID != 0 {
		extensionAssignedEvent := &etelecom.AssignedExtensionEvent{
			AccountID:       args.AccountID,
			UserID:          args.UserID,
			ExtensionNumber: externalExtensionResp.ExtensionNumber,
		}
		if err := a.eventBus.Publish(ctx, extensionAssignedEvent); err != nil {
			return nil, err
		}
	}

	return a.extensionStore(ctx).ID(ext.ID).GetExtension()
}

func (a *EtelecomAggregate) checkDuplicateExtensionNumber(ctx context.Context, tenantID dot.ID, extNumber string) error {
	if extNumber == "" {
		return nil
	}
	// make sure extension number does not duplicate in the same tenant
	_, err := a.extensionStore(ctx).TenantID(tenantID).ExtensionNumber(extNumber).GetExtension()
	switch cm.ErrorCode(err) {
	case cm.NoError:
		return cm.Errorf(cm.FailedPrecondition, nil, "Số máy nhánh đã tồn tại")
	case cm.NotFound:
		return nil
	default:
		return err
	}

}

func (a *EtelecomAggregate) CreateExtensionBySubscription(ctx context.Context, args *etelecom.CreateExtenstionBySubscriptionArgs) (_ *etelecom.Extension, _err error) {
	if err := args.Validate(); err != nil {
		return nil, err
	}

	// create by subscription plan
	if args.SubscriptionID == 0 && args.SubscriptionPlanID != 0 {
		subrID := dot.ID(0)
		invoiceID := dot.ID(0)
		// create subscription

		var res *etelecom.Extension
		var err error
		err = a.txDBMain.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
			cmd := &subscription.CreateSubscriptionCommand{
				AccountID: args.AccountID,
				Lines: []*subscription.SubscriptionLine{
					{
						PlanID:   args.SubscriptionPlanID,
						Quantity: 1,
					},
				},
			}
			if err = a.subrAggr.Dispatch(ctx, cmd); err != nil {
				return err
			}
			subrID = cmd.Result.ID

			// payment subscription
			if args.PaymentMethod != payment_method.Balance {
				return cm.Errorf(cm.InvalidArgument, nil, "Phương thức thanh toán không hợp lệ.")
			}
			queryShop := &identity.GetShopByIDQuery{
				ID: args.AccountID,
			}
			if err = a.identityQuery.Dispatch(ctx, queryShop); err != nil {
				return err
			}
			shop := queryShop.Result

			cmdInvoice := &invoicing.CreateInvoiceBySubrIDCommand{
				SubscriptionID: subrID,
				AccountID:      args.AccountID,
				Customer: &subscriptingtypes.CustomerInfo{
					FullName: shop.Name,
					Phone:    shop.Phone,
					Email:    shop.Email,
				},
				Description: "Thanh toán phí khởi tạo extension",
				Classify:    service_classify.Telecom,
			}
			if err = a.invoiceAggr.Dispatch(ctx, cmdInvoice); err != nil {
				return err
			}

			invoiceID = cmdInvoice.Result.ID
			cmdPayment := &invoicing.PaymentInvoiceCommand{
				InvoiceID:       invoiceID,
				AccountID:       args.AccountID,
				OwnerID:         args.OwnerID,
				TotalAmount:     cmdInvoice.Result.TotalAmount,
				PaymentMethod:   args.PaymentMethod,
				ServiceClassify: service_classify.Telecom.Wrap(),
			}
			if err = a.invoiceAggr.Dispatch(ctx, cmdPayment); err != nil {
				return err
			}

			// create extension with subscription
			querySubr := &subscription.GetSubscriptionByIDQuery{
				ID:        subrID,
				AccountID: args.AccountID,
			}
			if err = a.subrQuery.Dispatch(ctx, querySubr); err != nil {
				return err
			}
			cmdExt := &etelecom.CreateExtensionArgs{
				UserID:          args.UserID,
				AccountID:       args.AccountID,
				HotlineID:       args.HotlineID,
				OwnerID:         args.OwnerID,
				SubscriptionID:  subrID,
				ExpiresAt:       querySubr.Result.CurrentPeriodEndAt,
				ExtensionNumber: args.ExtensionNumber,
			}
			res, err = a.createExtension(ctx, cmdExt)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
		return res, nil
	}

	// create by subscription id
	// make sure:
	//    - SubrID is not belong to any extension
	//    - Subr has valid CurrentPeriodEndAt
	if args.SubscriptionID != 0 {
		_, err := a.extensionStore(ctx).AccountID(args.AccountID).OptionalSubscriptionID(args.SubscriptionID).GetExtension()
		if err == nil {
			return nil, cm.Errorf(cm.FailedPrecondition, nil, "Extention đã tồn tại. Vui lòng kiểm tra hoặc gia hạn extention.")
		}
		if err != nil && cm.ErrorCode(err) != cm.NotFound {
			return nil, err
		}

		query := &subscription.GetSubscriptionByIDQuery{
			ID:        args.SubscriptionID,
			AccountID: args.AccountID,
		}
		if err = a.subrQuery.Dispatch(ctx, query); err != nil {
			return nil, err
		}
		subr := query.Result
		if subr.CurrentPeriodEndAt.Before(time.Now()) {
			return nil, cm.Errorf(cm.FailedPrecondition, nil, "Subscription đã hết hạn. Vui lòng gia hạn trước khi khởi tạo extension.")
		}

		cmd := &etelecom.CreateExtensionArgs{
			UserID:          args.UserID,
			AccountID:       args.AccountID,
			HotlineID:       args.HotlineID,
			OwnerID:         args.OwnerID,
			SubscriptionID:  subr.ID,
			ExpiresAt:       subr.CurrentPeriodEndAt,
			ExtensionNumber: args.ExtensionNumber,
		}
		return a.createExtension(ctx, cmd)
	}
	return nil, nil
}

func (a *EtelecomAggregate) ExtendExtension(ctx context.Context, args *etelecom.ExtendExtensionArgs) (*etelecom.Extension, error) {
	if err := args.Validate(); err != nil {
		return nil, err
	}
	ext, err := a.extensionStore(ctx).ID(args.ExtensionID).AccountID(args.AccountID).OptionalUserID(args.UserID).GetExtension()
	if err != nil {
		return nil, err
	}

	if args.SubscriptionID == 0 && ext.SubscriptionID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing subscription ID")
	}
	subrID := args.SubscriptionID
	if subrID == 0 {
		subrID = ext.SubscriptionID
	}

	querySubr := &subscription.GetSubscriptionByIDQuery{
		ID:        subrID,
		AccountID: args.AccountID,
	}
	if err = a.subrQuery.Dispatch(ctx, querySubr); err != nil {
		return nil, err
	}
	subr := querySubr.Result

	if args.SubscriptionPlanID != 0 {
		subrLines := subr.Lines
		if len(subrLines) != 1 || subrLines[0].PlanID != args.SubscriptionPlanID {
			// update subscription
			update := &subscription.UpdateSubscriptionInfoCommand{
				ID:        subr.ID,
				AccountID: subr.AccountID,
				Lines: []*subscription.SubscriptionLine{
					{
						PlanID:         args.SubscriptionPlanID,
						SubscriptionID: subr.ID,
						Quantity:       1,
					},
				},
			}
			if err = a.subrAggr.Dispatch(ctx, update); err != nil {
				return nil, err
			}
		}
	}

	// payment
	err = a.txDBMain.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		queryShop := &identity.GetShopByIDQuery{
			ID: args.AccountID,
		}
		if err = a.identityQuery.Dispatch(ctx, queryShop); err != nil {
			return err
		}
		shop := queryShop.Result
		cmdInvoice := &invoicing.CreateInvoiceBySubrIDCommand{
			SubscriptionID: subrID,
			AccountID:      args.AccountID,
			Customer: &subscriptingtypes.CustomerInfo{
				FullName: shop.Name,
				Phone:    shop.Phone,
				Email:    shop.Email,
			},
			Description: "Thanh toán gia hạn extension",
			Classify:    service_classify.Telecom,
		}
		if err = a.invoiceAggr.Dispatch(ctx, cmdInvoice); err != nil {
			return err
		}

		inv := cmdInvoice.Result
		cmdPayment := &invoicing.PaymentInvoiceCommand{
			InvoiceID:       inv.ID,
			AccountID:       args.AccountID,
			OwnerID:         shop.OwnerID,
			TotalAmount:     inv.TotalAmount,
			PaymentMethod:   args.PaymentMethod,
			ServiceClassify: service_classify.Telecom.Wrap(),
		}
		if err = a.invoiceAggr.Dispatch(ctx, cmdPayment); err != nil {
			return err
		}

		query := &subscription.GetSubscriptionByIDQuery{
			ID:        subrID,
			AccountID: args.AccountID,
		}
		if err = a.subrQuery.Dispatch(ctx, query); err != nil {
			return err
		}
		subr = query.Result

		// update extension
		updateExt := &etelecom.Extension{
			ExpiresAt: subr.CurrentPeriodEndAt,
		}
		return a.extensionStore(ctx).ID(args.ExtensionID).AccountID(args.AccountID).UpdateExtension(updateExt)
	})
	if err != nil {
		return nil, err
	}
	return a.extensionStore(ctx).ID(args.ExtensionID).AccountID(args.AccountID).GetExtension()
}

func (a *EtelecomAggregate) getTenant(ctx context.Context, args *etelecom.CreateExtensionArgs) (*etelecom.Tenant, error) {
	hotline, err := a.hotlineStore(ctx).ID(args.HotlineID).GetHotline()
	if err != nil {
		return nil, err
	}
	if hotline.Status != status3.P {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Hotline không hợp lệ")
	}

	query := a.tenantStore(ctx).ConnectionID(hotline.ConnectionID)

	queryConn := &connectioning.GetConnectionByIDQuery{
		ID: hotline.ConnectionID,
	}
	if err = a.connectionQuery.Dispatch(ctx, queryConn); err != nil {
		return nil, err
	}
	conn := queryConn.Result

	ownerID := args.OwnerID
	if conn.ConnectionMethod != connection_type.ConnectionMethodBuiltin {
		if ownerID == 0 {
			queryShop := &identity.GetShopByIDQuery{
				ID: args.AccountID,
			}
			if err = a.identityQuery.Dispatch(ctx, queryShop); err != nil {
				return nil, err
			}
			ownerID = queryShop.Result.OwnerID
		}
		if hotline.OwnerID != ownerID {
			return nil, cm.Errorf(cm.FailedPrecondition, nil, "Hotline không thuộc chủ shop này")
		}

		query = query.OwnerID(args.OwnerID)
	}
	return query.GetTenant()
}

func (a *EtelecomAggregate) DeleteExtension(ctx context.Context, id dot.ID) error {
	_, err := a.extensionStore(ctx).ID(id).SoftDelete()
	return err
}

func (a *EtelecomAggregate) UpdateExternalExtensionInfo(ctx context.Context, args *etelecom.UpdateExternalExtensionInfoArgs) error {
	update := &etelecom.Extension{
		HotlineID:         args.HotlineID,
		ExtensionNumber:   args.ExtensionNumber,
		ExtensionPassword: args.ExtensionPassword,
		TenantDomain:      args.TenantDomain,
		ExternalData: &etelecom.ExtensionExternalData{
			ID: args.ExternalID,
		},
	}
	return a.extensionStore(ctx).ID(args.ID).UpdateExtension(update)
}

func (a *EtelecomAggregate) RemoveUserOfExtension(ctx context.Context, args *etelecom.RemoveUserOfExtensionArgs) (int, error) {
	queryUser := &identity.GetAccountUserQuery{
		UserID:    args.UserID,
		AccountID: args.AccountID,
	}
	if err := a.identityQuery.Dispatch(ctx, queryUser); err != nil {
		return 0, cm.Errorf(cm.ErrorCode(err), err, "Không tìm thấy nhân viên")
	}

	query := a.extensionStore(ctx).ID(args.ExtensionID).AccountID(args.AccountID)
	ext, err := query.GetExtension()
	if ext == nil || err != nil {
		return 0, cm.Errorf(cm.ErrorCode(err), err, "Không tìm thấy máy nhánh")
	}

	if ext.UserID != args.UserID {
		return 0, cm.Errorf(cm.InvalidArgument, nil, "Máy nhánh không hợp lệ")
	}

	update, err := a.extensionStore(ctx).AccountID(args.AccountID).UserID(args.UserID).ID(args.ExtensionID).RemoveUserID()
	if err != nil {
		return 0, err
	}

	removedUserOfExtensionEvent := &etelecom.RemovedUserOfExtensionEvent{
		ID:        ext.ID,
		AccountID: args.AccountID,
		UserID:    args.UserID,
	}
	if err = a.eventBus.Publish(ctx, removedUserOfExtensionEvent); err != nil {
		return 0, err
	}
	return update, nil
}

func (a *EtelecomAggregate) AssignUserToExtension(ctx context.Context, args *etelecom.AssignUserToExtensionArgs) error {
	queryUserExtension := a.extensionStore(ctx).AccountID(args.AccountID).UserID(args.UserID)
	_, err := queryUserExtension.GetExtension()
	switch cm.ErrorCode(err) {
	case cm.NotFound:
	case cm.NoError:
		return cm.Errorf(cm.InvalidArgument, nil, "Nhân viên này đã có máy nhánh")
	default:
		return err
	}
	query := a.extensionStore(ctx).ID(args.ExtensionID).AccountID(args.AccountID)
	ext, err := query.GetExtension()
	if err != nil && cm.ErrorCode(err) == cm.NotFound {
		return cm.Errorf(cm.NotFound, nil, "Không tìm thấy máy nhánh")
	}
	if err != nil {
		return err
	}
	if ext.UserID != 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Máy nhánh đã được gán với người dùng khác")
	}

	queryUser := &identity.GetAccountUserQuery{
		UserID:    args.UserID,
		AccountID: args.AccountID,
	}
	if err = a.identityQuery.Dispatch(ctx, queryUser); err != nil {
		return cm.Errorf(cm.ErrorCode(err), err, "Không tìm thấy nhân viên")
	}

	update := &etelecom.Extension{
		UserID: args.UserID,
	}
	if err = query.UpdateExtension(update); err != nil {
		return err
	}

	extensionAssignedEvent := &etelecom.AssignedExtensionEvent{
		AccountID:       args.AccountID,
		UserID:          args.UserID,
		ExtensionNumber: ext.ExtensionNumber,
	}
	if err := a.eventBus.Publish(ctx, extensionAssignedEvent); err != nil {
		return err
	}
	return nil
}

func (a *EtelecomAggregate) ImportExtensions(ctx context.Context, args *etelecom.ImportExtensionsArgs) error {
	created := 0
	mapTenant := map[dot.ID]*etelecom.Tenant{}
	err := a.txDBEtelecom.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		for _, imExt := range args.ImportExtensions {
			if _, ok := mapTenant[imExt.TenantID]; !ok {
				tenant, err := a.tenantStore(ctx).ID(imExt.TenantID).OwnerID(imExt.OwnerID).GetTenant()
				if err != nil {
					return err
				}
				mapTenant[tenant.ID] = tenant
			}
			tenant := mapTenant[imExt.TenantID]
			extension, _ := a.extensionStore(ctx).AccountID(imExt.AccountID).HotlineID(imExt.HotlineID).TenantID(imExt.TenantID).ExtensionNumber(imExt.ExtensionNumber).GetExtension()
			if extension != nil {
				updateExt := &etelecom.Extension{
					ExpiresAt: imExt.ExpiresAt,
				}
				err := a.extensionStore(ctx).ID(extension.ID).AccountID(extension.AccountID).UpdateExtension(updateExt)
				if err != nil {
					return err
				}
			}else {
				cmd := &etelecom.Extension{
					ID:              cm.NewID(),
					AccountID:       imExt.AccountID,
					HotlineID:       imExt.HotlineID,
					ExtensionNumber: imExt.ExtensionNumber,
					TenantID:        imExt.TenantID,
					ExpiresAt:       imExt.ExpiresAt,
				}
				ext, err := a.extensionStore(ctx).CreateExtension(cmd)
				if err != nil {
					return err
				}

				externalExtensionResp, _err := a.telecomManager.CreateExtension(ctx, ext)
				if _err != nil {
					return _err
				}
				updateExt := &etelecom.UpdateExternalExtensionInfoArgs{
					ID:                externalExtensionResp.ExtensionID,
					HotlineID:         externalExtensionResp.HotlineID,
					ExternalID:        externalExtensionResp.ExternalID,
					ExtensionNumber:   externalExtensionResp.ExtensionNumber,
					ExtensionPassword: externalExtensionResp.ExtensionPassword,
					TenantDomain:      tenant.Domain,
				}
				if err := a.UpdateExternalExtensionInfo(ctx, updateExt); err != nil {
					return err
				}
			}
			created++
		}
		return nil
	})
	ll.S.Infof("Import extension success: %v/%v", created, len(args.ImportExtensions))
	return err
}
