package security_test

import (
	"bytes"
	"testing"
	"text/template"

	"github.com/datadrivers/terraform-provider-nexus/internal/acceptance"
	"github.com/datadrivers/terraform-provider-nexus/internal/services/security"
	"github.com/datadrivers/terraform-provider-nexus/internal/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const resourceSecurityOIDCTemplateString = `
resource "nexus_security_oidc" "acceptance" {
	client_id         = "{{ .ClientID }}"
	client_secret     = "{{ .ClientSecret }}"
	authorization_url = "{{ .IdpAuthorizationURL }}"
	token_url         = "{{ .IdpTokenURL }}"
	jwks_url          = "{{ .IdpJwksURL }}"
	jws_algorithm     = "{{ .IdpJwsAlgorithm }}"
	username_claim    = "{{ .UsernameClaim }}"
	groups_claim      = "{{ .GroupsClaim }}"
	{{- if .IdpLogoutURL }}
	logout_url        = "{{ .IdpLogoutURL }}"
	{{- end }}
	{{- if .FirstNameClaim }}
	first_name_claim  = "{{ .FirstNameClaim }}"
	{{- end }}
	{{- if .LastNameClaim }}
	last_name_claim   = "{{ .LastNameClaim }}"
	{{- end }}
	{{- if .EmailClaim }}
	email_claim       = "{{ .EmailClaim }}"
	{{- end }}
	use_trust_store   = {{ .UseTrustStore }}
}`

func testAccResourceSecurityOIDC() security.OIDC {
	return security.OIDC{
		ClientID:            "nexus-acceptance",
		ClientSecret:        "s3cr3t",
		IdpAuthorizationURL: "https://idp.example.test/authorize",
		IdpTokenURL:         "https://idp.example.test/token",
		IdpJwksURL:          "https://idp.example.test/.well-known/jwks.json",
		IdpJwsAlgorithm:     "RS256",
		UsernameClaim:       "preferred_username",
		GroupsClaim:         "groups",
		IdpLogoutURL:        "https://idp.example.test/logout",
		FirstNameClaim:      "given_name",
		LastNameClaim:       "family_name",
		EmailClaim:          "email",
		UseTrustStore:       false,
	}
}

func TestAccResourceSecurityOIDC(t *testing.T) {
	if tools.GetEnv("SKIP_PRO_TESTS", "false") == "true" {
		t.Skip("Skipping Nexus Pro tests")
	}

	cfg := testAccResourceSecurityOIDC()
	resName := "nexus_security_oidc.acceptance"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSecurityOIDCConfig(cfg),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "client_id", cfg.ClientID),
					resource.TestCheckResourceAttr(resName, "authorization_url", cfg.IdpAuthorizationURL),
					resource.TestCheckResourceAttr(resName, "token_url", cfg.IdpTokenURL),
					resource.TestCheckResourceAttr(resName, "jwks_url", cfg.IdpJwksURL),
					resource.TestCheckResourceAttr(resName, "jws_algorithm", cfg.IdpJwsAlgorithm),
					resource.TestCheckResourceAttr(resName, "username_claim", cfg.UsernameClaim),
					resource.TestCheckResourceAttr(resName, "groups_claim", cfg.GroupsClaim),
					resource.TestCheckResourceAttr(resName, "logout_url", cfg.IdpLogoutURL),
					resource.TestCheckResourceAttr(resName, "first_name_claim", cfg.FirstNameClaim),
					resource.TestCheckResourceAttr(resName, "last_name_claim", cfg.LastNameClaim),
					resource.TestCheckResourceAttr(resName, "email_claim", cfg.EmailClaim),
				),
			},
			{
				ResourceName:            resName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"client_secret"},
			},
		},
	})
}

func testAccResourceSecurityOIDCConfig(o security.OIDC) string {
	buf := &bytes.Buffer{}
	t := template.Must(template.New("NexusSecurityOIDC").Funcs(acceptance.TemplateFuncMap).Parse(resourceSecurityOIDCTemplateString))
	if err := t.Execute(buf, o); err != nil {
		panic(err)
	}
	return buf.String()
}
