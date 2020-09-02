package api

import (
	"context"

	api "o.o/api/top/services/handler"
	pbcm "o.o/api/top/types/common"
)

type MiscService struct{}

func (s *MiscService) Clone() api.MiscService {
	res := *s
	return &res
}

func (s *MiscService) VersionInfo(context.Context, *pbcm.Empty) (*pbcm.VersionInfoResponse, error) {
	result := &pbcm.VersionInfoResponse{
		Service: "event-handler",
		Version: "1.0",
	}
	return result, nil
}
