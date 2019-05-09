package supplier

import (
	"context"

	etopP "etop.vn/backend/pb/etop"
	supplierP "etop.vn/backend/pb/etop/supplier"
	kiotvietP "etop.vn/backend/pb/integration/kiotviet"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	cmGrpc "etop.vn/backend/pkg/common/grpc"
	"etop.vn/backend/pkg/common/l"
	"etop.vn/backend/pkg/common/validate"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/backend/pkg/integration/kiotviet"
	supplierW "etop.vn/backend/wrapper/etop/supplier"
	kiotvietW "etop.vn/backend/wrapper/integration/kiotviet"
)

func init() {
	bus.AddHandler("api", ConnectKiotviet)
	bus.AddHandler("api", ConnectKiotvietTest)
}

var tokenKV string

func Init(token string) {
	if token == "" {
		ll.Fatal("Empty token")
	}
	tokenKV = token
}

func ConnectKiotvietTest(ctx context.Context, r *supplierW.ConnectKiotvietTestEndpoint) error {
	retailerID := r.Kiotviet.RetailerId
	clientID := r.Kiotviet.ClientId
	secret := r.Kiotviet.ClientSecret
	_, branches, err := kiotviet.NewConnection(ctx, retailerID, clientID, secret)
	if err != nil {
		return cm.Error(cm.Unknown, "Error requesting Kiotviet", err)
	}

	r.Result = &supplierP.ConnectKiotvietTestResponse{
		Branches: supplierP.PbBranches(kiotviet.ToModelBranches(branches)),
	}
	return nil
}

func ConnectKiotviet(ctx context.Context, r *supplierW.ConnectKiotvietEndpoint) error {
	var ok bool
	if r.Name, ok = validate.NormalizeName(r.Name); !ok {
		return cm.Error(cm.InvalidArgument, "Invalid name", nil)
	}
	if r.BranchId == "" {
		return cm.Error(cm.InvalidArgument, "Missing BranchID", nil)
	}

	retailerID := r.Kiotviet.RetailerId
	clientID := r.Kiotviet.ClientId
	secret := r.Kiotviet.ClientSecret
	if retailerID == "" {
		return cm.Error(cm.InvalidArgument, "Missing RetailerID", nil)
	}
	if clientID == "" || secret == "" {
		return cm.Error(cm.InvalidArgument, "Missing Kiotviet credential", nil)
	}

	psQuery := &model.GetProductSourceQuery{
		GetProductSourceProps: model.GetProductSourceProps{
			Type:        model.TypeKiotviet,
			ExternalKey: r.Kiotviet.RetailerId,
		},
	}
	if err := bus.Dispatch(ctx, psQuery); cm.ErrorCode(err) != cm.NotFound {
		if err == nil {
			return cm.Error(cm.AlreadyExists, "Kiotviet account already exist", nil)
		}
		return err
	}

	conn, branches, err := kiotviet.NewConnection(ctx, retailerID, clientID, secret)
	if err != nil {
		return cm.Error(cm.Unknown, "Error requesting Kiotviet", err)
	}

	var defaultBranchID string
	for _, b := range branches {
		if string(b.ID) == r.BranchId {
			defaultBranchID = r.BranchId
		}
	}
	if defaultBranchID == "" {
		return cm.Error(cm.NotFound, "BranchID not found", nil)
	}

	if r.UrlSlug != "" && !validate.URLSlug(r.UrlSlug) {
		return cm.Error(cm.InvalidArgument, "Thông tin url_slug không hợp lệ. Vui lòng kiểm tra lại.", nil)
	}

	cmd := &model.CreateSupplierKiotvietCommand{
		OwnerID: r.Context.UserID,
		IsTest:  r.Context.User.IsTest != 0,
		SupplierInfo: model.SupplierInfo{
			Name: r.Name,
		},
		Kiotviet: model.SupplierKiotviet{
			RetailerID:     retailerID,
			ClientID:       clientID,
			ClientSecret:   secret,
			ClientToken:    conn.TokenStr,
			ExpiresAt:      conn.ExpiresAt,
			KiotvietStatus: model.StatusCreated,
			Branches:       kiotviet.ToModelBranches(branches),
		},
		DefaultBranchID: defaultBranchID,
		URLSlug:         r.UrlSlug,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &supplierP.ConnectKiotvietResponse{
		Supplier: etopP.PbSupplier(cmd.Result.Supplier),
	}

	ctx = cmGrpc.AppendAccessToken(ctx, tokenKV)
	req := &kiotvietP.SyncProductSourceRequest{
		Id:            cmd.Result.ProductSource.ID,
		FromBeginning: true,
	}
	if _, err := kiotvietW.Client.SyncProductSource(ctx, req); err != nil {
		ll.Error("Sync Supplier", l.Error(err))
	}
	return nil
}
