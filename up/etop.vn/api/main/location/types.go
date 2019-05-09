package location

import locationv1 "etop.vn/api/main/location/v1"

const CountryVietnam = "Việt Nam"

type VietnamRegion = locationv1.VietnamRegion

const (
	North  VietnamRegion = locationv1.VietnamRegion_north
	Middle VietnamRegion = locationv1.VietnamRegion_middle
	South  VietnamRegion = locationv1.VietnamRegion_south
)

type UrbanType = locationv1.UrbanType

const (
	Urban     UrbanType = locationv1.UrbanType_urban
	Suburban1 UrbanType = locationv1.UrbanType_suburban1
	Suburban2 UrbanType = locationv1.UrbanType_suburban2
)

type Province = locationv1.Province

type District = locationv1.District

type Ward = locationv1.Ward

type LocationExtra = locationv1.Extra
