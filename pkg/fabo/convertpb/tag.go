package convertpb

import (
	"o.o/api/fabo/fbusering"
	"o.o/api/top/int/fabo"
)

func ConvertFbUserringTagToResponseTag(tag *fbusering.FbShopTag) *fabo.FbShopTagResponse {
	if tag == nil {
		return nil
	}

	return &fabo.FbShopTagResponse{
		ID:        tag.ID,
		Name:      tag.Name,
		Color:     tag.Color,
		ShopID:    tag.ShopID,
		CreatedAt: tag.CreatedAt,
		UpdatedAt: tag.UpdatedAt,
	}
}

func ConvertFbUserringTagsToResponseTags(tags []*fbusering.FbShopTag) []*fabo.FbShopTagResponse {
	if tags == nil {
		return nil
	}

	var result []*fabo.FbShopTagResponse
	for _, tag := range tags {
		result = append(result, ConvertFbUserringTagToResponseTag(tag))
	}

	return result
}
