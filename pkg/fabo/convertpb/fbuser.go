package convertpb

import (
	"etop.vn/api/fabo/fbusering"
	"etop.vn/api/top/int/fabo"
)

func PbFbUserCombined(m *fbusering.FbUserCombined) *fabo.FbUserCombined {
	if m == nil {
		return nil
	}
	return &fabo.FbUserCombined{
		ID:         m.FbUser.ID,
		ExternalID: m.FbUser.ExternalID,
		UserID:     m.FbUser.UserID,
		ExternalInfo: &fabo.ExternalFbUserInfo{
			Name:      m.FbUser.ExternalInfo.Name,
			FirstName: m.FbUser.ExternalInfo.FirstName,
			LastName:  m.FbUser.ExternalInfo.LastName,
			ShortName: m.FbUser.ExternalInfo.ShortName,
			ImageURL:  m.FbUser.ExternalInfo.ImageURL,
		},
		Status:    m.FbUser.Status,
		CreatedAt: m.FbUser.CreatedAt,
		UpdatedAt: m.FbUser.UpdatedAt,
	}
}

func PbFbUserCombineds(ms []*fbusering.FbUserCombined) []*fabo.FbUserCombined {
	res := make([]*fabo.FbUserCombined, len(ms))
	for i, m := range ms {
		res[i] = PbFbUserCombined(m)
	}
	return res
}

func PbFbUser(m *fbusering.FbUser) *fabo.FbUser {
	if m == nil {
		return nil
	}
	return &fabo.FbUser{
		ID:         m.ID,
		ExternalID: m.ExternalID,
		UserID:     m.UserID,
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

func PbFbUsers(ms []*fbusering.FbUser) []*fabo.FbUser {
	res := make([]*fabo.FbUser, len(ms))
	for i, m := range ms {
		res[i] = PbFbUser(m)
	}
	return res
}
