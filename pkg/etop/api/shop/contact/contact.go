package contact

import (
	"context"

	"o.o/api/main/contact"
	api "o.o/api/top/int/shop"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	convertpball "o.o/backend/pkg/etop/api/convertpb/_all"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/capi/dot"
)

type ContactService struct {
	session.Session

	ContactQuery contact.QueryBus
	ContactAggr  contact.CommandBus
}

func (s *ContactService) Clone() api.ContactService { res := *s; return &res }

func (s *ContactService) GetContact(
	ctx context.Context, req *api.GetContactRequest,
) (*api.Contact, error) {
	shopID := s.SS.Shop().ID
	query := &contact.GetContactByIDQuery{
		ID:     req.ID,
		ShopID: shopID,
	}
	if err := s.ContactQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}

	return convertpball.Convert_core_Contact_to_api_Contact(query.Result), nil
}

func (s *ContactService) GetContacts(
	ctx context.Context, req *api.GetContactsRequest,
) (*api.GetContactsResponse, error) {
	paging := cmapi.CMPaging(req.Paging)
	shopID := s.SS.Shop().ID
	var IDs []dot.ID
	var phone, name string
	if req.Filter != nil {
		IDs = req.Filter.IDs
		phone = req.Filter.Phone
		name = req.Filter.Name
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

	return &api.GetContactsResponse{
		Contacts: convertpball.PbContacts(getContactsQuery.Result.Contacts),
		Paging:   cmapi.PbPageInfo(paging),
	}, nil
}

func (s *ContactService) CreateContact(
	ctx context.Context, req *api.CreateContactRequest,
) (*api.Contact, error) {
	if req.Phone == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Số điện thoại không được để trống")
	}
	if req.FullName == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Tên không được để trống")
	}

	cmd := &contact.CreateContactCommand{
		ShopID:   s.SS.Shop().ID,
		FullName: req.FullName,
		Phone:    req.Phone,
	}
	if err := s.ContactAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}

	return convertpball.Convert_core_Contact_to_api_Contact(cmd.Result), nil
}

func (s *ContactService) UpdateContact(
	ctx context.Context, req *api.UpdateContactRequest,
) (*api.Contact, error) {
	shopID := s.SS.Shop().ID
	cmd := &contact.UpdateContactCommand{
		ID:       req.ID,
		ShopID:   shopID,
		FullName: req.FullName,
		Phone:    req.Phone,
	}
	if err := s.ContactAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}

	return convertpball.Convert_core_Contact_to_api_Contact(cmd.Result), nil
}

func (s *ContactService) DeleteContact(
	ctx context.Context, req *api.DeleteContactRequest,
) (*api.DeleteContactResponse, error) {
	shopID := s.SS.Shop().ID
	cmd := &contact.DeleteContactCommand{
		ID:     req.ID,
		ShopID: shopID,
	}
	if err := s.ContactAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return &api.DeleteContactResponse{Count: cmd.Result}, nil
}
