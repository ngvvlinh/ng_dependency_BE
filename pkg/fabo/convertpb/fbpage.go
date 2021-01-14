package convertpb

import (
	"o.o/api/fabo/fbpaging"
	"o.o/api/top/int/fabo"
)

func PbFbPageCombined(m *fbpaging.FbExternalPageCombined) *fabo.FbPageCombined {
	if m == nil || m.FbExternalPage == nil {
		return nil
	}
	externalCategoryList := make([]*fabo.ExternalCategory, 0, len(m.FbExternalPage.ExternalCategoryList))
	for _, category := range m.FbExternalPage.ExternalCategoryList {
		externalCategoryList = append(externalCategoryList, &fabo.ExternalCategory{
			ID:   category.ID,
			Name: category.Name,
		})
	}

	return &fabo.FbPageCombined{
		ID:                   m.FbExternalPage.ID,
		ExternalID:           m.FbExternalPage.ExternalID,
		ExternalUserID:       m.FbExternalPage.ExternalUserID,
		ShopID:               m.FbExternalPage.ShopID,
		ExternalName:         m.FbExternalPage.ExternalName,
		ExternalCategory:     m.FbExternalPage.ExternalCategory,
		ExternalCategoryList: externalCategoryList,
		ExternalTasks:        m.FbExternalPage.ExternalTasks,
		ExternalPermissions:  m.FbExternalPage.ExternalPermissions,
		ExternalImageURL:     m.FbExternalPage.ExternalImageURL,
		Status:               m.FbExternalPage.Status,
		ConnectionStatus:     m.FbExternalPage.ConnectionStatus,
		CreatedAt:            m.FbExternalPage.CreatedAt,
		UpdatedAt:            m.FbExternalPage.UpdatedAt,
	}
}

func PbFbPageCombineds(ms []*fbpaging.FbExternalPageCombined) []*fabo.FbPageCombined {
	res := make([]*fabo.FbPageCombined, len(ms))
	for i, m := range ms {
		res[i] = PbFbPageCombined(m)
	}
	return res
}

func PbFbPage(m *fbpaging.FbExternalPage) *fabo.FbPage {
	categoryList := make([]*fabo.ExternalCategory, 0, len(m.ExternalCategoryList))
	for _, category := range m.ExternalCategoryList {
		categoryList = append(categoryList, &fabo.ExternalCategory{
			ID:   category.ID,
			Name: category.Name,
		})
	}
	return &fabo.FbPage{
		ID:                   m.ID,
		ExternalID:           m.ExternalID,
		ShopID:               m.ShopID,
		ExternalName:         m.ExternalName,
		ExternalCategory:     m.ExternalCategory,
		ExternalCategoryList: categoryList,
		ExternalTasks:        m.ExternalTasks,
		ExternalPermissions:  m.ExternalPermissions,
		ExternalImageURL:     m.ExternalImageURL,
		Status:               m.Status,
		ConnectionStatus:     m.ConnectionStatus,
		CreatedAt:            m.CreatedAt,
		UpdatedAt:            m.UpdatedAt,
	}
}

func PbFbPages(ms []*fbpaging.FbExternalPage) []*fabo.FbPage {
	res := make([]*fabo.FbPage, len(ms))
	for i, m := range ms {
		res[i] = PbFbPage(m)
	}
	return res
}
