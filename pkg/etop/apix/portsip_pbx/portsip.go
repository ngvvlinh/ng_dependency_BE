package portsip_pbx

import (
	"gopkg.in/olivere/elastic.v5"
	"o.o/api/etelecom"
	"o.o/api/main/connectioning"
	"o.o/api/top/types/etc/account_type"
	identitymodel "o.o/backend/com/main/identity/model"
	identitymodelx "o.o/backend/com/main/identity/modelx"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/httpx"
	"o.o/backend/pkg/common/elasticsearch"
	"o.o/backend/pkg/common/headers"
	"o.o/backend/pkg/etc/idutil"
	"o.o/backend/pkg/etop/authorize/claims"
	"o.o/backend/pkg/etop/authorize/middleware"
	"time"
)

type PortsipService struct {
	SS            middleware.SessionStarter
	ESStore       elasticsearch.Store
	EtelecomQuery etelecom.QueryBus
}

func (s *PortsipService) GetCallLogs(c *httpx.Context) error {
	ctx := c.Context()
	tokenStr := headers.GetBearerTokenFromCtx(ctx)

	if tokenStr == "" {
		return cm.Errorf(cm.Unauthenticated, nil, "")
	}

	var claim *claims.Claim
	var err error
	var account identitymodel.AccountInterface
	claim, account, err = s.SS.VerifyAPIKey(ctx, tokenStr, account_type.Shop)
	if err != nil {
		return err
	}

	if !idutil.IsShopID(claim.AccountID) {
		return cm.ErrUnauthenticated
	}
	query := &identitymodelx.GetShopQuery{
		ShopID: claim.AccountID,
	}
	if err := s.SS.ShopStore.GetShop(ctx, query); err != nil {
		return cm.ErrUnauthenticated
	}
	account = query.Result
	acc := account.GetAccount()

	getTenantByConnectionQuery := &etelecom.GetTenantByConnectionQuery{
		OwnerID:      acc.OwnerID,
		ConnectionID: connectioning.DefaultDirectPortsipConnectionID,
	}

	if err := s.EtelecomQuery.Dispatch(ctx, getTenantByConnectionQuery); err != nil {
		return err
	}

	tenant := getTenantByConnectionQuery.Result
	values := c.Req.URL.Query()
	scrollId := values.Get("scroll_id")

	boolQuery := elastic.NewBoolQuery()
	rangeQuery := elastic.NewRangeQuery("ended_time")
	caller := values.Get("caller")
	if caller != "" {
		boolQuery.Must(elastic.NewTermsQuery("caller", caller))
	}

	callee := values.Get("callee")
	if callee != "" {
		boolQuery.Must(elastic.NewTermsQuery("callee", callee))
	}

	sessionID := values.Get("session_id")
	if sessionID != "" {
		boolQuery.Must(elastic.NewTermsQuery("session_id", sessionID))
	}

	startTime := values.Get("start_time")
	if startTime != "" {
		boolQuery.Filter(rangeQuery.From(time.Unix(ConvertStringToInt64(startTime), 0).Format(time.RFC3339)))
	}

	endTime := values.Get("end_time")
	if endTime != "" {
		boolQuery.Filter(rangeQuery.To(time.Unix(ConvertStringToInt64(endTime), 0).Format(time.RFC3339)))
	}

	boolQuery.Must(elastic.NewTermsQuery("tenant_id", tenant.ExternalData.ID))
	sort := elastic.NewFieldSort("start_time").Desc().SortMode("max")
	index := s.ESStore.GetIndex(tenant.ExternalData.ID)
	results, err := s.ESStore.Scroll(index, boolQuery, sort, scrollId)
	if err != nil {
		return err
	}

	resp := &SessionsResponse{
		ScrollId: results.ScrollId,
		Sesssion: Convert_core_SearchHits_To_api_PortSipCallLogs(results.Hits.Hits),
	}
	c.SetResult(resp)
	return nil
}
