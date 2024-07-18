package utils

import (
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func LogHttpRequest(ctx context.Context, req *http.Request) error {
	ctx = tflog.SetField(ctx, "http.url", req.URL.String())
	ctx = tflog.SetField(ctx, "http.method", req.Method)
	//ctx = tflog.SetField(ctx, "http.header", req.Header)

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

func LogHttpResponse(ctx context.Context, res *http.Response, body *[]byte) error {
	ctx = tflog.SetField(ctx, "http.url", res.Request.URL.String())
	ctx = tflog.SetField(ctx, "http.method", res.Request.Method)
	ctx = tflog.SetField(ctx, "http.status_code", res.StatusCode)
	//ctx = tflog.SetField(ctx, "http.header", res.Header)
	if body != nil {
		ctx = tflog.SetField(ctx, "http.request.body", string(*body))
	}

	tflog.Trace(ctx, "Receiving IDMC API response.")
	return nil

}
