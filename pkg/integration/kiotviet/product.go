package kiotviet

import (
	"context"
	"net/url"
	"time"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/httpreq"
	"etop.vn/backend/pkg/common/l"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/backend/pkg/integration/kiotviet/ssm"
)

func init() {
	bus.AddHandler("kiotviet", RequestUpdatedProducts)
}

func syncProducts(cmd *SyncProductsCommand) (_err error) {
	defer func() {
		err := recover()
		if err != nil {
			ll.Error("Panic", l.Any("err", err), l.Object("cmd", cmd), l.Stack())
			_err = cm.ErrorTrace(cm.RuntimePanic, cm.F("Panic: %v", err), nil)
		}
	}()

	for {
		state := cmd.SyncState
		ll.Info("Sync products", l.Int64("source", cmd.SourceID), l.Time("since", state.Since), l.Int("page", state.Page))

		ctx := bus.NewRootContext(context.Background())
		err := syncProductsStep(ctx, cmd)
		if err != nil {
			return err
		}

		ll.Info("Sync products", l.Int64("source", cmd.SourceID), l.Any("cmd", cmd))
		if cmd.Result.Done {
			ll.Info("Sync products: done", l.Int64("source", cmd.SourceID))
			return nil
		}

		cmd.SyncState = cmd.Result.NextState
		cmd.Result.Done = false
		cmd.Result.Response = nil
		cmd.Result.NextState = SyncState{}
	}
	return nil
}

func syncProductsStep(ctx context.Context, cmd *SyncProductsCommand) error {
	state := cmd.SyncState
	now := time.Now()
	reqCmd := &RequestUpdatedProductsCommand{
		Connection: cmd.Connection,
		UpdatedQuery: UpdatedQuery{
			LastUpdated:      state.Since,
			Page:             state.Page,
			IncludeRemoveIDs: true,
		},
	}
	if err := bus.Dispatch(ctx, reqCmd); err != nil {
		return err
	}

	result := reqCmd.Result
	data := result.Data
	// updateData := data

	var done bool
	var nextState SyncState
	if err := func() error {
		if len(data) == 0 {
			done = true
			nextState = ssm.AdvanceState(state, time.Millisecond,
				ssm.AdvanceArgs{
					Done: done,
					Size: 0,
				})
			return nil
		}

		if result.TotalResponse.PageSize <= 0 {
			ll.Error("Unexpected page size", l.Any("req", reqCmd))
			return cm.Error(cm.ExternalServiceError, "Unexpected page size from Kiotviet", nil)
		}
		if len(data) < int(result.TotalResponse.PageSize) {
			done = true
		}

		timeStart := data[0].UpdatedAt.ToTime()
		timeEnd := data[len(data)-1].UpdatedAt.ToTime()
		if timeStart.IsZero() || timeEnd.IsZero() || timeEnd.Before(timeStart) {
			// Some products have no ModifiedDate. We already handle this case
			// in AdvanceState.
		}

		// Only update just enough data, because the next request will continue
		// from timeEnd.
		if !done && timeStart.Before(timeEnd) {
			for i := len(data) - 1; ; i-- {
				ti := data[i].UpdatedAt.ToTime()
				if ti.Before(timeEnd) {
					// updateData = data[:i+1]
					break
				}
			}
		}
		nextState = ssm.AdvanceState(state, time.Millisecond,
			ssm.AdvanceArgs{
				Done:  done,
				Size:  len(data),
				Start: timeStart,
				End:   timeEnd,
			})
		return nil
	}(); err != nil {
		return err
	}

	savedSyncState := SavedSyncState{
		LastSyncAt: now,
		State:      &state,
		NextState:  &nextState,
		Done:       done,
	}
	updateCmd := &model.SyncUpdateProductsCommand{
		SourceID: cmd.SourceID,

		// Bug: Use data instead
		Data:       toExternalProducts(data, cmd.BranchID),
		DeletedIDs: httpreq.FromStrings(result.DeletedIDs),
		LastSyncAt: now,
		SyncState:  savedSyncState.ToJSON(),
	}
	if err := bus.Dispatch(ctx, updateCmd); err != nil {
		return err
	}

	cmd.Result.Response = reqCmd.Result
	cmd.Result.Done = done
	cmd.Result.NextState = nextState
	return nil
}

func RequestUpdatedProducts(ctx context.Context, cmd *RequestUpdatedProductsCommand) error {
	var result UpdatedProductsResponse

	reqURL, _ := url.Parse(BaseURLPublic + "products")
	urlQuery := cmd.UpdatedQuery.Build()
	urlQuery.Set("includeInventory", "1")
	reqURL.RawQuery = urlQuery.Encode()
	_, err := getx(cmd.Connection, reqURL.String(), &result)
	if err != nil {
		return cm.Error(cm.Unknown, "Error requesting Kiotviet", err)
	}

	cmd.Result = &result
	return nil
}

func toExternalProducts(data []*Product, branchID string) []*model.VariantExternalWithQuantity {
	result := make([]*model.VariantExternalWithQuantity, len(data))
	for i, p := range data {
		units := make([]*model.Unit, len(p.Units))
		for i, u := range p.Units {
			units[i] = &model.Unit{
				ID:       string(u.ID),
				Code:     u.Code,
				Name:     u.Name,
				FullName: u.FullName,
				Unit:     u.Unit,
				UnitConv: p.Conversionvalue,
				Price:    int(p.BasePrice),
			}
		}

		attrs := make([]model.ProductAttribute, len(p.Attributes))
		for i, a := range p.Attributes {
			attrs[i] = model.ProductAttribute{
				Name:  a.Name,
				Value: a.Value,
			}
		}

		masterProductID := string(p.MasterProductID)
		if masterProductID == "" {
			masterProductID = string(p.ID)
		}

		vx := &model.VariantExternal{
			ExternalProductID:  masterProductID,
			ExternalPrice:      int(p.BasePrice),
			ExternalBaseUnitID: string(p.MasterUnitID),
			ExternalUnitConv:   p.Conversionvalue,
			ExternalAttributes: attrs,
			ProductExternalCommon: model.ProductExternalCommon{
				ExternalID:         string(p.ID),
				ExternalCategoryID: string(p.CategoryID),
				ExternalCode:       p.Code,
				ExternalName:       p.FullName, // TODO(qv): Split the fullname and shortname

				ExternalUnit: string(p.Unit),
				// ExternalUnits:     units,
				ExternalImageURLs: p.Images,

				ExternalUpdatedAt: time.Time(p.UpdatedAt),
				// ExternalCreatedAt: (missing)
			},
		}
		vxq := &model.VariantExternalWithQuantity{
			Variant: vx,
		}
		for _, in := range p.Inventories {
			if string(in.BranchID) == branchID {
				vxq.QuantityOnHand = int(in.OnHand)
				vxq.QuantityReserved = int(in.Reserved)
			}
		}

		// Kiotviet returns IsActive and AllowSales:
		// - IsActive:   The product is allowed to sale from all channels.
		// - AllowSales: The product is allowed to sale on Kiotviet POS.
		if p.IsActive {
			vx.ExternalStatus = 1
		} else {
			vx.ExternalStatus = -1
		}
		result[i] = vxq
	}
	return result
}
