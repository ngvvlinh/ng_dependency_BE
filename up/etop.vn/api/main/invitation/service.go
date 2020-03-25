package invitation

import (
	"context"
	"time"

	"etop.vn/api/main/authorization"
	"etop.vn/api/meta"
	"etop.vn/api/shopping"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/capi/dot"
)

// +gen:api

type Aggregate interface {
	CreateInvitation(ctx context.Context, _ *CreateInvitationArgs) (*Invitation, error)
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
