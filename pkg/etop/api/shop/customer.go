package shop

import (
	"context"
	"fmt"
	"time"

	"o.o/api/main/ordering"
	"o.o/api/main/receipting"
	"o.o/api/shopping/customering"
	"o.o/api/top/int/shop"
	pbcm "o.o/api/top/types/common"
	"o.o/api/top/types/etc/receipt_type"
	"o.o/api/top/types/etc/status3"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/capi/dot"
)

func init() {
	bus.AddHandlers("api",
		customerService.CreateCustomer,
		customerService.UpdateCustomer,
		customerService.DeleteCustomer,
		customerService.GetCustomer,
		customerService.GetCustomers,
		customerService.GetCustomersByIDs,
		customerService.GetCustomerDetails,
		customerService.BatchSetCustomersStatus,

		customerGroupService.CreateCustomerGroup,
		customerGroupService.GetCustomerGroup,
		customerGroupService.GetCustomerGroups,
		customerGroupService.UpdateCustomerGroup,
		customerService.AddCustomersToGroup,
		customerService.RemoveCustomersFromGroup,
	)
}

func (s *CustomerService) CreateCustomer(ctx context.Context, r *CreateCustomerEndpoint) error {
	key := fmt.Sprintf("CreateCustomer %v-%v-%v-%v-%v",
		r.Context.Shop.ID, r.Context.UserID, r.Phone, r.FullName, r.Email)
	res, err := idempgroup.DoAndWrap(ctx, key, 15*time.Second,
		func() (interface{}, error) {
			cmd := &customering.CreateCustomerCommand{
				ShopID:   r.Context.Shop.ID,
				FullName: r.FullName,
				Gender:   r.Gender,
				Type:     r.Type,
				Birthday: r.Birthday,
				Note:     r.Note,
				Phone:    r.Phone,
				Email:    r.Email,
			}
			if err := customerAggr.Dispatch(ctx, cmd); err != nil {
				return nil, err
			}
			r.Result = convertpb.PbCustomer(cmd.Result)
			return r, nil
		}, "tạo khách hàng")

	if err != nil {
		return err
	}
	r.Result = res.(*CreateCustomerEndpoint).Result
	return nil
}

func (s *CustomerService) UpdateCustomer(ctx context.Context, r *UpdateCustomerEndpoint) error {
	cmd := &customering.UpdateCustomerCommand{
		ID:       r.Id,
		ShopID:   r.Context.Shop.ID,
		FullName: r.FullName,
		Gender:   r.Gender,
		Birthday: r.Birthday,
		Note:     r.Note,
		Phone:    r.Phone,
		Email:    r.Email,
		Type:     r.Type.Apply(0),
	}
	err := customerAggr.Dispatch(ctx, cmd)
	if err != nil {
		return err
	}

	r.Result = convertpb.PbCustomer(cmd.Result)
	return nil
}

func (s *CustomerService) BatchSetCustomersStatus(ctx context.Context, r *BatchSetCustomersStatusEndpoint) error {
	cmd := &customering.BatchSetCustomersStatusCommand{
		IDs:    r.Ids,
		ShopID: r.Context.Shop.ID,
		Status: int(r.Status),
	}
	if err := customerAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.UpdatedResponse{Updated: cmd.Result.Updated}
	return nil
}

func (s *CustomerService) DeleteCustomer(ctx context.Context, r *DeleteCustomerEndpoint) error {
	cmd := &customering.DeleteCustomerCommand{
		ID:     r.Id,
		ShopID: r.Context.Shop.ID,
	}
	if err := customerAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.DeletedResponse{Deleted: cmd.Result}
	return nil
}

func (s *CustomerService) GetCustomer(ctx context.Context, r *GetCustomerEndpoint) error {
	query := &customering.GetCustomerByIDQuery{
		ID:     r.Id,
		ShopID: r.Context.Shop.ID,
	}
	if err := customerQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = convertpb.PbCustomer(query.Result)
	if err := s.listLiabilities(ctx, r.Context.Shop.ID, []*shop.Customer{r.Result}); err != nil {
		return err
	}
	return nil
}

func (s *CustomerService) GetCustomers(ctx context.Context, r *GetCustomersEndpoint) error {
	paging := cmapi.CMPaging(r.Paging)
	if !r.GetAll.Valid {
		r.GetAll = dot.Bool(true)
	}
	switch r.GetAll.Bool {
	case true:
		if err := getAllCustomers(ctx, paging, r); err != nil {
			return err
		}
	case false:
		customers, err := getCustomers(ctx, paging, r)
		if err != nil {
			return err
		}
		r.Result.Customers = customers
	}
	if err := s.listLiabilities(ctx, r.Context.Shop.ID, r.Result.Customers); err != nil {
		return err
	}
	return nil
}

