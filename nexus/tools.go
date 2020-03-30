package nexus

import (
	"fmt"
	"os"

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

func interfaceSliceToStringSlice(data []interface{}) []string {
	result := make([]string, len(data))
	for i, v := range data {
		result[i] = v.(string)
	}
	return result
}

func stringSliceToInterfaceSlice(strings []string) []interface{} {
	s := make([]interface{}, len(strings))
	for i, v := range strings {
		s[i] = string(v)
	}
	return s
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}
