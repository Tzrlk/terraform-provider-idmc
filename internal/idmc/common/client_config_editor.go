package common

import (
	"context"
	"net/http"
	"terraform-provider-idmc/internal/utils"
)

// Typedefs ////////////////////////////////////////////////////////////////////

type EditorFn[T any] func(ctx context.Context, cfg *ClientConfig, subject *T) error

// RequestEditorFn  is the function signature for the RequestEditor callback function.
type RequestEditorFn = EditorFn[http.Request]

// ResponseEditorFn  is the function signature for the ResponseEditor callback function.
type ResponseEditorFn = EditorFn[http.Response]

// ApiResponseEditorFn are functions that inspect or alter api-wrapped http responses.
type ApiResponseEditorFn = EditorFn[ClientResponse]

// Editor //////////////////////////////////////////////////////////////////////

// ClientConfigEditor
// Combines some number of request, response, and/or apiResponse editors into
// one object that can be passed into api operations.
type ClientConfigEditor struct {
	RequestEditors     []RequestEditorFn
	ResponseEditors    []ResponseEditorFn
	ApiResponseEditors []ApiResponseEditorFn
}

func (c ClientConfigEditor) AsSlice(others ...ClientConfigEditor) []ClientConfigEditor {
	result := make([]ClientConfigEditor, 1+len(others))
	result[0] = c
	for index, item := range others {
		result[1+index] = item
	}
	return result
}

func (c ClientConfigEditor) Merge(others ...ClientConfigEditor) ClientConfigEditor {
	otherCount := len(others)
	if otherCount < 1 {
		return c
	}
	next := ClientConfigEditor{
		RequestEditors:     utils.NewSliceFrom(c.RequestEditors, others[0].RequestEditors),
		ResponseEditors:    utils.NewSliceFrom(c.ResponseEditors, others[0].ResponseEditors),
		ApiResponseEditors: utils.NewSliceFrom(c.ApiResponseEditors, others[0].ApiResponseEditors),
	}
	if otherCount < 2 {
		return next
	}
	return next.Merge(others[1:]...)
}

// EditHttpRequest
// Performs any needed manipulations to the api request before sending it.
func (c ClientConfigEditor) EditHttpRequest(ctx context.Context, cfg *ClientConfig, req *http.Request) error {
	for _, editor := range c.RequestEditors {
		if err := editor(ctx, cfg, req); err != nil {
			return err
		}
	}
	return nil
}

// EditHttpResponse
// Performs any needed manipulations to the api response after receiving it.
func (c ClientConfigEditor) EditHttpResponse(ctx context.Context, cfg *ClientConfig, res *http.Response) error {
	for _, editor := range c.ResponseEditors {
		if err := editor(ctx, cfg, res); err != nil {
			return err
		}
	}
	return nil
}

// EditApiResponse
// Performs any needed manipulations to the api response after parsing it.
func (c ClientConfigEditor) EditApiResponse(ctx context.Context, cfg *ClientConfig, apiRes *ClientResponse) error {
	for _, editor := range c.ApiResponseEditors {
		if err := editor(ctx, cfg, apiRes); err != nil {
			return err
		}
	}
	return nil
}
