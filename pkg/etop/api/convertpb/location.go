package convertpb

import (
	"o.o/api/main/location"
	"o.o/api/top/int/admin"
)

func PbCustomRegion(in *location.CustomRegion) *admin.CustomRegion {
	if in == nil {
		return nil
	}
	return &admin.CustomRegion{
		ID:            in.ID,
		Name:          in.Name,
		Description:   in.Description,
		ProvinceCodes: in.ProvinceCodes,
		CreatedAt:     in.CreatedAt,
		UpdatedAt:     in.UpdatedAt,
	}
}

func PbCustomRegions(items []*location.CustomRegion) []*admin.CustomRegion {
	result := make([]*admin.CustomRegion, len(items))
	for i, item := range items {
		result[i] = PbCustomRegion(item)
	}
	return result
}
