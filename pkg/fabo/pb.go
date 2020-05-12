package fabo

import (
	"time"

	"o.o/api/fabo/fbusering"
	"o.o/api/top/int/fabo"
)

func PbFbUser(args *fbusering.FbExternalUser) *fabo.FbUser {
	return &fabo.FbUser{
		ExternalID:   "",
		ExternalInfo: nil,
		Status:       0,
		CreatedAt:    time.Time{},
		UpdatedAt:    time.Time{},
	}
}
