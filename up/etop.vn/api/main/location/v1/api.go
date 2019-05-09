package v1

func (g VietnamRegion) Name() string {
	switch g {
	case VietnamRegion_north:
		return "Miền Bắc"
	case VietnamRegion_middle:
		return "Miền Trung"
	case VietnamRegion_south:
		return "Miền Nam"
	default:
		return "?"
	}
}

func (a UrbanType) Name() string {
	switch a {
	case UrbanType_urban:
		return "Nội thành"
	case UrbanType_suburban1:
		return "Ngoại thành 1"
	case UrbanType_suburban2:
		return "Ngoại thành 2"
	default:
		return "?"
	}
}
