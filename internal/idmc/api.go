package idmc

import (
	"net/http"
)

type Api struct {
	Client    *http.Client
	BaseUrl   string
	SessionId string
	V2        *ApiV2
	V3        *ApiV3
}

func NewApi(baseUrl string) *Api {
	api := &Api{
		Client:    http.DefaultClient,
		BaseUrl:   baseUrl,
		SessionId: "",
		V2:        nil,
		V3:        nil,
	}
	api.V2 = &ApiV2 {
		Root: api,
	}
	api.V3 = &ApiV3 {
		Root: api,
	}

	return api
}

type ApiV2 struct {
	Root *Api
}

type ApiV3 struct {
	Root *Api
}
