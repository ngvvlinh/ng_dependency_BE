package provider

import "o.o/capi/dot"

type CreateExtensionResponse struct {
	ExtensionID       dot.ID
	ExtensionNumber   string
	ExtensionPassword string
	ExternalID        string
	HotlineID         dot.ID
}
