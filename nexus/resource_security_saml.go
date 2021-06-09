/*
Use this resource to create a Nexus Security SAML

Example Usage

```hcl
resource "nexus_security_saml" "example" {
  idp_metadata                  = "<EntityDescriptor ...>...</EntityDescriptor>"
  entity_id                     = "http://nexus.example/service/rest/v1/security/saml/metadata"
  validate_response_signature   = true
  validate_assertion_signature  = true
  username_attribute            = "username"
  first_name_attribute          = "firstName"
  last_name_attribute           = "lastName
  email_attribute               = "email
  groups_attribute              = "groups"
}
```
*/
package nexus

import (
	nexus "github.com/datadrivers/go-nexus-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceSecuritySAML() *schema.Resource {
	return &schema.Resource{
		Create: resourceSecuritySAMLUpdate,
		Read:   resourceSecuritySAMLRead,
		Update: resourceSecuritySAMLUpdate,
		Delete: resourceSecuritySAMLDelete,
		Exists: resourceSecuritySAMLExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"idp_metadata": {
				Description: "SAML Identity Provider Metadata XML",
				Required:    true,
				Type:        schema.TypeString,
			},
			"entity_id": {
				Description: "Entity ID URI",
				Optional:    true,
				Type:        schema.TypeString,
			},
			"validate_response_signature": {
				Description: "By default, if a signing key is found in the IdP metadata, then NXRM will attempt to validate signatures on the response.",
				Optional:    true,
				Type:        schema.TypeBool,
			},
			"validate_assertion_signature": {
				Description: "By default, if a signing key is found in the IdP metadata, then NXRM will attempt to validate signatures on the assertions.",
				Optional:    true,
				Type:        schema.TypeBool,
			},
			"username_attribute": {
				Description: "IdP field mappings for username",
				Optional:    true,
				Type:        schema.TypeString,
			},
			"first_name_attribute": {
				Description: "IdP field mappings for user's given name",
				Optional:    true,
				Type:        schema.TypeString,
			},
			"last_name_attribute": {
				Description: "IdP field mappings for user's family name",
				Optional:    true,
				Type:        schema.TypeString,
			},
			"email_attribute": {
				Description: "IdP field mappings for user's email address",
				Optional:    true,
				Type:        schema.TypeString,
			},
			"groups_attribute": {
				Description: "IdP field mappings for user's groups",
				Optional:    true,
				Type:        schema.TypeString,
			},
		},
	}
}

func resourceSecuritySAMLRead(d *schema.ResourceData, m interface{}) error {
	client := m.(nexus.Client)

	saml, err := client.SAMLRead()
	if err != nil {
		return err
	}

	if saml == nil {
		d.SetId("")
		return nil
	}

	return setSecuritySAMLToResourceData(saml, d)
}

func resourceSecuritySAMLUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(nexus.Client)

	saml := getSecuritySAMLFromResourceData(d)

	if err := client.SAMLApply(saml); err != nil {
		return err
	}

	if err := setSecuritySAMLToResourceData(&saml, d); err != nil {
		return err
	}

	return resourceSecuritySAMLRead(d, m)
}

func resourceSecuritySAMLDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(nexus.Client)

	return client.SAMLDelete()
}

func resourceSecuritySAMLExists(d *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(nexus.Client)

	saml, _ := client.SAMLRead()

	return saml != nil, nil
}

func setSecuritySAMLToResourceData(saml *nexus.SAML, d *schema.ResourceData) error {
	d.SetId("saml")
	d.Set("idp_metadata", saml.IdpMetadata)
	d.Set("entity_id", saml.EntityId)
	d.Set("validate_response_signature", saml.ValidateResponseSignature)
	d.Set("validate_assertion_signature", saml.ValidateAssertionSignature)
	d.Set("username_attribute", saml.UsernameAttribute)
	d.Set("first_name_attribute", saml.FirstNameAttribute)
	d.Set("last_name_attribute", saml.LastNameAttribute)
	d.Set("email_attribute", saml.EmailAttribute)
	d.Set("groups_attribute", saml.GroupsAttribute)

	return nil
}

func getSecuritySAMLFromResourceData(d *schema.ResourceData) nexus.SAML {
	saml := nexus.SAML{
		IdpMetadata:                d.Get("idp_metadata").(string),
		EntityId:                   d.Get("entity_id").(string),
		ValidateResponseSignature:  d.Get("validate_response_signature").(bool),
		ValidateAssertionSignature: d.Get("validate_assertion_signature").(bool),
		UsernameAttribute:          d.Get("username_attribute").(string),
		FirstNameAttribute:         d.Get("first_name_attribute").(string),
		LastNameAttribute:          d.Get("last_name_attribute").(string),
		EmailAttribute:             d.Get("email_attribute").(string),
		GroupsAttribute:            d.Get("groups_attribute").(string),
	}

	return saml
}
