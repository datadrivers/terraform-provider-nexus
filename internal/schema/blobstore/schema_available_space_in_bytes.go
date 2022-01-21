package blobstore

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	ResourceAvailableSpaceInBytes = &schema.Schema{
		Computed:    true,
		Description: "Available space in Bytes",
		Type:        schema.TypeInt,
	}

	DataSourceAvailableSpaceInBytes = ResourceAvailableSpaceInBytes
)
