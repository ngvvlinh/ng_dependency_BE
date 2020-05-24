package query

import (
	"context"

	"o.o/api/main/authorization"
	"o.o/backend/com/main/authorization/convert"
	identitymodelx "o.o/backend/com/main/identity/modelx"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/capi/dot"
)

var _ authorization.QueryService = &AuthorizationQuery{}

type AuthorizationQuery struct{}

func NewAuthorizationQuery() *AuthorizationQuery {
	return &AuthorizationQuery{}
}

func AuthorizationQueryMessageBus(a *AuthorizationQuery) authorization.QueryBus {
	b := bus.New()
	return authorization.NewQueryServiceHandler(a).RegisterHandlers(b)
}

func (a *AuthorizationQuery) GetAuthorization(
	ctx context.Context, accountID, userID dot.ID,
) (auth *authorization.Authorization, _ error) {
	getAccountUserQuery := &identitymodelx.GetAccountUserExtendedQuery{
		UserID:    userID,
		AccountID: accountID,
	}
	if err := bus.Dispatch(ctx, getAccountUserQuery); err != nil {
		return nil, cm.MapError(err).
			Wrap(cm.NotFound, "Authorization not found").
			Throw()
	}
	auth = convert.ConvertAccountUserExtendedToAuthorization(&getAccountUserQuery.Result)
	return auth, nil
}

func (a *AuthorizationQuery) GetAccountAuthorization(
	ctx context.Context, accountID dot.ID,
) (auths []*authorization.Authorization, _ error) {
	getAccountUsersQuery := &identitymodelx.GetAccountUserExtendedsQuery{
		AccountIDs: []dot.ID{accountID},
	}
	if err := bus.Dispatch(ctx, getAccountUsersQuery); err != nil {
		return nil, err
	}
	for _, accountUser := range getAccountUsersQuery.Result.AccountUsers {
		auths = append(auths, convert.ConvertAccountUserExtendedToAuthorization(accountUser))
	}
	return auths, nil
}

func (a *AuthorizationQuery) GetRelationships(
	ctx context.Context, args *authorization.GetRelationshipsArgs,
) (relationships []*authorization.Relationship, _ error) {
	getAccountUsersQuery := &identitymodelx.GetAccountUserExtendedsQuery{
		AccountIDs: []dot.ID{args.AccountID},
		Filters:    args.Filters,
	}
	if err := bus.Dispatch(ctx, getAccountUsersQuery); err != nil {
		return nil, err
	}
	for _, accountUser := range getAccountUsersQuery.Result.AccountUsers {
		relationships = append(relationships, convert.ConvertAccountUserToRelationship(accountUser.AccountUser))
	}
	return relationships, nil
}
