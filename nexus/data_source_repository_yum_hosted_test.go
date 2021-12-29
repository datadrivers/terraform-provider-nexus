package nexus

import (
	"fmt"
)

func testAccDataSourceRepositoryYumHostedConfig(name string) string {
	return fmt.Sprintf(`
data "nexus_repository_yum_hosted" "acceptance" {
	name   = "%s"
}`, name)
}
