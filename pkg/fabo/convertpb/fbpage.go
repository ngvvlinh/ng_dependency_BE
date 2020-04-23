package convertpb

import (
	"etop.vn/api/fabo/fbpaging"
	"etop.vn/api/top/int/fabo"
)

func PbFbPageCombined(m *fbpaging.FbPageCombined) *fabo.FbPageCombined {
	if m == nil || m.FbPage == nil {
		return nil
	}
	externalCategoryList := make([]*fabo.ExternalCategory, 0, len(m.FbPage.ExternalCategoryList))
	for _, category := range m.FbPage.ExternalCategoryList {
		externalCategoryList = append(externalCategoryList, &fabo.ExternalCategory{
			ID:   category.ID,
			Name: category.Name,
		})
	}

	return &fabo.FbPageCombined{
		ID:                   m.FbPage.ID,
		ExternalID:           m.FbPage.ExternalID,
		FbUserID:             m.FbPage.FbUserID,
		ShopID:               m.FbPage.ShopID,
		UserID:               m.FbPage.UserID,
		ExternalName:         m.FbPage.ExternalName,
		ExternalCategory:     m.FbPage.ExternalCategory,
		ExternalCategoryList: externalCategoryList,
		ExternalTasks:        m.FbPage.ExternalTasks,
		ExternalImageURL:     m.FbPage.ExternalImageURL,
		Status:               m.FbPage.Status,
		ConnectionStatus:     m.FbPage.ConnectionStatus,
		CreatedAt:            m.FbPage.CreatedAt,
		UpdatedAt:            m.FbPage.UpdatedAt,
	}
}

func PbFbPageCombineds(ms []*fbpaging.FbPageCombined) []*fabo.FbPageCombined {
	res := make([]*fabo.FbPageCombined, len(ms))
	for i, m := range ms {
		res[i] = PbFbPageCombined(m)
	}
	return res
}

func PbFbPage(m *fbpaging.FbPage) *fabo.FbPage {
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
		FbUserID:             m.FbUserID,
		ShopID:               m.ShopID,
		UserID:               m.UserID,
		ExternalName:         m.ExternalName,
		ExternalCategory:     m.ExternalCategory,
		ExternalCategoryList: categoryList,
		ExternalTasks:        m.ExternalTasks,
		ExternalImageURL:     m.ExternalImageURL,
		Status:               m.Status,
		ConnectionStatus:     m.ConnectionStatus,
		CreatedAt:            m.CreatedAt,
		UpdatedAt:            m.UpdatedAt,
	}
}

func PbFbPages(ms []*fbpaging.FbPage) []*fabo.FbPage {
	res := make([]*fabo.FbPage, len(ms))
	for i, m := range ms {
		res[i] = PbFbPage(m)
	}
	return res
}
