package convertpb

import (
	"o.o/api/fabo/fbusering"
	"o.o/api/top/int/fabo"
)

func Convert_core_FbShopUserTag_To_api_FbShopUserTag(tag *fbusering.FbShopUserTag) *fabo.FbShopUserTag {
	if tag == nil {
		return nil
	}

	return &fabo.FbShopUserTag{
		ID:        tag.ID,
		Name:      tag.Name,
		Color:     tag.Color,
		ShopID:    tag.ShopID,
		CreatedAt: tag.CreatedAt,
		UpdatedAt: tag.UpdatedAt,
	}
}

func Convert_core_FbShopUserTags_To_api_FbShopUserTags(tags []*fbusering.FbShopUserTag) []*fabo.FbShopUserTag {
	if tags == nil {
		return nil
	}

	var result []*fabo.FbShopUserTag
	for _, tag := range tags {
		result = append(result, Convert_core_FbShopUserTag_To_api_FbShopUserTag(tag))
	}
	return result
}
