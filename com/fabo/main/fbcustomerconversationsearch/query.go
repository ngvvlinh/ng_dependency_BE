package fbcustomerconversationsearch

import (
	"context"

	"o.o/api/fabo/fbcustomerconversationsearch"
	"o.o/backend/com/fabo/main/fbcustomerconversationsearch/sqlstore"
	com "o.o/backend/com/main"
	"o.o/backend/pkg/common/bus"
)

var _ fbcustomerconversationsearch.QueryService = &FbSearchService{}

type FbSearchService struct {
	fbCommentSearchStore        sqlstore.FbExternalCommentSearchStoreFactory
	fbMessageSearchStore        sqlstore.FbExternalMessageSearchStoreFactory
	fbCustomerConversationStore sqlstore.FbExternalCustomerConversationSearchStoreFactory
}

func NewFbSearchServiceQuery(db com.MainDB) *FbSearchService {
	return &FbSearchService{
		fbCommentSearchStore:        sqlstore.NewFbExternalCommentSearchStore(db),
		fbMessageSearchStore:        sqlstore.NewFbExternalMessageSearchStore(db),
		fbCustomerConversationStore: sqlstore.NewFbExternalCustomerConversationSearchStore(db),
	}
}

func FbSearchQueryMessageBus(q *FbSearchService) fbcustomerconversationsearch.QueryBus {
	b := bus.New()
	return fbcustomerconversationsearch.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (q *FbSearchService) ListFbExternalCommentSearch(
	ctx context.Context, pageIDs []string, externalMsg string,
) ([]*fbcustomerconversationsearch.FbExternalCommentSearch, error) {
	return q.
		fbCommentSearchStore(ctx).
		ByPageIDs(pageIDs).
		ByExternalMessageNorm(externalMsg).
		ListExternalCommentSearch()
}

func (q *FbSearchService) ListFbExternalMessageSearch(
	ctx context.Context, pageIDs []string, externalMsg string,
) ([]*fbcustomerconversationsearch.FbExternalMessageSearch, error) {
	return q.
		fbMessageSearchStore(ctx).
		ByPageIDs(pageIDs).
		ByExternalMessageNorm(externalMsg).
		ListExternalMessageSearch()
}

func (q *FbSearchService) ListFbExternalConversationSearch(
	ctx context.Context, pageIDs []string, extUserName string,
) ([]*fbcustomerconversationsearch.FbCustomerConversationSearch, error) {
	return q.
		fbCustomerConversationStore(ctx).
		ByPageIDs(pageIDs).
		ByExternalUserNameNorm(extUserName).
		ListFbCustomerConversationSearch()
}
