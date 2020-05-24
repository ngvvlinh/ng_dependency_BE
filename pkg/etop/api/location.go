package api

import (
	"context"

	"o.o/api/main/location"
	apietop "o.o/api/top/int/etop"
	"o.o/backend/pkg/etop/api/convertpb"
)

type LocationService struct {
	LocationQuery location.QueryBus
}

func (s *LocationService) Clone() *LocationService {
	res := *s
	return &res
}

func (s *LocationService) GetProvinces(ctx context.Context, q *GetProvincesEndpoint) error {
	query := &location.GetAllLocationsQuery{All: true}
	if err := s.LocationQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &apietop.GetProvincesResponse{
		Provinces: convertpb.PbProvinces(query.Result.Provinces),
	}
	return nil
}

func (s *LocationService) GetDistricts(ctx context.Context, q *GetDistrictsEndpoint) error {
	query := &location.GetAllLocationsQuery{All: true}
	if err := s.LocationQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &apietop.GetDistrictsResponse{
		Districts: convertpb.PbDistricts(query.Result.Districts),
	}
	return nil
}

func (s *LocationService) GetDistrictsByProvince(ctx context.Context, q *GetDistrictsByProvinceEndpoint) error {
	query := &location.GetAllLocationsQuery{ProvinceCode: q.ProvinceCode}
	if err := s.LocationQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &apietop.GetDistrictsResponse{
		Districts: convertpb.PbDistricts(query.Result.Districts),
	}
	return nil
}

func (s *LocationService) GetWards(ctx context.Context, q *GetWardsEndpoint) error {
	query := &location.GetAllLocationsQuery{All: true}
	if err := s.LocationQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &apietop.GetWardsResponse{
		Wards: convertpb.PbWards(query.Result.Wards),
	}
	return nil
}

func (s *LocationService) GetWardsByDistrict(ctx context.Context, q *GetWardsByDistrictEndpoint) error {
	query := &location.GetAllLocationsQuery{DistrictCode: q.DistrictCode}
	if err := s.LocationQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &apietop.GetWardsResponse{
		Wards: convertpb.PbWards(query.Result.Wards),
	}
	return nil
}

func (s *LocationService) ParseLocation(ctx context.Context, q *ParseLocationEndpoint) error {
	query := &location.FindLocationQuery{
		Province: q.ProvinceName,
		District: q.DistrictName,
		Ward:     q.WardName,
	}
	if err := s.LocationQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	loc := query.Result
	res := &apietop.ParseLocationResponse{}
	if loc.Province != nil {
		res.Province = convertpb.PbProvince(loc.Province)
	}
	if loc.District != nil {
		res.District = convertpb.PbDistrict(loc.District)
	}
	if loc.Ward != nil {
		res.Ward = convertpb.PbWard(loc.Ward)
	}
	q.Result = res
	return nil
}
