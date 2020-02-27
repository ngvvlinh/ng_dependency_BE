// +build !generator

// Code generated by generator apix. DO NOT EDIT.

package partner_proto

import (
	context "context"
	fmt "fmt"
	http "net/http"

	externaltypes "etop.vn/api/top/external/types"
	common "etop.vn/api/top/types/common"
	capi "etop.vn/capi"
	httprpc "etop.vn/capi/httprpc"
)

type Server interface {
	http.Handler
	PathPrefix() string
}

type CustomerAddressServiceServer struct {
	inner CustomerAddressService
}

func NewCustomerAddressServiceServer(svc CustomerAddressService) Server {
	return &CustomerAddressServiceServer{
		inner: svc,
	}
}

const CustomerAddressServicePathPrefix = "/partner.CustomerAddress/"

func (s *CustomerAddressServiceServer) PathPrefix() string {
	return CustomerAddressServicePathPrefix
}

func (s *CustomerAddressServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	serve, err := httprpc.ParseRequestHeader(req)
	if err != nil {
		httprpc.WriteError(ctx, resp, err)
		return
	}
	reqMsg, exec, err := s.parseRoute(req.URL.Path)
	if err != nil {
		httprpc.WriteError(ctx, resp, err)
		return
	}
	serve(ctx, resp, req, reqMsg, exec)
}

func (s *CustomerAddressServiceServer) parseRoute(path string) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/partner.CustomerAddress/CreateAddress":
		msg := &externaltypes.CreateCustomerAddressRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.CreateAddress(ctx, msg)
		}
		return msg, fn, nil
	case "/partner.CustomerAddress/DeleteAddress":
		msg := &common.IDRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.DeleteAddress(ctx, msg)
		}
		return msg, fn, nil
	case "/partner.CustomerAddress/GetAddress":
		msg := &externaltypes.OrderIDRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.GetAddress(ctx, msg)
		}
		return msg, fn, nil
	case "/partner.CustomerAddress/ListAddresses":
		msg := &externaltypes.ListCustomerAddressesRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.ListAddresses(ctx, msg)
		}
		return msg, fn, nil
	case "/partner.CustomerAddress/UpdateAddress":
		msg := &externaltypes.UpdateCustomerAddressRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.UpdateAddress(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}

type CustomerGroupRelationshipServiceServer struct {
	inner CustomerGroupRelationshipService
}

func NewCustomerGroupRelationshipServiceServer(svc CustomerGroupRelationshipService) Server {
	return &CustomerGroupRelationshipServiceServer{
		inner: svc,
	}
}

const CustomerGroupRelationshipServicePathPrefix = "/partner.CustomerGroupRelationship/"

func (s *CustomerGroupRelationshipServiceServer) PathPrefix() string {
	return CustomerGroupRelationshipServicePathPrefix
}

func (s *CustomerGroupRelationshipServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	serve, err := httprpc.ParseRequestHeader(req)
	if err != nil {
		httprpc.WriteError(ctx, resp, err)
		return
	}
	reqMsg, exec, err := s.parseRoute(req.URL.Path)
	if err != nil {
		httprpc.WriteError(ctx, resp, err)
		return
	}
	serve(ctx, resp, req, reqMsg, exec)
}

func (s *CustomerGroupRelationshipServiceServer) parseRoute(path string) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/partner.CustomerGroupRelationship/CreateRelationship":
		msg := &externaltypes.AddCustomerRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.CreateRelationship(ctx, msg)
		}
		return msg, fn, nil
	case "/partner.CustomerGroupRelationship/DeleteRelationship":
		msg := &externaltypes.RemoveCustomerRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.DeleteRelationship(ctx, msg)
		}
		return msg, fn, nil
	case "/partner.CustomerGroupRelationship/ListRelationships":
		msg := &externaltypes.ListCustomerGroupRelationshipsRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.ListRelationships(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}

type CustomerGroupServiceServer struct {
	inner CustomerGroupService
}

func NewCustomerGroupServiceServer(svc CustomerGroupService) Server {
	return &CustomerGroupServiceServer{
		inner: svc,
	}
}

const CustomerGroupServicePathPrefix = "/partner.CustomerGroup/"

func (s *CustomerGroupServiceServer) PathPrefix() string {
	return CustomerGroupServicePathPrefix
}

