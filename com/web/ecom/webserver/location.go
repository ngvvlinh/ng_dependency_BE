package webserver

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"o.o/api/main/location"
	apietop "o.o/api/top/int/etop"
	"o.o/backend/pkg/etop/api/convertpb"
)

type GetDistrictsRequest struct {
	ProvinceCode string `json:"province_code"`
}

type GetWardsRequest struct {
	DistrictCode string `json:"district_code"`
}

func (s *Server) Provinces(c echo.Context) error {
	query := &location.GetAllLocationsQuery{All: true}
	if err := locationBus.Dispatch(c.Request().Context(), query); err != nil {
		return err
	}
	result := &apietop.GetProvincesResponse{
		Provinces: convertpb.PbProvinces(query.Result.Provinces),
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response()).Encode(result)
}

func (s *Server) Districts(c echo.Context) error {
	var b []byte
	_, _ = c.Request().Body.Read(b)
	u := new(GetDistrictsRequest)
	if err := c.Bind(u); err != nil {
		return err
	}
	query := &location.GetAllLocationsQuery{ProvinceCode: u.ProvinceCode}
	if err := locationBus.Dispatch(c.Request().Context(), query); err != nil {
		return err
	}
	result := &apietop.GetDistrictsResponse{
		Districts: convertpb.PbDistricts(query.Result.Districts),
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response()).Encode(result)
}

func (s *Server) Wards(c echo.Context) error {
	u := new(GetWardsRequest)
	if err := c.Bind(u); err != nil {
		return err
	}
	query := &location.GetAllLocationsQuery{DistrictCode: u.DistrictCode}
	if err := locationBus.Dispatch(c.Request().Context(), query); err != nil {
		return err
	}
	result := &apietop.GetWardsResponse{
		Wards: convertpb.PbWards(query.Result.Wards),
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response()).Encode(result)
}
