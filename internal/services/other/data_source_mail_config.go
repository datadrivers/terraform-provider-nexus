package other

import (
	"github.com/dre2004/terraform-provider-nexus/internal/schema/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceMailConfig() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to query the mail config",

		Read: DataSourceMailConfigRead,
		Schema: map[string]*schema.Schema{
			"id": common.DataSourceID,
			"host": {
				Description: "hostname",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"from_address": {
				Description: "fromAddress",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"port": {
				Description: "port",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"enabled": {
				Description: "Whether the config is enabled or not",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"username": {
				Description: "Username",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"subject_prefix": {
				Description: "Subject prefix",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"start_tls_enabled": {
				Description: "Star TLS Enabled",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"start_tls_required": {
				Description: "Star TLS required",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"ssl_on_connect_enabled": {
				Description: "SSL on connect enabled",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"ssl_server_identity_check_enabled": {
				Description: "SSL on connect enabled",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"nexus_trust_store_enabled": {
				Description: "SSL on connect enabled",
				Type:        schema.TypeBool,
				Computed:    true,
			},
		},
	}
}

func DataSourceMailConfigRead(d *schema.ResourceData, m interface{}) error {
	d.SetId(MailConfigId)
	return resourceMailConfigRead(d, m)
}
