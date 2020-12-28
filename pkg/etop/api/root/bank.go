package root

import (
	"context"

	api "o.o/api/top/int/etop"
	pbcm "o.o/api/top/types/common"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/backend/pkg/integration/bank"
)

type BankService struct {
	session.Session
}

func (s *BankService) Clone() api.BankService {
	res := *s
	return &res
}

func (s *BankService) GetBanks(ctx context.Context, q *pbcm.Empty) (*api.GetBanksResponse, error) {
	result := &api.GetBanksResponse{
		Banks: convertpb.PbBanks(bank.Banks),
	}
	return result, nil
}

func (s *BankService) GetProvincesByBank(ctx context.Context, q *api.GetProvincesByBankResquest) (*api.GetBankProvincesResponse, error) {
	query := &bank.BankQuery{
		Code: q.BankCode,
		Name: q.BankName,
	}

	provinces := bank.GetProvinceByBank(query)
	result := &api.GetBankProvincesResponse{
		Provinces: convertpb.PbBankProvinces(provinces),
	}
	return result, nil
}

func (s *BankService) GetBankProvinces(ctx context.Context, q *api.GetBankProvincesRequest) (*api.GetBankProvinceResponse, error) {
	query := &bank.BankQuery{
		Code: q.BankCode,
		Name: q.BankName,
	}

	provinces := bank.GetBankProvinces(query, q.All)

	return &api.GetBankProvinceResponse{
		Provinces: convertpb.PbBankProvinces(provinces),
	}, nil
}

func (s *BankService) GetBranchesByBankProvince(ctx context.Context, q *api.GetBranchesByBankProvinceResquest) (*api.GetBranchesByBankProvinceResponse, error) {
	bankQuery := &bank.BankQuery{
		Code: q.BankCode,
		Name: q.BankName,
	}
	provinceQuery := &bank.ProvinceQuery{
		Code: q.ProvinceCode,
		Name: q.ProvinceName,
	}

	branches := bank.GetBranchByBankProvince(bankQuery, provinceQuery)
	result := &api.GetBranchesByBankProvinceResponse{
		Branches: convertpb.PbBankBranches(branches),
	}
	return result, nil
}

func (s *BankService) GetBankBranches(ctx context.Context, q *api.GetBankBranchesRequest) (*api.GetBankBranchesResponse, error) {
	bankQuery := &bank.BankQuery{
		Code: q.BankCode,
		Name: q.BankName,
	}
	branches := bank.GetBankBranches(bankQuery, q.All)

	return &api.GetBankBranchesResponse{
		Branches: convertpb.PbBankBranches(branches),
	}, nil
}
