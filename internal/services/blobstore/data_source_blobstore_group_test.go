package blobstore_test

import (
	"fmt"
	"testing"

	"github.com/datadrivers/terraform-provider-nexus/internal/acceptance"
	"github.com/datadrivers/terraform-provider-nexus/internal/tools"
	"github.com/gcroucher/go-nexus-client/nexus3/schema/blobstore"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceBlobstoreGroup(t *testing.T) {
	if tools.GetEnv("SKIP_PRO_TESTS", "false") == "true" {
		t.Skip("Skipping Nexus Pro tests")
	}

	dataSourceName := "data.nexus_blobstore_group.acceptance"

	memberBlobStore := blobstore.File{
		Name: fmt.Sprintf("test_file_%s", acctest.RandString(5)),
		Path: fmt.Sprintf("/nexus-data/test-file-%s", acctest.RandString(5)),
	}
	bs := blobstore.Group{
		Name:       fmt.Sprintf("test-blobstore-%s", acctest.RandString(5)),
		FillPolicy: blobstore.GroupFillPolicyWriteToFirst,
		Members: []string{
			memberBlobStore.Name,
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceBlobstoreFileConfig(memberBlobStore) + testAccResourceBlobstoreGroupConfig(bs) + testAccDataSourceBlobstoreTypeGroupConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "id", bs.Name),
					resource.TestCheckResourceAttr(dataSourceName, "name", bs.Name),
					resource.TestCheckResourceAttr(dataSourceName, "fill_policy", bs.FillPolicy),
					resource.TestCheckResourceAttr(dataSourceName, "members.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "members.0", bs.Members[0]),
					resource.TestCheckResourceAttrSet(dataSourceName, "blob_count"),
					resource.TestCheckResourceAttrSet(dataSourceName, "total_size_in_bytes"),
					resource.TestCheckResourceAttrSet(dataSourceName, "available_space_in_bytes"),
				),
			},
		},
	})
}

func testAccDataSourceBlobstoreTypeGroupConfig() string {
	return `
data "nexus_blobstore_group" "acceptance" {
	name = nexus_blobstore_group.acceptance.name
}`
}
