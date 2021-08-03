package postage

import (
	"o.o/api/etelecom/mobile_network"
	"o.o/backend/pkg/common/validate"
)

var MobileNetworkList = map[mobile_network.MobileNetwork][]string{
	mobile_network.MobiFone: {
		"089", "090", "093", "070", "079", "077",
		"076", "078",
	},
	mobile_network.Vinaphone: {
		"088", "091", "094", "083", "084", "085",
		"081", "082",
	},
	mobile_network.Viettel: {
		"086", "096", "097", "098", "032", "033", "034",
		"035", "036", "037", "038", "039",
	},
}

var PrefixPhoneNumberNetWork = map[string]mobile_network.MobileNetwork{}

func init() {
	for network, phoneLists := range MobileNetworkList {
		for _, prefixPhone := range phoneLists {
			PrefixPhoneNumberNetWork[prefixPhone] = network
		}
	}
}

func GetPhoneNetwork(phone string) mobile_network.MobileNetwork {
	phoneNorm, ok := validate.NormalizePhone(phone)
	if !ok {
		return mobile_network.Other
	}
	prefixPhoneNumber := string(phoneNorm)[:3]
	network, ok := PrefixPhoneNumberNetWork[prefixPhoneNumber]
	if !ok {
		return mobile_network.Other
	}
	return network
}
