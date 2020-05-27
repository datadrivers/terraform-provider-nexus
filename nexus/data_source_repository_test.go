package nexus

import (
	"fmt"
)

func testAccDataSourceRepository(name string) string {
	return fmt.Sprintf(`
data "nexus_repository" "acceptance" {
	name   = "%s"
}`, name)
}
