package customer

import (
	"context"
	"fmt"
	"time"

	"o.o/api/main/location"
	"o.o/api/main/ordering"
	"o.o/api/main/receipting"
	"o.o/api/shopping/addressing"
	"o.o/api/shopping/customering"
	"o.o/api/top/int/shop"
	api "o.o/api/top/int/shop"
	pbcm "o.o/api/top/types/common"
	"o.o/api/top/types/etc/receipt_type"
	"o.o/api/top/types/etc/status3"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/etop/api/convertpb"
	shop2 "o.o/backend/pkg/etop/api/shop"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/capi/dot"
	"o.o/capi/filter"
)

type CustomerService struct {
	session.Session

	LocationQuery location.QueryBus
	CustomerQuery customering.QueryBus
	CustomerAggr  customering.CommandBus
	AddressAggr   addressing.CommandBus
	AddressQuery  addressing.QueryBus
	OrderQuery    ordering.QueryBus
	ReceiptQuery  receipting.QueryBus
}

func (s *CustomerService) Clone() api.CustomerService { res := *s; return &res }

func (s *CustomerService) CreateCustomer(ctx context.Context, r *api.CreateCustomerRequest) (*api.Customer, error) {
	key := fmt.Sprintf("CreateCustomer %v-%v-%v-%v-%v",
		s.SS.Shop().ID, s.SS.Claim().UserID, r.Phone, r.FullName, r.Email)
	res, _, err := shop2.Idempgroup.DoAndWrap(
		ctx, key, 15*time.Second, "tạo khách hàng",
		func() (interface{}, error) {
			cmd := &customering.CreateCustomerCommand{
				ShopID:   s.SS.Shop().ID,
				FullName: r.FullName,
				Gender:   r.Gender,
				Type:     r.Type,
				Birthday: r.Birthday,
				Note:     r.Note,
				Phone:    r.Phone,
				Email:    r.Email,
			}
			if err := s.CustomerAggr.Dispatch(ctx, cmd); err != nil {
				return nil, err
			}
			result := convertpb.PbCustomer(cmd.Result)
			return result, nil
		})

	if err != nil {
		return nil, err
	}
	result := res.(*shop.Customer)
	return result, nil
}

func (s *CustomerService) UpdateCustomer(ctx context.Context, r *api.UpdateCustomerRequest) (*api.Customer, error) {
	cmd := &customering.UpdateCustomerCommand{
		ID:       r.Id,
		ShopID:   s.SS.Shop().ID,
		FullName: r.FullName,
		Gender:   r.Gender,
		Birthday: r.Birthday,
		Note:     r.Note,
		Phone:    r.Phone,
		Email:    r.Email,
		Type:     r.Type.Apply(0),
	}
	err := s.CustomerAggr.Dispatch(ctx, cmd)
	if err != nil {
		return nil, err
	}

	result := convertpb.PbCustomer(cmd.Result)
	return result, nil
}

func (s *CustomerService) BatchSetCustomersStatus(ctx context.Context, r *api.SetCustomersStatusRequest) (*pbcm.UpdatedResponse, error) {
	cmd := &customering.BatchSetCustomersStatusCommand{
		IDs:    r.Ids,
		ShopID: s.SS.Shop().ID,
		Status: int(r.Status),
	}
	if err := s.CustomerAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.UpdatedResponse{Updated: cmd.Result.Updated}
	return result, nil
}

func (s *CustomerService) DeleteCustomer(ctx context.Context, r *pbcm.IDRequest) (*pbcm.DeletedResponse, error) {
	cmd := &customering.DeleteCustomerCommand{
		ID:     r.Id,
		ShopID: s.SS.Shop().ID,
	}
	if err := s.CustomerAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.DeletedResponse{Deleted: cmd.Result}
	return result, nil
}

func (s *CustomerService) GetCustomer(ctx context.Context, r *pbcm.IDRequest) (*api.Customer, error) {
	query := &customering.GetCustomerByIDQuery{
		ID:     r.Id,
		ShopID: s.SS.Shop().ID,
	}
	if err := s.CustomerQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := convertpb.PbCustomer(query.Result)
	if err := s.listLiabilities(ctx, s.SS.Shop().ID, []*shop.Customer{result}); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *CustomerService) GetCustomers(ctx context.Context, r *api.GetCustomersRequest) (resp *api.CustomersResponse, err error) {
	paging := cmapi.CMPaging(r.Paging)
	if !r.GetAll.Valid {
		r.GetAll = dot.Bool(true)
	}
	if r.GetAll.Apply(false) {
		resp, err = s.getAllCustomers(ctx, paging, r)
	} else {
		resp, err = s.getCustomers(ctx, paging, r)
	}
	if err != nil {
		return nil, err
	}
	err = s.listLiabilities(ctx, s.SS.Shop().ID, resp.Customers)
	return resp, err
}

