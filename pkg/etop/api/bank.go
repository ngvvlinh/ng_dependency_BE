package api

import (
	"context"

	apietop "o.o/api/top/int/etop"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/integration/bank"
)

type BankService struct{}

func (s *BankService) Clone() *BankService {
	res := *s
	return &res
}

func (s *BankService) GetBanks(ctx context.Context, q *GetBanksEndpoint) error {
	q.Result = &apietop.GetBanksResponse{
		Banks: convertpb.PbBanks(bank.Banks),
	}
	return nil
}

func (s *BankService) GetProvincesByBank(ctx context.Context, q *GetProvincesByBankEndpoint) error {
	query := &bank.BankQuery{
		Code: q.BankCode,
		Name: q.BankName,
	}

	provinces := bank.GetProvinceByBank(query)
	q.Result = &apietop.GetBankProvincesResponse{
		Provinces: convertpb.PbBankProvinces(provinces),
	}
	return nil
}

func (s *BankService) GetBranchesByBankProvince(ctx context.Context, q *GetBranchesByBankProvinceEndpoint) error {
	bankQuery := &bank.BankQuery{
		Code: q.BankCode,
		Name: q.BankName,
	}
	provinceQuery := &bank.ProvinceQuery{
		Code: q.ProvinceCode,
		Name: q.ProvinceName,
	}

	branches := bank.GetBranchByBankProvince(bankQuery, provinceQuery)
	q.Result = &apietop.GetBranchesByBankProvinceResponse{
		Branches: convertpb.PbBankBranches(branches),
	}
	return nil
}
