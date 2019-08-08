package pm

import (
	"etop.vn/api/main/identity"
)

type ProcessManager struct {
	identityQuery identity.QueryService
}

func NewProcessManager(identityQuery identity.QueryService) *ProcessManager {
	return &ProcessManager{
		identityQuery: identityQuery,
	}
}
