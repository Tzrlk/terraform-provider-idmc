package common

import (
	"net/http"
)

var _ HttpRequestDoer = &HttpRequestDoerSimple{}

//var _ v2.HttpRequestDoer = &HttpRequestDoerSimple{}
//var _ v3.HttpRequestDoer = &HttpRequestDoerSimple{}

type HttpRequestDoerSimple struct {
	doer func(req *http.Request) (*http.Response, error)
}

func NewHttpRequestDoerSimple(doer func(req *http.Request) (*http.Response, error)) HttpRequestDoerSimple {
	return HttpRequestDoerSimple{
		doer: doer,
	}
}

func (s HttpRequestDoerSimple) Do(req *http.Request) (*http.Response, error) {
	return s.doer(req)
}
