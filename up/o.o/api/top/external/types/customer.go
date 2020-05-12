package types

import (
	"o.o/api/top/int/etop"
	"o.o/api/top/types/common"
	"o.o/api/top/types/etc/customer_type"
	"o.o/api/top/types/etc/gender"
	"o.o/api/top/types/etc/status3"
	"o.o/capi/dot"
	"o.o/capi/filter"
	"o.o/common/jsonx"
)

type Customer struct {
	ExternalId   dot.NullString `json:"external_id"`
	ExternalCode dot.NullString `json:"external_code"`

	Id        dot.ID                         `json:"id"`
	ShopId    dot.ID                         `json:"shop_id"`
	FullName  dot.NullString                 `json:"full_name"`
	Code      dot.NullString                 `json:"code"`
	Note      dot.NullString                 `json:"note"`
	Phone     dot.NullString                 `json:"phone"`
	Email     dot.NullString                 `json:"email"`
	Gender    gender.NullGender              `json:"gender"`
	Type      customer_type.NullCustomerType `json:"type"`
	Birthday  dot.NullString                 `json:"birthday"`
	CreatedAt dot.Time                       `json:"created_at"`
	UpdatedAt dot.Time                       `json:"updated_at"`
	Status    status3.NullStatus             `json:"status"`
	Deleted   bool                           `json:"deleted"`
}

func (m *Customer) String() string { return jsonx.MustMarshalToString(m) }

func (m *Customer) HasChanged() bool {
	return m.FullName.Valid ||
		m.Code.Valid ||
		m.Note.Valid ||
		m.Phone.Valid ||
		m.Email.Valid ||
		m.Gender.Valid ||
		m.Type.Valid ||
		m.Birthday.Valid ||
		m.Status.Valid
}

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
	Filter         CustomerFilter       `json:"filter"`
	Paging         *common.CursorPaging `json:"paging"`
	IncludeDeleted bool                 `json:"include_deleted"`
}

func (m *ListCustomersRequest) String() string { return jsonx.MustMarshalToString(m) }

type CustomerAddress struct {
	Id           dot.ID            `json:"id"`
	CustomerID   dot.ID            `json:"customer_id"`
	Province     dot.NullString    `json:"province"`
	ProvinceCode dot.NullString    `json:"province_code"`
	District     dot.NullString    `json:"district"`
	DistrictCode dot.NullString    `json:"district_code"`
	Ward         dot.NullString    `json:"ward"`
	WardCode     dot.NullString    `json:"ward_code"`
	Address1     dot.NullString    `json:"address1"`
	Address2     dot.NullString    `json:"address2"`
	FullName     dot.NullString    `json:"full_name"`
	Company      dot.NullString    `json:"company"`
	Phone        dot.NullString    `json:"phone"`
	Email        dot.NullString    `json:"email"`
	Position     dot.NullString    `json:"position"`
	Coordinates  *etop.Coordinates `json:"coordinates"`
	Deleted      bool              `json:"deleted"`
}

func (m *CustomerAddress) String() string { return jsonx.MustMarshalToString(m) }

func (m *CustomerAddress) HasChanged() bool {
	return m.WardCode.Valid ||
		m.Ward.Valid ||
		m.DistrictCode.Valid ||
		m.District.Valid ||
		m.Position.Valid ||
		m.Email.Valid ||
		m.Phone.Valid ||
		m.Company.Valid ||
		m.FullName.Valid ||
		m.Address1.Valid ||
		m.Address2.Valid ||
		m.Province.Valid ||
		m.ProvinceCode.Valid
}

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
	Filter         CustomerAddressFilter `json:"filter"`
	Paging         *common.CursorPaging  `json:"paging"`
	IncludeDeleted bool                  `json:"include_deleted"`
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