func (s *CustomerService) getAllCustomers(ctx context.Context, paging *cm.Paging, r *api.GetCustomersRequest) (*api.CustomersResponse, error) {
	queryCustomerIndenpendent := &customering.GetCustomerIndependentQuery{}
	if err := s.CustomerQuery.Dispatch(ctx, queryCustomerIndenpendent); err != nil {
		return nil, err
	}
	var customers []*shop.Customer
	customers = append(customers, convertpb.PbCustomer(queryCustomerIndenpendent.Result))

	result := &api.CustomersResponse{
		Paging: cmapi.PbPageInfo(paging),
	}
	if paging.Limit == 1 && paging.Offset == 0 {
		result.Customers = customers
		return result, nil
	}
	if paging.Offset == 0 {
		paging.Limit--
		resp, err := s.getCustomers(ctx, paging, r)
		if err != nil {
			return nil, err
		}
		customers = append(customers, resp.Customers...)
		result.Customers = customers
	} else {
		paging.Offset--
		return s.getCustomers(ctx, paging, r)
	}
	return result, nil
}

func (s *CustomerService) getCustomers(ctx context.Context, paging *cm.Paging, r *api.GetCustomersRequest) (*api.CustomersResponse, error) {
	var fullTextSearch filter.FullTextSearch = ""
	if r.Filter != nil {
		fullTextSearch = r.Filter.FullName
	}
	query := &customering.ListCustomersQuery{
		ShopID:  s.SS.Shop().ID,
		Paging:  *paging,
		Filters: cmapi.ToFilters(r.Filters),
		Name:    fullTextSearch,
	}
	if err := s.CustomerQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := &api.CustomersResponse{
		Customers: convertpb.PbCustomers(query.Result.Customers),
		Paging:    cmapi.PbPageInfo(paging),
	}
	return result, nil
}

func (s *CustomerService) GetCustomersByIDs(ctx context.Context, r *pbcm.IDsRequest) (*api.CustomersResponse, error) {
	query := &customering.ListCustomersByIDsQuery{
		IDs:    r.Ids,
		ShopID: s.SS.Shop().ID,
	}
	if err := s.CustomerQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := &api.CustomersResponse{
		Customers: convertpb.PbCustomers(query.Result.Customers),
	}
	if err := s.listLiabilities(ctx, s.SS.Shop().ID, result.Customers); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *CustomerService) GetCustomerDetails(ctx context.Context, r *pbcm.IDRequest) (*api.CustomerDetailsResponse, error) {
	return nil, cm.ErrTODO
}

func (s *CustomerService) AddCustomersToGroup(ctx context.Context, r *api.AddCustomerToGroupRequest) (*pbcm.UpdatedResponse, error) {
	cmd := &customering.AddCustomersToGroupCommand{
		ShopID:      s.SS.Shop().ID,
		GroupID:     r.GroupId,
		CustomerIDs: r.CustomerIds,
	}
	if err := s.CustomerAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.UpdatedResponse{Updated: cmd.Result}
	return result, nil
}

func (s *CustomerService) RemoveCustomersFromGroup(ctx context.Context, r *api.RemoveCustomerOutOfGroupRequest) (*pbcm.RemovedResponse, error) {
	cmd := &customering.RemoveCustomersFromGroupCommand{
		ShopID:      s.SS.Shop().ID,
		GroupID:     r.GroupId,
		CustomerIDs: r.CustomerIds,
	}
	if err := s.CustomerAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.RemovedResponse{Removed: cmd.Result}
	return result, nil
}

func (s *CustomerService) listLiabilities(ctx context.Context, shopID dot.ID, customers []*shop.Customer) error {
	var customerIDs []dot.ID
	mapCustomerIDAndNumberOfOrders := make(map[dot.ID]int)
	mapCustomerIDAndTotalAmountOrders := make(map[dot.ID]int)
	mapCustomerIDAndTotalAmountReceipts := make(map[dot.ID]int)

	for _, customer := range customers {
		customerIDs = append(customerIDs, customer.Id)
	}

	getOrdersByCustomerIDs := &ordering.ListOrdersByCustomerIDsQuery{
		CustomerIDs: customerIDs,
		ShopID:      shopID,
	}
	if err := s.OrderQuery.Dispatch(ctx, getOrdersByCustomerIDs); err != nil {
		return err
	}
	for _, order := range getOrdersByCustomerIDs.Result.Orders {
		mapCustomerIDAndNumberOfOrders[order.CustomerID] += 1
		mapCustomerIDAndTotalAmountOrders[order.CustomerID] += order.TotalAmount
	}

	getReceiptsByCustomerIDs := &receipting.ListReceiptsByTraderIDsAndStatusesQuery{
		ShopID:    shopID,
		TraderIDs: customerIDs,
		Statuses:  []status3.Status{status3.P},
	}
	if err := s.ReceiptQuery.Dispatch(ctx, getReceiptsByCustomerIDs); err != nil {
		return err
	}
	for _, receipt := range getReceiptsByCustomerIDs.Result.Receipts {
		switch receipt.Type {
		case receipt_type.Receipt:
			mapCustomerIDAndTotalAmountReceipts[receipt.TraderID] += receipt.Amount
		case receipt_type.Payment:
			mapCustomerIDAndTotalAmountReceipts[receipt.TraderID] -= receipt.Amount
		}
	}

	for _, customer := range customers {
		customer.Liability = &shop.CustomerLiability{
			TotalOrders:    mapCustomerIDAndNumberOfOrders[customer.Id],
			TotalAmount:    mapCustomerIDAndTotalAmountOrders[customer.Id],
			ReceivedAmount: mapCustomerIDAndTotalAmountReceipts[customer.Id],
			Liability:      mapCustomerIDAndTotalAmountOrders[customer.Id] - mapCustomerIDAndTotalAmountReceipts[customer.Id],
		}
	}
	return nil
}
