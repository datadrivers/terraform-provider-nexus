package security

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceSecuritySAML() *schema.Resource {
	return &schema.Resource{
		Description: `~> PRO Feature

Use this data source to get the saml configuration.`,

		Read: dataSourceSecuritySamlRead,
		Schema: map[string]*schema.Schema{
			"idp_metadata": {
				Description: "SAML Identity Provider Metadata XML",
				Computed:    true,
				Type:        schema.TypeString,
			},
			"entity_id": {
				Description: "Entity ID URI",
				Computed:    true,
				Type:        schema.TypeString,
			},
			"validate_response_signature": {
				Description: "By default, if a signing key is found in the IdP metadata, then NXRM will attempt to validate signatures on the response.",
				Computed:    true,
				Type:        schema.TypeBool,
			},
			"validate_assertion_signature": {
				Description: "By default, if a signing key is found in the IdP metadata, then NXRM will attempt to validate signatures on the assertions.",
				Computed:    true,
				Type:        schema.TypeBool,
			},
			"username_attribute": {
				Description: "IdP field mappings for username",
				Computed:    true,
				Type:        schema.TypeString,
			},
			"first_name_attribute": {
				Description: "IdP field mappings for user's given name",
				Computed:    true,
				Type:        schema.TypeString,
			},
			"last_name_attribute": {
				Description: "IdP field mappings for user's family name",
				Computed:    true,
				Type:        schema.TypeString,
			},
			"email_attribute": {
				Description: "IdP field mappings for user's email address",
				Computed:    true,
				Type:        schema.TypeString,
			},
			"groups_attribute": {
				Description: "IdP field mappings for user's groups",
				Computed:    true,
				Type:        schema.TypeString,
			},
		},
	}
}

func dataSourceSecuritySamlRead(d *schema.ResourceData, m interface{}) error {
	return resourceSecuritySAMLRead(d, m)
}
