package repository

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const repositoryReadMaxAttempts = 10

var repositoryReadRetryDelay = 30 * time.Second

type repositoryReadFunc func(*schema.ResourceData, interface{}) error

func retryRepositoryRead(read repositoryReadFunc, resourceData *schema.ResourceData, m interface{}) error {
	var err error

	for attempt := 1; attempt <= repositoryReadMaxAttempts; attempt++ {
		err = read(resourceData, m)
		if err == nil {
			return nil
		}

		if attempt < repositoryReadMaxAttempts {
			time.Sleep(repositoryReadRetryDelay)
		}
	}

	return err
}
