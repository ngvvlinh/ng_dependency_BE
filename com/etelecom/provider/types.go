package provider

import "o.o/capi/dot"

type CreateExtensionResponse struct {
	ExtensionID       dot.ID
	ExtensionNumber   string
	ExtensionPassword string
	ExternalID        string
	HotlineID         dot.ID
}

type CreateTenantResponse struct {
	TenantID         dot.ID
	ExternalTenantID string
}

type CreateOutboundRuleRequest struct {
	// trunk provider: aarenat provider id in portsip
	TrunkProviderID string
	OwnerID         dot.ID
	ConnectionID    dot.ID
}
