package shop

import (
	"context"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/etop/api/export"
	"etop.vn/backend/pkg/etop/authorize/auth"
	"etop.vn/backend/tools/pkg/acl"
)

func init() {
	bus.AddHandlers("export",
		exportService.RequestExport,
		exportService.GetExports,
	)
}

//- khi truyền "export_type": "shop/orders" Cần quyền "order:export"
//- khi truyền "export_type": "shop/fulfillments" và filter theo:
//    - "fulfillments.id" Cần quyền "fulfillment:export"
//    - "money_transaction.id" Cần quyền "money_transaction:export"
func (s *ExportService) RequestExport(ctx context.Context, r *RequestExportEndpoint) (_err error) {
	claim := r.Context
	authorization := auth.New()
	isTest := 0
	if claim.Shop != nil {
		isTest = claim.Shop.IsTest
	}
	switch r.ExportType {
	case "shop/orders":
		// Do not check permission for 3rd party requests
		if claim.AuthPartnerID == 0 && !authorization.Check(claim.Roles, string(acl.ShopOrderExport), isTest) {
			return cm.Error(cm.PermissionDenied, "", nil)
		}
	case "shop/fulfillments":
		if claim.AuthPartnerID == 0 && !authorization.Check(claim.Roles, string(acl.ShopFulfillmentExport), isTest) {
			return cm.Error(cm.PermissionDenied, "", nil)
		}

		for _, filter := range r.Filters {
			if filter.Name == "money_transaction.id" {
				if claim.AuthPartnerID == 0 && !authorization.Check(claim.Roles, string(acl.ShopMoneyTransactionExport), isTest) {
					return cm.Error(cm.PermissionDenied, "", nil)
				}
				break
			}
		}
	}

	resp, err := export.ServiceImpl.RequestExport(ctx, r.Context, r.Context.Shop, r.Context.UserID, r.RequestExportRequest)
	r.Result = resp
	return err
}

func (s *ExportService) GetExports(ctx context.Context, r *GetExportsEndpoint) error {
	resp, err := export.ServiceImpl.GetExports(ctx, r.Context.Shop.ID, r.GetExportsRequest)
	r.Result = resp
	return err
}
