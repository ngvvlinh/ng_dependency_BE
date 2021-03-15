package model

import (
	"time"

	"o.o/capi/dot"
)

// +sqlgen
type Ticket struct {
	ID                 string
	AssignedUserID     string
	TicketStatus       string `sq:"'ticketstatus'"`
	TicketCategories   string `sq:"'ticketcategories'"`
	ModifiedBy         string `sq:"'modifiedby'"`
	TicketTitle        string
	Description        string
	ContactID          string
	Note               string
	FfmCode            string
	FfmUrl             string
	FfmID              string
	EtopID             string
	OrderID            dot.ID
	OrderCode          string
	Company            string
	Provider           string
	FromApp            string
	Code               string
	OldValue           string
	NewValue           string
	SubStatus          string `sq:"'substatus'"`
	EtopNote           string
	Reason             string
	CreatedAt          time.Time
	CreatedTime        time.Time `sq:"'createdtime'"`
	ModifiedAt         time.Time
	ConfirmedAt        time.Time
	ConfirmedBy        string
	ClosedAt           time.Time
	ClosedBy           string
	TicketNo           string
	EtopAccountID      string
	ExternalID         string
	ExternalStatus     string
	ExternalCSProvider string
	XData              string
	XError             string
	OData              string
}

// +sqlgen
type TicketComment struct {
	ID        string
	Content   string
	CreatorID string
	TicketID  string
	ThreadID  string
	CreatedAt time.Time
	OData     string `sq:"'o_data'"`
	XID       string `sq:"'x_id'"`
	XError    string `sq:"'x_error'"`
	From      string
}

// +sqlgen
type VtigerAccount struct {
	ID             string
	UserName       string
	FirstName      string
	LastName       string
	RoleID         string `sq:"'roleid'"`
	Email1         string `sq:"'email1'"`
	Email2         string `sq:"'email2'"`
	SecondaryEmail string `sq:"'secondaryemail'"`
	Status         string
	OData          string
	CreatedAt      string
	UpdatedAt      string
}
