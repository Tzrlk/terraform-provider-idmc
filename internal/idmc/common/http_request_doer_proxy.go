package common

import (
	"net/http"
	"terraform-provider-idmc/internal/utils"
)

var _ HttpRequestDoer = &HttpRequestDoerProxy{}
//var _ v2.HttpRequestDoer = &HttpRequestDoerProxy{}
//var _ v3.HttpRequestDoer = &HttpRequestDoerProxy{}

type HttpRequestDoerProxy struct {
	target      *HttpRequestDoer
	interceptor *HttpInspector
}

func NewHttpRequestDoerProxy(target *HttpRequestDoer, interceptor *HttpInspector) HttpRequestDoerProxy {
	return HttpRequestDoerProxy{
		target:      target,
		interceptor: interceptor,
	}
}

func (h HttpRequestDoerProxy) Do(req *http.Request) (*http.Response, error) {

	// First apply the onRequest hooks
	for _, onRequest := range h.interceptor.onRequest {
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
	for _, onResponse := range h.interceptor.onResponse {
		respHandleErr := onResponse(resp)
		if respHandleErr != nil {
			return nil, respHandleErr
		}
	}

	// Finally return the response
	return resp, nil

}

