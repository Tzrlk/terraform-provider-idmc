package provider

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"terraform-provider-idmc/internal/idmc/common"
	"terraform-provider-idmc/internal/utils"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	. "github.com/onsi/gomega"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"idmc": providerserver.NewProtocol6WithError(New("test")()),
}

func TestDoLogin(t *testing.T) {
	RegisterTestingT(t)

	authHost := gofakeit.DomainName()
	authUser := gofakeit.LetterN(8)
	authPass := gofakeit.LetterN(8)

	// Case inputs
	ctx := context.TODO()

	// Case outputs
	fakeApiUrl := fmt.Sprintf("https://%s/saas", gofakeit.DomainName())
	fakeSessionId := gofakeit.LetterN(8)
	fakeBody := fmt.Sprintf(
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

	baseApiUrl, sessionId, loginErr := doLogin(
		ctx, authHost, authUser, authPass,
		common.NewHttpRequestDoerSimple(func(req *http.Request) (*http.Response, error) {
			return utils.OkPtr(&http.Response{
				Status:        "200 OK",
				StatusCode:    200,
				Proto:         "HTTP/1.1",
				ProtoMajor:    1,
				ProtoMinor:    1,
				Body:          io.NopCloser(bytes.NewBufferString(fakeBody)),
				ContentLength: int64(len(fakeBody)),
				Request:       req,
				Header: http.Header{
					"Content-Type": {"application/json"},
				},
			})
		}),
	)

	Expect(loginErr).To(BeNil())
	Expect(baseApiUrl).To(Equal(fakeApiUrl))
	Expect(sessionId).To(Equal(fakeSessionId))

}
