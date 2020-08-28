package accountshipnow

import (
	"context"
	"encoding/json"
	"time"

	"o.o/capi/dot"
)

// +gen:api

type Aggregate interface {
	CreateExternalAccountAhamove(context.Context, *CreateExternalAccountAhamoveArgs) (*ExternalAccountAhamove, error)

	RequestVerifyExternalAccountAhamove(context.Context, *RequestVerifyExternalAccountAhamoveArgs) (*RequestVerifyExternalAccountAhamoveResult, error)

	UpdateVerifiedExternalAccountAhamove(context.Context, *UpdateVerifiedExternalAccountAhamoveArgs) (*ExternalAccountAhamove, error)

	UpdateExternalAccountAhamoveVerification(context.Context, *UpdateExternalAccountAhamoveVerificationArgs) (*ExternalAccountAhamove, error)

	UpdateExternalAccountAhamoveExternalInfo(context.Context, *UpdateXAccountAhamoveExternalInfoArgs) error

	DeleteExternalAccountAhamove(context.Context, *DeleteXAccountAhamoveArgs) error
}

type QueryService interface {
	GetExternalAccountAhamove(context.Context, *GetExternalAccountAhamoveArgs) (*ExternalAccountAhamove, error)

	GetExternalAccountAhamoveByExternalID(context.Context, *GetExternalAccountAhamoveByExternalIDQueryArgs) (*ExternalAccountAhamove, error)

	GetAccountShipnow(context.Context, *GetAccountShipnowArgs) (*ExternalAccountAhamove, error)
}

type CreateExternalAccountAhamoveArgs struct {
	ShopID       dot.ID
	OwnerID      dot.ID // user id
	Phone        string
	Name         string
	ConnectionID dot.ID
}

type RequestVerifyExternalAccountAhamoveArgs struct {
	OwnerID      dot.ID
	Phone        string
	ConnectionID dot.ID
}

type RequestVerifyExternalAccountAhamoveResult struct {
}

type UpdateVerifiedExternalAccountAhamoveArgs struct {
	OwnerID      dot.ID
	Phone        string
	ConnectionID dot.ID
}

type UpdateExternalAccountAhamoveVerificationArgs struct {
	OwnerID             dot.ID
	Phone               string
	IDCardFrontImg      string
	IDCardBackImg       string
	PortraitImg         string
	WebsiteURL          string
	FanpageURL          string
	CompanyImgs         []string
	BusinessLicenseImgs []string
}

type UpdateXAccountAhamoveExternalInfoArgs struct {
	ID                   dot.ID
	ExternalID           string
	ExternalCreatedAt    time.Time
	ExternalVerified     bool
	ExternalToken        string
	ExternalTicketID     string
	LastSendVerifiedAt   time.Time
	ExternalDataVerified json.RawMessage
}

type DeleteXAccountAhamoveArgs struct {
	ID           dot.ID
	OwnerID      dot.ID
	ConnectionID dot.ID
}

type GetExternalAccountAhamoveArgs struct {
	OwnerID      dot.ID
	Phone        string
	ConnectionID dot.ID
}

type GetExternalAccountAhamoveByExternalIDQueryArgs struct {
	ExternalID string
}

type GetAccountShipnowArgs struct {
	Phone        string
	OwnerID      dot.ID
	ConnectionID dot.ID
}
