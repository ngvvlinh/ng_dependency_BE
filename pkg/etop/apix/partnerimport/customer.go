package partnerimport

import (
	"context"

	"o.o/api/top/external/whitelabel"
	"o.o/api/top/types/etc/customer_type"
	"o.o/backend/com/shopping/customering/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/capi/dot"
)

func (s *ImportService) Customers(ctx context.Context, r *CustomersEndpoint) error {
	if len(r.Customers) > MaximumItems {
		return cm.Errorf(cm.InvalidArgument, nil, "cannot handle rather than 100 items at once")
	}

	var ids []dot.ID
	for _, customer := range r.Customers {
		if customer.ExternalID == "" {
			return cm.Errorf(cm.InvalidArgument, nil, "id should not be null")
		}
		if customer.Type != customer_type.Individual && customer.Type != customer_type.Organization {
			return cm.Errorf(cm.InvalidArgument, nil, "type is invalid")
		}

		shopCustomer := &model.ShopCustomer{
			ShopID:       r.Context.Shop.ID,
			PartnerID:    r.Context.AuthPartnerID,
			ExternalID:   customer.ExternalID,
			ExternalCode: customer.ExternalCode,
			Code:         customer.ExternalCode,
			FullName:     customer.FullName,
			Gender:       customer.Gender,
			Type:         customer.Type,
			Birthday:     customer.Birthday,
			Note:         customer.Note,
			Phone:        customer.Phone,
			Email:        customer.Email,
			Status:       0,
			CreatedAt:    customer.CreatedAt.ToTime(),
			UpdatedAt:    customer.UpdatedAt.ToTime(),
			DeletedAt:    customer.DeletedAt.ToTime(),
		}

		oldShopCustomer, err := s.customerStoreFactory(ctx).ExternalID(customer.ExternalID).GetCustomerDB()
		switch cm.ErrorCode(err) {
		case cm.NotFound:
			id := cm.NewID()
			ids = append(ids, id)
			shopCustomer.ID = id
			if _err := s.customerStoreFactory(ctx).CreateCustomer(shopCustomer); _err != nil {
				return _err
			}
		case cm.NoError:
			shopCustomer.ID = oldShopCustomer.ID
			ids = append(ids, oldShopCustomer.ID)
			if _err := s.customerStoreFactory(ctx).UpdateCustomerDB(shopCustomer); _err != nil {
				return _err
			}
		default:
			return err
		}
	}

	modelCustomers, err := s.customerStoreFactory(ctx).IDs(ids...).ListCustomersDB()
	if err != nil {
		return err
	}

	var customersResponse []*whitelabel.Customer
	for _, customer := range modelCustomers {
		customersResponse = append(customersResponse, &whitelabel.Customer{
			ExternalId:   customer.ExternalID,
			ExternalCode: customer.ExternalCode,
			ID:           customer.ID,
			PartnerID:    customer.PartnerID,
			ShopID:       customer.ShopID,
			FullName:     customer.FullName,
			Gender:       customer.Gender,
			Birthday:     customer.Birthday,
			Type:         customer.Type,
			Note:         customer.Note,
			Phone:        customer.Phone,
			Email:        customer.Email,
			CreatedAt:    cmapi.PbTime(customer.CreatedAt),
			UpdatedAt:    cmapi.PbTime(customer.UpdatedAt),
			DeletedAt:    cmapi.PbTime(customer.DeletedAt),
		})
	}
	r.Result = &whitelabel.ImportCustomersResponse{Customers: customersResponse}
	return nil
}
