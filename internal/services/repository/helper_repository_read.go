package repository

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const repositoryReadMaxAttempts = 10

var repositoryReadRetryDelay = 30 * time.Second

type repositoryReadFunc func(*schema.ResourceData, interface{}) error

func retryRepositoryRead(read repositoryReadFunc, resourceData *schema.ResourceData, m interface{}) error {
	var err error
	id := ""
	if resourceData != nil {
		id = resourceData.Id()
	}

	for attempt := 1; attempt <= repositoryReadMaxAttempts; attempt++ {
		err = read(resourceData, m)
		if err == nil && id != "" && resourceData.Id() == "" {
			// The read found no repository just after a write; on a cluster the write may not have propagated yet.
			resourceData.SetId(id)
			err = fmt.Errorf("repository %q not found after it was written", id)
		}
		if err == nil {
			return nil
		}

		if attempt < repositoryReadMaxAttempts {
			time.Sleep(repositoryReadRetryDelay)
		}
	}

	return err
}
