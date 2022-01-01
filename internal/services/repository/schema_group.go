package repository

import (
	"strings"

	"github.com/datadrivers/go-nexus-client/nexus3/schema/repository"
	"github.com/datadrivers/terraform-provider-nexus/internal/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func getResourceGroupSchema() *schema.Schema {
	return &schema.Schema{
		Description: "Configuration for repository group",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"member_names": {
					Description: "Member repositories names",
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
					Required: true,
					Set: func(v interface{}) int {
						return schema.HashString(strings.ToLower(v.(string)))
					},
					Type: schema.TypeSet,
				},
			},
		},
		MaxItems: 1,
		Optional: true,
		Type:     schema.TypeList,
	}
}

func flattenRepositoryGroup(group *repository.Group) []map[string]interface{} {
	if group == nil {
		return nil
	}
	data := map[string]interface{}{
		"member_names": tools.StringSliceToInterfaceSlice(group.MemberNames),
	}
	return []map[string]interface{}{data}
}
