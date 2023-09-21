package security

import (
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceSecurityPrivilegeRepositoryView() *schema.Resource {
	return &schema.Resource{
		Description: `Use this data source to get a privilege for a repository view`,

		Read: dataSourceSecurityPrivilegeRepositoryViewRead,
		Schema: map[string]*schema.Schema{
			"id": common.DataSourceID,
			"name": {
				Description: "Name of the application privilege",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "Description of the application privilege",
				Computed:    true,
				Type:        schema.TypeString,
			},
			"readonly": {
				Description: "Whether it is readonly or not",
				Computed:    true,
				Type:        schema.TypeBool,
			},
			"repository": {
				Description: "Name of the repository the privilege applies to",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"format": {
				Description: "The format of the referenced Repository",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"actions": {
				Description: "A list of allowed actions. For a list of applicable values see https://help.sonatype.com/repomanager3/nexus-repository-administration/access-control/privileges#Privileges-PrivilegeTypes",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"ADD", "READ", "DELETE", "BROWSE", "EDIT"}, false),
				},
			},
		},
	}
}

func dataSourceSecurityPrivilegeRepositoryViewRead(d *schema.ResourceData, m interface{}) error {
	d.SetId(d.Get("name").(string))
	return resourceSecurityPrivilegeRepositoryViewRead(d, m)
}
