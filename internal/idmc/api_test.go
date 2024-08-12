package idmc

import (
	"net/http"
)

type FakeHttpRequestDoer struct {
	OnDo func(req *http.Request) (*http.Response, error)
}

func (f FakeHttpRequestDoer) Do(req *http.Request) (*http.Response, error) {
	return f.OnDo(req)
}
