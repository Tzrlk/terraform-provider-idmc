package common

import (
	"net/http"
)

type OnRequestHandler  = func(r *http.Request) error
type OnResponseHandler = func(r *http.Response) error

type HttpInspector struct {
	onRequest  []OnRequestHandler
	onResponse []OnResponseHandler
}

func NewHttpInspector() *HttpInspector {
	return &HttpInspector{
		onRequest:  make([]OnRequestHandler, 0),
		onResponse: make([]OnResponseHandler, 0),
	}
}

func (h *HttpInspector) Wrap(target *HttpRequestDoer) HttpRequestDoerProxy {
	return NewHttpRequestDoerProxy(target, h)
}

func (h *HttpInspector) OnRequest(handlers ...OnRequestHandler) *HttpInspector {
	h.onRequest = append(h.onRequest, handlers...)
	return h
}

func (h *HttpInspector) OnResponse(handlers ...OnResponseHandler) *HttpInspector {
	h.onResponse = append(h.onResponse, handlers...)
	return h
}


