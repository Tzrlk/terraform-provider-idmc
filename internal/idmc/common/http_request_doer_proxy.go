package common

import (
	"net/http"
	"terraform-provider-idmc/internal/utils"
)

var _ HttpRequestDoer = &HttpRequestDoerProxy{}

type HttpRequestDoerProxy struct {
	target    *HttpRequestDoer
	inspector *HttpInspector
}

func NewHttpRequestDoerProxy(target *HttpRequestDoer, inspector *HttpInspector) HttpRequestDoerProxy {
	return HttpRequestDoerProxy{
		target:    target,
		inspector: inspector,
	}
}

func (h HttpRequestDoerProxy) Do(req *http.Request) (*http.Response, error) {

	// First apply the onRequest hooks
	for _, onRequest := range h.inspector.onRequest {
		reqHandleErr := onRequest(req)
		if reqHandleErr != nil {
			return nil, reqHandleErr
		}
	}

	// Then pass it on to the actual doer to get the response
	resp, respErr := utils.Val(h.target).Do(req)
	if respErr != nil {
		return nil, respErr
	}

	// Now apply the onResponse hooks
	for _, onResponse := range h.inspector.onResponse {
		respHandleErr := onResponse(resp)
		if respHandleErr != nil {
			return nil, respHandleErr
		}
	}

	// Finally return the response
	return resp, nil

}
