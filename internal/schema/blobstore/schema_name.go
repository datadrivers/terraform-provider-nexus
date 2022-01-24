package blobstore

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	ResourceName = &schema.Schema{
		Description: "Blobstore name",
		Required:    true,
		Type:        schema.TypeString,
	}

	DataSourceName = ResourceName
)
