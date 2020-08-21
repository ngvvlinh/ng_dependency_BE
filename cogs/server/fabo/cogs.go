package fabo

import (
	"net/http"

	"o.o/backend/pkg/common/apifw/httpx"
	cmservice "o.o/backend/pkg/common/apifw/service"
	"o.o/backend/res/dl/fabo"
)

type FaboImageHandler httpx.Server

func BuildFaboImageHandler() FaboImageHandler {
	mux := http.NewServeMux()

	mux.Handle("/dl/fabo/default_avatar.png",
		cmservice.ServeFaboAssets(
			fabo.AssetDefaultAvatarPath,
			cmservice.MINEPNG,
		),
	)
	return httpx.MakeServer("/dl/fabo/", mux)
}
