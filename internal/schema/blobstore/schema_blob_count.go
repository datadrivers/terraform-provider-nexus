package blobstore

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	ResourceBlobCount = &schema.Schema{
		Computed:    true,
		Description: "Count of blobs",
		Type:        schema.TypeInt,
	}

	DataSourceBlobCount = ResourceBlobCount
)
