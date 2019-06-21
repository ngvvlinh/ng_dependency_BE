package ahamove

import (
	"context"
	"fmt"
	"math/rand"
	"strings"

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
	Code          string
	Name          string
	MinStopPoints int
	MaxStopPoints int
	MaxWeight     int
	City          CityCode
	Description   string
}

type Service struct {
	Code      ServiceCode
	ShortCode byte
	Name      string
}

var (
	BIKE    ServiceCode = "BIKE"
	POOL    ServiceCode = "POOL"
	SAMEDAY ServiceCode = "SAMEDAY"

	SGNCode CityCode = "SGN"
	HANCode CityCode = "HAN"

	ServicesIndexCode      = make(map[ServiceCode]*Service)
	ServicesIndexShortCode = make(map[byte]*Service)
)

var Services = []*Service{
	{
		Code:      BIKE,
		ShortCode: 'B',
		Name:      "Siêu Tốc",
	}, {
		Code:      POOL,
		ShortCode: 'P',
		Name:      "Siêu Rẻ",
	}, {
		Code:      SAMEDAY,
		ShortCode: 'S',
		Name:      "Trong ngày",
	},
}

type serviceIDGenerator struct {
	rd *rand.Rand
}

func init() {
	for _, service := range Services {
		ServicesIndexCode[service.Code] = service
		ServicesIndexShortCode[service.ShortCode] = service
	}
}

func newServiceIDGenerator(seed int64) serviceIDGenerator {
	src := rand.NewSource(seed)
	rd := rand.New(src)
	return serviceIDGenerator{rd}
}

// GenerateServiceID generate new service id for using with ahamove. The generated
// id is always 8 character in length.
func (c serviceIDGenerator) GenerateServiceID(serviceCode string) (string, error) {
	city, code, err := ParseInfoAhamoveServiceCode(serviceCode)
	if err != nil {
		return "", err
	}

	service := ServicesIndexCode[code]
	if service == nil {
		return "", cm.Errorf(cm.InvalidArgument, nil, "Ahamove: Service not found")
	}

	n := c.rd.Uint64()
	v := gencode.Alphabet32.EncodeReverse(n, 8)
	v = v[:8]

	switch city {
	case SGNCode:
		v[6] = 'S'
	case HANCode:
		v[6] = 'H'
	default:
		return "", cm.Errorf(cm.InvalidArgument, nil, "Ahamove: Invalid city code")
	}
	v[7] = service.ShortCode

	return string(v), nil
}

func DecodeShippingServiceName(code string) (name string, ok bool) {
	if len(code) != 8 {
		return "", false
	}
	service := ServicesIndexShortCode[code[6]]
	if service == nil {
		return "", false
	}

	return service.Name, true
}

func (c *Carrier) GetServiceName(code string) (serviceName string, ok bool) {
	return DecodeShippingServiceName(code)
}

func (c *Carrier) ParseServiceCode(code string) (serviceCode string, _err error) {
	sCode, err := parseServiceCode(code)
	return sCode, err
}

func (c *Carrier) GetAvailableServices(ctx context.Context, points []*ahamoveClient.DeliveryPointRequest) (res []*ShippingService) {
	pointCount := len(points) - 1
	// min stop point = 1
	if pointCount < 1 {
		return nil
	}

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

	cmd := &GetServiceCommand{
		Request: &ahamoveClient.GetServicesRequest{
			CityID: string(cityCode),
		},
	}
	if err := c.GetServices(ctx, cmd); err != nil {
		return nil
	}
	for _, service := range cmd.Result {
		s := ToService(service)
		if s != nil && pointCount >= s.MinStopPoints && pointCount <= s.MaxStopPoints {
			res = append(res, s)
		}
	}
	return res
}

func ToService(service *ahamoveClient.ServiceType) *ShippingService {
	if service == nil {
		return nil
	}

	return &ShippingService{
		Code:          service.ID,
		Name:          service.NameViVn,
		MinStopPoints: 1,
		MaxStopPoints: service.MaxStopPoints,
		City:          CityCode(service.CityID),
		Description:   service.DescriptionViVn,
	}
}

func parseServiceCode(code string) (serviceCode string, err error) {
	if code == "" {
		err = cm.Errorf(cm.InvalidArgument, nil, "Ahamove: Missing service code")
		return
	}
	if len(code) != 8 {
		err = cm.Errorf(cm.InvalidArgument, nil, "Ahamove: Invalid service code")
		return
	}

	service := ServicesIndexShortCode[code[7]]
	if service == nil {
		err = cm.Errorf(cm.InvalidArgument, nil, "Ahamove: Invalid service code")
	}

	var city CityCode
	switch {
	case code[6] == 'S':
		city = SGNCode
	case code[6] == 'H':
		city = HANCode
	default:
		return "", cm.Errorf(cm.InvalidArgument, nil, "Ahamove: Invalid city code")
	}

	return AhamoveServiceCodeFormat(city, service.Code), nil
}

func AhamoveServiceCodeFormat(city CityCode, code ServiceCode) string {
	return fmt.Sprintf("%v-%v", city, code)
}

func ParseInfoAhamoveServiceCode(code string) (CityCode, ServiceCode, error) {
	// Ahamove service code format: SGN-BIKE
	if code == "" {
		return "", "", cm.Errorf(cm.InvalidArgument, nil, "Ahamove: Missing service ID")
	}

	arr := strings.Split(code, "-")
	if len(arr) != 2 {
		return "", "", cm.Errorf(cm.InvalidArgument, nil, "Ahamove: Invalid service code")
	}
	city, code := arr[0], arr[1]
	return CityCode(city), ServiceCode(code), nil
}
