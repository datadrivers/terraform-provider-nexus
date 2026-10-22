package repository

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// RequireRepositoryFound makes a repository data source fail when the repository does not exist.
// The shared resource read clears the ID for a missing repository, so that a resource is recreated;
// a data source must report it instead of returning empty attributes.
func RequireRepositoryFound(dataSource *schema.Resource) *schema.Resource {
	read := dataSource.Read
	dataSource.Read = func(d *schema.ResourceData, m interface{}) error {
		if err := read(d, m); err != nil {
			return err
		}
		if d.Id() == "" {
			return fmt.Errorf("repository %q not found", d.Get("name"))
		}
		return nil
	}
	return dataSource
}
