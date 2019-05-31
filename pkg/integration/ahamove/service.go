package ahamove

import (
	"math/rand"

	"etop.vn/backend/pkg/services/location"

	ahamoveClient "etop.vn/backend/pkg/integration/ahamove/client"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/gencode"
)

type (
	ServiceCode string
	CityCode    string
)
type ShippingService struct {
	ID            ServiceCode
	Name          string
	MinStopPoints int
	MaxStopPoints int
	MaxWeight     int
	City          CityCode
	ShortCode     string
}

var (
	SGNBIKE ServiceCode = "SGN-BIKE"
	SGNPOOL ServiceCode = "SGN-POOL"
	SGNDG   ServiceCode = "SGN-DG"
	HANBIKE ServiceCode = "HAN-BIKE"
	HANPOOL ServiceCode = "HAN-POOL"
	HANDG   ServiceCode = "HAN-DG"

	SGNCode CityCode = "SGN"
	HANCode CityCode = "HAN"

	ServicesIndexID        = make(map[ServiceCode]*ShippingService)
	ServicesIndexShortCode = make(map[string]*ShippingService)
	ServicesIndexCity      = make(map[CityCode][]*ShippingService)
)
var Services = []*ShippingService{
	{
		ID:            SGNBIKE,
		Name:          "Siêu tốc",
		MinStopPoints: 1,
		MaxStopPoints: 5,
		MaxWeight:     0,
		City:          SGNCode,
		ShortCode:     "SB",
	}, {
		ID:            SGNPOOL,
		Name:          "Siêu rẻ",
		MinStopPoints: 1,
		MaxStopPoints: 1,
		MaxWeight:     0,
		City:          SGNCode,
		ShortCode:     "SP",
	}, {
		ID:            SGNDG,
		Name:          "Đồng giá",
		MinStopPoints: 4,
		MaxStopPoints: 10,
		MaxWeight:     0,
		City:          SGNCode,
		ShortCode:     "SD",
	}, {
		ID:            HANBIKE,
		Name:          "Siêu tốc",
		MinStopPoints: 1,
		MaxStopPoints: 5,
		MaxWeight:     0,
		City:          HANCode,
		ShortCode:     "HB",
	}, {
		ID:            HANPOOL,
		Name:          "Siêu rẻ",
		MinStopPoints: 1,
		MaxStopPoints: 1,
		MaxWeight:     0,
		City:          HANCode,
		ShortCode:     "HP",
	}, {
		ID:            HANDG,
		Name:          "Đồng giá",
		MinStopPoints: 3,
		MaxStopPoints: 10,
		MaxWeight:     0,
		City:          HANCode,
		ShortCode:     "HD",
	},
}

type serviceIDGenerator struct {
	rd *rand.Rand
}

func init() {
	for _, service := range Services {
		ServicesIndexID[service.ID] = service
		ServicesIndexShortCode[service.ShortCode] = service
		ServicesIndexCity[service.City] = append(ServicesIndexCity[service.City], service)
	}
}

func newServiceIDGenerator(seed int64) serviceIDGenerator {
	src := rand.NewSource(seed)
	rd := rand.New(src)
	return serviceIDGenerator{rd}
}

// GenerateServiceID generate new service id for using with ahamove. The generated
// id is always 8 character in length.
func (c serviceIDGenerator) GenerateServiceID(serviceID ServiceCode) (string, error) {
	if serviceID == "" {
		return "", cm.Errorf(cm.InvalidArgument, nil, "ahamove: Missing service ID")
	}
	service := ServicesIndexID[serviceID]
	if service == nil {
		return "", cm.Errorf(cm.InvalidArgument, nil, "ahamove: Service not found")
	}
	n := c.rd.Uint64()
	v := gencode.Alphabet32.EncodeReverse(n, 8)

	code := string(v)
	code = code[:6] + service.ShortCode
	return code, nil
}

func DecodeShippingServiceName(code string) (name string, ok bool) {
	if len(code) != 8 {
		return "", false
	}
	shortCode := name[6:]
	service := ServicesIndexShortCode[shortCode]
	if service == nil {
		return "", false
	}
	return service.Name, false
}

func (c *Carrier) GetServiceName(code string) (serviceName string, ok bool) {
	return DecodeShippingServiceName(code)
}

func (c *Carrier) ParseServiceCode(code string) (serviceID string, _err error) {
	sID, err := ParseServiceID(code)
	return string(sID), err
}

func ParseServiceID(code string) (serviceID ServiceCode, err error) {
	if code == "" {
		err = cm.Errorf(cm.InvalidArgument, nil, "ahamove: missing service id")
		return
	}
	if len(code) != 8 {
		err = cm.Errorf(cm.InvalidArgument, nil, "ahamove: invalid service id")
		return
	}

	shortCode := code[6:]
	service := ServicesIndexShortCode[shortCode]

	if service == nil {
		err = cm.Errorf(cm.InvalidArgument, nil, "ahamove: invalid service id")
	}
	serviceID = service.ID
	return
}

func GetAvailableServices(points []*ahamoveClient.DeliveryPointRequest) []*ShippingService {
	pointCount := len(points) - 1
	provinceCode := ""
	for i, point := range points {
		if i == 0 {
			continue
		}
		if provinceCode == "" {
			provinceCode = point.ProvinceCode
			continue
		}
		if provinceCode != point.ProvinceCode {
			return nil
		}
	}
	var cityCode CityCode
	switch provinceCode {
	case location.HCMProvinceCode:
		cityCode = SGNCode
	case location.HNProvinceCode:
		cityCode = HANCode
	default:
		return nil
	}
	services := ServicesIndexCity[cityCode]
	var result = make([]*ShippingService, 0, len(services))
	for _, s := range services {
		if pointCount >= s.MinStopPoints && pointCount <= s.MaxStopPoints {
			result = append(result, s)
		}
	}
	return result
}
