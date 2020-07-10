package affiliate

import (
	"context"

	"o.o/api/top/int/affiliate"
	pbcm "o.o/api/top/types/common"
)

type MiscService struct{}

func (s *MiscService) Clone() affiliate.MiscService {
	res := *s
	return &res
}

func (s *MiscService) VersionInfo(ctx context.Context, q *pbcm.Empty) (*pbcm.VersionInfoResponse, error) {
	result := &pbcm.VersionInfoResponse{
		Service: "etop.affiliate",
		Version: "0.1",
	}
	return result, nil
}