func (s *CustomerGroupServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	serve, err := httprpc.ParseRequestHeader(req)
	if err != nil {
		httprpc.WriteError(ctx, resp, err)
		return
	}
	reqMsg, exec, err := s.parseRoute(req.URL.Path)
	if err != nil {
		httprpc.WriteError(ctx, resp, err)
		return
	}
	serve(ctx, resp, req, reqMsg, exec)
}

func (s *CustomerGroupServiceServer) parseRoute(path string) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/partner.CustomerGroup/CreateGroup":
		msg := &externaltypes.CreateCustomerGroupRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.CreateGroup(ctx, msg)
		}
		return msg, fn, nil
	case "/partner.CustomerGroup/DeleteGroup":
		msg := &common.IDRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.DeleteGroup(ctx, msg)
		}
		return msg, fn, nil
	case "/partner.CustomerGroup/GetGroup":
		msg := &common.IDRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.GetGroup(ctx, msg)
		}
		return msg, fn, nil
	case "/partner.CustomerGroup/ListGroups":
		msg := &externaltypes.ListCustomerGroupsRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.ListGroups(ctx, msg)
		}
		return msg, fn, nil
	case "/partner.CustomerGroup/UpdateGroup":
		msg := &externaltypes.UpdateCustomerGroupRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.UpdateGroup(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}

type CustomerServiceServer struct {
	inner CustomerService
}

func NewCustomerServiceServer(svc CustomerService) Server {
	return &CustomerServiceServer{
		inner: svc,
	}
}

const CustomerServicePathPrefix = "/partner.Customer/"

func (s *CustomerServiceServer) PathPrefix() string {
	return CustomerServicePathPrefix
}

func (s *CustomerServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	serve, err := httprpc.ParseRequestHeader(req)
	if err != nil {
		httprpc.WriteError(ctx, resp, err)
		return
	}
	reqMsg, exec, err := s.parseRoute(req.URL.Path)
	if err != nil {
		httprpc.WriteError(ctx, resp, err)
		return
	}
	serve(ctx, resp, req, reqMsg, exec)
}

func (s *CustomerServiceServer) parseRoute(path string) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/partner.Customer/CreateCustomer":
		msg := &externaltypes.CreateCustomerRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.CreateCustomer(ctx, msg)
		}
		return msg, fn, nil
	case "/partner.Customer/DeleteCustomer":
		msg := &externaltypes.DeleteCustomerRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.DeleteCustomer(ctx, msg)
		}
		return msg, fn, nil
	case "/partner.Customer/GetCustomer":
		msg := &externaltypes.GetCustomerRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.GetCustomer(ctx, msg)
		}
		return msg, fn, nil
	case "/partner.Customer/ListCustomers":
		msg := &externaltypes.ListCustomersRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.ListCustomers(ctx, msg)
		}
		return msg, fn, nil
	case "/partner.Customer/UpdateCustomer":
		msg := &externaltypes.UpdateCustomerRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.UpdateCustomer(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}

type FulfillmentServiceServer struct {
	inner FulfillmentService
}

func NewFulfillmentServiceServer(svc FulfillmentService) Server {
	return &FulfillmentServiceServer{
		inner: svc,
	}
}

const FulfillmentServicePathPrefix = "/partner.Fulfillment/"

func (s *FulfillmentServiceServer) PathPrefix() string {
	return FulfillmentServicePathPrefix
}

func (s *FulfillmentServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	serve, err := httprpc.ParseRequestHeader(req)
	if err != nil {
		httprpc.WriteError(ctx, resp, err)
		return
	}
	reqMsg, exec, err := s.parseRoute(req.URL.Path)
	if err != nil {
		httprpc.WriteError(ctx, resp, err)
		return
	}
	serve(ctx, resp, req, reqMsg, exec)
}

func (s *FulfillmentServiceServer) parseRoute(path string) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/partner.Fulfillment/CancelFulfillment":
		msg := &externaltypes.CancelFulfillmentRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.CancelFulfillment(ctx, msg)
		}
		return msg, fn, nil
	case "/partner.Fulfillment/CreateFulfillment":
		msg := &externaltypes.CreateFulfillmentRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.CreateFulfillment(ctx, msg)
		}
		return msg, fn, nil
	case "/partner.Fulfillment/GetFulfillment":
		msg := &externaltypes.FulfillmentIDRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.GetFulfillment(ctx, msg)
		}
		return msg, fn, nil
	case "/partner.Fulfillment/ListFulfillments":
		msg := &externaltypes.ListFulfillmentsRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.ListFulfillments(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}

