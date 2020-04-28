package sqlstore

import (
	"context"

	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/etop/model"
)

func init() {
	bus.AddHandlers("sql",
		GetShippingSources,
		GetShippingSourceInternal,
		UpdateOrCreateShippingSourceInternal,
		CreateShippingSource,
		GetShippingSource,
	)
}

func GetShippingSources(ctx context.Context, query *model.GetShippingSources) error {
	s := x.Table("shipping_source")
	if query.Type != 0 {
		s = s.Where("type = ?", query.Type)
	}
	if len(query.Names) != 0 {
		s = s.In("name", query.Names)
	}
	err := s.Find((*model.ShippingSources)(&query.Result))
	return err
}

func GetShippingSource(ctx context.Context, query *model.GetShippingSource) error {
	s := x.Table("shipping_source")
	if query.ID != 0 {
		s = s.Where("id = ?", query.ID)
	}
	if query.Name != "" {
		s = s.Where("name = ?", query.Name)
	}
	if query.Username != "" {
		s = s.Where("username = ?", query.Username)
	}
	if query.Type != 0 {
		s = s.Where("type = ?", query.Type)
	}
	var shippingSource = new(model.ShippingSource)
	if err := s.ShouldGet(shippingSource); err != nil {
		return err
	}
	query.Result.ShippingSource = shippingSource
	querySourceInternal := &model.GetShippingSourceInternal{
		ID: shippingSource.ID,
	}
	if err := bus.Dispatch(ctx, querySourceInternal); err != nil {
		return err
	}
	query.Result.ShippingSourceInternal = querySourceInternal.Result
	return nil
}

func CreateShippingSource(ctx context.Context, cmd *model.CreateShippingSource) error {
	if cmd.Name == "" {
		return cm.Error(cm.InvalidArgument, "Missing name", nil)
	}
	if cmd.Type == 0 {
		return cm.Error(cm.InvalidArgument, "Missing shipping provider type", nil)
	}

	query := &model.GetShippingSource{
		Name:     cmd.Name,
		Username: cmd.Username,
		Type:     cmd.Type,
	}
	bus.Dispatch(ctx, query)

	err := inTransaction(func(s Qx) error {
		newID := cm.NewID()
		ss := &model.ShippingSource{
			ID:       newID,
			Name:     cmd.Name,
			Username: cmd.Username,
			Type:     cmd.Type.String(),
		}
		if query.Result.ShippingSource != nil {
			ss.ID = query.Result.ShippingSource.ID
			if err := s.Table("shipping_source").Where("id = ?", ss.ID).ShouldUpdate(ss); err != nil {
				return err
			}
		} else {
			if err := s.Table("shipping_source").ShouldInsert(ss); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	if err = bus.Dispatch(ctx, query); err != nil {
		return err
	}
	cmd.Result = query.Result
	return nil
}

/*
func VTPostRequestInternalInfo(ctx context.Context, cmd *model.VTPostRequestInternalInfo) error {
	return vtpostRequestInternalInfo(ctx, x, cmd)
}

func vtpostRequestInternalInfo(ctx context.Context, x Qx, cmd *model.VTPostRequestInternalInfo) error {
	if cmd.ShippingSourceID == 0 || cmd.Username == "" || cmd.Password == "" {
		return cm.Error(cm.InvalidArgument, "Missing arguments", nil)
	}
	loginCmd := &vtpost.LoginCommand{
		Request: &vtpostClient.LoginRequest{
			Username: cmd.Username,
			Password: cmd.Password,
		},
	}
	if err := bus.Dispatch(ctx, loginCmd); err != nil {
		return err
	}
	res := loginCmd.Result.Data
	customerID := res.UserId
	token := res.ApiKey
	// Try getting expiresAt from JWT first, then from response
	expiresAt := cm.GetJWTExpires(token)
	if expiresAt.IsZero() {
		now := time.Now()
		expiresAt = now.AddDate(0, 1, 0)
	}

	update := &model.UpdateOrCreateShippingSourceInternal{
		ID:          cmd.ShippingSourceID,
		AccessToken: token,
		ExpiresAt:   expiresAt,
		Secret: &model.ShippingSourceSecret{
			CustomerID: customerID,
			Username:   cmd.Username,
			Password:   cm.EncodePassword(cmd.Password),
		},
	}
	if err := updateOrCreateShippingSourceInternal(ctx, x, update); err != nil {
		return err
	}
	return nil
}
*/

func GetShippingSourceInternal(ctx context.Context, query *model.GetShippingSourceInternal) error {
	if query.ID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing ID", nil)
	}

	var shippingSourceState = new(model.ShippingSourceInternal)
	if _, err := x.Table("shipping_source_internal").
		Where("id = ?", query.ID).Get(shippingSourceState); err != nil {
		return err
	}
	query.Result = shippingSourceState
	return nil
}

func UpdateOrCreateShippingSourceInternal(ctx context.Context, cmd *model.UpdateOrCreateShippingSourceInternal) error {
	return updateOrCreateShippingSourceInternal(ctx, x, cmd)
}

func updateOrCreateShippingSourceInternal(ctx context.Context, x Qx, cmd *model.UpdateOrCreateShippingSourceInternal) error {
	if cmd.ID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing ID", nil)
	}
	shippingSourceInternal := &model.ShippingSourceInternal{
		ID:          cmd.ID,
		LastSyncAt:  cmd.LastSyncAt,
		AccessToken: cmd.AccessToken,
		ExpiresAt:   cmd.ExpiresAt,
		Secret:      cmd.Secret,
	}

	var ssState = new(model.ShippingSourceInternal)
	if has, err := x.Table("shipping_source_internal").Where("id = ?", cmd.ID).
		Get(ssState); err != nil {
		return nil
	} else if !has {
		if err := x.Table("shipping_source_internal").ShouldInsert(shippingSourceInternal); err != nil {
			return err
		}
	} else if has {
		if err := x.Table("shipping_source_internal").Where("id = ?", cmd.ID).
			ShouldUpdate(shippingSourceInternal); err != nil {
			return err
		}
	}

	cmd.Result.Updated = 1
	return nil
}
