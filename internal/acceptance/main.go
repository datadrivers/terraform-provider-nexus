package acceptance

import (
	"fmt"
	"os"
	"testing"
	"text/template"

	"github.com/datadrivers/terraform-provider-nexus/internal/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	TestAccProviders map[string]*schema.Provider
	TestAccProvider  *schema.Provider
	TemplateFuncMap  = template.FuncMap{
		"deref": func(data interface{}) string {
			switch v := data.(type) {
			case *string:
				return *v
			case *int:
				return fmt.Sprintf("%d", *v)
			default:
				return fmt.Sprintf("%v", v)
			}
		},
	}
)

func init() {
	TestAccProvider = provider.Provider()
	TestAccProviders = map[string]*schema.Provider{
		"nexus": TestAccProvider,
	}
}

func AccPreCheck(t *testing.T) {
	if v := os.Getenv("NEXUS_URL"); v == "" {
		t.Fatalf("NEXUS_URL must be set for acceptance tests")
	}
	if v := os.Getenv("NEXUS_USERNAME"); v == "" {
		t.Fatalf("NEXUS_USERNAME must be set for acceptance tests")
	}
	if v := os.Getenv("NEXUS_PASSWORD"); v == "" {
		t.Fatalf("NEXUS_PASSWORD must be set for acceptance tests")
	}
}
