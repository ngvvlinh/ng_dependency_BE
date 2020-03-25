package model

import (
	"time"

	"etop.vn/capi/dot"
)

// +sqlgen
type Address struct {
	ID        dot.ID `json:"id"`
	FullName  string `json:"full_name"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Position  string `json:"position"`
	Email     string `json:"email"`

	Country  string `json:"country"`
	City     string `json:"city"`
	Province string `json:"province"`
	District string `json:"district"`
	Ward     string `json:"ward"` // Ward may be non-empty while WardCode is empty
	Zip      string `json:"zip"`

	DistrictCode string `json:"district_code"`
	ProvinceCode string `json:"province_code"`
	WardCode     string `json:"ward_code"`

	Company     string       `json:"company"`
	Address1    string       `json:"address1"`
	Address2    string       `json:"address2"`
	Type        string       `json:"type"`
	AccountID   dot.ID       `json:"account_id"`
	Notes       *AddressNote `json:"notes"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Coordinates *Coordinates `json:"coordinates"`

	Rid dot.ID `json:"rid"`
}

type AddressNote struct {
	Note       string `json:"note"`
	OpenTime   string `json:"open_time"`
	LunchBreak string `json:"lunch_break"`
	Other      string `json:"other"`
}

type Coordinates struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}
