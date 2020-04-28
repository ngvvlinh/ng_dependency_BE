package services

import (
	"o.o/api/top/types/etc/shipping_provider"
	ghnclient "o.o/backend/pkg/integration/shipping/ghn/client"
	ghtkclient "o.o/backend/pkg/integration/shipping/ghtk/client"
	vtpostclient "o.o/backend/pkg/integration/shipping/vtpost/client"
)

type ShipmentService struct {
	ServiceID string
	Name      string
}

var shipmentServicesByCarrier = map[shipping_provider.ShippingProvider][]ShipmentService{
	shipping_provider.GHN: {
		{
			ServiceID: string(ghnclient.ServiceFee6Hours),
			Name:      "Gói 6 giờ",
		}, {
			ServiceID: string(ghnclient.ServiceFee1Day),
			Name:      "Gói 1 ngày",
		}, {
			ServiceID: string(ghnclient.ServiceFee2Days),
			Name:      "Gói 2 ngày",
		}, {
			ServiceID: string(ghnclient.ServiceFee3Days),
			Name:      "Gói 3 ngày",
		}, {
			ServiceID: string(ghnclient.ServiceFee4Days),
			Name:      "Gói 4 ngày",
		}, {
			ServiceID: string(ghnclient.ServiceFee5Days),
			Name:      "Gói 5 ngày",
		}, {
			ServiceID: string(ghnclient.ServiceFee6Days),
			Name:      "Gói 6 ngày",
		},
	},
	shipping_provider.GHTK: {
		{
			ServiceID: string(ghtkclient.TransportRoad),
			Name:      "Đường bộ",
		}, {
			ServiceID: string(ghtkclient.TransportFly),
			Name:      "Đường hàng không",
		},
	},
	shipping_provider.VTPost: {
		{
			ServiceID: string(vtpostclient.OrderServiceCodeSCOD),
			Name:      "Nhanh - SCOD Giao hàng thu tiền",
		}, {
			ServiceID: string(vtpostclient.OrderServiceCodeVCN),
			Name:      "Nhanh - VCN Chuyển phát nhanh - Express dilivery",
		}, {
			ServiceID: string(vtpostclient.OrderServiceCodeVTK),
			Name:      "Chậm - VTK - VTK Tiết kiệm - Express Saver",
		}, {
			ServiceID: string(vtpostclient.OrderServiceCodePHS),
			Name:      "Chậm - PHS Phát hôm sau nội tỉnh",
		}, {
			ServiceID: string(vtpostclient.OrderServiceCodeVVT),
			Name:      "Chậm - VVT Dịch vụ vận tải",
		}, {
			ServiceID: string(vtpostclient.OrderServiceCodeVHT),
			Name:      "Nhanh - VHT Phát Hỏa tốc",
		}, {
			ServiceID: string(vtpostclient.OrderServiceCodePTN),
			Name:      "Nhanh - PTN Phát trong ngày nội tỉnh",
		}, {
			ServiceID: string(vtpostclient.OrderServiceCodePHT),
			Name:      "Nhanh - PHT Phát hỏa tốc nội tỉnh",
		}, {
			ServiceID: string(vtpostclient.OrderServiceCodeVBS),
			Name:      "Nhanh - VBS Nhanh theo hộp",
		}, {
			ServiceID: string(vtpostclient.OrderServiceCodeVBE),
			Name:      "Chậm - VBE Tiết kiệm theo hộp",
		},
	},
}

func GetServicesByCarrier(carrier shipping_provider.ShippingProvider) []ShipmentService {
	res, ok := shipmentServicesByCarrier[carrier]
	if !ok {
		return []ShipmentService{}
	}
	return res
}
