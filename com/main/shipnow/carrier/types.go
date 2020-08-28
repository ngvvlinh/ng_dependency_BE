package carrier

import "o.o/capi/dot"

type ShopConnectionSignInArgs struct {
	Phone        string
	ConnectionID dot.ID
	OwnerID      dot.ID
	ShopID       dot.ID
}
