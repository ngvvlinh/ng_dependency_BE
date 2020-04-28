package ahamove

import (
	"context"
	"fmt"
	"math/rand"
	"strings"

	"o.o/backend/com/main/location"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/code/gencode"
	"o.o/backend/pkg/integration/shipnow/ahamove/client"
)

type (
	ServiceCode string
	CityCode    string
)

type AhamoveShippingService struct {
	Code          string // SGN-BIKE
	Name          string
	MinStopPoints int
	MaxStopPoints int
	MaxWeight     int
	MaxCOD        int
	City          CityCode
	Description   string
}

type Service struct {
	Code      ServiceCode
	ShortCode byte
	Name      string
}

var (
	BIKE      ServiceCode = "BIKE"
	POOL      ServiceCode = "POOL"
	SAMEDAY   ServiceCode = "SAMEDAY"
	SAMEPRICE ServiceCode = "DG"

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
		Name:      "Trong Ngày",
	}, {
		Code:      SAMEPRICE,
		ShortCode: 'D',
		Name:      "Đồng Giá 25",
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
func (c serviceIDGenerator) generateServiceID(serviceCode string) (string, error) {
	city, code, err := parseInfoAhamoveServiceCode(serviceCode)
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

func decodeShippingServiceName(code string) (name string, ok bool) {
	if len(code) != 8 {
		return "", false
	}
	service := ServicesIndexShortCode[code[6]]
	if service == nil {
		return "", false
	}

	return service.Name, true
}

func (c *Carrier) getAvailableServices(ctx context.Context, points []*client.DeliveryPointRequest) (res []*AhamoveShippingService) {
	pointCount := len(points) - 1
	totalCOD := 0
	// min stop point = 1
	if pointCount < 1 {
		return nil
	}

	provinceCode := ""
	for i, point := range points {
		if i == 0 {
			continue
		}
		totalCOD += point.COD
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
	request := &client.GetServicesRequest{
		CityID: string(cityCode),
	}
	services, err := c.client.GetServices(ctx, request)
	if err != nil {
		return nil
	}
	for _, service := range services {
		s := toService(service)
		if validateService(s, pointCount, totalCOD) {
			res = append(res, s)
		}
	}
	return res
}

func validateService(s *AhamoveShippingService, pointCount int, totalCOD int) bool {
	if s == nil {
		return false
	}
	if totalCOD >= s.MaxCOD {
		return false
	}
	if pointCount < s.MinStopPoints || pointCount > s.MaxStopPoints {
		return false
	}
	return true
}

func toService(service *client.ServiceType) *AhamoveShippingService {
	if service == nil {
		return nil
	}
	city, _, err := parseInfoAhamoveServiceCode(service.ID)
	if err != nil {
		return nil
	}
	minStopPoints := service.MinStopPoints
	if minStopPoints == 0 {
		minStopPoints = 1
	}

	return &AhamoveShippingService{
		Code:          service.ID, // keep the original service code
		Name:          service.NameViVn,
		MinStopPoints: int(minStopPoints),
		MaxStopPoints: int(service.MaxStopPoints),
		City:          city,
		Description:   service.DescriptionViVn,
		MaxCOD:        cm.CoalesceInt(int(service.MaxCOD), int(service.COD)),
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
		return
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

	return ahamoveServiceCodeFormat(city, service.Code), nil
}

func ahamoveServiceCodeFormat(city CityCode, code ServiceCode) string {
	return fmt.Sprintf("%v-%v", city, code)
}

func parseInfoAhamoveServiceCode(code string) (CityCode, ServiceCode, error) {
	// Ahamove service code format: SGN-BIKE
	if code == "" {
		return "", "", cm.Errorf(cm.InvalidArgument, nil, "Ahamove: Missing service ID")
	}

	arr := strings.Split(code, "-")
	if len(arr) != 2 {
		return "", "", cm.Errorf(cm.InvalidArgument, nil, "Ahamove: Invalid service code")
	}
	city, code := arr[0], arr[1]

	serviceCode := ServiceCode(code)
	service := ServicesIndexCode[serviceCode]
	if service == nil {
		return "", "", cm.Errorf(cm.InvalidArgument, nil, "Ahamove: Invalid service code (%v)", serviceCode)
	}
	return CityCode(city), serviceCode, nil
}
