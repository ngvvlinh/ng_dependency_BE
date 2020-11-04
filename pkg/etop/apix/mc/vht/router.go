package vht

import "o.o/capi/httprpc"

type Servers []httprpc.Server

func NewServers(
	userService *VHTUserService,
) Servers {
	servers := httprpc.MustNewServers(
		userService.Clone,
	)
	return servers
}
