package api

import (
	"context"

	"o.o/api/main/authorization"
	"o.o/api/main/invitation"
	apietop "o.o/api/top/int/etop"
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
	"o.o/backend/pkg/etop/sqlstore"
	"o.o/capi/dot"
)

type AccountRelationshipService struct {
	InvitationAggr    invitation.CommandBus
	InvitationQuery   invitation.QueryBus
	AuthorizationAggr authorization.CommandBus
}

func (s *AccountRelationshipService) Clone() *AccountRelationshipService {
	res := *s
	return &res
}

func (s *AccountRelationshipService) CreateInvitation(ctx context.Context, q *AccountRelationshipCreateInvitationEndpoint) error {
	if q.Email == "" && q.Phone == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "email and phone must not be null")
	}

	var roles []authorization.Role
	for _, role := range q.Roles {
		roles = append(roles, authorization.Role(role))
	}
	cmd := &invitation.CreateInvitationCommand{
		AccountID: q.Context.Shop.ID,
		Email:     q.Email,
		Phone:     q.Phone,
		FullName:  q.FullName,
		ShortName: q.ShortName,
		Position:  q.Position,
		Roles:     roles,
		Status:    status3.Z,
		InvitedBy: q.Context.UserID,
	}
	if err := s.InvitationAggr.Dispatch(ctx, cmd); err != nil {
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
	if err := s.InvitationQuery.Dispatch(ctx, query); err != nil {
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
	if err := s.InvitationAggr.Dispatch(ctx, cmd); err != nil {
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
	if err := s.AuthorizationAggr.Dispatch(ctx, cmd); err != nil {
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
	if err := s.AuthorizationAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = convertpb.PbRelationship(cmd.Result)
	return nil
}

func (s *AccountRelationshipService) GetRelationships(ctx context.Context, q *AccountRelationshipGetRelationshipsEndpoint) error {
	paging := cmapi.CMPaging(q.Paging)
	query := &identitymodelx.GetAccountUserExtendedsQuery{
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
	mapUser := make(map[dot.ID]*identitymodel.User)
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
		relationship.Phone = mapUser[relationship.UserID].Phone
	}

	return nil
}

func (s *AccountRelationshipService) RemoveUser(ctx context.Context, q *AccountRelationshipRemoveUserEndpoint) error {
	cmd := &authorization.RemoveUserCommand{
		AccountID:     q.Context.Shop.ID,
		CurrentUserID: q.Context.UserID,
		UserID:        q.UserID,
	}
	if err := s.AuthorizationAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.UpdatedResponse{Updated: cmd.Result}
	return nil
}