type HistoryServiceServer struct {
	inner HistoryService
}

func NewHistoryServiceServer(svc HistoryService) Server {
	return &HistoryServiceServer{
		inner: svc,
	}
}

const HistoryServicePathPrefix = "/partner.History/"

func (s *HistoryServiceServer) PathPrefix() string {
	return HistoryServicePathPrefix
}

func (s *HistoryServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	serve, err := httprpc.ParseRequestHeader(req)
	if err != nil {
		httprpc.WriteError(ctx, resp, err)
		return
	}
	reqMsg, exec, err := s.parseRoute(req.URL.Path)
	if err != nil {
		httprpc.WriteError(ctx, resp, err)
		return
	}
	serve(ctx, resp, req, reqMsg, exec)
}

func (s *HistoryServiceServer) parseRoute(path string) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/partner.History/GetChanges":
		msg := &common.Empty{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.GetChanges(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}

type InventoryServiceServer struct {
	inner InventoryService
}

func NewInventoryServiceServer(svc InventoryService) Server {
	return &InventoryServiceServer{
		inner: svc,
	}
}

const InventoryServicePathPrefix = "/partner.Inventory/"

func (s *InventoryServiceServer) PathPrefix() string {
	return InventoryServicePathPrefix
}

func (s *InventoryServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	serve, err := httprpc.ParseRequestHeader(req)
	if err != nil {
		httprpc.WriteError(ctx, resp, err)
		return
	}
	reqMsg, exec, err := s.parseRoute(req.URL.Path)
	if err != nil {
		httprpc.WriteError(ctx, resp, err)
		return
	}
	serve(ctx, resp, req, reqMsg, exec)
}

func (s *InventoryServiceServer) parseRoute(path string) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/partner.Inventory/ListInventoryLevels":
		msg := &externaltypes.ListInventoryLevelsRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.ListInventoryLevels(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}

type MiscServiceServer struct {
	inner MiscService
}

func NewMiscServiceServer(svc MiscService) Server {
	return &MiscServiceServer{
		inner: svc,
	}
}

const MiscServicePathPrefix = "/partner.Misc/"

func (s *MiscServiceServer) PathPrefix() string {
	return MiscServicePathPrefix
}

func (s *MiscServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	serve, err := httprpc.ParseRequestHeader(req)
	if err != nil {
		httprpc.WriteError(ctx, resp, err)
		return
	}
	reqMsg, exec, err := s.parseRoute(req.URL.Path)
	if err != nil {
		httprpc.WriteError(ctx, resp, err)
		return
	}
	serve(ctx, resp, req, reqMsg, exec)
}

func (s *MiscServiceServer) parseRoute(path string) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/partner.Misc/CurrentAccount":
		msg := &common.Empty{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.CurrentAccount(ctx, msg)
		}
		return msg, fn, nil
	case "/partner.Misc/GetLocationList":
		msg := &common.Empty{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.GetLocationList(ctx, msg)
		}
		return msg, fn, nil
	case "/partner.Misc/VersionInfo":
		msg := &common.Empty{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.VersionInfo(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}

type OrderServiceServer struct {
	inner OrderService
}

func NewOrderServiceServer(svc OrderService) Server {
	return &OrderServiceServer{
		inner: svc,
	}
}

const OrderServicePathPrefix = "/partner.Order/"

func (s *OrderServiceServer) PathPrefix() string {
	return OrderServicePathPrefix
}

func (s *OrderServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	serve, err := httprpc.ParseRequestHeader(req)
	if err != nil {
		httprpc.WriteError(ctx, resp, err)
		return
	}
	reqMsg, exec, err := s.parseRoute(req.URL.Path)
	if err != nil {
		httprpc.WriteError(ctx, resp, err)
		return
	}
	serve(ctx, resp, req, reqMsg, exec)
}

func (s *OrderServiceServer) parseRoute(path string) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/partner.Order/CancelOrder":
		msg := &externaltypes.CancelOrderRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.CancelOrder(ctx, msg)
		}
		return msg, fn, nil
	case "/partner.Order/ConfirmOrder":
		msg := &externaltypes.ConfirmOrderRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.ConfirmOrder(ctx, msg)
		}
		return msg, fn, nil
	case "/partner.Order/CreateOrder":
		msg := &externaltypes.CreateOrderRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.CreateOrder(ctx, msg)
		}
		return msg, fn, nil
	case "/partner.Order/GetOrder":
		msg := &externaltypes.OrderIDRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.GetOrder(ctx, msg)
		}
		return msg, fn, nil
	case "/partner.Order/ListOrders":
		msg := &externaltypes.ListOrdersRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.ListOrders(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}

