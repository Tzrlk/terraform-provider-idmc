package idmc

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit/v7"

	v3 "terraform-provider-idmc/internal/idmc/admin/v3"

	. "github.com/onsi/gomega"
)

type FakeHttpRequestDoer struct {
	OnDo func(req *http.Request) (*http.Response, error)
}
func (f FakeHttpRequestDoer) Do(req *http.Request) (*http.Response, error) {
	return f.OnDo(req)
}

func TestDoLogin(t *testing.T) {
	RegisterTestingT(t)

	authHost := gofakeit.DomainName()
	authUser := gofakeit.LetterN(8)
	authPass := gofakeit.LetterN(8)

	// Case inputs
	ctx  := context.TODO()

	// Case outputs
	fakeApiUrl := fmt.Sprintf("https://%s/saas", gofakeit.DomainName())
	fakeSessionId := gofakeit.LetterN(8)

	baseApiUrl, sessionId, loginErr := doLogin(
		ctx, authHost, authUser, authPass,
		func(client *v3.Client) error {
			client.Client = FakeHttpRequestDoer {
				OnDo: func(req *http.Request) (*http.Response, error) {
					body := fmt.Sprintf(
						`{
 	"products": [
 		{
 			"name": "Integration Cloud",
 			"baseApiUrl": "%s"
 		}
 	],
 	"userInfo": {
 		"sessionId": "%s",
 		"id": "9L1GFroXSDHe2IIg7QhBaT",
 		"name": "user",
 		"parentOrgId": "52ZSTB0IDK6dXxaEQLUaQu",
 		"orgId": "0cuQSDTq5sikvN7x8r1xm1",
 		"orgName": "MyOrg_INFA",
 		"groups": {},
 		"status": "Active"
 	}
}`,
						fakeApiUrl,
						fakeSessionId,
					)
					return &http.Response{
						Status:        "200 OK",
						StatusCode:    200,
						Proto:         "HTTP/1.1",
						ProtoMajor:    1,
						ProtoMinor:    1,
						Body:          io.NopCloser(bytes.NewBufferString(body)),
						ContentLength: int64(len(body)),
						Request:       req,
						Header:        http.Header {
							"Content-Type": {"application/json"},
						},
					}, nil
				},
			}
			return nil
		},
	)

	Expect(loginErr).To(BeNil())
	Expect(*baseApiUrl).To(Equal(fakeApiUrl))
	Expect(*sessionId).To(Equal(fakeSessionId))

}
