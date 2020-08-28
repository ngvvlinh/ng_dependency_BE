package driverv2

import (
	"context"
)

type SupportedGHNDriver interface {
	// update setting of shipment price list for merchant
	AddClientContract(ctx context.Context, clientID int) error
}
