package repository

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestRequireRepositoryFound(t *testing.T) {
	newDataSource := func(found bool) *schema.Resource {
		return RequireRepositoryFound(&schema.Resource{
			Schema: map[string]*schema.Schema{"name": {Type: schema.TypeString, Required: true}},
			Read: func(d *schema.ResourceData, _ interface{}) error {
				d.SetId(d.Get("name").(string))
				if !found {
					d.SetId("")
				}
				return nil
			},
		})
	}

	t.Run("passes a found repository through", func(t *testing.T) {
		dataSource := newDataSource(true)
		d := schema.TestResourceDataRaw(t, dataSource.Schema, map[string]interface{}{"name": "example"})
		assert.NoError(t, dataSource.Read(d, nil))
		assert.Equal(t, "example", d.Id())
	})

	t.Run("fails for a missing repository", func(t *testing.T) {
		dataSource := newDataSource(false)
		d := schema.TestResourceDataRaw(t, dataSource.Schema, map[string]interface{}{"name": "example"})
		assert.EqualError(t, dataSource.Read(d, nil), `repository "example" not found`)
	})
}
