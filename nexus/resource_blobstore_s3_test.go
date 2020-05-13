package nexus

import (
	"fmt"
	"os"
	"testing"

	nexus "github.com/datadrivers/go-nexus-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccResourceBlobstoreS3(t *testing.T) {
	if os.Getenv("SKIP_S3_TESTS") != "" {
		t.Skip("Skipping S3 tests")
	}
	awsAccessKeyID := getEnv("AWS_ACCESS_KEY_ID", "")
	awsSecretAccessKey := getEnv("AWS_SECRET_ACCESS_KEY", "")
	bsName := fmt.Sprintf("test-blobstore-s3-%d", acctest.RandIntRange(0, 99))
	bsType := nexus.BlobstoreTypeS3
	bucketName := getEnv("AWS_BUCKET_NAME", "terraform-provider-nexus-s3-test")
	bucketRegion := getEnv("AWS_DEFAULT_REGION", "eu-central-1")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBlobstoreResourceS3Minimal(bsName, bsType, bucketName, bucketRegion, awsAccessKeyID, awsSecretAccessKey),
				// FIXME: Increase test coverage
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("nexus_blobstore.acceptance", "name", bsName),
					resource.TestCheckResourceAttr("nexus_blobstore.acceptance", "type", bsType),
				),
			},
			{
				ResourceName:      "nexus_blobstore.acceptance",
				ImportState:       true,
				ImportStateId:     bsName,
				ImportStateVerify: true,
				// available_space_in_bytes changes too frequently.
				ImportStateVerifyIgnore: []string{"available_space_in_bytes", "bucket_configuration.0.bucket_security.0.secret_access_key"},
			},
		},
	})
}

func testAccBlobstoreResourceS3Minimal(name string, bsType string, bucketName string, bucketRegion string, awsAccessKeyID string, awsSecretAccessKey string) string {
	return fmt.Sprintf(`
resource "nexus_blobstore" "acceptance" {
	name = "%s"
	type = "%s"

	bucket_configuration {
		bucket {
		  name   = "%s"
		  region = "%s"
		}

		bucket_security {
		  access_key_id     = "%s"
		  secret_access_key = "%s"
		}
	}
}`, name, bsType, bucketName, bucketRegion, awsAccessKeyID, awsSecretAccessKey)
}
