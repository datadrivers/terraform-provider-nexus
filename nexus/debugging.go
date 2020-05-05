package nexus

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Helper function to pretty-print tf state
// Usage: put it as first function inside resource.ComposeTestCheckFunc()
// it will print the statee during test
func printState(s *terraform.State) error {
	res, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(res))
	return nil
}
