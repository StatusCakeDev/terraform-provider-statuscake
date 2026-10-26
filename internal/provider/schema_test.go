package provider_test

import (
	"context"
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"

	"github.com/StatusCakeDev/terraform-provider-statuscake/v2/internal/provider"
)

var update = flag.Bool("update", false, "rewrite the golden files")

// The schema is described as Terraform receives it from the provider, which is
// what users' configurations and states depend on. It does not depend on how
// the provider is built, so it is the same before and after a resource is moved
// to the Plugin Framework. Descriptions are left out because the documentation
// is generated from them.
type attributeShape struct {
	Type       string `json:"type"`
	Required   bool   `json:"required,omitempty"`
	Optional   bool   `json:"optional,omitempty"`
	Computed   bool   `json:"computed,omitempty"`
	Sensitive  bool   `json:"sensitive,omitempty"`
	Deprecated bool   `json:"deprecated,omitempty"`
}

type blockShape struct {
	Nesting    string                    `json:"nesting"`
	MinItems   int64                     `json:"min_items,omitempty"`
	MaxItems   int64                     `json:"max_items,omitempty"`
	Attributes map[string]attributeShape `json:"attributes,omitempty"`
	Blocks     map[string]blockShape     `json:"blocks,omitempty"`
}

type schemaShape struct {
	Attributes map[string]attributeShape `json:"attributes,omitempty"`
	Blocks     map[string]blockShape     `json:"blocks,omitempty"`
}

type providerShape struct {
	Provider    schemaShape            `json:"provider"`
	Resources   map[string]schemaShape `json:"resources"`
	DataSources map[string]schemaShape `json:"data_sources"`
}

func describeBlock(b *tfprotov5.SchemaBlock) schemaShape {
	var s schemaShape
	for _, a := range b.Attributes {
		if s.Attributes == nil {
			s.Attributes = map[string]attributeShape{}
		}
		s.Attributes[a.Name] = attributeShape{
			Type:       a.Type.String(),
			Required:   a.Required,
			Optional:   a.Optional,
			Computed:   a.Computed,
			Sensitive:  a.Sensitive,
			Deprecated: a.Deprecated,
		}
	}
	for _, nb := range b.BlockTypes {
		if s.Blocks == nil {
			s.Blocks = map[string]blockShape{}
		}
		inner := describeBlock(nb.Block)
		s.Blocks[nb.TypeName] = blockShape{
			Nesting:    strings.TrimPrefix(nb.Nesting.String(), "SchemaNestedBlockNestingMode"),
			MinItems:   nb.MinItems,
			MaxItems:   nb.MaxItems,
			Attributes: inner.Attributes,
			Blocks:     inner.Blocks,
		}
	}
	return s
}

func getSchema(t *testing.T) *tfprotov5.GetProviderSchemaResponse {
	t.Helper()

	server, err := provider.NewProtoV5Server(context.Background(), "test", "")
	if err != nil {
		t.Fatal(err)
	}
	res, err := server().GetProviderSchema(context.Background(), &tfprotov5.GetProviderSchemaRequest{})
	if err != nil {
		t.Fatalf("failed to get the provider schema: %v", err)
	}
	for _, d := range res.Diagnostics {
		t.Errorf("%s: %s: %s", d.Severity, d.Summary, d.Detail)
	}
	return res
}

// TestSchemaIsUnchanged guards the contract with existing users. Their
// configurations and states depend on every name, type and flag here, so a
// change to any of them must be deliberate. When one is, review the difference
// and refresh the file with:
//
//	go test ./internal/provider -run TestSchemaIsUnchanged -update
func TestSchemaIsUnchanged(t *testing.T) {
	t.Parallel()

	if developmentBuild {
		t.Skip("this is a development build, whose schema has settings that a release does not")
	}

	res := getSchema(t)
	shape := providerShape{
		Provider:    describeBlock(res.Provider.Block),
		Resources:   map[string]schemaShape{},
		DataSources: map[string]schemaShape{},
	}
	for name, s := range res.ResourceSchemas {
		shape.Resources[name] = describeBlock(s.Block)
	}
	for name, s := range res.DataSourceSchemas {
		shape.DataSources[name] = describeBlock(s.Block)
	}

	got, err := json.MarshalIndent(shape, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	got = append(got, '\n')

	path := filepath.Join("testdata", "schema.golden.json")
	if *update {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, got, 0o644); err != nil {
			t.Fatal(err)
		}
		return
	}

	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read golden file: %v (create it with -update)", err)
	}
	if diff := cmp.Diff(string(want), string(got)); diff != "" {
		t.Errorf("the provider's schema has changed (-golden +current):\n%s", diff)
	}
}

// The providers that are served as one must declare the same settings, or
// Terraform cannot be served by both of them. This also holds when the API
// token comes from the environment, which makes it optional.
func TestProvidersAgreeOnTheirSettings(t *testing.T) {
	tests := map[string]struct {
		token    string
		required bool
	}{
		"token in configuration": {
			token:    "",
			required: true,
		},
		"token in environment": {
			token:    testToken,
			required: false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Setenv("STATUSCAKE_API_TOKEN", tt.token)

			res := getSchema(t)

			var token *tfprotov5.SchemaAttribute
			for _, a := range res.Provider.Block.Attributes {
				if a.Name == "api_token" {
					token = a
				}
			}
			if token == nil {
				t.Fatal("api_token is not a setting")
			}
			if token.Required != tt.required || token.Optional == tt.required {
				t.Errorf("api_token required = %v, optional = %v; want required = %v", token.Required, token.Optional, tt.required)
			}

			if !slices.Contains(attributeNames(res.Provider.Block), "rps") {
				t.Error("rps is not a setting")
			}
		})
	}
}

func attributeNames(b *tfprotov5.SchemaBlock) []string {
	var names []string
	for _, a := range b.Attributes {
		names = append(names, a.Name)
	}
	return names
}
