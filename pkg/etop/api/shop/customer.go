package shop

import (
	"context"

	"etop.vn/api/main/etop"
	"etop.vn/api/main/ordering"
	"etop.vn/api/main/receipting"
	pbcm "etop.vn/api/pb/common"
	pbshop "etop.vn/api/pb/etop/shop"
	"etop.vn/api/shopping/customering"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmapi"
	"etop.vn/backend/pkg/etop/api/convertpb"
	. "etop.vn/capi/dot"
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
	cmd := &customering.CreateCustomerCommand{
		ShopID:   r.Context.Shop.ID,
		FullName: r.FullName,
		Gender:   r.Gender,
		Type:     customering.CustomerType(r.Type),
		Birthday: r.Birthday,
		Note:     r.Note,
		Phone:    r.Phone,
		Email:    r.Email,
	}
	err := customerAggr.Dispatch(ctx, cmd)
	if err != nil {
		return err
	}

	r.Result = convertpb.PbCustomer(cmd.Result)
	return nil
}

func (s *CustomerService) UpdateCustomer(ctx context.Context, r *UpdateCustomerEndpoint) error {
	cmd := &customering.UpdateCustomerCommand{
		ID:       r.Id,
		ShopID:   r.Context.Shop.ID,
		FullName: PString(r.FullName),
		Gender:   PString(r.Gender),
		Birthday: PString(r.Birthday),
		Note:     PString(r.Note),
		Phone:    PString(r.Phone),
		Email:    PString(r.Email),
	}
	if r.Type != nil {
		cmd.Type = customering.CustomerType(*r.Type)
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
		Status: int32(r.Status),
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
	r.Result = &pbcm.DeletedResponse{Deleted: int32(cmd.Result)}
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
	if err := s.listLiabilities(ctx, r.Context.Shop.ID, []*pbshop.Customer{r.Result}); err != nil {
		return err
	}
	return nil
}

func (s *CustomerService) GetCustomers(ctx context.Context, r *GetCustomersEndpoint) error {
	paging := cmapi.CMPaging(r.Paging)
	query := &customering.ListCustomersQuery{
		ShopID:  r.Context.Shop.ID,
		Paging:  *paging,
		Filters: cmapi.ToFilters(r.Filters),
	}
	if err := customerQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = &pbshop.CustomersResponse{
		Customers: convertpb.PbCustomers(query.Result.Customers),
		Paging:    cmapi.PbPageInfo(paging, query.Result.Count),
	}
	if err := s.listLiabilities(ctx, r.Context.Shop.ID, r.Result.Customers); err != nil {
		return err
	}
	return nil
}

func (s *CustomerService) GetCustomersByIDs(ctx context.Context, r *GetCustomersByIDsEndpoint) error {
	query := &customering.ListCustomersByIDsQuery{
		IDs:    r.Ids,
		ShopID: r.Context.Shop.ID,
	}
	if err := customerQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = &pbshop.CustomersResponse{
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
		Name: r.Name,
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
	q.Result = &pbshop.CustomerGroupsResponse{
		Paging:         cmapi.PbPageInfo(paging, query.Result.Count),
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
		GroupID:     r.GroupId,
		CustomerIDs: r.CustomerIds,
	}
	if err := customerAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.UpdatedResponse{Updated: int32(cmd.Result)}
	return nil
}

func (s *CustomerService) RemoveCustomersFromGroup(ctx context.Context, r *RemoveCustomersFromGroupEndpoint) error {
	cmd := &customering.RemoveCustomersFromGroupCommand{
		GroupID:     r.GroupId,
		CustomerIDs: r.CustomerIds,
	}
	if err := customerAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.RemovedResponse{Removed: int32(cmd.Result)}
	return nil
}

func (s *CustomerService) listLiabilities(ctx context.Context, shopID int64, customers []*pbshop.Customer) error {
	var customerIDs []int64
	mapCustomerIDAndNumberOfOrders := make(map[int64]int)
	mapCustomerIDAndTotalAmountOrders := make(map[int64]int64)
	mapCustomerIDAndTotalAmountReceipts := make(map[int64]int64)

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
		mapCustomerIDAndTotalAmountOrders[order.CustomerID] += int64(order.TotalAmount)
	}

	getReceiptsByCustomerIDs := &receipting.ListReceiptsByTraderIDsAndStatusesQuery{
		ShopID:    shopID,
		TraderIDs: customerIDs,
		Statuses:  []etop.Status3{etop.S3Positive},
	}
	if err := receiptQuery.Dispatch(ctx, getReceiptsByCustomerIDs); err != nil {
		return err
	}
	for _, receipt := range getReceiptsByCustomerIDs.Result.Receipts {
		switch receipt.Type {
		case receipting.ReceiptTypeReceipt:
			mapCustomerIDAndTotalAmountReceipts[receipt.TraderID] += int64(receipt.Amount)
		case receipting.ReceiptTypePayment:
			mapCustomerIDAndTotalAmountReceipts[receipt.TraderID] -= int64(receipt.Amount)
		}
	}

	for _, customer := range customers {
		customer.Liability = &pbshop.CustomerLiability{
			TotalOrders:    int32(mapCustomerIDAndNumberOfOrders[customer.Id]),
			TotalAmount:    mapCustomerIDAndTotalAmountOrders[customer.Id],
			ReceivedAmount: mapCustomerIDAndTotalAmountReceipts[customer.Id],
			Liability:      mapCustomerIDAndTotalAmountOrders[customer.Id] - mapCustomerIDAndTotalAmountReceipts[customer.Id],
		}
	}

	return nil
}
