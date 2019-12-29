package api

import (
	"context"

	"etop.vn/api/main/authorization"
	"etop.vn/api/main/invitation"
	"etop.vn/api/main/location"
	apietop "etop.vn/api/top/int/etop"
	pbcm "etop.vn/api/top/types/common"
	"etop.vn/api/top/types/etc/status3"
	authorizationconvert "etop.vn/backend/com/main/authorization/convert"
	"etop.vn/backend/com/main/invitation/convert"
	servicelocation "etop.vn/backend/com/main/location"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/apifw/cmapi"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/etop/api/convertpb"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/backend/pkg/etop/sqlstore"
	"etop.vn/backend/pkg/integration/bank"
	"etop.vn/capi/dot"
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
		userRelationshipService.AcceptInvitation,
		userRelationshipService.RejectInvitation,
		userRelationshipService.GetInvitationByToken,
		userRelationshipService.GetInvitations,
		accountRelationshipService.UpdatePermission,
	)
}

var ll = l.New()
var locationBus = servicelocation.New().MessageBus()

type MiscService struct{}
type LocationService struct{}
type BankService struct{}
type AddressService struct{}
type UserRelationshipService struct{}
type AccountRelationshipService struct{}

var miscService = &MiscService{}
var locationService = &LocationService{}
var bankService = &BankService{}
var addressService = &AddressService{}
var userRelationshipService = &UserRelationshipService{}
var accountRelationshipService = &AccountRelationshipService{}

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
	q.Result = &apietop.GetProvincesResponse{
		Provinces: convertpb.PbProvinces(query.Result.Provinces),
	}
	return nil
}

func (s *LocationService) GetDistricts(ctx context.Context, q *GetDistrictsEndpoint) error {
	query := &location.GetAllLocationsQuery{All: true}
	if err := locationBus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &apietop.GetDistrictsResponse{
		Districts: convertpb.PbDistricts(query.Result.Districts),
	}
	return nil
}

func (s *LocationService) GetDistrictsByProvince(ctx context.Context, q *GetDistrictsByProvinceEndpoint) error {
	query := &location.GetAllLocationsQuery{ProvinceCode: q.ProvinceCode}
	if err := locationBus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &apietop.GetDistrictsResponse{
		Districts: convertpb.PbDistricts(query.Result.Districts),
	}
	return nil
}

func (s *LocationService) GetWards(ctx context.Context, q *GetWardsEndpoint) error {
	query := &location.GetAllLocationsQuery{All: true}
	if err := locationBus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &apietop.GetWardsResponse{
		Wards: convertpb.PbWards(query.Result.Wards),
	}
	return nil
}

