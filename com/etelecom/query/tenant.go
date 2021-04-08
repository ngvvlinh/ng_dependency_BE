package query

import (
	"context"

	"o.o/api/etelecom"
	cm "o.o/backend/pkg/common"
)

func (q *QueryService) GetTenant(ctx context.Context, args *etelecom.GetTenantArgs) (*etelecom.Tenant, error) {
	query := q.tenantStore(ctx)
	count := 0
	if args.ID != 0 {
		query = query.ID(args.ID)
		count++
	}
	if args.OwnerID != 0 {
		query = query.OwnerID(args.OwnerID)
		count++
	}
	if args.ConnectionID != 0 {
		query = query.ConnectionID(args.ConnectionID)
		count++
	}
	if args.ID == 0 && count < 2 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing required params")
	}

	return query.GetTenant()
}

func (q *QueryService) ListTenants(ctx context.Context, args *etelecom.ListTenantsArgs) (*etelecom.ListTenantsResponse, error) {
	query := q.tenantStore(ctx).WithPaging(args.Paging)
	if args.OwnerID != 0 {
		query = query.OwnerID(args.OwnerID)
	}
	if args.ConnectionID != 0 {
		query = query.ConnectionID(args.ConnectionID)
	}
	res, err := query.ListTenants()
	if err != nil {
		return nil, err
	}
	return &etelecom.ListTenantsResponse{
		Tenants: res,
		Paging:  query.GetPaging(),
	}, nil
}
