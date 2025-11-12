package repository

import (
	"github.com/gcroucher/go-nexus-client/nexus3/schema/repository"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var (
	ResourceMaven = &schema.Schema{
		Description: "Maven contains additional data of maven repository",
		Type:        schema.TypeList,
		Required:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"version_policy": {
					Description: "What type of artifacts does this repository store? Possible Value: `RELEASE`, `SNAPSHOT` or `MIXED`",
					Required:    true,
					Type:        schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						string(repository.MavenVersionPolicyRelease),
						string(repository.MavenVersionPolicySnapshot),
						string(repository.MavenVersionPolicyMixed),
					}, false),
				},
				"layout_policy": {
					Description: "Validate that all paths are maven artifact or metadata paths. Possible Value: `STRICT` or `PERMISSIVE`",
					Required:    true,
					Type:        schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						string(repository.MavenLayoutPolicyStrict),
						string(repository.MavenLayoutPolicyPermissive),
					}, false),
				},
				"content_disposition": {
					Description: "Add Content-Disposition header as 'Attachment' to disable some content from being inline in a browse. Possible Value: `INLINE` or `ATTACHMENT`",
					Optional:    true,
					Type:        schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						string(repository.MavenContentDispositionInline),
						string(repository.MavenContentDispositionAttachment),
					}, false),
				},
			},
		},
	}
	DataSourceMaven = &schema.Schema{
		Description: "Maven contains additional data of maven repository",
		Type:        schema.TypeList,
		Computed:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"version_policy": {
					Description: "What type of artifacts does this repository store?",
					Computed:    true,
					Type:        schema.TypeString,
				},
				"layout_policy": {
					Description: "Validate that all paths are maven artifact or metadata paths",
					Computed:    true,
					Type:        schema.TypeString,
				},
				"content_disposition": {
					Description: "Add Content-Disposition header as 'Attachment' to disable some content from being inline in a browse",
					Computed:    true,
					Type:        schema.TypeString,
				},
			},
		},
	}
)
