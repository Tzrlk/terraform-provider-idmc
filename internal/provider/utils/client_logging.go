package utils

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"io"
	"net/http"
	"strings"
	"terraform-provider-idmc/internal/idmc/common"
	"terraform-provider-idmc/internal/utils"
)

func ReformatHeaders(header http.Header) map[string]string {
	return utils.TransformMapValues(header, func(key string, val []string) string {
		return strings.Join(val, " | ")
	})
}

func GetHttpRequestCtx(ctx context.Context, req *http.Request) (context.Context, error) {
	var err error
	if req == nil {
		return ctx, fmt.Errorf("unable to get context for nil http request")
	}

	// Enrich context from reasonable request properties.
	ctx = tflog.SetField(ctx, "http.url", req.URL.String())
	ctx = tflog.SetField(ctx, "http.method", req.Method)
	ctx = tflog.SetField(ctx, "http.request.headers", ReformatHeaders(req.Header))

	// Attempt to get a copy of the request body.
	var bodyReadCloser io.ReadCloser
	bodyReadCloser, err = req.GetBody()
	if err != nil {
		return ctx, err
	}

	// Attempt to re-serialise the request body copy.
	copyBuffer := new(strings.Builder)
	if _, err = io.Copy(copyBuffer, bodyReadCloser); err != nil {
		return ctx, err
	}

	// Finally enrich the context with the body info and return.
	ctx = tflog.SetField(ctx, "http.request.body", copyBuffer.String())
	return ctx, nil
}

func GetHttpResponseCtx(ctx context.Context, res *http.Response) (context.Context, error) {
	if res == nil {
		return ctx, fmt.Errorf("unable to get context for nil http response")
	}

	ctx = tflog.SetField(ctx, "http.status_code", res.StatusCode)
	ctx = tflog.SetField(ctx, "http.response.headers", ReformatHeaders(res.Header))

	return GetHttpRequestCtx(ctx, res.Request)
}

func GetApiResponseCtx(ctx context.Context, apiRes common.ApiResponse) (context.Context, error) {
	if apiRes == nil {
		return ctx, fmt.Errorf("unable to get context for nil api response")
	}

	resBody := apiRes.BodyData()
	ctx = tflog.SetField(ctx, "http.request.body", string(resBody))

	return GetHttpResponseCtx(ctx, apiRes.HttpResponse())
}

func LogHttpRequest(ctx context.Context, req *http.Request) error {
	reqCtx, ctxErr := GetHttpRequestCtx(ctx, req)
	if ctxErr == nil {
		tflog.Trace(reqCtx, "Sending IDMC HTTP request.")
	}
	return ctxErr
}

func LogHttpResponse(ctx context.Context, res *http.Response) error {
	resCtx, ctxErr := GetHttpResponseCtx(ctx, res)
	if ctxErr == nil {
		tflog.Trace(resCtx, "Receiving IDMC HTTP response.")
	}
	return ctxErr
}

func LogApiResponse(ctx context.Context, apiRes common.ApiResponse) error {
	apiResCtx, apiResCtxErr := GetApiResponseCtx(ctx, apiRes)
	if apiResCtxErr == nil {
		tflog.Trace(apiResCtx, "Receiving IDMC API response.")
	}
	return apiResCtxErr
}
