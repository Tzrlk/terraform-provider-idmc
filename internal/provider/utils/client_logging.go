package utils

import (
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"terraform-provider-idmc/internal/idmc/common"
)

func LogHttpRequest(ctx context.Context, req *http.Request) error {
	ctx = tflog.SetField(ctx, "http.method", req.Method)
	ctx = tflog.SetField(ctx, "http.url", req.URL.String())
	ctx = tflog.SetField(ctx, "http.user_agent", req.UserAgent())
	ctx = tflog.SetField(ctx, "http.header", req.Header)

	// Attempt to re-serialise the request body.
	bodyReadCloser, bodyReadErr := req.GetBody()
	if bodyReadErr == nil {
		copyBuffer := new(strings.Builder)
		_, copyErr := io.Copy(copyBuffer, bodyReadCloser)
		if copyErr == nil {
			ctx = tflog.SetField(ctx, "http.request.body", copyBuffer.String())
		}
	}

	tflog.Trace(ctx, "Sending IDMC API request.")
	return nil

}

func WithRequestLogger(tfCtx context.Context) common.RequestEditorFn {
	return func(httpCtx context.Context, req *http.Request) error {
		return LogHttpRequest(tfCtx, req)
	}
}

func LogHttpResponse(ctx context.Context, res *http.Response) error {
	ctx = tflog.SetField(ctx, "http.status_code", res.StatusCode)
	ctx = tflog.SetField(ctx, "http.header", res.Header)

	copyBuffer := new(strings.Builder)
	_, copyErr := io.Copy(copyBuffer, res.Body)
	if copyErr == nil {
		ctx = tflog.SetField(ctx, "http.request.body", copyBuffer.String())
	} else {
		ctx = tflog.SetField(ctx, "http.request.body.err", copyErr)
	}

	tflog.Trace(ctx, "Receiving IDMC API response.")
	return nil

}
