package model

import "time"

type Address struct {
	Country      string
	CountryCode  string
	CountryName  string
	ProvinceCode string
	District     string
	DistrictCode string
	Ward         string
	WardCode     string
	Address1     string
	Address2     string
	Zip          string
	City         string
	Phone        string
	Name         string
}

type ShippingRate struct {
	ServiceID       string
	ServiceName     string
	ServiceCode     string
	Currency        string
	TotalPrice      int32
	PhoneRequired   bool
	MinDeliveryDate time.Time
	MaxDeliveryDate time.Time
	Description     string
}
