package partner

import (
	"context"

	cm "o.o/backend/pkg/common"
)

type HistoryService struct{}

func (s *HistoryService) Clone() *HistoryService { res := *s; return &res }

func (s *HistoryService) GetChanges(ctx context.Context, r *GetChangesEndpoint) error {
	return cm.ErrTODO
}
