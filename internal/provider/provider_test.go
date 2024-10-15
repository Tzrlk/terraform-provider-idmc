package provider

import (
	"errors"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/joho/godotenv"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProviders = map[string]func() (tfprotov6.ProviderServer, error){
	"idmc": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) func() {
	return func() {

		// <editor-fold desc="DOTENV"> /////////////////////////////////////////
		// Check if a dotenv file exists at the project root.
		_, err := os.Stat("../../.env")
		if err == nil {
			// If it exists, try to load it for later.
			err = godotenv.Load("../../.env")
			if err != nil {
				// The file exists, but failed to load.
				t.Fatal("error loading .env file", err)
			}
		} else if !errors.Is(err, os.ErrNotExist) {
			// The file has some other issue apart from existence.
			t.Fatal("failed to stat dotenv file", err)
		} else {
			// The file just doesn't exist.
			t.Log("dotenv file not found at project root")
		}
		// </editor-fold> //////////////////////////////////////////////////////

		// <editor-fold desc="IDMC AUTH"> //////////////////////////////////////
		// Ensure required environment variables are set.
		failed := false
		for _, key := range []string{"HOST", "USER", "PASS"} {
			if val := os.Getenv("IDMC_AUTH_" + key); val == "" {
				t.Log("Missing environment variable: " + key)
				failed = true
			}
		}
		if failed {
			t.FailNow()
		}
		// </editor-fold> //////////////////////////////////////////////////////

	}
}