type ProductCollectionRelationshipServiceServer struct {
	inner ProductCollectionRelationshipService
}

func NewProductCollectionRelationshipServiceServer(svc ProductCollectionRelationshipService) Server {
	return &ProductCollectionRelationshipServiceServer{
		inner: svc,
	}
}

const ProductCollectionRelationshipServicePathPrefix = "/partner.ProductCollectionRelationship/"

func (s *ProductCollectionRelationshipServiceServer) PathPrefix() string {
	return ProductCollectionRelationshipServicePathPrefix
}

func (s *ProductCollectionRelationshipServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	serve, err := httprpc.ParseRequestHeader(req)
	if err != nil {
		httprpc.WriteError(ctx, resp, err)
		return
	}
	reqMsg, exec, err := s.parseRoute(req.URL.Path)
	if err != nil {
		httprpc.WriteError(ctx, resp, err)
		return
	}
	serve(ctx, resp, req, reqMsg, exec)
}

func (s *ProductCollectionRelationshipServiceServer) parseRoute(path string) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/partner.ProductCollectionRelationship/CreateRelationship":
		msg := &externaltypes.CreateProductCollectionRelationshipRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.CreateRelationship(ctx, msg)
		}
		return msg, fn, nil
	case "/partner.ProductCollectionRelationship/DeleteRelationship":
		msg := &externaltypes.RemoveProductCollectionRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.DeleteRelationship(ctx, msg)
		}
		return msg, fn, nil
	case "/partner.ProductCollectionRelationship/ListRelationships":
		msg := &externaltypes.ListProductCollectionRelationshipsRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.ListRelationships(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}

type ProductCollectionServiceServer struct {
	inner ProductCollectionService
}

func NewProductCollectionServiceServer(svc ProductCollectionService) Server {
	return &ProductCollectionServiceServer{
		inner: svc,
	}
}

const ProductCollectionServicePathPrefix = "/partner.ProductCollection/"

func (s *ProductCollectionServiceServer) PathPrefix() string {
	return ProductCollectionServicePathPrefix
}

func (s *ProductCollectionServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	serve, err := httprpc.ParseRequestHeader(req)
	if err != nil {
		httprpc.WriteError(ctx, resp, err)
		return
	}
	reqMsg, exec, err := s.parseRoute(req.URL.Path)
	if err != nil {
		httprpc.WriteError(ctx, resp, err)
		return
	}
	serve(ctx, resp, req, reqMsg, exec)
}

func (s *ProductCollectionServiceServer) parseRoute(path string) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/partner.ProductCollection/CreateCollection":
		msg := &externaltypes.CreateCollectionRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.CreateCollection(ctx, msg)
		}
		return msg, fn, nil
	case "/partner.ProductCollection/DeleteCollection":
		msg := &externaltypes.GetCollectionRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.DeleteCollection(ctx, msg)
		}
		return msg, fn, nil
	case "/partner.ProductCollection/GetCollection":
		msg := &externaltypes.GetCollectionRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.GetCollection(ctx, msg)
		}
		return msg, fn, nil
	case "/partner.ProductCollection/ListCollections":
		msg := &externaltypes.ListCollectionsRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.ListCollections(ctx, msg)
		}
		return msg, fn, nil
	case "/partner.ProductCollection/UpdateCollection":
		msg := &externaltypes.UpdateCollectionRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.UpdateCollection(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}

type ProductServiceServer struct {
	inner ProductService
}

func NewProductServiceServer(svc ProductService) Server {
	return &ProductServiceServer{
		inner: svc,
	}
}

const ProductServicePathPrefix = "/partner.Product/"

func (s *ProductServiceServer) PathPrefix() string {
	return ProductServicePathPrefix
}

func (s *ProductServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	serve, err := httprpc.ParseRequestHeader(req)
	if err != nil {
		httprpc.WriteError(ctx, resp, err)
		return
	}
	reqMsg, exec, err := s.parseRoute(req.URL.Path)
	if err != nil {
		httprpc.WriteError(ctx, resp, err)
		return
	}
	serve(ctx, resp, req, reqMsg, exec)
}

