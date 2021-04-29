package query

import (
	"context"

	"o.o/api/etelecom"
	"o.o/api/main/connectioning"
	"o.o/api/top/types/etc/connection_type"
	cm "o.o/backend/pkg/common"
	"o.o/capi/dot"
)

func (q *QueryService) GetTenantByConnection(ctx context.Context, args *etelecom.GetTenantByConnectionArgs) (*etelecom.Tenant, error) {
	if args.ConnectionID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing connection ID")
	}
	query := q.tenantStore(ctx).ConnectionID(args.ConnectionID)

	queryConn := &connectioning.GetConnectionByIDQuery{
		ID: args.ConnectionID,
	}
	if err := q.connectionQS.Dispatch(ctx, queryConn); err != nil {
		return nil, err
	}
	conn := queryConn.Result
	if conn.ConnectionMethod != connection_type.ConnectionMethodBuiltin {
		if args.OwnerID == 0 {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing owner ID")
		}
		query = query.OwnerID(args.OwnerID)
	}

	return query.GetTenant()
}

func (q *QueryService) GetTenantByID(ctx context.Context, id dot.ID) (*etelecom.Tenant, error) {
	if id == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing tenant ID")
	}
	return q.tenantStore(ctx).ID(id).GetTenant()
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
