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

func TestAccResourceBlobstoreAzure(t *testing.T) {
	if tools.GetEnv("SKIP_AZURE_TESTS", "false") == "true" {
		t.Skip("Skipping Nexus blobstore for Azure tests")
	}
	if tools.GetEnv("SKIP_PRO_TESTS", "false") == "true" {
		t.Skip("Skipping Nexus Pro tests")
	}

	resourceName := "nexus_blobstore_azure.acceptance"

	bs := blobstore.Azure{
		Name: fmt.Sprintf("test-blobstore-azure-%s", acctest.RandString(5)),
		BucketConfiguration: blobstore.AzureBucketConfiguration{
			AccountName: "terraformprovidernexus",
			Authentication: blobstore.AzureBucketConfigurationAuthentication{
				AuthenticationMethod: blobstore.AzureAuthenticationMethodAccountKey,
				AccountKey:           tools.GetEnv("AZURE_STORAGE_ACCOUNT_KEY", "test-key"),
			},
			ContainerName: "acceptance",
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceBlobstoreTypeAzureConfig(bs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", bs.Name),
					resource.TestCheckResourceAttrSet(resourceName, "blob_count"),
					resource.TestCheckResourceAttrSet(resourceName, "total_size_in_bytes"),
					resource.TestCheckResourceAttr(resourceName, "bucket_configuration.0.account_name", bs.BucketConfiguration.AccountName),
					resource.TestCheckResourceAttr(resourceName, "bucket_configuration.0.container_name", bs.BucketConfiguration.ContainerName),
					resource.TestCheckResourceAttr(resourceName, "bucket_configuration.0.authentication.0.authentication_method", string(bs.BucketConfiguration.Authentication.AuthenticationMethod)),
					resource.TestCheckResourceAttr(resourceName, "bucket_configuration.0.authentication.0.account_key", bs.BucketConfiguration.Authentication.AccountKey),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateId:           bs.Name,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"bucket_configuration.0.authentication.0.account_key"},
			},
		},
	})
}

func testAccResourceBlobstoreTypeAzureConfig(bs blobstore.Azure) string {
	return fmt.Sprintf(`
resource "nexus_blobstore_azure" "acceptance" {
	name = "%s"

	bucket_configuration {
		account_name = "%s"
		authentication {
			authentication_method			= "%s"
			account_key	= "%s"
		}
		container_name = "%s"
	}
}`, bs.Name, bs.BucketConfiguration.AccountName, bs.BucketConfiguration.Authentication.AuthenticationMethod, bs.BucketConfiguration.Authentication.AccountKey, bs.BucketConfiguration.ContainerName)
}
