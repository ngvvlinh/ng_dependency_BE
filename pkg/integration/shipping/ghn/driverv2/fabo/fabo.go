package fabo

import (
	"context"

	ghnclient "o.o/backend/pkg/integration/shipping/ghn/clientv2"
	"o.o/backend/pkg/integration/shipping/ghn/driverv2"
)

var _ driverv2.SupportedGHNDriver = &FaboSupportedGHNDriver{}

type FaboSupportedGHNDriver struct {
	client *ghnclient.Client
}

func NewFaboSupportedGHNDriver(env string, cfg ghnclient.GHNAccountCfg) *FaboSupportedGHNDriver {
	return &FaboSupportedGHNDriver{client: ghnclient.New(env, cfg)}
}

func (f *FaboSupportedGHNDriver) AddClientContract(ctx context.Context, clientID int) error {
	return f.client.AddClientContract(ctx, &ghnclient.AddClientContractRequest{ClientID: clientID})
}
