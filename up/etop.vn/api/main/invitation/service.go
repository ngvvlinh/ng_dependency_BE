package invitation

import (
	"context"
	"time"

	"etop.vn/api/main/etop"
	"etop.vn/api/meta"
	"etop.vn/api/shopping"
)

// +gen:api

type Aggregate interface {
	CreateInvitation(ctx context.Context, _ *CreateInvitationArgs) (*Invitation, error)
	AcceptInvitation(ctx context.Context, userID int64, token string) (updated int, _ error)
	RejectInvitation(ctx context.Context, userID int64, token string) (updated int, _ error)
}

type QueryService interface {
	GetInvitation(ctx context.Context, ID int64) (*Invitation, error)
	GetInvitationByToken(ctx context.Context, token string) (*Invitation, error)

	ListInvitations(context.Context, *shopping.ListQueryShopArgs) (*InvitationsResponse, error)
	ListInvitationsByEmail(context.Context, *ListInvitationsByEmailArgs) (*InvitationsResponse, error)
	ListInvitationsAcceptedByEmail(ctx context.Context, email string) (*InvitationsResponse, error)
}

//-- queries --//

type InvitationsResponse struct {
	Invitations []*Invitation
	Count       int32
	Paging      meta.PageInfo
}

type ListInvitationsByEmailArgs struct {
	Email   string
	Paging  meta.Paging
	Filters meta.Filters
}

//-- commands --//

// +convert:create=Invitation
type CreateInvitationArgs struct {
	AccountID int64
	Email     string
	Roles     []Role
	Status    etop.Status3
	InvitedBy int64
	CreatedBy time.Time
}
