package contact

import (
	"context"

	"o.o/api/main/contact"
	api "o.o/api/top/int/shop"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/authorize/session"
)

type ContactService struct {
	session.Session

	ContactQuery contact.QueryBus
	ContactAggr  contact.CommandBus
}

func (s *ContactService) Clone() api.ContactService { res := *s; return &res }

func (s *ContactService) GetContactByID(
	ctx context.Context, req *api.GetContactByIDRequest,
) (*api.Contact, error) {
	shopID := s.SS.Shop().ID
	query := &contact.GetContactByIDQuery{
		ID:     req.ID,
		ShopID: shopID,
	}
	if err := s.ContactQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}

	return convertpb.Convert_core_Contact_to_api_Contact(query.Result), nil
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

	return convertpb.Convert_core_Contact_to_api_Contact(cmd.Result), nil
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

	return convertpb.Convert_core_Contact_to_api_Contact(cmd.Result), nil
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
