/*
Use this data source to get an existing yum repository

Example Usage

```hcl
data "nexus_repository_yum_hosted" "yummy" {
  name = "yummy"
}
```
*/
package nexus

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceRepositoryYumHosted() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRepositoryYumHostedRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "A unique identifier for this repository",
				Required:    true,
				Type:        schema.TypeString,
			},
			"online": {
				Description: "Whether this repository accepts incoming requests",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"repodata_depth": {
				Description: "Specifies the repository depth where repodata folder(s) are created. Possible values: 0-5",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"deploy_policy": {
				Description: "Validate that all paths are RPMs or yum metadata. Possible values: `STRICT` or `PERMISSIVE`",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"cleanup": {
				Description: "Cleanup policies",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy_names": {
							Description: "List of policy names",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional: true,
							Set: func(v interface{}) int {
								return schema.HashString(strings.ToLower(v.(string)))
							},
							Type: schema.TypeSet,
						},
					},
				},
			},
			"storage": {
				Description: "The storage configuration of the repository",
				Type:        schema.TypeList,
				Computed:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"blob_store_name": {
							Description: "Blob store used to store repository contents",
							Computed:    true,
							Type:        schema.TypeString,
						},
						"strict_content_type_validation": {
							Description: "Whether to validate uploaded content's MIME type appropriate for the repository format",
							Computed:    true,
							Type:        schema.TypeBool,
						},
						"write_policy": {
							Description: "Controls if deployments of and updates to assets are allowed",
							Computed:    true,
							Type:        schema.TypeString,
						},
					},
				},
			},
		},
	}
}

func dataSourceRepositoryYumHostedRead(d *schema.ResourceData, m interface{}) error {
	d.SetId(d.Get("name").(string))

	return resourceYumHostedRepositoryRead(d, m)
}
