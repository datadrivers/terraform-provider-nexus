package repository

import (
	"strings"

	"github.com/datadrivers/go-nexus-client/nexus3/schema/repository"
	"github.com/datadrivers/terraform-provider-nexus/internal/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func getResourceCleanupSchema() *schema.Schema {
	return &schema.Schema{
		DefaultFunc: cleanupDefault,
		Description: "Cleanup policies",
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    true,
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
	}
}

func getDataSourceCleanupSchema() *schema.Schema {
	return &schema.Schema{
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
					Computed: true,
					Set: func(v interface{}) int {
						return schema.HashString(strings.ToLower(v.(string)))
					},
					Type: schema.TypeSet,
				},
			},
		},
	}
}

func cleanupDefault() (interface{}, error) {
	data := map[string]interface{}{
		"policy_names": []string{},
	}
	return []map[string]interface{}{data}, nil
}

func flattenCleanup(cleanup *repository.Cleanup) []map[string]interface{} {
	if cleanup == nil {
		return nil
	}
	data := map[string]interface{}{
		"policy_names": tools.StringSliceToInterfaceSlice(cleanup.PolicyNames),
	}

	return []map[string]interface{}{data}
}
