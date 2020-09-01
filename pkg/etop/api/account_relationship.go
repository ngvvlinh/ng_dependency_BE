package api

import (
	"context"
	"fmt"
	"time"

	"o.o/api/main/authorization"
	"o.o/api/main/invitation"
	api "o.o/api/top/int/etop"
	pbcm "o.o/api/top/types/common"
	"o.o/api/top/types/etc/status3"
	authorizationconvert "o.o/backend/com/main/authorization/convert"
	identitymodel "o.o/backend/com/main/identity/model"
	identitymodelx "o.o/backend/com/main/identity/modelx"
	"o.o/backend/com/main/invitation/convert"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/backend/pkg/etop/sqlstore"
	"o.o/capi/dot"
)

type AccountRelationshipService struct {
	session.Session

	InvitationAggr    invitation.CommandBus
	InvitationQuery   invitation.QueryBus
	AuthorizationAggr authorization.CommandBus
}

func (s *AccountRelationshipService) Clone() api.AccountRelationshipService {
	res := *s
	return &res
}

func (s *AccountRelationshipService) ResendInvitation(ctx context.Context, q *api.ResendInvitationRequest) (*api.Invitation, error) {
	key := fmt.Sprintf("resend-invitation:%v-%v-%v-%v",
		s.SS.Shop().ID, s.SS.Claim().UserID, q.Email, q.Phone)
	resp, _, err := idempgroup.DoAndWrap(
		ctx, key, 10*time.Minute, "Resend invitation",
		func() (interface{}, error) { return s.resendInvitation(ctx, q) })
	if err != nil {
		return nil, err
	}
	result := convertpb.PbInvitation(resp.(*invitation.Invitation))
	return result, nil
}

func (s *AccountRelationshipService) resendInvitation(ctx context.Context, q *api.ResendInvitationRequest) (*invitation.Invitation, error) {
	if q.Email == "" && q.Phone == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "email and phone must not be null")
	}

	cmd := &invitation.ResendInvitationCommand{
		AccountID: s.SS.Shop().ID,
		ResendBy:  s.SS.Claim().UserID,
		Email:     q.Email,
		Phone:     q.Phone,
	}
	if err := s.InvitationAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}

	return cmd.Result, nil
}

func (s *AccountRelationshipService) CreateInvitation(ctx context.Context, q *api.CreateInvitationRequest) (*api.Invitation, error) {
	if q.Email == "" && q.Phone == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "email and phone must not be null")
	}

	var roles []authorization.Role
	for _, role := range q.Roles {
		roles = append(roles, authorization.Role(role))
	}
	cmd := &invitation.CreateInvitationCommand{
		AccountID: s.SS.Shop().ID,
		Email:     q.Email,
		Phone:     q.Phone,
		FullName:  q.FullName,
		ShortName: q.ShortName,
		Position:  q.Position,
		Roles:     roles,
		Status:    status3.Z,
		InvitedBy: s.SS.Claim().UserID,
	}
	if err := s.InvitationAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}

	result := convertpb.PbInvitation(cmd.Result)
	return result, nil
}

func (s *AccountRelationshipService) GetInvitations(ctx context.Context, q *api.GetInvitationsRequest) (*api.InvitationsResponse, error) {
	paging := cmapi.CMPaging(q.Paging)
	query := &invitation.ListInvitationsQuery{
		ShopID:  s.SS.Shop().ID,
		Paging:  *paging,
		Filters: cmapi.ToFilters(q.Filters),
	}
	if err := s.InvitationQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := &api.InvitationsResponse{
		Invitations: convertpb.PbInvitations(query.Result.Invitations),
		Paging:      cmapi.PbPageInfo(paging),
	}
	return result, nil
}

func (s *AccountRelationshipService) DeleteInvitation(ctx context.Context, q *api.DeleteInvitationRequest) (*pbcm.UpdatedResponse, error) {
	cmd := &invitation.DeleteInvitationCommand{
		UserID:    s.SS.Claim().UserID,
		AccountID: s.SS.Shop().ID,
		Token:     q.Token,
	}
	if err := s.InvitationAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.UpdatedResponse{Updated: cmd.Result}
	return result, nil
}

func (s *AccountRelationshipService) UpdatePermission(ctx context.Context, q *api.UpdateAccountUserPermissionRequest) (*api.Relationship, error) {
	cmd := &authorization.UpdatePermissionCommand{
		AccountID:  s.SS.Shop().ID,
		CurrUserID: s.SS.Claim().UserID,
		UserID:     q.UserID,
		Roles:      convert.ConvertStringsToRoles(q.Roles),
	}
	if err := s.AuthorizationAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}

	result := convertpb.PbRelationship(cmd.Result)
	return result, nil
}

func (s *AccountRelationshipService) UpdateRelationship(ctx context.Context, q *api.UpdateRelationshipRequest) (*api.Relationship, error) {
	cmd := &authorization.UpdateRelationshipCommand{
		AccountID: s.SS.Shop().ID,
		UserID:    q.UserID,
		FullName:  q.FullName,
		ShortName: q.ShortName,
		Position:  q.Position,
	}
	if err := s.AuthorizationAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := convertpb.PbRelationship(cmd.Result)
	return result, nil
}

func (s *AccountRelationshipService) GetRelationships(ctx context.Context, q *api.GetRelationshipsRequest) (*api.RelationshipsResponse, error) {
	paging := cmapi.CMPaging(q.Paging)
	query := &identitymodelx.GetAccountUserExtendedsQuery{
		AccountIDs:     []dot.ID{s.SS.Shop().ID},
		Paging:         paging,
		Filters:        cmapi.ToFilters(q.Filters),
		IncludeDeleted: true,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return nil, err
	}

	var relationships []*authorization.Relationship
	for _, accountUser := range query.Result.AccountUsers {
		relationships = append(relationships, authorizationconvert.ConvertAccountUserToRelationship(s.SS.Authorizer(), accountUser.AccountUser))
	}

	result := &api.RelationshipsResponse{Relationships: convertpb.PbRelationships(relationships)}

	var userIDs []dot.ID
	mapUser := make(map[dot.ID]*identitymodel.User)
	for _, relationship := range result.Relationships {
		userIDs = append(userIDs, relationship.UserID)
	}

	users, err := sqlstore.User(ctx).IDs(userIDs...).List()
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		mapUser[user.ID] = user
	}
	for _, relationship := range result.Relationships {
		if relationship.FullName == "" {
			relationship.FullName = mapUser[relationship.UserID].FullName
		}

		relationship.Email = mapUser[relationship.UserID].Email
		relationship.Phone = mapUser[relationship.UserID].Phone
	}
	return result, nil
}

func (s *AccountRelationshipService) RemoveUser(ctx context.Context, q *api.RemoveUserRequest) (*pbcm.UpdatedResponse, error) {
	cmd := &authorization.RemoveUserCommand{
		AccountID:     s.SS.Shop().ID,
		CurrentUserID: s.SS.Claim().UserID,
		UserID:        q.UserID,
	}
	if err := s.AuthorizationAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.UpdatedResponse{Updated: cmd.Result}
	return result, nil
}
