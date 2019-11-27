package carrier

import "etop.vn/capi/dot"

type ConnectionSignInArgs struct {
	ConnectionID dot.ID
	Email        string
	Password     string
}

type ShopConnectionSignInArgs struct {
	ConnectionID dot.ID
	ShopID       dot.ID
	Email        string
	Password     string
}

type ConnectionSignUpArgs struct {
	ConnectionID dot.ID
	Name         string
	Email        string
	Password     string
	Phone        string
	Province     string
	District     string
	Address      string
}

type ShopConnectionSignUpArgs struct {
	ConnectionID dot.ID
	ShopID       dot.ID
	Name         string
	Email        string
	Password     string
	Phone        string
	Province     string
	District     string
	Address      string
}
