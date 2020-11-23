package contact

import (
	"time"

	"o.o/capi/dot"
)

type Contact struct {
	ID          dot.ID
	ShopID      dot.ID
	FullName    string
	Phone       string
	PhoneNorm   string
	WLPartnerID dot.ID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
