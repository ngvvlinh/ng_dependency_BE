package invitation

import (
	"context"
	"time"

	"o.o/api/main/authorization"
	"o.o/api/meta"
	"o.o/api/shopping"
	"o.o/api/top/types/etc/status3"
	"o.o/capi/dot"
)

// +gen:api

type Aggregate interface {
	CreateInvitation(ctx context.Context, _ *CreateInvitationArgs) (*Invitation, error)
	ResendInvitation(ctx context.Context, _ *ResendInvitationArgs) (*Invitation, error)
	AcceptInvitation(ctx context.Context, userID dot.ID, token string) (updated int, _ error)
	RejectInvitation(ctx context.Context, userID dot.ID, token string) (updated int, _ error)
	DeleteInvitation(ctx context.Context, userID, accountID dot.ID, token string) (updated int, _ error)
}

type QueryService interface {
	GetInvitation(ctx context.Context, ID dot.ID) (*Invitation, error)
	GetInvitationByToken(ctx context.Context, token string) (*Invitation, error)

	ListInvitations(context.Context, *shopping.ListQueryShopArgs) (*InvitationsResponse, error)
	ListInvitationsByEmailAndPhone(context.Context, *ListInvitationsByEmailAndPhoneArgs) (*InvitationsResponse, error)
	ListInvitationsAcceptedByEmail(ctx context.Context, email string) (*InvitationsResponse, error)
}

//-- queries --//

type InvitationsResponse struct {
	Invitations []*Invitation
	Paging      meta.PageInfo
}

type ListInvitationsByEmailAndPhoneArgs struct {
	Email   string
	Phone   string
	Paging  meta.Paging
	Filters meta.Filters
}

//-- commands --//

// +convert:create=Invitation
type CreateInvitationArgs struct {
	AccountID dot.ID
	Email     string
	Phone     string
	FullName  string
	ShortName string
	Position  string
	Roles     []authorization.Role
	Status    status3.Status
	InvitedBy dot.ID
	CreatedBy time.Time
}

type ResendInvitationArgs struct {
	AccountID dot.ID
	ResendBy  dot.ID
	Email     string
	Phone     string
}
