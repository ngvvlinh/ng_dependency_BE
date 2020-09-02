package api

import (
	"context"

	pbcm "o.o/api/top/types/common"
)

type MiscService struct{}

func (s *MiscService) Clone() *MiscService {
	res := *s
	return &res
}

func (s *MiscService) VersionInfo(ctx context.Context, q *pbcm.Empty) (*pbcm.VersionInfoResponse, error) {
	result := &pbcm.VersionInfoResponse{
		Service: "pgevent-forwarder",
		Version: "0.1",
	}
	return result, nil
}
