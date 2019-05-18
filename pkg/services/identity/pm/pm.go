package pm

import (
	"context"

	"etop.vn/api/main/identity"
)

type ProcessManager struct {
	identityQuery identity.QueryService
}

func NewProcessManager(identityQuery identity.QueryService) *ProcessManager {
	return &ProcessManager{
		identityQuery: identityQuery,
	}
}

func (pm *ProcessManager) GetShopByID(ctx context.Context, args *identity.GetShopByIDQueryArgs) (*identity.Shop, error) {
	return pm.identityQuery.GetShopByID(ctx, args)
}
