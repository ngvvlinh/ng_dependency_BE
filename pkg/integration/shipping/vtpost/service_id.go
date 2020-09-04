package vtpost

import (
	cm "o.o/backend/pkg/common"
	vtpostclient "o.o/backend/pkg/integration/shipping/vtpost/client"
)

func getLast3Character(code vtpostclient.VTPostOrderServiceCode) string {
	return string(code[len(code)-3:])
}

func ParseServiceID(id string) (clientCode byte, orderService vtpostclient.VTPostOrderServiceCode, err error) {
	if id == "" {
		err = cm.Errorf(cm.InvalidArgument, nil, "missing service id")
		return
	}
	if len(id) != 8 {
		err = cm.Errorf(cm.InvalidArgument, nil, "invalid service id")
		return
	}
	if id[1] == 'D' {
		clientCode = VTPostCodePublic
	}

	code := id[len(id)-3:]
	switch code {
	case getLast3Character(vtpostclient.OrderServiceCodeSCOD):
		orderService = vtpostclient.OrderServiceCodeSCOD

	case vtpostclient.OrderServiceCodeVCN.String(),
		vtpostclient.OrderServiceCodeVTK.String(),
		vtpostclient.OrderServiceCodePHS.String(),
		vtpostclient.OrderServiceCodeVVT.String(),
		vtpostclient.OrderServiceCodeVHT.String(),
		vtpostclient.OrderServiceCodePTN.String(),
		vtpostclient.OrderServiceCodePHT.String(),
		vtpostclient.OrderServiceCodeVBS.String(),
		vtpostclient.OrderServiceCodeVBE.String():
		orderService = vtpostclient.VTPostOrderServiceCode(code)

	default:
		// Backwark compatible: the old service ids has the following format:
		// id[4] = D => vtpostclient.OrderServiceCodeSCOD
		// ...
		if id[4] == 'D' {
			orderService = vtpostclient.OrderServiceCodeSCOD
		}
		if id[5] == 'N' {
			if orderService != "" {
				err = cm.Errorf(cm.InvalidArgument, nil, "invalid service id")
			}
			orderService = vtpostclient.OrderServiceCodeVCN
		}
		if id[6] == 'K' {
			if orderService != "" {
				err = cm.Errorf(cm.InvalidArgument, nil, "invalid service id")
			}
			orderService = vtpostclient.OrderServiceCodeVTK
		}
		if id[7] == 'S' {
			if orderService != "" {
				err = cm.Errorf(cm.InvalidArgument, nil, "invalid service id")
			}
			orderService = vtpostclient.OrderServiceCodePHS
		}
	}

	if clientCode == 0 || orderService == "" {
		err = cm.Errorf(cm.InvalidArgument, nil, "invalid service id")
	}
	return
}
