package blobstore

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	ResourceTotalSizeInBytes = &schema.Schema{
		Computed:    true,
		Description: "The total size of the blobstore in Bytes",
		Type:        schema.TypeInt,
	}

	DataSourceTotalSizeInBytes = ResourceTotalSizeInBytes
)
