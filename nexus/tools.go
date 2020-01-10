package nexus

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceDataStringSlice(d *schema.ResourceData, attribute string) []string {
	n := d.Get(fmt.Sprintf("%s.#", attribute)).(int)
	data := make([]string, n)
	for i := 0; i < n; i++ {
		data[i] = d.Get(fmt.Sprintf("%s.%d", attribute, i)).(string)
	}
	return data
}
