package pm

import (
	"o.o/api/main/address"
)

type ProcessManager struct {
	addressQuery address.QueryService
}

func NewProcessManager(addressQuery address.QueryService) *ProcessManager {
	return &ProcessManager{
		addressQuery: addressQuery,
	}
}
