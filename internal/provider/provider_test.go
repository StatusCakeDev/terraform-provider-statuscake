package provider_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/StatusCakeDev/terraform-provider-statuscake/v2/internal/provider"
)

// These tests use the provider only through what main does: the server that
// Terraform talks to. Tests that need an API to talk to run only in a
// development build (`-tags dev`), which has the setting that points the
// provider at one. See CONTRIBUTING.md.

// testToken satisfies the provider's validation of the API token. It is not a
// credential.
const testToken = "0123456789abcdefghij0123"

// TestMain removes the provider's settings from the environment, so that the
// tests do not depend on those of whoever runs them. A test that needs one
// sets it.
func TestMain(m *testing.M) {
	for _, kv := range os.Environ() {
		if name, _, _ := strings.Cut(kv, "="); strings.HasPrefix(name, "STATUSCAKE_") {
			os.Unsetenv(name)
		}
	}
	os.Exit(m.Run())
}

// providerFactories serves the provider in process, as Terraform would reach it
// over gRPC.
func providerFactories() map[string]func() (tfprotov5.ProviderServer, error) {
	return map[string]func() (tfprotov5.ProviderServer, error){
		"statuscake": func() (tfprotov5.ProviderServer, error) {
			server, err := provider.NewProtoV5Server(context.Background(), "test", "abc123")
			if err != nil {
				return nil, err
			}
			return server(), nil
		},
	}
}

// providerBlock is the provider configuration for the tests, which does not
// limit or retry requests.
func providerBlock() string {
	return fmt.Sprintf(`
		provider "statuscake" {
			api_token = %q
			rps       = 1000
			retries   = 0
		}
	`, testToken)
}

// requireCLI skips the test when there is no Terraform (or compatible)
// command-line tool to run it with. In CI, where the tool is installed, its
// absence is a failure rather than a skip.
func requireCLI(t *testing.T) {
	t.Helper()

	if os.Getenv("TF_ACC_TERRAFORM_PATH") != "" {
		return
	}
	if _, err := exec.LookPath("terraform"); err == nil {
		return
	}
	if os.Getenv("CI") != "" {
		t.Fatal("terraform is not installed, and these tests need it")
	}
	t.Skip("terraform is not installed: set TF_ACC_TERRAFORM_PATH or add it to PATH to run these tests")
}

// runSteps runs the steps with the provider, which is not given an API to talk
// to. For steps that Terraform stops before reaching the API.
func runSteps(t *testing.T, steps ...resource.TestStep) {
	t.Helper()
	requireCLI(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: providerFactories(),
		Steps:                    steps,
	})
}
