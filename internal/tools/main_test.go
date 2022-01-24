package tools

import (
	"hash/fnv"
	"sort"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestInterfaceSliceToStringSlice(t *testing.T) {
	input := []interface{}{"foo", "bar"}
	output := InterfaceSliceToStringSlice(input)

	for i := range input {
		assert.Equal(t, input[i], output[i])
	}

}

func TestStringSliceToInterfaceSlice(t *testing.T) {
	input := []string{"foo", "bar"}
	output := StringSliceToInterfaceSlice(input)

	for i := range input {
		assert.Equal(t, input[i], output[i])
	}
}

func TestConvertStringSet(t *testing.T) {
	testStrings := []string{
		"test",
		"blub",
		"bla",
	}
	set := schema.NewSet(func(s interface{}) int {
		h := fnv.New32a()
		h.Write([]byte(s.(string)))
		return int(h.Sum32())
	}, []interface{}{})
	for _, value := range testStrings {
		set.Add(value)
	}
	convertedSet := ConvertStringSet(set)
	sort.Strings(testStrings)

	assert.Equal(t, testStrings, convertedSet)
}
