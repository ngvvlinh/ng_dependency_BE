package xshop

import (
	"context"

	"o.o/api/main/contact"
	api "o.o/api/top/external/shop"
	externaltypes "o.o/api/top/external/types"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/etop/apix/convertpb"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/capi/dot"
)

type ContactService struct {
	session.Session
	ContactQuery contact.QueryBus
	ContactAggr  contact.CommandBus
}

func (s *ContactService) Clone() api.ContactService { res := *s; return &res }

func (s *ContactService) ListContacts(ctx context.Context, r *externaltypes.ListContactsRequest) (*externaltypes.ContactsResponse, error) {
	paging, err := cmapi.CMCursorPaging(r.Paging)
	if err != nil {
		return nil, err
	}
	shopID := s.SS.Shop().ID
	var IDs []dot.ID
	var phone, name string
	if r.Filter != nil {
		IDs = r.Filter.IDs
		phone = r.Filter.Phone
		name = r.Filter.Name
	}

	getContactsQuery := &contact.GetContactsQuery{
		ShopID: shopID,
		IDs:    IDs,
		Phone:  phone,
		Name:   name,
		Paging: *paging,
	}
	if err := s.ContactQuery.Dispatch(ctx, getContactsQuery); err != nil {
		return nil, err
	}

	return &externaltypes.ContactsResponse{
		Contacts: convertpb.Convert_core_Contacts_To_apix_Contacts(getContactsQuery.Result.Contacts),
		Paging:   cmapi.PbCursorPageInfo(paging, &getContactsQuery.Result.Paging),
	}, nil
}

func (s *ContactService) CreateContact(ctx context.Context, r *externaltypes.CreateContactRequest) (*externaltypes.Contact, error) {
	if r.Phone == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Số điện thoại không được để trống")
	}
	if r.FullName == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Tên không được để trống")
	}

	cmd := &contact.CreateContactCommand{
		ShopID:   s.SS.Shop().ID,
		FullName: r.FullName,
		Phone:    r.Phone,
	}
	if err := s.ContactAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}

	return convertpb.Convert_core_Contact_To_apix_Contact(cmd.Result), nil
}
