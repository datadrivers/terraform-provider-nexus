package repository

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	ResourceYumSigning = &schema.Schema{
		Description: "Contains signing data of repositores",
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"keypair": {
					Description: "PGP signing key pair (armored private key e.g. gpg --export-secret-key --armor)",
					Type:        schema.TypeString,
					Required:    true,
					Sensitive:   true,
				},
				"passphrase": {
					Description: "Passphrase to access PGP signing key",
					Type:        schema.TypeString,
					Optional:    true,
					Sensitive:   true,
				},
			},
		},
	}
	DataSourceYumSigning = &schema.Schema{
		Description: "Contains signing data of repositores",
		Type:        schema.TypeList,
		Computed:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"keypair": {
					Description: "PGP signing key pair (armored private key e.g. gpg --export-secret-key --armor)",
					Type:        schema.TypeString,
					Computed:    true,
					Sensitive:   true,
				},
				"passphrase": {
					Description: "Passphrase to access PGP signing key",
					Type:        schema.TypeString,
					Computed:    true,
					Sensitive:   true,
				},
			},
		},
	}
)