func (s *ProductServiceServer) parseRoute(path string) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/partner.Product/CreateProduct":
		msg := &externaltypes.CreateProductRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.CreateProduct(ctx, msg)
		}
		return msg, fn, nil
	case "/partner.Product/DeleteProduct":
		msg := &externaltypes.GetProductRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.DeleteProduct(ctx, msg)
		}
		return msg, fn, nil
	case "/partner.Product/GetProduct":
		msg := &externaltypes.GetProductRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.GetProduct(ctx, msg)
		}
		return msg, fn, nil
	case "/partner.Product/ListProducts":
		msg := &externaltypes.ListProductsRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.ListProducts(ctx, msg)
		}
		return msg, fn, nil
	case "/partner.Product/UpdateProduct":
		msg := &externaltypes.UpdateProductRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.UpdateProduct(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}

type ShipmentConnectionServiceServer struct {
	inner ShipmentConnectionService
}

func NewShipmentConnectionServiceServer(svc ShipmentConnectionService) Server {
	return &ShipmentConnectionServiceServer{
		inner: svc,
	}
}

const ShipmentConnectionServicePathPrefix = "/partner.ShipmentConnection/"

func (s *ShipmentConnectionServiceServer) PathPrefix() string {
	return ShipmentConnectionServicePathPrefix
}

func (s *ShipmentConnectionServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	serve, err := httprpc.ParseRequestHeader(req)
	if err != nil {
		httprpc.WriteError(ctx, resp, err)
		return
	}
	reqMsg, exec, err := s.parseRoute(req.URL.Path)
	if err != nil {
		httprpc.WriteError(ctx, resp, err)
		return
	}
	serve(ctx, resp, req, reqMsg, exec)
}

func (s *ShipmentConnectionServiceServer) parseRoute(path string) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/partner.ShipmentConnection/CreateConnection":
		msg := &CreateConnectionRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.CreateConnection(ctx, msg)
		}
		return msg, fn, nil
	case "/partner.ShipmentConnection/DeleteConnection":
		msg := &common.IDRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.DeleteConnection(ctx, msg)
		}
		return msg, fn, nil
	case "/partner.ShipmentConnection/GetConnections":
		msg := &common.Empty{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.GetConnections(ctx, msg)
		}
		return msg, fn, nil
	case "/partner.ShipmentConnection/UpdateConnection":
		msg := &UpdateConnectionRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.UpdateConnection(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}

type ShipmentServiceServer struct {
	inner ShipmentService
}

func NewShipmentServiceServer(svc ShipmentService) Server {
	return &ShipmentServiceServer{
		inner: svc,
	}
}

const ShipmentServicePathPrefix = "/partner.Shipment/"

func (s *ShipmentServiceServer) PathPrefix() string {
	return ShipmentServicePathPrefix
}

func (s *ShipmentServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	serve, err := httprpc.ParseRequestHeader(req)
	if err != nil {
		httprpc.WriteError(ctx, resp, err)
		return
	}
	reqMsg, exec, err := s.parseRoute(req.URL.Path)
	if err != nil {
		httprpc.WriteError(ctx, resp, err)
		return
	}
	serve(ctx, resp, req, reqMsg, exec)
}

func (s *ShipmentServiceServer) parseRoute(path string) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/partner.Shipment/UpdateFulfillment":
		msg := &UpdateFulfillmentRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.UpdateFulfillment(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}

type ShippingServiceServer struct {
	inner ShippingService
}

func NewShippingServiceServer(svc ShippingService) Server {
	return &ShippingServiceServer{
		inner: svc,
	}
}

const ShippingServicePathPrefix = "/partner.Shipping/"

func (s *ShippingServiceServer) PathPrefix() string {
	return ShippingServicePathPrefix
}

func (s *ShippingServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	serve, err := httprpc.ParseRequestHeader(req)
	if err != nil {
		httprpc.WriteError(ctx, resp, err)
		return
	}
	reqMsg, exec, err := s.parseRoute(req.URL.Path)
	if err != nil {
		httprpc.WriteError(ctx, resp, err)
		return
	}
	serve(ctx, resp, req, reqMsg, exec)
}

