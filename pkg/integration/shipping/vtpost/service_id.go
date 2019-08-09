package vtpost

import (
	"math/rand"

	vtpostclient2 "etop.vn/backend/pkg/integration/shipping/vtpost/client"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/gencode"
	"etop.vn/backend/pkg/etop/model"
)

type serviceIDGenerator struct {
	rd *rand.Rand
}

func newServiceIDGenerator(seed int64) serviceIDGenerator {
	src := rand.NewSource(seed)
	rd := rand.New(src)
	return serviceIDGenerator{rd}
}

// GenerateServiceID generate new service id for using with ghtk. The generated
// id is always 8 character in length.
func (c serviceIDGenerator) GenerateServiceID(clientCode byte, orderService vtpostclient2.VTPostOrderServiceCode) (string, error) {
	n := c.rd.Uint64()
	v := gencode.Alphabet32.EncodeReverse(n, 8)
	v = v[:5]

	switch clientCode {
	case 'D':
		v[1] = 'D'
		v[2] = blacklist(v[2], 'D')
	default:
		return "", cm.Errorf(cm.Internal, nil, "invalid code")
	}

	switch orderService.Name() {
	case model.ShippingServiceNameStandard:
		v[3] = 'S'
		v[4] = blacklist(v[4], 'S', 'F')
	case model.ShippingServiceNameFaster:
		v[3] = blacklist(v[3], 'S', 'F')
		v[4] = 'F'
	default:
		return "", cm.Errorf(cm.Internal, nil, "invalid code")
	}

	// Get 3 last characters of service
	_orderService := orderService[len(orderService)-3:]
	code := string(v) + string(_orderService)
	return code, nil
}

func DecodeShippingServiceName(code string) (name string, ok bool) {
	if len(code) != 8 {
		return "", false
	}
	switch {
	case code[3] == 'S':
		return model.ShippingServiceNameStandard, true
	case code[4] == 'F':
		return model.ShippingServiceNameFaster, true
	}
	return "", false
}

func (c *Carrier) ParseServiceCode(code string) (serviceName string, ok bool) {
	return DecodeShippingServiceName(code)
}

func blacklist(current byte, blacks ...byte) byte {
	for _, b := range blacks {
		if current == b {
			// return an arbitrary character which does not collide with blacklist values
			return 'V'
		}
	}
	return current
}

func getLast3Character(code vtpostclient2.VTPostOrderServiceCode) string {
	return string(code[len(code)-3:])
}

func ParseServiceID(id string) (clientCode byte, orderService vtpostclient2.VTPostOrderServiceCode, err error) {
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
	case getLast3Character(vtpostclient2.OrderServiceCodeSCOD):
		orderService = vtpostclient2.OrderServiceCodeSCOD

	case string(vtpostclient2.OrderServiceCodeVCN),
		string(vtpostclient2.OrderServiceCodeVTK),
		string(vtpostclient2.OrderServiceCodePHS),
		string(vtpostclient2.OrderServiceCodeVVT),
		string(vtpostclient2.OrderServiceCodeVHT),
		string(vtpostclient2.OrderServiceCodePTN),
		string(vtpostclient2.OrderServiceCodePHT),
		string(vtpostclient2.OrderServiceCodeVBS),
		string(vtpostclient2.OrderServiceCodeVBE):
		orderService = vtpostclient2.VTPostOrderServiceCode(code)

	default:
		// Backwark compatible: the old service ids has the following format:
		// id[4] = D => vtpostclient.OrderServiceCodeSCOD
		// ...
		if id[4] == 'D' {
			orderService = vtpostclient2.OrderServiceCodeSCOD
		}
		if id[5] == 'N' {
			if orderService != "" {
				err = cm.Errorf(cm.InvalidArgument, nil, "invalid service id")
			}
			orderService = vtpostclient2.OrderServiceCodeVCN
		}
		if id[6] == 'K' {
			if orderService != "" {
				err = cm.Errorf(cm.InvalidArgument, nil, "invalid service id")
			}
			orderService = vtpostclient2.OrderServiceCodeVTK
		}
		if id[7] == 'S' {
			if orderService != "" {
				err = cm.Errorf(cm.InvalidArgument, nil, "invalid service id")
			}
			orderService = vtpostclient2.OrderServiceCodePHS
		}
	}

	if clientCode == 0 || orderService == "" {
		err = cm.Errorf(cm.InvalidArgument, nil, "invalid service id")
	}
	return
}
