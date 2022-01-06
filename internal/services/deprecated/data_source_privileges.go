package deprecated

import (
	"fmt"

	nexus "github.com/datadrivers/go-nexus-client/nexus3"
	"github.com/datadrivers/go-nexus-client/nexus3/schema/repository"
	"github.com/datadrivers/go-nexus-client/nexus3/schema/security"
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	"github.com/datadrivers/terraform-provider-nexus/internal/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourcePrivileges() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to work with privileges.",

		Read: dataSourcePrivilegesRead,
		Schema: map[string]*schema.Schema{
			"id": common.DataSourceID,
			"domain": {
				Description:  "The domain of the privilege",
				ForceNew:     true,
				Optional:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice(security.PrivilegeDomains, false),
			},
			"format": {
				Description:  "The format of the privilege",
				ForceNew:     true,
				Optional:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice(repository.RepositoryFormats, false),
			},
			"name": {
				Description: "The name of the privilege",
				ForceNew:    true,
				Optional:    true,
				Type:        schema.TypeString,
			},
			"privileges": {
				Computed:    true,
				Description: "List of privileges",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"actions": {
							Description: "Actions for the privilege (browse, read, edit, add, delete, all and run)",
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Type:        schema.TypeSet,
						},
						"content_selector": {
							Description: "The content selector for the privilege",
							Computed:    true,
							Type:        schema.TypeString,
						},
						"description": {
							Description: "A description of the privilege",
							Computed:    true,
							Type:        schema.TypeString,
						},
						"domain": {
							Description: "The domain of the privilege",
							Computed:    true,
							Type:        schema.TypeString,
						},
						"format": {
							Description: "The format of the privilege",
							Computed:    true,
							Type:        schema.TypeString,
						},
						"name": {
							Description: "The name of the privilege",
							Computed:    true,
							Type:        schema.TypeString,
						},
						"pattern": {
							Description: "The wildcard privilege pattern",
							Computed:    true,
							Type:        schema.TypeBool,
						},
						"read_only": {
							Description: "Indicates whether the privilege can be changed. External values supplied to this will be ignored by the system.",
							Computed:    true,
							Type:        schema.TypeBool,
						},
						"repository": {
							Description: "The repository of the privilege",
							Computed:    true,
							Type:        schema.TypeString,
						},
						"type": {
							Description: "The type of the privilege",
							Computed:    true,
							Type:        schema.TypeString,
						},
					},
				},
				Type: schema.TypeList,
			},
			"repository": {
				Description: "The repository of the privilege",
				ForceNew:    true,
				Optional:    true,
				Type:        schema.TypeString,
			},
			"type": {
				Description: "The type of the privilege",
				ForceNew:    true,
				Optional:    true,
				Type:        schema.TypeString,
				// ValidateFunc: validation.StringInSlice(nexus.PrivilegeTypes, false),
			},
		},
	}
}

func dataSourcePrivilegesRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	privileges, err := client.Security.Privilege.List()
	if err != nil {
		return err
	}

	dsDomain := d.Get("domain").(string)
	dsFormat := d.Get("format").(string)
	dsName := d.Get("name").(string)
	dsRepository := d.Get("repository").(string)
	dsType := d.Get("type").(string)

	d.SetId(fmt.Sprintf("%s-%s-%s-%s-%s", dsDomain, dsFormat, dsName, dsRepository, dsType))
	d.Set("domain", dsDomain)
	d.Set("format", dsFormat)
	d.Set("name", dsName)
	d.Set("repository", dsRepository)
	d.Set("type", dsType)

	var filteredPrivileges []security.Privilege
	if filteredPrivileges, err = filterPrivileges(privileges, dsDomain, dsFormat, dsName, dsRepository, dsType); err != nil {
		return err
	}

	if err := d.Set("privileges", flattenPrivileges(filteredPrivileges)); err != nil {
		return err
	}

	return nil
}

func filterPrivileges(privileges []security.Privilege, filterDomain, filterFormat, filterName, filterRepository, filterType string) ([]security.Privilege, error) {
	domains := make([]int, 0)
	formats := make([]int, 0)
	repositories := make([]int, 0)
	types := make([]int, 0)

	// Filter
	// Golang OR operator precedence is left to right
	// https://www.tutorialspoint.com/go/go_operators_precedence.htm
	for i, priv := range privileges {
		// filter by domain
		if filterDomain == "" || priv.Domain == filterDomain {
			domains = append(domains, i)
		}

		// filter by format
		if filterFormat == "" || priv.Format == filterFormat {
			formats = append(formats, i)
		}

		// filter by repository
		if filterRepository == "" || priv.Repository == filterRepository {
			repositories = append(repositories, i)
		}

		// filter by type
		if filterType == "" || priv.Type == filterType {
			types = append(types, i)
		}
	}

	intSlice := tools.Intersection(tools.Intersection(tools.Intersection(domains, repositories), formats), types)

	result := make([]security.Privilege, len(intSlice))
	for i, v := range intSlice {
		result[i] = privileges[v]
	}

	return result, nil
}

func flattenPrivileges(privileges []security.Privilege) []map[string]interface{} {
	if privileges == nil {
		return nil
	}

	data := make([]map[string]interface{}, len(privileges))
	for i, priv := range privileges {
		data[i] = map[string]interface{}{
			"actions":          priv.Actions,
			"content_selector": priv.ContentSelector,
			"description":      priv.Description,
			"domain":           priv.Domain,
			"format":           priv.Format,
			"name":             priv.Name,
			"read_only":        priv.ReadOnly,
			"repository":       priv.Repository,
			"type":             priv.Type,
		}
	}

	return data
}
