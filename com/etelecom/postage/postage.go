package postage

import (
	"math"

	"o.o/api/etelecom/call_direction"
	"o.o/api/etelecom/mobile_network"
)

type PriceList struct {
	Direction    call_direction.CallDirection
	Network      mobile_network.MobileNetwork
	FeePerSecond float64
}

var PriceListCall = []PriceList{
	{
		Direction:    call_direction.Out,
		Network:      mobile_network.MobiFone,
		FeePerSecond: 7.5,
	},
	{
		Direction:    call_direction.Out,
		Network:      mobile_network.Viettel,
		FeePerSecond: 7.5,
	},
	{
		Direction:    call_direction.Out,
		Network:      mobile_network.Vinaphone,
		FeePerSecond: 7.5,
	},
	{
		Direction:    call_direction.Out,
		Network:      mobile_network.Other,
		FeePerSecond: 13.4,
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
	Phone          string
	Direction      call_direction.CallDirection
	DurationSecond int
}

func CalcPostage(args CalcPostageArgs) int {
	network := GetPhoneNetwork(args.Phone)
	pricelists := PriceListCallByNetWork[network]
	for _, pl := range pricelists {
		if pl.Network == network && pl.Direction == args.Direction {
			fee := float64(args.DurationSecond) * pl.FeePerSecond
			return int(math.Ceil(fee))
		}
	}
	return 0
}
