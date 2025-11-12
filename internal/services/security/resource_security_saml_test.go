package security_test

import (
	"bytes"
	"os"
	"strconv"
	"testing"
	"text/template"

	"github.com/datadrivers/terraform-provider-nexus/internal/acceptance"
	"github.com/datadrivers/terraform-provider-nexus/internal/tools"
	"github.com/gcroucher/go-nexus-client/nexus3/schema/security"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/assert"
)

const (
	resourceSecuritySamlTemplateString = `
resource "nexus_security_saml" "acceptance" {
	idp_metadata                 = <<-EOF
{{ .IdpMetadata }}EOF
	entity_id                    = "{{ .EntityId }}"
	username_attribute           = "{{ .UsernameAttribute }}"
	{{- if .ValidateAssertionSignature }}
	validate_assertion_signature = {{ .ValidateAssertionSignature }}
	{{- end }}
	{{- if .ValidateResponseSignature }}
	validate_response_signature  = {{ .ValidateResponseSignature }}
	{{- end }}
	{{- if .FirstNameAttribute }}
	first_name_attribute         = "{{ .FirstNameAttribute }}"
	{{- end }}
	{{- if .LastNameAttribute }}
	last_name_attribute          = "{{ .LastNameAttribute }}"
	{{- end }}
	{{- if .EmailAttribute }}
	email_attribute              = "{{ .EmailAttribute }}"
	{{- end }}
	{{- if .GroupsAttribute }}
	groups_attribute             = "{{ .GroupsAttribute }}"
	{{- end }}
}`
)

func testAccResourceSecuritySAML() (*security.SAML, error) {
	dat, err := os.ReadFile("../../../examples/saml-testconfig.xml")
	if err != nil {
		return nil, err
	}

	validateAssertionSignature := false
	validateResponseSignature := false
	firstNameAttribute := "firstName2"
	lastNameAttribute := "lastName2"
	emailAttribute := "email2"
	groupsAttribute := "groups2"

	return &security.SAML{
		// https://samltest.id/saml/idp
		IdpMetadata:                string(dat),
		EntityId:                   "http://example.test/client/rest/v1/security/saml/metadata",
		ValidateAssertionSignature: &validateAssertionSignature,
		ValidateResponseSignature:  &validateResponseSignature,
		UsernameAttribute:          "username2",
		FirstNameAttribute:         &firstNameAttribute,
		LastNameAttribute:          &lastNameAttribute,
		EmailAttribute:             &emailAttribute,
		GroupsAttribute:            &groupsAttribute,
	}, nil
}

func TestAccResourceSecuritySAML(t *testing.T) {
	if tools.GetEnv("SKIP_PRO_TESTS", "false") == "true" {
		t.Skip("Skipping Nexus Pro tests")
	}

	saml, err := testAccResourceSecuritySAML()
	assert.Nil(t, err)
	resName := "nexus_security_saml.acceptance"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSecuritySAMLConfig(*saml),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "idp_metadata", saml.IdpMetadata),
					resource.TestCheckResourceAttr(resName, "entity_id", saml.EntityId),
					resource.TestCheckResourceAttr(resName, "validate_response_signature", strconv.FormatBool(*saml.ValidateResponseSignature)),
					resource.TestCheckResourceAttr(resName, "validate_assertion_signature", strconv.FormatBool(*saml.ValidateAssertionSignature)),
					resource.TestCheckResourceAttr(resName, "username_attribute", saml.UsernameAttribute),
					resource.TestCheckResourceAttr(resName, "first_name_attribute", *saml.FirstNameAttribute),
					resource.TestCheckResourceAttr(resName, "last_name_attribute", *saml.LastNameAttribute),
					resource.TestCheckResourceAttr(resName, "email_attribute", *saml.EmailAttribute),
					resource.TestCheckResourceAttr(resName, "groups_attribute", *saml.GroupsAttribute),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResourceSecuritySAMLConfig(saml security.SAML) string {
	buf := &bytes.Buffer{}
	resourceTemplate := template.Must(template.New("NexusSecuritySaml").Funcs(acceptance.TemplateFuncMap).Parse(resourceSecuritySamlTemplateString))
	if err := resourceTemplate.Execute(buf, saml); err != nil {
		panic(err)
	}
	return buf.String()
}
