package api

import (
	"context"

	"etop.vn/api/main/invitation"
	"etop.vn/api/main/location"
	servicelocation "etop.vn/backend/com/main/location"
	pbcm "etop.vn/backend/pb/common"
	pbetop "etop.vn/backend/pb/etop"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/backend/pkg/integration/bank"
	"etop.vn/common/l"
)

func init() {
	bus.AddHandlers("api",
		miscService.VersionInfo,
		locationService.GetProvinces,
		locationService.GetDistricts,
		locationService.GetDistrictsByProvince,
		locationService.GetWards,
		locationService.GetWardsByDistrict,
		locationService.ParseLocation,
		bankService.GetBanks,
		bankService.GetProvincesByBank,
		bankService.GetBranchesByBankProvince,
		addressService.CreateAddress,
		addressService.GetAddresses,
		addressService.UpdateAddress,
		addressService.RemoveAddress,
		invitationService.AcceptInvitation,
		invitationService.RejectInvitation,
		invitationService.GetInvitationByToken,
		invitationService.GetInvitations,
	)
}

var ll = l.New()
var locationBus = servicelocation.New().MessageBus()

type MiscService struct{}
type LocationService struct{}
type BankService struct{}
type AddressService struct{}
type InvitationService struct{}

var miscService = &MiscService{}
var locationService = &LocationService{}
var bankService = &BankService{}
var addressService = &AddressService{}
var invitationService = &InvitationService{}

func (s *MiscService) VersionInfo(ctx context.Context, q *VersionInfoEndpoint) error {
	q.Result = &pbcm.VersionInfoResponse{
		Service: "etop",
		Version: "0.1",
	}
	return nil
}

func (s *LocationService) GetProvinces(ctx context.Context, q *GetProvincesEndpoint) error {
	query := &location.GetAllLocationsQuery{All: true}
	if err := locationBus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &pbetop.GetProvincesResponse{
		Provinces: pbetop.PbProvinces(query.Result.Provinces),
	}
	return nil
}

func (s *LocationService) GetDistricts(ctx context.Context, q *GetDistrictsEndpoint) error {
	query := &location.GetAllLocationsQuery{All: true}
	if err := locationBus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &pbetop.GetDistrictsResponse{
		Districts: pbetop.PbDistricts(query.Result.Districts),
	}
	return nil
}

func (s *LocationService) GetDistrictsByProvince(ctx context.Context, q *GetDistrictsByProvinceEndpoint) error {
	query := &location.GetAllLocationsQuery{ProvinceCode: q.ProvinceCode}
	if err := locationBus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &pbetop.GetDistrictsResponse{
		Districts: pbetop.PbDistricts(query.Result.Districts),
	}
	return nil
}

func (s *LocationService) GetWards(ctx context.Context, q *GetWardsEndpoint) error {
	query := &location.GetAllLocationsQuery{All: true}
	if err := locationBus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &pbetop.GetWardsResponse{
		Wards: pbetop.PbWards(query.Result.Wards),
	}
	return nil
}

func (s *LocationService) GetWardsByDistrict(ctx context.Context, q *GetWardsByDistrictEndpoint) error {
	query := &location.GetAllLocationsQuery{DistrictCode: q.DistrictCode}
	if err := locationBus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &pbetop.GetWardsResponse{
		Wards: pbetop.PbWards(query.Result.Wards),
	}
	return nil
}

func (s *LocationService) ParseLocation(ctx context.Context, q *ParseLocationEndpoint) error {
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

func (s *BankService) GetBanks(ctx context.Context, q *GetBanksEndpoint) error {
	q.Result = &pbetop.GetBanksResponse{
		Banks: pbetop.PbBanks(bank.Banks),
	}
	return nil
}

func (s *BankService) GetProvincesByBank(ctx context.Context, q *GetProvincesByBankEndpoint) error {
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
	q.Result = &pbetop.GetBranchesByBankProvinceResponse{
		Branches: pbetop.PbBankBranches(branches),
	}
	return nil
}

func (s *AddressService) CreateAddress(ctx context.Context, q *CreateAddressEndpoint) error {
	address, err := pbetop.PbCreateAddressToModel(q.Context.AccountID, q.CreateAddressRequest)
	if err != nil {
		return err
	}
	cmd := &model.CreateAddressCommand{
		Address: address,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = pbetop.PbAddress(cmd.Result)
	return nil
}

func (s *AddressService) GetAddresses(ctx context.Context, q *GetAddressesEndpoint) error {
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

func (s *AddressService) UpdateAddress(ctx context.Context, q *UpdateAddressEndpoint) error {
	accountID := q.Context.AccountID
	address, err := pbetop.PbUpdateAddressToModel(accountID, q.UpdateAddressRequest)
	if err != nil {
		return err
	}
	cmd := &model.UpdateAddressCommand{
		Address: address,
	}

	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = pbetop.PbAddress(cmd.Result)
	return nil
}

func (s *AddressService) RemoveAddress(ctx context.Context, q *RemoveAddressEndpoint) error {
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

func (s *InvitationService) AcceptInvitation(ctx context.Context, q *AcceptInvitationEndpoint) error {
	cmd := &invitation.AcceptInvitationCommand{
		UserID: q.Context.UserID,
		Token:  q.Token,
	}
	if err := invitationAggregate.Dispatch(ctx, cmd); err != nil {
		return err
	}

	q.Result = &pbcm.UpdatedResponse{Updated: int32(cmd.Result)}
	return nil
}

func (s *InvitationService) RejectInvitation(ctx context.Context, q *RejectInvitationEndpoint) error {
	cmd := &invitation.AcceptInvitationCommand{
		UserID: q.Context.UserID,
		Token:  q.Token,
	}
	if err := invitationAggregate.Dispatch(ctx, cmd); err != nil {
		return err
	}

	q.Result = &pbcm.UpdatedResponse{Updated: int32(cmd.Result)}
	return nil
}

func (s *InvitationService) GetInvitationByToken(ctx context.Context, q *GetInvitationByTokenEndpoint) error {
	query := &invitation.GetInvitationByTokenQuery{
		Token: q.Token,
	}
	if err := invitationQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = pbetop.PbInvitation(query.Result)
	return nil
}

func (s *InvitationService) GetInvitations(ctx context.Context, q *GetInvitationsEndpoint) error {
	paging := q.Paging.CMPaging()
	query := &invitation.ListInvitationsByEmailQuery{
		Email:   q.Context.User.Email,
		Paging:  *paging,
		Filters: pbcm.ToFilters(q.Filters),
	}
	if err := invitationQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &pbetop.InvitationsResponse{
		Invitations: pbetop.PbInvitations(query.Result.Invitations),
		Paging:      pbcm.PbPageInfo(paging, query.Result.Count),
	}
	return nil
}
