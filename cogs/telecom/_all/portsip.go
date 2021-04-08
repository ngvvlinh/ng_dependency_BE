package _all

import (
	telecomtypes "o.o/backend/com/etelecom/provider/types"
	portsipclient "o.o/backend/pkg/integration/telecom/portsip/client"
	portsipadmindriver "o.o/backend/pkg/integration/telecom/portsip/driver/admin"
)

type AdminPortsipConfig struct {
	Name                      string `yaml:"name"`
	Password                  string `yaml:"password"`
	AarenetTrunkingProviderID string `yaml:"aarenet_trunking_provider_id"`
}

func SupportAdminPortsipDriver(cfg AdminPortsipConfig) telecomtypes.AdministratorTelecom {
	accountCfg := portsipclient.PortsipAccountCfg{
		Username: cfg.Name,
		Password: cfg.Password,
	}
	driver := portsipadmindriver.NewAdminAccount(accountCfg)
	return telecomtypes.AdministratorTelecom{
		Driver:                 driver,
		TrunkProviderDefaultID: cfg.AarenetTrunkingProviderID,
	}
}
