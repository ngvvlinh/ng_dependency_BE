package authorization

import (
	"context"

	"o.o/api/meta"
	"o.o/capi/dot"
)

// +gen:api

type Aggregate interface {
	UpdatePermission(ctx context.Context, _ *UpdatePermissionArgs) (*Relationship, error)
	UpdateRelationship(ctx context.Context, _ *UpdateRelationshipArgs) (*Relationship, error)
	LeaveAccount(ctx context.Context, userID, accountID dot.ID) (updated int, _ error)
	RemoveUser(ctx context.Context, _ *RemoveUserArgs) (update int, _ error)
}

type QueryService interface {
	GetAuthorization(ctx context.Context, accountID, userID dot.ID) (*Authorization, error)
	GetAccountAuthorization(ctx context.Context, accountID dot.ID) ([]*Authorization, error)
	GetRelationships(ctx context.Context, _ *GetRelationshipsArgs) ([]*Relationship, error)
}

//-- command --//
type UpdatePermissionArgs struct {
	AccountID  dot.ID
	CurrUserID dot.ID
	UserID     dot.ID
	Roles      []Role
}

type UpdateRelationshipArgs struct {
	AccountID dot.ID
	UserID    dot.ID
	FullName  dot.NullString
	ShortName dot.NullString
	Position  dot.NullString
}

type RemoveUserArgs struct {
	AccountID     dot.ID
	CurrentUserID dot.ID
	UserID        dot.ID
	Roles         []Role
}

type GetRelationshipsArgs struct {
	AccountID dot.ID
	Paging    meta.Paging
	Filters   meta.Filters
}
