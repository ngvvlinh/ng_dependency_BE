package convertpb

import (
	"o.o/api/fabo/fbusering"
	"o.o/api/top/int/fabo"
)

func PbFbUserCombined(m *fbusering.FbExternalUserCombined) *fabo.FbUserCombined {
	if m == nil {
		return nil
	}
	return &fabo.FbUserCombined{
		ExternalID: m.FbExternalUser.ExternalID,
		ExternalInfo: &fabo.ExternalFbUserInfo{
			Name:      m.FbExternalUser.ExternalInfo.Name,
			FirstName: m.FbExternalUser.ExternalInfo.FirstName,
			LastName:  m.FbExternalUser.ExternalInfo.LastName,
			ShortName: m.FbExternalUser.ExternalInfo.ShortName,
			ImageURL:  m.FbExternalUser.ExternalInfo.ImageURL,
		},
		Status:    m.FbExternalUser.Status,
		CreatedAt: m.FbExternalUser.CreatedAt,
		UpdatedAt: m.FbExternalUser.UpdatedAt,
	}
}

func PbFbUserCombineds(ms []*fbusering.FbExternalUserCombined) []*fabo.FbUserCombined {
	res := make([]*fabo.FbUserCombined, len(ms))
	for i, m := range ms {
		res[i] = PbFbUserCombined(m)
	}
	return res
}

func PbFbUser(m *fbusering.FbExternalUser) *fabo.FbUser {
	if m == nil {
		return nil
	}
	return &fabo.FbUser{
		ExternalID: m.ExternalID,
		ExternalInfo: &fabo.ExternalFbUserInfo{
			Name:      m.ExternalInfo.Name,
			FirstName: m.ExternalInfo.FirstName,
			LastName:  m.ExternalInfo.LastName,
			ShortName: m.ExternalInfo.ShortName,
			ImageURL:  m.ExternalInfo.ImageURL,
		},
		Status:    m.Status,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func PbFbUsers(ms []*fbusering.FbExternalUser) []*fabo.FbUser {
	res := make([]*fabo.FbUser, len(ms))
	for i, m := range ms {
		res[i] = PbFbUser(m)
	}
	return res
}
