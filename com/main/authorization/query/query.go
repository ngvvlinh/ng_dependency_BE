package query

import (
	"context"

	"o.o/api/main/authorization"
	"o.o/backend/com/main/authorization/convert"
	identitymodelx "o.o/backend/com/main/identity/modelx"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/etop/authorize/auth"
	"o.o/backend/pkg/etop/sqlstore"
	"o.o/capi/dot"
)

var _ authorization.QueryService = &AuthorizationQuery{}

type AuthorizationQuery struct {
	Auth *auth.Authorizer

	AccountUserStore sqlstore.AccountUserStoreInterface
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
	if err := a.AccountUserStore.GetAccountUserExtended(ctx, getAccountUserQuery); err != nil {
		return nil, cm.MapError(err).
			Wrap(cm.NotFound, "Authorization not found").
			Throw()
	}
	auth = convert.ConvertAccountUserExtendedToAuthorization(a.Auth, &getAccountUserQuery.Result)
	return auth, nil
}

func (a *AuthorizationQuery) GetAccountAuthorization(
	ctx context.Context, accountID dot.ID,
) (auths []*authorization.Authorization, _ error) {
	getAccountUsersQuery := &identitymodelx.GetAccountUserExtendedsQuery{
		AccountIDs: []dot.ID{accountID},
	}
	if err := a.AccountUserStore.GetAccountUserExtendeds(ctx, getAccountUsersQuery); err != nil {
		return nil, err
	}
	for _, accountUser := range getAccountUsersQuery.Result.AccountUsers {
		auths = append(auths, convert.ConvertAccountUserExtendedToAuthorization(a.Auth, accountUser))
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
	if err := a.AccountUserStore.GetAccountUserExtendeds(ctx, getAccountUsersQuery); err != nil {
		return nil, err
	}
	for _, accountUser := range getAccountUsersQuery.Result.AccountUsers {
		relationships = append(relationships, convert.ConvertAccountUserToRelationship(a.Auth, accountUser.AccountUser))
	}
	return relationships, nil
}
