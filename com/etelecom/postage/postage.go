package postage

import (
	"o.o/api/etelecom/call_log_direction"
	"o.o/api/etelecom/mobile_network"
)

type PriceList struct {
	Direction    call_log_direction.CallLogDirection
	Network      mobile_network.MobileNetwork
	FeePerMinute int
}

var PriceListCall = []PriceList{
	{
		Direction:    call_log_direction.Out,
		Network:      mobile_network.Mobiphone,
		FeePerMinute: 450,
	},
	{
		Direction:    call_log_direction.Out,
		Network:      mobile_network.Viettel,
		FeePerMinute: 450,
	},
	{
		Direction:    call_log_direction.Out,
		Network:      mobile_network.Vinaphone,
		FeePerMinute: 450,
	},
	{
		Direction:    call_log_direction.Out,
		Network:      mobile_network.Other,
		FeePerMinute: 800,
	},
}

var (
	PriceListCallByNetWork = make(map[mobile_network.MobileNetwork][]PriceList)
)

func init() {
	for _, pl := range PriceListCall {
		PriceListCallByNetWork[pl.Network] = append(PriceListCallByNetWork[pl.Network], pl)
	}
}

type CalcPostageArgs struct {
	Phone     string
	Direction call_log_direction.CallLogDirection
	// Duration: minute
	DurationPostage int
}

func CalcPostage(args CalcPostageArgs) int {
	network := GetPhoneNetwork(args.Phone)
	pricelists := PriceListCallByNetWork[network]
	for _, pl := range pricelists {
		if pl.Network == network && pl.Direction == args.Direction {
			return args.DurationPostage * pl.FeePerMinute
		}
	}
	return 0
}
