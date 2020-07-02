package shop

import (
	"context"

	api "o.o/api/top/int/shop"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/etop/api/export"
	"o.o/backend/pkg/etop/authorize/auth"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/backend/tools/pkg/acl"
)

type ExportService struct {
	session.Session

	ExportInner *export.Service
}

func (s *ExportService) Clone() api.ExportService { res := *s; return &res }

//- khi truyền "export_type": "shop/orders" Cần quyền "order:export"
//- khi truyền "export_type": "shop/fulfillments" và filter theo:
//    - "fulfillments.id" Cần quyền "fulfillment:export"
//    - "money_transaction.id" Cần quyền "money_transaction:export"
func (s *ExportService) RequestExport(ctx context.Context, r *api.RequestExportRequest) (resp *api.RequestExportResponse, _err error) {
	claim := s.SS.Claim()
	authorization := auth.New()
	shop, roles := s.SS.Shop(), s.SS.Permission().Roles
	isTest := 0
	if shop != nil {
		isTest = shop.IsTest
	}
	switch r.ExportType {
	case "shop/orders":
		// Do not check permission for 3rd party requests
		if claim.AuthPartnerID == 0 && !authorization.Check(roles, string(acl.ShopOrderExport), isTest) {
			return nil, cm.Error(cm.PermissionDenied, "", nil)
		}
	case "shop/fulfillments":
		if claim.AuthPartnerID == 0 && !authorization.Check(roles, string(acl.ShopFulfillmentExport), isTest) {
			return nil, cm.Error(cm.PermissionDenied, "", nil)
		}

		for _, filter := range r.Filters {
			if filter.Name == "money_transaction.id" {
				if claim.AuthPartnerID == 0 && !authorization.Check(roles, string(acl.ShopMoneyTransactionExport), isTest) {
					return nil, cm.Error(cm.PermissionDenied, "", nil)
				}
				break
			}
		}
	}

	resp, err := s.ExportInner.RequestExport(ctx, s.SS.Claim(), s.SS.Shop(), s.SS.Claim().UserID, r)
	return resp, err
}

func (s *ExportService) GetExports(ctx context.Context, r *api.GetExportsRequest) (*api.GetExportsResponse, error) {
	resp, err := s.ExportInner.GetExports(ctx, s.SS.Shop().ID, r)
	return resp, err
}
