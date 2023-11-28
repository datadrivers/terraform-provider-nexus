package repository

import (
	nexus "github.com/dre2004/go-nexus-client/nexus3"
	"github.com/dre2004/terraform-provider-nexus/internal/schema/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceRepositoryList() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get a list with all repositories.",

		Read: dataSourceRepositoryList,
		Schema: map[string]*schema.Schema{
			"id": common.DataSourceID,
			"items": {
				Description: "A List of all repositories",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "A unique identifier for this repository",
							Computed:    true,
							Type:        schema.TypeString,
						},
						"format": {
							Description: "Repository format",
							Computed:    true,
							Type:        schema.TypeString,
						},
						"type": {
							Description: "Repository type",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"url": {
							Computed:    true,
							Description: "The URL of the repository",
							Type:        schema.TypeString,
						},
					},
				},
			},
		},
	}
}

func dataSourceRepositoryList(dataSource *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	items := []map[string]string{}
	repositories, err := client.Repository.List()
	if err != nil {
		return err
	}

	for _, repository := range repositories {
		items = append(items, map[string]string{
			"name":   repository.Name,
			"format": repository.Format,
			"type":   repository.Type,
			"url":    repository.URL,
		})
	}
	if err := dataSource.Set("items", items); err != nil {
		return err
	}
	dataSource.SetId("repositoryList")
	return nil
}
