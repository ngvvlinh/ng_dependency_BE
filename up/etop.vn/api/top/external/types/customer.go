package types

import (
	"etop.vn/api/shopping/customering/customer_type"
	"etop.vn/api/top/int/etop"
	"etop.vn/api/top/types/common"
	"etop.vn/api/top/types/etc/gender"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/capi/dot"
	"etop.vn/capi/filter"
	"etop.vn/common/jsonx"
)

type Customer struct {
	Id           dot.ID         `json:"id"`
	ShopId       dot.ID         `json:"shop_id"`
	ExternalId   string         `json:"external_id"`
	ExternalCode string         `json:"external_code"`
	FullName     string         `json:"full_name"`
	Code         string         `json:"code"`
	Note         string         `json:"note"`
	Phone        string         `json:"phone"`
	Email        string         `json:"email"`
	Gender       string         `json:"gender"`
	Type         string         `json:"type"`
	Birthday     string         `json:"birthday"`
	CreatedAt    dot.Time       `json:"created_at"`
	UpdatedAt    dot.Time       `json:"updated_at"`
	Status       status3.Status `json:"status"`
}

func (m *Customer) String() string { return jsonx.MustMarshalToString(m) }

type CustomersResponse struct {
	Customers []*Customer            `json:"customers"`
	Paging    *common.CursorPageInfo `json:"paging"`
}

func (m *CustomersResponse) String() string { return jsonx.MustMarshalToString(m) }

type CreateCustomerRequest struct {
	ExternalId   string `json:"external_id"`
	ExternalCode string `json:"external_code"`

	// @required
	FullName string        `json:"full_name"`
	Gender   gender.Gender `json:"gender"`
	Birthday string        `json:"birthday"`
	// enum ('independent', 'individual', 'organization')
	Type customer_type.CustomerType `json:"type"`
	Note string                     `json:"note"`
	// @required
	Phone string `json:"phone"`
	Email string `json:"email"`
}

func (m *CreateCustomerRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateCustomerRequest struct {
	Id       dot.ID            `json:"id"`
	FullName dot.NullString    `json:"full_name"`
	Gender   gender.NullGender `json:"gender"`
	Birthday dot.NullString    `json:"birthday"`
	// enum ('individual', 'organization','independent')
	Type  dot.NullString `json:"type"`
	Note  dot.NullString `json:"note"`
	Phone dot.NullString `json:"phone"`
	Email dot.NullString `json:"email"`
}

func (m *UpdateCustomerRequest) String() string { return jsonx.MustMarshalToString(m) }

type DeleteCustomerRequest struct {
	Id         dot.ID `json:"id"`
	Code       string `json:"code"`
	ExternalId string `json:"external_id"`
}

func (m *DeleteCustomerRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetCustomerRequest struct {
	Id         dot.ID `json:"id"`
	Code       string `json:"code"`
	ExternalId string `json:"external_id"`
}

func (m *GetCustomerRequest) String() string { return jsonx.MustMarshalToString(m) }

type CustomerFilter struct {
	ID filter.IDs `json:"id"`
}

func (m *CustomerFilter) String() string { return jsonx.MustMarshalToString(m) }

type ListCustomersRequest struct {
	Filter CustomerFilter       `json:"filter"`
	Paging *common.CursorPaging `json:"paging"`
}

func (m *ListCustomersRequest) String() string { return jsonx.MustMarshalToString(m) }

type CustomerAddress struct {
	Id           dot.ID            `json:"id"`
	CustomerID   dot.ID            `json:"customer_id"`
	Province     string            `json:"province"`
	ProvinceCode string            `json:"province_code"`
	District     string            `json:"district"`
	DistrictCode string            `json:"district_code"`
	Ward         string            `json:"ward"`
	WardCode     string            `json:"ward_code"`
	Address1     string            `json:"address1"`
	Address2     string            `json:"address2"`
	FullName     string            `json:"full_name"`
	Company      string            `json:"company"`
	Phone        string            `json:"phone"`
	Email        string            `json:"email"`
	Position     string            `json:"position"`
	Coordinates  *etop.Coordinates `json:"coordinates"`
}

func (m *CustomerAddress) String() string { return jsonx.MustMarshalToString(m) }

type CustomerAddressesResponse struct {
	CustomerAddresses []*CustomerAddress     `json:"addresses"`
	Paging            *common.CursorPageInfo `json:"paging"`
}

func (m *CustomerAddressesResponse) String() string { return jsonx.MustMarshalToString(m) }

type CustomerAddressFilter struct {
	CustomerId filter.IDs `json:"customer_id"`
}

func (m *CustomerAddressFilter) String() string { return jsonx.MustMarshalToString(m) }

type ListCustomerAddressesRequest struct {
	Filter CustomerAddressFilter `json:"filter"`
	Paging *common.CursorPaging  `json:"paging"`
}

func (m *ListCustomerAddressesRequest) String() string { return jsonx.MustMarshalToString(m) }

type CreateCustomerAddressRequest struct {
	CustomerId   dot.ID       `json:"customer_id"`
	ProvinceCode string       `json:"province_code"`
	DistrictCode string       `json:"district_code"`
	WardCode     string       `json:"ward_code"`
	Address1     string       `json:"address1"`
	Address2     string       `json:"address2"`
	Country      string       `json:"country"`
	FullName     string       `json:"full_name"`
	Company      string       `json:"company"`
	Phone        string       `json:"phone"`
	Email        string       `json:"email"`
	Position     string       `json:"position"`
	Coordinates  *Coordinates `json:"coordinates"`
}

func (m *CreateCustomerAddressRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateCustomerAddressRequest struct {
	Id           dot.ID            `json:"id"`
	ProvinceCode dot.NullString    `json:"province_code"`
	DistrictCode dot.NullString    `json:"district_code"`
	WardCode     dot.NullString    `json:"ward_code"`
	Address1     dot.NullString    `json:"address1"`
	Address2     dot.NullString    `json:"address2"`
	Country      dot.NullString    `json:"country"`
	FullName     dot.NullString    `json:"full_name"`
	Phone        dot.NullString    `json:"phone"`
	Email        dot.NullString    `json:"email"`
	Position     dot.NullString    `json:"position"`
	Company      dot.NullString    `json:"company"`
	Coordinates  *etop.Coordinates `json:"coordinates"`
}

func (m *UpdateCustomerAddressRequest) String() string { return jsonx.MustMarshalToString(m) }
