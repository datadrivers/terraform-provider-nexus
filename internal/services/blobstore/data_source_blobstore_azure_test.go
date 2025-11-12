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

func TestAccDataSourceBlobstoreAzure(t *testing.T) {
	if tools.GetEnv("SKIP_AZURE_TESTS", "false") == "true" {
		t.Skip("Skipping Nexus blobstore for Azure tests")
	}
	if tools.GetEnv("SKIP_PRO_TESTS", "false") == "true" {
		t.Skip("Skipping Nexus Pro tests")
	}

	dataSourceName := "data.nexus_blobstore_azure.acceptance"

	bs := blobstore.Azure{
		Name: fmt.Sprintf("test-blobstore-azure-%s", acctest.RandString(5)),
		BucketConfiguration: blobstore.AzureBucketConfiguration{
			AccountName: "terraformprovidernexus",
			Authentication: blobstore.AzureBucketConfigurationAuthentication{
				AuthenticationMethod: blobstore.AzureAuthenticationMethodAccountKey,
				AccountKey:           tools.GetEnv("AZURE_STORAGE_ACCOUNT_KEY", "test-key"),
			},
			ContainerName: "datasource-acceptance",
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceBlobstoreTypeAzureConfig(bs) + testAccDataSourceBlobstoreTypeAzureConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", bs.Name),
					resource.TestCheckResourceAttr(dataSourceName, "bucket_configuration.0.account_name", bs.BucketConfiguration.AccountName),
					resource.TestCheckResourceAttr(dataSourceName, "bucket_configuration.0.container_name", bs.BucketConfiguration.ContainerName),
					resource.TestCheckResourceAttr(dataSourceName, "bucket_configuration.0.authentication.0.authentication_method", string(bs.BucketConfiguration.Authentication.AuthenticationMethod)),
					resource.TestCheckResourceAttrSet(dataSourceName, "blob_count"),
					resource.TestCheckResourceAttrSet(dataSourceName, "total_size_in_bytes"),
				),
			},
		},
	})
}

func testAccDataSourceBlobstoreTypeAzureConfig() string {
	return `
data "nexus_blobstore_azure" "acceptance" {
	name = nexus_blobstore_azure.acceptance.name
}`
}
