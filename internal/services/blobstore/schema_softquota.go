package blobstore

import (
	"github.com/datadrivers/go-nexus-client/nexus3/schema/blobstore"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func getResourceSoftQuotaSchema() *schema.Schema {
	return &schema.Schema{
		Description: "Soft quota of the blobstore",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"limit": {
					Description:  "The limit in Bytes. Minimum value is 1000000",
					Required:     true,
					Type:         schema.TypeInt,
					ValidateFunc: validation.IntAtLeast(100000),
				},
				"type": {
					Description:  "The type to use such as spaceRemainingQuota, or spaceUsedQuota",
					Required:     true,
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"spaceRemainingQuota", "spaceUsedQuota"}, false),
				},
			},
		},
		MaxItems: 1,
		Optional: true,
		Type:     schema.TypeList,
	}
}

func getDataSourceSoftQuotaSchema() *schema.Schema {
	return &schema.Schema{
		Description: "Soft quota of the blobstore",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"limit": {
					Description: "The limit in Bytes. Minimum value is 1000000",
					Type:        schema.TypeInt,
					Computed:    true,
				},
				"type": {
					Description: "The type to use such as spaceRemainingQuota, or spaceUsedQuota",
					Type:        schema.TypeString,
					Computed:    true,
				},
			},
		},
		Computed: true,
		Type:     schema.TypeList,
	}
}

func flattenSoftQuota(softQuota *blobstore.SoftQuota) []map[string]interface{} {
	if softQuota == nil {
		return nil
	}
	data := map[string]interface{}{
		"limit": softQuota.Limit,
		"type":  softQuota.Type,
	}
	return []map[string]interface{}{data}
}
