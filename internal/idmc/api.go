package idmc

import (
	"net/http"
)

type Api struct {
	Client    *http.Client
	BaseUrl   string
	SessionId string
}

func NewApi(baseUrl string) *Api {
	return &Api{
		Client:    http.DefaultClient,
		BaseUrl:   baseUrl,
		SessionId: "",
	}
}