func (s *ShippingServiceServer) parseRoute(path string) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/partner.Shipping/CancelOrder":
		msg := &externaltypes.CancelOrderRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.CancelOrder(ctx, msg)
		}
		return msg, fn, nil
	case "/partner.Shipping/CreateAndConfirmOrder":
		msg := &externaltypes.CreateAndConfirmOrderRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.CreateAndConfirmOrder(ctx, msg)
		}
		return msg, fn, nil
	case "/partner.Shipping/GetFulfillment":
		msg := &externaltypes.FulfillmentIDRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.GetFulfillment(ctx, msg)
		}
		return msg, fn, nil
	case "/partner.Shipping/GetOrder":
		msg := &externaltypes.OrderIDRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.GetOrder(ctx, msg)
		}
		return msg, fn, nil
	case "/partner.Shipping/GetShippingServices":
		msg := &externaltypes.GetShippingServicesRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.GetShippingServices(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}

type ShopServiceServer struct {
	inner ShopService
}

func NewShopServiceServer(svc ShopService) Server {
	return &ShopServiceServer{
		inner: svc,
	}
}

const ShopServicePathPrefix = "/partner.Shop/"

func (s *ShopServiceServer) PathPrefix() string {
	return ShopServicePathPrefix
}

func (s *ShopServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	serve, err := httprpc.ParseRequestHeader(req)
	if err != nil {
		httprpc.WriteError(ctx, resp, err)
		return
	}
	reqMsg, exec, err := s.parseRoute(req.URL.Path)
	if err != nil {
		httprpc.WriteError(ctx, resp, err)
		return
	}
	serve(ctx, resp, req, reqMsg, exec)
}

func (s *ShopServiceServer) parseRoute(path string) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/partner.Shop/AuthorizeShop":
		msg := &AuthorizeShopRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.AuthorizeShop(ctx, msg)
		}
		return msg, fn, nil
	case "/partner.Shop/CurrentShop":
		msg := &common.Empty{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.CurrentShop(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}

type VariantServiceServer struct {
	inner VariantService
}

func NewVariantServiceServer(svc VariantService) Server {
	return &VariantServiceServer{
		inner: svc,
	}
}

const VariantServicePathPrefix = "/partner.Variant/"

func (s *VariantServiceServer) PathPrefix() string {
	return VariantServicePathPrefix
}

func (s *VariantServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	serve, err := httprpc.ParseRequestHeader(req)
	if err != nil {
		httprpc.WriteError(ctx, resp, err)
		return
	}
	reqMsg, exec, err := s.parseRoute(req.URL.Path)
	if err != nil {
		httprpc.WriteError(ctx, resp, err)
		return
	}
	serve(ctx, resp, req, reqMsg, exec)
}

func (s *VariantServiceServer) parseRoute(path string) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/partner.Variant/CreateVariant":
		msg := &externaltypes.CreateVariantRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.CreateVariant(ctx, msg)
		}
		return msg, fn, nil
	case "/partner.Variant/DeleteVariant":
		msg := &externaltypes.GetVariantRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.DeleteVariant(ctx, msg)
		}
		return msg, fn, nil
	case "/partner.Variant/GetVariant":
		msg := &externaltypes.GetVariantRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.GetVariant(ctx, msg)
		}
		return msg, fn, nil
	case "/partner.Variant/ListVariants":
		msg := &externaltypes.ListVariantsRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.ListVariants(ctx, msg)
		}
		return msg, fn, nil
	case "/partner.Variant/UpdateVariant":
		msg := &externaltypes.UpdateVariantRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.UpdateVariant(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}

type WebhookServiceServer struct {
	inner WebhookService
}

func NewWebhookServiceServer(svc WebhookService) Server {
	return &WebhookServiceServer{
		inner: svc,
	}
}

const WebhookServicePathPrefix = "/partner.Webhook/"

func (s *WebhookServiceServer) PathPrefix() string {
	return WebhookServicePathPrefix
}

func (s *WebhookServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	serve, err := httprpc.ParseRequestHeader(req)
	if err != nil {
		httprpc.WriteError(ctx, resp, err)
		return
	}
	reqMsg, exec, err := s.parseRoute(req.URL.Path)
	if err != nil {
		httprpc.WriteError(ctx, resp, err)
		return
	}
	serve(ctx, resp, req, reqMsg, exec)
}

func (s *WebhookServiceServer) parseRoute(path string) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/partner.Webhook/CreateWebhook":
		msg := &externaltypes.CreateWebhookRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.CreateWebhook(ctx, msg)
		}
		return msg, fn, nil
	case "/partner.Webhook/DeleteWebhook":
		msg := &externaltypes.DeleteWebhookRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.DeleteWebhook(ctx, msg)
		}
		return msg, fn, nil
	case "/partner.Webhook/GetWebhooks":
		msg := &common.Empty{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.GetWebhooks(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}
