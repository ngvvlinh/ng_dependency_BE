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
	bus.AddHandler("kiotviet", RequestUpdatedCategories)
}

func syncCategories(cmd *SyncCategoriesCommand) (_err error) {
	defer func() {
		err := recover()
		if err != nil {
			ll.Error("Panic", l.Any("err", err), l.Object("cmd", cmd), l.Stack())
			_err = cm.ErrorTrace(cm.RuntimePanic, cm.F("Panic: %v", err), nil)
		}
	}()

	for {
		state := cmd.SyncState
		ll.Info("Sync categories", l.Int64("source", cmd.SourceID), l.Time("since", state.Since), l.Int("page", state.Page))

		ctx := bus.NewRootContext(context.Background())
		err := syncCategoriesStep(ctx, cmd)
		ll.Info("Sync categories", l.Int64("source", cmd.SourceID), l.Any("cmd", cmd))
		if err != nil {
			return err
		}
		if cmd.Result.Done {
			ll.Info("Sync categories: done", l.Int64("source", cmd.SourceID))
			return nil
		}

		cmd.SyncState = cmd.Result.NextState
		cmd.Result.Done = false
		cmd.Result.Response = nil
		cmd.Result.NextState = SyncState{}
	}
}

func syncCategoriesStep(ctx context.Context, cmd *SyncCategoriesCommand) error {
	state := cmd.SyncState
	now := time.Now()
	reqCmd := &RequestUpdatedCategoriesCommand{
		Connection: cmd.Connection,
		UpdatedQuery: UpdatedQuery{
			LastUpdated:      cmd.SyncState.Since,
			Page:             cmd.SyncState.Page,
			IncludeRemoveIDs: true,
		},
	}
	if err := bus.Dispatch(ctx, reqCmd); err != nil {
		return err
	}

	result := reqCmd.Result
	data := result.Data
	updateData := data

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
			// Some categories have no ModifiedDate. We already handle this case
			// in AdvanceState.
		}

		// Only update just enough data, because the next request will continue
		// from timeEnd.
		if !done && timeStart.Before(timeEnd) {
			for i := len(data) - 1; ; i-- {
				ti := data[i].UpdatedAt.ToTime()
				if ti.Before(timeEnd) {
					updateData = data[:i+1]
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
	updateCmd := &model.SyncUpdateCategoriesCommand{
		SourceID:   cmd.SourceID,
		Data:       toExternalCategories(updateData),
		DeletedIDs: httpreq.FromStrings(reqCmd.Result.DeletedIDs),
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

func RequestUpdatedCategories(ctx context.Context, cmd *RequestUpdatedCategoriesCommand) error {
	var result UpdatedCategoriesResponse

	reqURL, _ := url.Parse(BaseURLPublic + "categories")
	reqURL.RawQuery = cmd.UpdatedQuery.Build().Encode()
	_, err := getx(cmd.Connection, reqURL.String(), &result)
	if err != nil {
		return cm.Error(cm.Unknown, "Error requesting Kiotviet", err)
	}

	cmd.Result = &result
	return nil
}

func toExternalCategories(data []*Category) []*model.ProductSourceCategoryExternal {
	result := make([]*model.ProductSourceCategoryExternal, len(data))
	for i, c := range data {
		result[i] = &model.ProductSourceCategoryExternal{
			ExternalID:        string(c.ID),
			ExternalParentID:  string(c.ParentID),
			ExternalName:      c.Name,
			ExternalUpdatedAt: c.UpdatedAt.ToTime(),
			ExternalCreatedAt: c.CreatedAt.ToTime(),

			// ExternalCode: (missing)
		}
	}
	return result
}
