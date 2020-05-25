package shop

import (
	"context"

	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/etop/api/export"
	"o.o/backend/pkg/etop/authorize/auth"
	"o.o/backend/tools/pkg/acl"
)

type ExportService struct{}

func (s *ExportService) Clone() *ExportService { res := *s; return &res }

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
