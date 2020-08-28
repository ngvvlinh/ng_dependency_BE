package etop

import (
	"context"

	ghnclient "o.o/backend/pkg/integration/shipping/ghn/clientv2"
	"o.o/backend/pkg/integration/shipping/ghn/driverv2"
)

var _ driverv2.SupportedGHNDriver = &EtopSupportedGHNDriver{}

type EtopSupportedGHNDriver struct {
	client *ghnclient.Client
}

func NewEtopSupportedGHNDriver(env string, cfg ghnclient.GHNAccountCfg) *EtopSupportedGHNDriver {
	return &EtopSupportedGHNDriver{client: ghnclient.New(env, cfg)}
}

func (f *EtopSupportedGHNDriver) AddClientContract(ctx context.Context, clientID int) error {
	return nil
}
