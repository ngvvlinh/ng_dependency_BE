package _main

import (
	"o.o/backend/pkg/common/apifw/httpx"
	"o.o/capi/httprpc"
)

type IntHandlers []httprpc.Server

type ExtHandlers []httprpc.Server

type AuthxHandler httpx.Server

type OIDCHandler httpx.Server

type PortSipHandler httpx.Server
