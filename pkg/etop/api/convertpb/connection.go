package convertpb

import (
	"etop.vn/api/main/connectioning"
	"etop.vn/api/top/int/types"
	"etop.vn/backend/pkg/common/apifw/cmapi"
)

func PbShopConnection(c *connectioning.ShopConnection) *types.ShopConnection {
	if c == nil {
		return nil
	}
	res := &types.ShopConnection{
		ShopID:       c.ShopID,
		ConnectionID: c.ConnectionID,
		Status:       c.Status,
		CreatedAt:    cmapi.PbTime(c.CreatedAt),
		UpdatedAt:    cmapi.PbTime(c.UpdatedAt),
		DeletedAt:    cmapi.PbTime(c.DeletedAt),
		IsGlobal:     c.IsGlobal,
	}
	if c.ExternalData != nil {
		res.ExternalData = &types.ShopConnectionExternalData{
			Email:  c.ExternalData.Email,
			UserID: c.ExternalData.UserID,
		}
	}
	return res
}

func PbShopConnections(items []*connectioning.ShopConnection) []*types.ShopConnection {
	result := make([]*types.ShopConnection, len(items))
	for i, item := range items {
		result[i] = PbShopConnection(item)
	}
	return result
}

func PbConnection(c *connectioning.Connection) *types.Connection {
	if c == nil {
		return nil
	}
	return &types.Connection{
		ID:                 c.ID,
		Name:               c.Name,
		Status:             c.Status,
		CreatedAt:          cmapi.PbTime(c.CreatedAt),
		UpdatedAt:          cmapi.PbTime(c.UpdatedAt),
		ConnectionType:     c.ConnectionType,
		ConnectionSubtype:  c.ConnectionSubtype,
		ConnectionMethod:   c.ConnectionMethod,
		ConnectionProvider: c.ConnectionProvider,
	}
}

func PbConnections(items []*connectioning.Connection) []*types.Connection {
	result := make([]*types.Connection, len(items))
	for i, item := range items {
		result[i] = PbConnection(item)
	}
	return result
}

func PbConnectionService(in *connectioning.ConnectionService) *types.ConnectionService {
	if in == nil {
		return nil
	}
	return &types.ConnectionService{
		ServiceID: in.ServiceID,
		Name:      in.Name,
	}
}

func PbConnectionServices(items []*connectioning.ConnectionService) []*types.ConnectionService {
	result := make([]*types.ConnectionService, len(items))
	for i, item := range items {
		result[i] = PbConnectionService(item)
	}
	return result
}
