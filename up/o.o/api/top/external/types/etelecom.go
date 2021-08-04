package types

import (
	"time"

	"o.o/common/jsonx"
	"o.o/common/xerrors"
)

type GetExtensionInfoRequest struct {
	ExtensionNumber string `json:"extension_number"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
}

func (m *GetExtensionInfoRequest) String() string { return jsonx.MustMarshalToString(m) }

func (r *GetExtensionInfoRequest) Validate() error {
	if r.ExtensionNumber == "" && r.Email == "" && r.Phone == "" {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "Missing required params")
	}
	return nil
}

type ExtensionInfo struct {
	ExtensionNumber   string    `json:"extension_number"`
	ExtensionPassword string    `json:"extension_password"`
	TenantDomain      string    `json:"tenant_domain"`
	ExpiresAt         time.Time `json:"expires_at"`
}

func (m *ExtensionInfo) String() string { return jsonx.MustMarshalToString(m) }
