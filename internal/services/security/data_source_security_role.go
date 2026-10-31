package security

import (
	"fmt"
	"strings"

	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceSecurityRole() *schema.Resource {
	return &schema.Resource{
		Description: "Use this to get a specified secrity role.",

		Read: dataSourceSecurityRoleRead,
		Schema: map[string]*schema.Schema{
			"id": common.DataSourceID,
			"roleid": {
				Description: "The id of the role.",
				Required:    true,
				Type:        schema.TypeString,
			},
			"name": {
				Description: "The name of the role.",
				Computed:    true,
				Type:        schema.TypeString,
			},
			"description": {
				Description: "The description of this role.",
				Computed:    true,
				Type:        schema.TypeString,
			},
			"privileges": {
				Description: "The privileges of this role.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
				Set: func(v interface{}) int {
					return schema.HashString(strings.ToLower(v.(string)))
				},
				Type: schema.TypeSet,
			},
			"roles": {
				Description: "The roles of this role.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
				Set: func(v interface{}) int {
					return schema.HashString(strings.ToLower(v.(string)))
				},
				Type: schema.TypeSet,
			},
		},
	}
}

func dataSourceSecurityRoleRead(d *schema.ResourceData, m interface{}) error {
	roleID := d.Get("roleid").(string)
	d.SetId(roleID)

	if err := resourceSecurityRoleRead(d, m); err != nil {
		return err
	}
	if d.Id() == "" {
		return fmt.Errorf("role %q not found", roleID)
	}
	return nil
}
