package pm

import (
	"etop.vn/api/main/address"
)

type ProcessManager struct {
	addressQuery address.QueryService
}

func NewProcessManager(addressQuery address.QueryService) *ProcessManager {
	return &ProcessManager{
		addressQuery: addressQuery,
	}
}
