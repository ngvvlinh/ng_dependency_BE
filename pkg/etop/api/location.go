package api

import (
	"context"

	"o.o/api/main/location"
	api "o.o/api/top/int/etop"
	pbcm "o.o/api/top/types/common"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/authorize/session"
)

type LocationService struct {
	session.Session

	LocationQuery location.QueryBus
}

func (s *LocationService) Clone() api.LocationService {
	res := *s
	return &res
}

func (s *LocationService) GetProvinces(ctx context.Context, q *pbcm.Empty) (*api.GetProvincesResponse, error) {
	query := &location.GetAllLocationsQuery{All: true}
	if err := s.LocationQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := &api.GetProvincesResponse{
		Provinces: convertpb.PbProvinces(query.Result.Provinces),
	}
	return result, nil
}

func (s *LocationService) GetDistricts(ctx context.Context, q *pbcm.Empty) (*api.GetDistrictsResponse, error) {
	query := &location.GetAllLocationsQuery{All: true}
	if err := s.LocationQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := &api.GetDistrictsResponse{
		Districts: convertpb.PbDistricts(query.Result.Districts),
	}
	return result, nil
}

func (s *LocationService) GetDistrictsByProvince(ctx context.Context, q *api.GetDistrictsByProvinceRequest) (*api.GetDistrictsResponse, error) {
	query := &location.GetAllLocationsQuery{ProvinceCode: q.ProvinceCode}
	if err := s.LocationQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := &api.GetDistrictsResponse{
		Districts: convertpb.PbDistricts(query.Result.Districts),
	}
	return result, nil
}

func (s *LocationService) GetWards(ctx context.Context, q *pbcm.Empty) (*api.GetWardsResponse, error) {
	query := &location.GetAllLocationsQuery{All: true}
	if err := s.LocationQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := &api.GetWardsResponse{
		Wards: convertpb.PbWards(query.Result.Wards),
	}
	return result, nil
}

func (s *LocationService) GetWardsByDistrict(ctx context.Context, q *api.GetWardsByDistrictRequest) (*api.GetWardsResponse, error) {
	query := &location.GetAllLocationsQuery{DistrictCode: q.DistrictCode}
	if err := s.LocationQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := &api.GetWardsResponse{
		Wards: convertpb.PbWards(query.Result.Wards),
	}
	return result, nil
}

func (s *LocationService) ParseLocation(ctx context.Context, q *api.ParseLocationRequest) (*api.ParseLocationResponse, error) {
	query := &location.FindLocationQuery{
		Province: q.ProvinceName,
		District: q.DistrictName,
		Ward:     q.WardName,
	}
	if err := s.LocationQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	loc := query.Result
	res := &api.ParseLocationResponse{}
	if loc.Province != nil {
		res.Province = convertpb.PbProvince(loc.Province)
	}
	if loc.District != nil {
		res.District = convertpb.PbDistrict(loc.District)
	}
	if loc.Ward != nil {
		res.Ward = convertpb.PbWard(loc.Ward)
	}
	result := res
	return result, nil
}
