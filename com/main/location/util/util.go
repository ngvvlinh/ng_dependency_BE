package util

import "etop.vn/api/main/location"

func CheckUrbanHCMHN(district *location.District) bool {
	// HCM code: "79"
	// HN code: "01"
	if district.ProvinceCode != "79" && district.ProvinceCode != "01" {
		return false
	}
	if district.UrbanType != location.Urban {
		return false
	}
	return true
}
