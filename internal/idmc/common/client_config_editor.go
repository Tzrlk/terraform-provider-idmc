package common

import (
	"context"
	"net/http"
	"terraform-provider-idmc/internal/utils"
)

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// ResponseEditorFn  is the function signature for the ResponseEditor callback function
type ResponseEditorFn func(ctx context.Context, res *http.Response) error

// ApiResponseEditorFn are functions that inspect or alter api-wrapped http responses.
type ApiResponseEditorFn func(ctx context.Context, apiRes *ClientResponse) error

// ClientConfigEditor
// Combines some number of request, response, and/or apiResponse editors into
// one object that can be passed into api operations.
type ClientConfigEditor struct {
	RequestEditors     []RequestEditorFn
	ResponseEditors    []ResponseEditorFn
	ApiResponseEditors []ApiResponseEditorFn
}

func (c ClientConfigEditor) Merge(other ...ClientConfigEditor) ClientConfigEditor {
	otherCount := len(other)
	if otherCount < 1 {
		return c
	}
	next := ClientConfigEditor{
		RequestEditors:     utils.NewSliceFrom(c.RequestEditors, other[0].RequestEditors),
		ResponseEditors:    utils.NewSliceFrom(c.ResponseEditors, other[0].ResponseEditors),
		ApiResponseEditors: utils.NewSliceFrom(c.ApiResponseEditors, other[0].ApiResponseEditors),
	}
	if otherCount < 2 {
		return next
	}
	return c.Merge(other[1:]...)
}

// EditHttpRequest
// Performs any needed manipulations to the api request before sending it.
func (c ClientConfigEditor) EditHttpRequest(ctx context.Context, req *http.Request) error {
	for _, editor := range c.RequestEditors {
		if err := editor(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// EditHttpResponse
// Performs any needed manipulations to the api response after receiving it.
func (c ClientConfigEditor) EditHttpResponse(ctx context.Context, res *http.Response) error {
	for _, editor := range c.ResponseEditors {
		if err := editor(ctx, res); err != nil {
			return err
		}
	}
	return nil
}

// EditApiResponse
// Performs any needed manipulations to the api response after parsing it.
func (c ClientConfigEditor) EditApiResponse(ctx context.Context, apiRes *ClientResponse) error {
	for _, editor := range c.ApiResponseEditors {
		if err := editor(ctx, apiRes); err != nil {
			return err
		}
	}
	return nil
}