func getAllCustomers(ctx context.Context, paging *cm.Paging, r *GetCustomersEndpoint) error {
	queryCustomerIndenpendent := &customering.GetCustomerIndependentQuery{}
	if err := customerQuery.Dispatch(ctx, queryCustomerIndenpendent); err != nil {
		return err
	}
	var customers []*shop.Customer
	customers = append(customers, convertpb.PbCustomer(queryCustomerIndenpendent.Result))

	if paging.Limit == 1 && paging.Offset == 0 {
		r.Result.Customers = customers
		return nil
	}
	if paging.Offset == 0 {
		paging.Limit--
		cts, err := getCustomers(ctx, paging, r)
		if err != nil {
			return err
		}
		customers = append(customers, cts...)
		r.Result.Customers = customers
		r.Result.Paging.Limit++
	} else {
		paging.Offset--
		_, err := getCustomers(ctx, paging, r)
		if err != nil {
			return err
		}
	}
	return nil
}

func getCustomers(ctx context.Context, paging *cm.Paging, r *GetCustomersEndpoint) ([]*shop.Customer, error) {
	query := &customering.ListCustomersQuery{
		ShopID:  r.Context.Shop.ID,
		Paging:  *paging,
		Filters: cmapi.ToFilters(r.Filters),
	}
	if err := customerQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	r.Result = &shop.CustomersResponse{
		Customers: convertpb.PbCustomers(query.Result.Customers),
		Paging:    cmapi.PbPageInfo(paging),
	}
	return convertpb.PbCustomers(query.Result.Customers), nil
}

func (s *CustomerService) GetCustomersByIDs(ctx context.Context, r *GetCustomersByIDsEndpoint) error {
	query := &customering.ListCustomersByIDsQuery{
		IDs:    r.Ids,
		ShopID: r.Context.Shop.ID,
	}
	if err := customerQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = &shop.CustomersResponse{
		Customers: convertpb.PbCustomers(query.Result.Customers),
	}
	if err := s.listLiabilities(ctx, r.Context.Shop.ID, r.Result.Customers); err != nil {
		return err
	}
	return nil
}

func (s *CustomerService) GetCustomerDetails(ctx context.Context, r *GetCustomerDetailsEndpoint) error {
	return cm.ErrTODO
}

func (s *CustomerGroupService) CreateCustomerGroup(ctx context.Context, r *CreateCustomerGroupEndpoint) error {
	cmd := &customering.CreateCustomerGroupCommand{
		ShopID: r.Context.Shop.ID,
		Name:   r.Name,
	}
	if err := customerAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = convertpb.PbCustopmerGroup(cmd.Result)
	return nil
}

func (s *CustomerGroupService) GetCustomerGroup(ctx context.Context, q *GetCustomerGroupEndpoint) error {
	query := &customering.GetCustomerGroupQuery{
		ID: q.Id,
	}
	if err := customerQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = convertpb.PbCustopmerGroup(query.Result)
	return nil
}

func (s *CustomerGroupService) GetCustomerGroups(ctx context.Context, q *GetCustomerGroupsEndpoint) error {
	paging := cmapi.CMPaging(q.Paging)
	query := &customering.ListCustomerGroupsQuery{
		Paging:  *paging,
		Filters: cmapi.ToFilters(q.Filters),
	}
	if err := customerQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &shop.CustomerGroupsResponse{
		Paging:         cmapi.PbPageInfo(paging),
		CustomerGroups: convertpb.PbCustomerGroups(query.Result.CustomerGroups),
	}
	return nil
}

func (s *CustomerGroupService) UpdateCustomerGroup(ctx context.Context, r *UpdateCustomerGroupEndpoint) error {
	cmd := &customering.UpdateCustomerGroupCommand{
		ID:   r.GroupId,
		Name: r.Name,
	}
	if err := customerAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = convertpb.PbCustopmerGroup(cmd.Result)
	return nil
}

func (s *CustomerService) AddCustomersToGroup(ctx context.Context, r *AddCustomersToGroupEndpoint) error {
	cmd := &customering.AddCustomersToGroupCommand{
		ShopID:      r.Context.Shop.ID,
		GroupID:     r.GroupId,
		CustomerIDs: r.CustomerIds,
	}
	if err := customerAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.UpdatedResponse{Updated: cmd.Result}
	return nil
}

func (s *CustomerService) RemoveCustomersFromGroup(ctx context.Context, r *RemoveCustomersFromGroupEndpoint) error {
	cmd := &customering.RemoveCustomersFromGroupCommand{
		ShopID:      r.Context.Shop.ID,
		GroupID:     r.GroupId,
		CustomerIDs: r.CustomerIds,
	}
	if err := customerAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.RemovedResponse{Removed: cmd.Result}
	return nil
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
	if err := orderQuery.Dispatch(ctx, getOrdersByCustomerIDs); err != nil {
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
	if err := receiptQuery.Dispatch(ctx, getReceiptsByCustomerIDs); err != nil {
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
