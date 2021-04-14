package aggregate

import (
	"context"
	"time"

	"o.o/api/etelecom"
	"o.o/api/main/identity"
	"o.o/api/main/invoicing"
	"o.o/api/subscripting/subscription"
	subscriptingtypes "o.o/api/subscripting/types"
	"o.o/api/top/types/etc/payment_method"
	"o.o/api/top/types/etc/service_classify"
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

func (a *EtelecomAggregate) createExtension(ctx context.Context, args *etelecom.CreateExtensionArgs) (*etelecom.Extension, error) {
	if err := args.Validate(); err != nil {
		return nil, err
	}
	tenantDomain, err := a.getTenantDomain(ctx, args)
	if err != nil {
		return nil, err
	}

	ext, err := a.extensionStore(ctx).UserID(args.UserID).AccountID(args.AccountID).HotlineID(args.HotlineID).GetExtension()
	switch cm.ErrorCode(err) {
	case cm.NoError:
		if ext.ExtensionNumber != "" {
			return nil, cm.Errorf(cm.FailedPrecondition, nil, "Extension đã được tạo cho người dùng này.")
		}
	case cm.NotFound:
		// create new one
		var extension etelecom.Extension
		if err = scheme.Convert(args, &extension); err != nil {
			return nil, err
		}

		extension.ID = cm.NewID()
		ext, err = a.extensionStore(ctx).CreateExtension(&extension)
		if err != nil {
			return nil, err
		}
	default:
		return nil, err
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
		TenantDomain:      tenantDomain,
	}
	if err = a.UpdateExternalExtensionInfo(ctx, updateExt); err != nil {
		return nil, err
	}

	return a.extensionStore(ctx).ID(ext.ID).GetExtension()
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
				UserID:         args.UserID,
				AccountID:      args.AccountID,
				HotlineID:      args.HotlineID,
				OwnerID:        args.OwnerID,
				SubscriptionID: subrID,
				ExpiresAt:      querySubr.Result.CurrentPeriodEndAt,
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
			UserID:         args.UserID,
			AccountID:      args.AccountID,
			HotlineID:      args.HotlineID,
			OwnerID:        args.OwnerID,
			SubscriptionID: subr.ID,
			ExpiresAt:      subr.CurrentPeriodEndAt,
		}
		return a.createExtension(ctx, cmd)
	}
	return nil, nil
}

func (a *EtelecomAggregate) ExtendExtension(ctx context.Context, args *etelecom.ExtendExtensionArgs) (*etelecom.Extension, error) {
	if err := args.Validate(); err != nil {
		return nil, err
	}
	ext, err := a.extensionStore(ctx).ID(args.ExtensionID).AccountID(args.AccountID).UserID(args.UserID).GetExtension()
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

func (a *EtelecomAggregate) getTenantDomain(ctx context.Context, args *etelecom.CreateExtensionArgs) (tenantDomain string, _ error) {
	hotline, err := a.hotlineStore(ctx).ID(args.HotlineID).GetHotline()
	if err != nil {
		return "", err
	}

	// get ownerID
	ownerID := args.OwnerID
	if ownerID == 0 {
		shopQuery := &identity.GetShopByIDQuery{
			ID: args.AccountID,
		}
		if err = a.identityQuery.Dispatch(ctx, shopQuery); err != nil {
			return "", err
		}
		ownerID = shopQuery.Result.OwnerID
	}

	_, shopConn, err := a.telecomManager.GetTelecomConnection(ctx, hotline.ConnectionID, ownerID)
	if err != nil {
		return "", err
	}
	return shopConn.TelecomData.TenantDomain, nil
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
