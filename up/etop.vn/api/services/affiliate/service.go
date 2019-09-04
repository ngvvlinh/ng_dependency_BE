package affiliate

import (
	"context"
)

// +gen:api

type Aggregate interface {
	CreateOrUpdateCommissionSetting(context.Context, *CreateCommissionSettingArgs) (*CommissionSetting, error)
}

type CreateCommissionSettingArgs struct {
	ProductID int64
	AccountID int64
	Amount    int32
	Unit      string
	Type      string
}

type QueryService interface {
	GetCommissionByProductIDs(context.Context, *GetCommissionByProductIDsArgs) ([]*CommissionSetting, error)
}

type GetCommissionByProductIDsArgs struct {
	AccountID  int64
	ProductIDs []int64
}
