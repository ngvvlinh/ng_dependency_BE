// +build !generator

// Code generated by generator api. DO NOT EDIT.

package fbcustomerconversationsearch

import (
	context "context"

	capi "o.o/capi"
)

type QueryBus struct{ bus capi.Bus }

func NewQueryBus(bus capi.Bus) QueryBus { return QueryBus{bus} }

func (b QueryBus) Dispatch(ctx context.Context, msg interface{ query() }) error {
	return b.bus.Dispatch(ctx, msg)
}

type ListFbExternalCommentSearchQuery struct {
	PageIDs     []string
	ExternalMsg string

	Result []*FbExternalCommentSearch `json:"-"`
}

func (h QueryServiceHandler) HandleListFbExternalCommentSearch(ctx context.Context, msg *ListFbExternalCommentSearchQuery) (err error) {
	msg.Result, err = h.inner.ListFbExternalCommentSearch(msg.GetArgs(ctx))
	return err
}

type ListFbExternalConversationSearchQuery struct {
	PageIDs     []string
	ExtUserName string

	Result []*FbCustomerConversationSearch `json:"-"`
}

func (h QueryServiceHandler) HandleListFbExternalConversationSearch(ctx context.Context, msg *ListFbExternalConversationSearchQuery) (err error) {
	msg.Result, err = h.inner.ListFbExternalConversationSearch(msg.GetArgs(ctx))
	return err
}

type ListFbExternalMessageSearchQuery struct {
	PageIDs     []string
	ExternalMsg string

	Result []*FbExternalMessageSearch `json:"-"`
}

func (h QueryServiceHandler) HandleListFbExternalMessageSearch(ctx context.Context, msg *ListFbExternalMessageSearchQuery) (err error) {
	msg.Result, err = h.inner.ListFbExternalMessageSearch(msg.GetArgs(ctx))
	return err
}

// implement interfaces

func (q *ListFbExternalCommentSearchQuery) query()      {}
func (q *ListFbExternalConversationSearchQuery) query() {}
func (q *ListFbExternalMessageSearchQuery) query()      {}

// implement conversion

func (q *ListFbExternalCommentSearchQuery) GetArgs(ctx context.Context) (_ context.Context, pageIDs []string, externalMsg string) {
	return ctx,
		q.PageIDs,
		q.ExternalMsg
}

func (q *ListFbExternalConversationSearchQuery) GetArgs(ctx context.Context) (_ context.Context, pageIDs []string, extUserName string) {
	return ctx,
		q.PageIDs,
		q.ExtUserName
}

func (q *ListFbExternalMessageSearchQuery) GetArgs(ctx context.Context) (_ context.Context, pageIDs []string, externalMsg string) {
	return ctx,
		q.PageIDs,
		q.ExternalMsg
}

// implement dispatching

type QueryServiceHandler struct {
	inner QueryService
}

func NewQueryServiceHandler(service QueryService) QueryServiceHandler {
	return QueryServiceHandler{service}
}

func (h QueryServiceHandler) RegisterHandlers(b interface {
	capi.Bus
	AddHandler(handler interface{})
}) QueryBus {
	b.AddHandler(h.HandleListFbExternalCommentSearch)
	b.AddHandler(h.HandleListFbExternalConversationSearch)
	b.AddHandler(h.HandleListFbExternalMessageSearch)
	return QueryBus{b}
}
