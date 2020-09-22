package service

import (
	"o.o/api/top/external/mc/vnp"
)

func clone(dst *vnp.ShipnowService, src vnp.ShipnowService) {
	*dst = src.(interface{ Clone() vnp.ShipnowService }).Clone()
}
