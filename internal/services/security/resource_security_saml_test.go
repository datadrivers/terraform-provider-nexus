package security_test

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/datadrivers/go-nexus-client/nexus3/schema/security"
	"github.com/datadrivers/terraform-provider-nexus/internal/acceptance"
	"github.com/datadrivers/terraform-provider-nexus/internal/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/assert"
)

func testAccResourceSecuritySAML() (*security.SAML, error) {
	dat, err := os.ReadFile("../../../examples/saml-testconfig.xml")
	if err != nil {
		return nil, err
	}

	return &security.SAML{
		// https://samltest.id/saml/idp
		IdpMetadata:                string(dat),
		EntityId:                   "http://example.test/client/rest/v1/security/saml/metadata",
		ValidateAssertionSignature: false,
		ValidateResponseSignature:  true,
		UsernameAttribute:          "username2",
		FirstNameAttribute:         "firstName2",
		LastNameAttribute:          "lastName2",
		EmailAttribute:             "email2",
		GroupsAttribute:            "groups2",
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
					resource.TestCheckResourceAttr(resName, "validate_response_signature", strconv.FormatBool(saml.ValidateResponseSignature)),
					resource.TestCheckResourceAttr(resName, "validate_assertion_signature", strconv.FormatBool(saml.ValidateAssertionSignature)),
					resource.TestCheckResourceAttr(resName, "username_attribute", saml.UsernameAttribute),
					resource.TestCheckResourceAttr(resName, "first_name_attribute", saml.FirstNameAttribute),
					resource.TestCheckResourceAttr(resName, "last_name_attribute", saml.LastNameAttribute),
					resource.TestCheckResourceAttr(resName, "email_attribute", saml.EmailAttribute),
					resource.TestCheckResourceAttr(resName, "groups_attribute", saml.GroupsAttribute),
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
	return fmt.Sprintf(`
resource "nexus_security_saml" "acceptance" {
	idp_metadata                 = <<-EOF
%sEOF
	entity_id                    = "%s"
	validate_response_signature  = "%t"
	validate_assertion_signature = "%t"
	username_attribute           = "%s"
	first_name_attribute         = "%s"
	last_name_attribute          = "%s"
	email_attribute              = "%s"
	groups_attribute             = "%s"
}
`, saml.IdpMetadata, saml.EntityId, saml.ValidateResponseSignature, saml.ValidateAssertionSignature, saml.UsernameAttribute, saml.FirstNameAttribute, saml.LastNameAttribute, saml.EmailAttribute, saml.GroupsAttribute)
}