func (s *LocationService) GetWardsByDistrict(ctx context.Context, q *GetWardsByDistrictEndpoint) error {
	query := &location.GetAllLocationsQuery{DistrictCode: q.DistrictCode}
	if err := locationBus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &apietop.GetWardsResponse{
		Wards: convertpb.PbWards(query.Result.Wards),
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
	res := &apietop.ParseLocationResponse{}
	if loc.Province != nil {
		res.Province = convertpb.PbProvince(loc.Province)
	}
	if loc.District != nil {
		res.District = convertpb.PbDistrict(loc.District)
	}
	if loc.Ward != nil {
		res.Ward = convertpb.PbWard(loc.Ward)
	}
	q.Result = res
	return nil
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

func (s *AddressService) CreateAddress(ctx context.Context, q *CreateAddressEndpoint) error {
	address, err := convertpb.PbCreateAddressToModel(q.Context.AccountID, q.CreateAddressRequest)
	if err != nil {
		return err
	}
	cmd := &model.CreateAddressCommand{
		Address: address,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = convertpb.PbAddress(cmd.Result)
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
	q.Result = &apietop.GetAddressResponse{
		Addresses: convertpb.PbAddresses(query.Result.Addresses),
	}
	return nil
}

func (s *AddressService) UpdateAddress(ctx context.Context, q *UpdateAddressEndpoint) error {
	accountID := q.Context.AccountID
	address, err := convertpb.PbUpdateAddressToModel(accountID, q.UpdateAddressRequest)
	if err != nil {
		return err
	}
	cmd := &model.UpdateAddressCommand{
		Address: address,
	}

	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = convertpb.PbAddress(cmd.Result)
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

func (s *UserRelationshipService) AcceptInvitation(ctx context.Context, q *UserRelationshipAcceptInvitationEndpoint) error {
	cmd := &invitation.AcceptInvitationCommand{
		UserID: q.Context.UserID,
		Token:  q.Token,
	}
	if err := invitationAggregate.Dispatch(ctx, cmd); err != nil {
		return err
	}

	q.Result = &pbcm.UpdatedResponse{Updated: cmd.Result}
	return nil
}

func (s *UserRelationshipService) RejectInvitation(ctx context.Context, q *UserRelationshipRejectInvitationEndpoint) error {
	cmd := &invitation.RejectInvitationCommand{
		UserID: q.Context.UserID,
		Token:  q.Token,
	}
	if err := invitationAggregate.Dispatch(ctx, cmd); err != nil {
		return err
	}

	q.Result = &pbcm.UpdatedResponse{Updated: cmd.Result}
	return nil
}

func (s *UserRelationshipService) GetInvitationByToken(ctx context.Context, q *UserRelationshipGetInvitationByTokenEndpoint) error {
	query := &invitation.GetInvitationByTokenQuery{
		Token: q.Token,
	}
	if err := invitationQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = convertpb.PbInvitation(query.Result)

	getAccountQuery := &model.GetShopQuery{
		ShopID: query.Result.AccountID,
	}
	if err := bus.Dispatch(ctx, getAccountQuery); err != nil {
		return err
	}
	q.Result.ShopShort = &apietop.ShopShort{
		ID:       getAccountQuery.Result.ID,
		Name:     getAccountQuery.Result.Name,
		Code:     getAccountQuery.Result.Code,
		ImageUrl: getAccountQuery.Result.ImageURL,
	}

	getUserQuery := &model.GetUserByEmailOrPhoneQuery{
		Email: query.Result.Email,
	}
	err := bus.Dispatch(ctx, getUserQuery)
	switch cm.ErrorCode(err) {
	case cm.NotFound:
	// no-op
	case cm.NoError:
		q.Result.UserId = getUserQuery.Result.ID
	default:
		return err
	}

	getInvitedByUserQuery := &model.GetUserByIDQuery{
		UserID: query.Result.InvitedBy,
	}
	if err := bus.Dispatch(ctx, getInvitedByUserQuery); err != nil {
		return err
	}
	q.Result.InvitedByUser = getInvitedByUserQuery.Result.FullName

	return nil
}

func (s *UserRelationshipService) GetInvitations(ctx context.Context, q *UserRelationshipGetInvitationsEndpoint) error {
	paging := cmapi.CMPaging(q.Paging)
	query := &invitation.ListInvitationsByEmailQuery{
		Email:   q.Context.User.Email,
		Paging:  *paging,
		Filters: cmapi.ToFilters(q.Filters),
	}
	if err := invitationQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &apietop.InvitationsResponse{
		Invitations: convertpb.PbInvitations(query.Result.Invitations),
		Paging:      cmapi.PbPageInfo(paging),
	}

	var accountIDs []dot.ID
	var hasAccountID bool
	for _, invitationEl := range query.Result.Invitations {
		hasAccountID = false
		for _, accountID := range accountIDs {
			if accountID == invitationEl.AccountID {
				hasAccountID = true
			}
		}
		if !hasAccountID {
			accountIDs = append(accountIDs, invitationEl.AccountID)
		}
	}

	getAccountsQuery := &model.GetShopsQuery{
		ShopIDs: accountIDs,
	}
	if err := bus.Dispatch(ctx, getAccountsQuery); err != nil {
		return err
	}
	mapShop := make(map[dot.ID]*model.Shop)
	for _, shop := range getAccountsQuery.Result.Shops {
		mapShop[shop.ID] = shop
	}

	for _, invitationEl := range q.Result.Invitations {
		invitationEl.ShopShort = &apietop.ShopShort{
			ID:       invitationEl.ShopId,
			Name:     mapShop[invitationEl.ShopId].Name,
			Code:     mapShop[invitationEl.ShopId].Code,
			ImageUrl: mapShop[invitationEl.ShopId].ImageURL,
		}
	}

	return nil
}

func (s *UserRelationshipService) LeaveAccount(ctx context.Context, q *UserRelationshipLeaveAccountEndpoint) error {
	cmd := &authorization.LeaveAccountCommand{
		UserID:    q.Context.UserID,
		AccountID: q.AccountID,
	}
	if err := authorizationAggregate.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.UpdatedResponse{Updated: cmd.Result}
	return nil
}

func (s *AccountRelationshipService) CreateInvitation(ctx context.Context, q *AccountRelationshipCreateInvitationEndpoint) error {
	var roles []authorization.Role
	for _, role := range q.Roles {
		roles = append(roles, authorization.Role(role))
	}
	cmd := &invitation.CreateInvitationCommand{
		AccountID: q.Context.Shop.ID,
		Email:     q.Email,
		FullName:  q.FullName,
		ShortName: q.ShortName,
		Position:  q.Position,
		Roles:     roles,
		Status:    status3.Z,
		InvitedBy: q.Context.UserID,
	}
	if err := invitationAggregate.Dispatch(ctx, cmd); err != nil {
		return err
	}

	q.Result = convertpb.PbInvitation(cmd.Result)
	return nil
}

func (s *AccountRelationshipService) GetInvitations(ctx context.Context, q *AccountRelationshipGetInvitationsEndpoint) error {
	paging := cmapi.CMPaging(q.Paging)
	query := &invitation.ListInvitationsQuery{
		ShopID:  q.Context.Shop.ID,
		Paging:  *paging,
		Filters: cmapi.ToFilters(q.Filters),
	}
	if err := invitationQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &apietop.InvitationsResponse{
		Invitations: convertpb.PbInvitations(query.Result.Invitations),
		Paging:      cmapi.PbPageInfo(paging),
	}
	return nil
}

func (s *AccountRelationshipService) DeleteInvitation(ctx context.Context, q *AccountRelationshipDeleteInvitationEndpoint) error {
	cmd := &invitation.DeleteInvitationCommand{
		UserID:    q.Context.UserID,
		AccountID: q.Context.Shop.ID,
		Token:     q.Token,
	}
	if err := invitationAggregate.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.UpdatedResponse{Updated: cmd.Result}
	return nil
}

func (s *AccountRelationshipService) UpdatePermission(ctx context.Context, q *AccountRelationshipUpdatePermissionEndpoint) error {
	cmd := &authorization.UpdatePermissionCommand{
		AccountID:  q.Context.Shop.ID,
		CurrUserID: q.Context.UserID,
		UserID:     q.UserID,
		Roles:      convert.ConvertStringsToRoles(q.Roles),
	}
	if err := authorizationAggregate.Dispatch(ctx, cmd); err != nil {
		return err
	}

	q.Result = convertpb.PbRelationship(cmd.Result)
	return nil
}

func (s *AccountRelationshipService) UpdateRelationship(ctx context.Context, q *AccountRelationshipUpdateRelationshipEndpoint) error {
	cmd := &authorization.UpdateRelationshipCommand{
		AccountID: q.Context.Shop.ID,
		UserID:    q.UserID,
		FullName:  q.FullName,
		ShortName: q.ShortName,
		Position:  q.Position,
	}
	if err := authorizationAggregate.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = convertpb.PbRelationship(cmd.Result)
	return nil
}

func (s *AccountRelationshipService) GetRelationships(ctx context.Context, q *AccountRelationshipGetRelationshipsEndpoint) error {
	paging := cmapi.CMPaging(q.Paging)
	query := &model.GetAccountUserExtendedsQuery{
		AccountIDs:     []dot.ID{q.Context.Shop.ID},
		Paging:         paging,
		Filters:        cmapi.ToFilters(q.Filters),
		IncludeDeleted: true,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}

	var relationships []*authorization.Relationship
	for _, accountUser := range query.Result.AccountUsers {
		relationships = append(relationships, authorizationconvert.ConvertAccountUserToRelationship(accountUser.AccountUser))
	}

	q.Result = &apietop.RelationshipsResponse{Relationships: convertpb.PbRelationships(relationships)}

	var userIDs []dot.ID
	mapUser := make(map[dot.ID]*model.User)
	for _, relationship := range q.Result.Relationships {
		userIDs = append(userIDs, relationship.UserID)
	}

	users, err := sqlstore.User(ctx).IDs(userIDs...).List()
	if err != nil {
		return err
	}

	for _, user := range users {
		mapUser[user.ID] = user
	}
	for _, relationship := range q.Result.Relationships {
		if relationship.FullName == "" {
			relationship.FullName = mapUser[relationship.UserID].FullName
		}

		relationship.Email = mapUser[relationship.UserID].Email
	}

	return nil
}

func (s *AccountRelationshipService) RemoveUser(ctx context.Context, q *AccountRelationshipRemoveUserEndpoint) error {
	cmd := &authorization.RemoveUserCommand{
		AccountID:     q.Context.Shop.ID,
		CurrentUserID: q.Context.UserID,
		UserID:        q.UserID,
	}
	if err := authorizationAggregate.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.UpdatedResponse{Updated: cmd.Result}
	return nil
}
