package xshop

import (
	"context"

	api "o.o/api/top/external/shop"
	externaltypes "o.o/api/top/external/types"
	pbcm "o.o/api/top/types/common"
	cm "o.o/backend/pkg/common"
)

type HistoryService struct{}

func (s *HistoryService) Clone() api.HistoryService { res := *s; return &res }

func (s *HistoryService) GetChanges(ctx context.Context, r *pbcm.Empty) (*externaltypes.Callback, error) {
	return nil, cm.ErrTODO
}
