package server

import (
	"context"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"

	"o.o/api/main/accountshipnow"
	"o.o/api/main/identity"
	"o.o/api/main/ordering"
	"o.o/api/main/shipnow"
	shippingcarrier "o.o/backend/com/main/shipping/carrier"
	"o.o/backend/pkg/common/apifw/httpx"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/integration/shipnow/ahamove"
	"o.o/backend/pkg/integration/shipnow/ahamove/webhook"
)

const PathAhamoveUserVerification = "/ahamove/user_verification"

type WebhookConfig cc.HTTP

type AhamoveWebhookServer *http.Server

func NewAhamoveWebhookServer(
	cfg WebhookConfig,
	shipmentManager *shippingcarrier.ShipmentManager,
	ahamoveCarrier *ahamove.Carrier,
	identityQuery identity.QueryBus,
	shipnowQuery shipnow.QueryBus,
	shipnowAggr shipnow.CommandBus,
	orderAggr ordering.CommandBus,
	orderQuery ordering.QueryBus,
	fileServer AhamoveVerificationFileServer,
	webhook *webhook.Webhook,
) AhamoveWebhookServer {

	mux := http.NewServeMux()
	{
		rt := httpx.New()
		rt.Use(httpx.RecoverAndLog(true))
		webhook.Register(rt)

		mux.Handle("/webhook/", rt)
	}

	// serve ahamove verification files
	mux.Handle(PathAhamoveUserVerification+"/", (*httpx.Router)(fileServer))

	svr := &http.Server{
		Addr:    cc.HTTP(cfg).Address(),
		Handler: mux,
	}
	return svr
}

type AhamoveVerificationFileServer *httpx.Router

func NewAhamoveVerificationFileServer(ctx context.Context, accountshipnowQS accountshipnow.QueryBus) AhamoveVerificationFileServer {
	// path: <UploadDirAhamoveVerification>/<originname>/<filename>.jpg
	// filepath:
	// user_id_front_<user.id>_<user.create_time>.jpg
	// user_portrait_<user.id>_<user.create_time>.jpg
	regex := regexp.MustCompile(`([0-9]+)_([0-9]+)`)

	rt := httpx.New()
	path := PathAhamoveUserVerification + "/:originname/:filename"

	rt.Router.GET(path, func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		originname := ps.ByName("originname")
		fileName := ps.ByName("filename")
		parts := regex.FindStringSubmatch(fileName)
		if len(parts) == 0 {
			http.NotFound(w, req)
			return
		}
		userID := parts[1]
		createTime := parts[2]

		query := &accountshipnow.GetExternalAccountAhamoveByExternalIDQuery{
			ExternalID: userID,
		}
		if err := accountshipnowQS.Dispatch(ctx, query); err != nil {
			http.NotFound(w, req)
			return
		}
		accountAhamove := query.Result
		xCreatedAt := accountAhamove.ExternalCreatedAt
		if strconv.FormatInt(xCreatedAt.Unix(), 10) != createTime {
			http.NotFound(w, req)
			return
		}

		url := ""
		if strings.Contains(fileName, "user_id_front") {
			url = accountAhamove.IDCardFrontImg
		} else if strings.Contains(fileName, "user_id_back") {
			url = accountAhamove.IDCardBackImg
		} else if strings.Contains(fileName, "user_portrait") {
			url = accountAhamove.PortraitImg
		}
		if strings.Contains(url, originname) {
			http.Redirect(w, req, url, http.StatusSeeOther)
			return
		}
		http.NotFound(w, req)
	})
	return rt
}
