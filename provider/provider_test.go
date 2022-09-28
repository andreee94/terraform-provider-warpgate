package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"warpgate": providerserver.NewProtocol6WithError(New("0.0.2")()),
}

func testAccPreCheck(t *testing.T) {
	checkEnvNotNull(t, "WARPGATE_HOST")
	checkEnvNotNull(t, "WARPGATE_PORT")
	checkEnvNotNull(t, "WARPGATE_USERNAME")
	checkEnvNotNull(t, "WARPGATE_PASSWORD")
	checkEnvNotNull(t, "WARPGATE_INSECURE_SKIP_VERIFY")
}

func checkEnvNotNull(t *testing.T, env string) {
	if v := os.Getenv(env); v == "" {
		t.Fatalf("%s must be set for acceptance tests", env)
	}
}

// package provider

// import (
// 	"os"
// 	"testing"

// 	"github.com/hashicorp/terraform-plugin-framework/providerserver"
// 	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
// )

// // testAccProtoV6ProviderFactories are used to instantiate a provider during
// // acceptance testing. The factory function will be invoked for every Terraform
// // CLI command executed to create a provider server to which the CLI can
// // reattach.
// var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
// 	"warpgate": providerserver.NewProtocol6WithError(New()()),
// }

// func testAccPreCheck(t *testing.T) {

// 	checkEnvNotNull(t, "WARPGATE_HOST")
// 	checkEnvNotNull(t, "WARPGATE_PORT")
// 	checkEnvNotNull(t, "WARPGATE_USERNAME")
// 	checkEnvNotNull(t, "WARPGATE_PASSWORD")
// 	checkEnvNotNull(t, "WARPGATE_INSECURE_SKIP_VERIFY")

// 	// provider, err := testAccProtoV6ProviderFactories["warpgate"]()

// 	// provider.ConfigureProvider(context.Background(),)

// 	// err := testAccProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(nil))
// 	// if err != nil {
// 	// 	t.Fatal(err)
// 	// }
// }
