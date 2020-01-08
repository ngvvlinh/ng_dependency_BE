package types

import (
	"etop.vn/api/top/types/common"
	"etop.vn/capi/dot"
	"etop.vn/capi/filter"
	"etop.vn/common/jsonx"
)

type AddCustomerRequest struct {
	GroupID    dot.ID `json:"group_id"`
	CustomerID dot.ID `json:"customer_id"`
}

func (m *AddCustomerRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetCustomerAddressRequest struct {
	ID dot.ID `json:"id"`
}

func (m *GetCustomerAddressRequest) String() string { return jsonx.MustMarshalToString(m) }

type RemoveCustomerRequest struct {
	GroupID    dot.ID `json:"group_id"`
	CustomerID dot.ID `json:"customer_id"`
}

func (m *RemoveCustomerRequest) String() string { return jsonx.MustMarshalToString(m) }

type CustomerGroup struct {
	Id   dot.ID `json:"id"`
	Name string `json:"name"`
}

func (m *CustomerGroup) String() string { return jsonx.MustMarshalToString(m) }

type CustomerGroupsResponse struct {
	CustomerGroups []*CustomerGroup       `json:"customer_groups"`
	Paging         *common.CursorPageInfo `json:"paging"`
}

func (m *CustomerGroupsResponse) String() string { return jsonx.MustMarshalToString(m) }

type CustomerGroupFilter struct {
	ID filter.IDs `json:"id"`
}

func (m *CustomerGroupFilter) String() string { return jsonx.MustMarshalToString(m) }

type ListCustomerGroupsRequest struct {
	Filter CustomerGroupFilter  `json:"filter"`
	Paging *common.CursorPaging `json:"paging"`
}

func (m *ListCustomerGroupsRequest) String() string { return jsonx.MustMarshalToString(m) }

type CreateCustomerGroupRequest struct {
	Name string `json:"name"`
}

func (m *CreateCustomerGroupRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateCustomerGroupRequest struct {
	GroupId dot.ID         `json:"id"`
	Name    dot.NullString `json:"name"`
}

func (m *UpdateCustomerGroupRequest) String() string { return jsonx.MustMarshalToString(m) }

type CustomerGroupRelationship struct {
	CustomerID dot.ID `json:"customer_id"`
	GroupID    dot.ID `json:"group_id"`
}

func (m *CustomerGroupRelationship) String() string { return jsonx.MustMarshalToString(m) }

type CustomerGroupRelationshipFilter struct {
	CustomerID filter.IDs `json:"customer_id"`
	GroupID    filter.IDs `json:"group_id"`
}

func (m *CustomerGroupRelationshipFilter) String() string { return jsonx.MustMarshalToString(m) }

type ListCustomerGroupRelationshipsRequest struct {
	Filter CustomerGroupRelationshipFilter `json:"filter"`
	Paging *common.CursorPaging            `json:"paging"`
}

func (m *ListCustomerGroupRelationshipsRequest) String() string { return jsonx.MustMarshalToString(m) }

type CustomerGroupRelationshipsResponse struct {
	Relationships []*CustomerGroupRelationship `json:"relationship"`
	Paging        *common.CursorPageInfo       `json:"paging"`
}

func (m *CustomerGroupRelationshipsResponse) String() string { return jsonx.MustMarshalToString(m) }