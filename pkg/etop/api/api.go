package api

import (
	"context"

	"etop.vn/api/main/location"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/l"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/backend/pkg/integration/bank"
	servicelocation "etop.vn/backend/pkg/services/location"

	pbcm "etop.vn/backend/pb/common"
	pbetop "etop.vn/backend/pb/etop"
	wrapetop "etop.vn/backend/wrapper/etop"
)

func init() {
	bus.AddHandlers("api",
		VersionInfo,
		GetProvinces,
		GetDistricts,
		GetDistrictsByProvince,
		GetWards,
		GetWardsByDistrict,
		ParseLocation,
		GetBanks,
		GetProvincesByBank,
		GetBranchesByBankProvince,
		CreateAddress,
		GetAddresses,
		UpdateAddress,
		RemoveAddress,
	)
}

var ll = l.New()
var locationBus = servicelocation.New().MessageBus()

func VersionInfo(ctx context.Context, q *wrapetop.VersionInfoEndpoint) error {
	q.Result = &pbcm.VersionInfoResponse{
		Service: "etop",
		Version: "0.1",
	}
	return nil
}

func GetProvinces(ctx context.Context, q *wrapetop.GetProvincesEndpoint) error {
	query := &location.GetAllLocationsQuery{All: true}
	if err := locationBus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &pbetop.GetProvincesResponse{
		Provinces: pbetop.PbProvinces(query.Result.Provinces),
	}
	return nil
}

func GetDistricts(ctx context.Context, q *wrapetop.GetDistrictsEndpoint) error {
	query := &location.GetAllLocationsQuery{All: true}
	if err := locationBus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &pbetop.GetDistrictsResponse{
		Districts: pbetop.PbDistricts(query.Result.Districts),
	}
	return nil
}

func GetDistrictsByProvince(ctx context.Context, q *wrapetop.GetDistrictsByProvinceEndpoint) error {
	query := &location.GetAllLocationsQuery{ProvinceCode: q.ProvinceCode}
	if err := locationBus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &pbetop.GetDistrictsResponse{
		Districts: pbetop.PbDistricts(query.Result.Districts),
	}
	return nil
}

func GetWards(ctx context.Context, q *wrapetop.GetWardsEndpoint) error {
	query := &location.GetAllLocationsQuery{All: true}
	if err := locationBus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &pbetop.GetWardsResponse{
		Wards: pbetop.PbWards(query.Result.Wards),
	}
	return nil
}

func GetWardsByDistrict(ctx context.Context, q *wrapetop.GetWardsByDistrictEndpoint) error {
	query := &location.GetAllLocationsQuery{DistrictCode: q.DistrictCode}
	if err := locationBus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &pbetop.GetWardsResponse{
		Wards: pbetop.PbWards(query.Result.Wards),
	}
	return nil
}

func ParseLocation(ctx context.Context, q *wrapetop.ParseLocationEndpoint) error {
	query := &location.FindLocationQuery{
		Province: q.ProvinceName,
		District: q.DistrictName,
		Ward:     q.WardName,
	}
	if err := locationBus.Dispatch(ctx, query); err != nil {
		return err
	}
	loc := query.Result
	res := &pbetop.ParseLocationResponse{}
	if loc.Province != nil {
		res.Province = pbetop.PbProvince(loc.Province)
	}
	if loc.District != nil {
		res.District = pbetop.PbDistrict(loc.District)
	}
	if loc.Ward != nil {
		res.Ward = pbetop.PbWard(loc.Ward)
	}
	q.Result = res
	return nil
}

func GetBanks(ctx context.Context, q *wrapetop.GetBanksEndpoint) error {
	q.Result = &pbetop.GetBanksResponse{
		Banks: pbetop.PbBanks(bank.Banks),
	}
	return nil
}

func GetProvincesByBank(ctx context.Context, q *wrapetop.GetProvincesByBankEndpoint) error {
	query := &bank.BankQuery{
		Code: q.BankCode,
		Name: q.BankName,
	}

	provinces := bank.GetProvinceByBank(query)
	q.Result = &pbetop.GetBankProvincesResponse{
		Provinces: pbetop.PbBankProvinces(provinces),
	}
	return nil
}

func GetBranchesByBankProvince(ctx context.Context, q *wrapetop.GetBranchesByBankProvinceEndpoint) error {
	bankQuery := &bank.BankQuery{
		Code: q.BankCode,
		Name: q.BankName,
	}
	provinceQuery := &bank.ProvinceQuery{
		Code: q.ProvinceCode,
		Name: q.ProvinceName,
	}

	branches := bank.GetBranchByBankProvince(bankQuery, provinceQuery)
	q.Result = &pbetop.GetBranchesByBankProvinceResponse{
		Branches: pbetop.PbBankBranches(branches),
	}
	return nil
}

func CreateAddress(ctx context.Context, q *wrapetop.CreateAddressEndpoint) error {
	cmd := &model.CreateAddressCommand{
		Address: pbetop.PbCreateAddressToModel(q.Context.AccountID, q.CreateAddressRequest),
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = pbetop.PbAddress(cmd.Result)
	return nil
}

func GetAddresses(ctx context.Context, q *wrapetop.GetAddressesEndpoint) error {
	accountID := q.Context.AccountID
	query := &model.GetAddressesQuery{
		AccountID: accountID,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return nil
	}
	q.Result = &pbetop.GetAddressResponse{
		Addresses: pbetop.PbAddresses(query.Result.Addresses),
	}
	return nil
}

func UpdateAddress(ctx context.Context, q *wrapetop.UpdateAddressEndpoint) error {
	accountID := q.Context.AccountID
	cmd := &model.UpdateAddressCommand{
		Address: pbetop.PbUpdateAddressToModel(accountID, q.UpdateAddressRequest),
	}

	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = pbetop.PbAddress(cmd.Result)
	return nil
}

func RemoveAddress(ctx context.Context, q *wrapetop.RemoveAddressEndpoint) error {
	accountID := q.Context.AccountID
	cmd := &model.DeleteAddressCommand{
		ID:        q.Id,
		AccountID: accountID,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.Empty{}
	return nil
}
