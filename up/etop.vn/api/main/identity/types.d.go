// Code generated by gen-cmd-query. DO NOT EDIT.

package identity

import (
	context "context"
	unsafe "unsafe"
)

type GetShopByIDQuery struct {
	ID int64

	Result *Shop `json:"-"`
}

// implement conversion

func (q *GetShopByIDQuery) GetArgs() *GetShopByIDQueryArgs {
	return (*GetShopByIDQueryArgs)(unsafe.Pointer(q))
}

// implement dispatching

type QueryServiceHandler struct {
	inner QueryService
}

func NewQueryServiceHandler(service QueryService) QueryServiceHandler {
	return QueryServiceHandler{service}
}

func (h QueryServiceHandler) RegisterHandlers(b interface {
	AddHandler(handler interface{})
}) {
	b.AddHandler(h.HandleGetShopByID)
}

func (h QueryServiceHandler) HandleGetShopByID(ctx context.Context, query *GetShopByIDQuery) error {
	result, err := h.inner.GetShopByID(ctx, query.GetArgs())
	query.Result = result
	return err
}
