package tools

import (
	"testing"

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
