package nexus

import (
	"fmt"
	"testing"

	nexus "github.com/datadrivers/go-nexus-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccResourceBlobstoreFile(t *testing.T) {
	bsName := fmt.Sprintf("test-blobstore-%d", acctest.RandIntRange(0, 99))
	bsType := nexus.BlobstoreTypeFile
	bsPath := fmt.Sprintf("/nexus-data/%s", bsName)
	quotaLimit := acctest.RandIntRange(100, 300)
	quotaType := "spaceRemainingQuota"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBlobstoreResourceFile(bsName, bsType, bsPath, quotaLimit, quotaType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("nexus_blobstore.acceptance", "name", bsName),
					resource.TestCheckResourceAttr("nexus_blobstore.acceptance", "type", bsType),
					resource.TestCheckResourceAttr("nexus_blobstore.acceptance", "path", bsPath),
				),
			},
			{
				ResourceName:      "nexus_blobstore.acceptance",
				ImportState:       true,
				ImportStateId:     bsName,
				ImportStateVerify: true,
				// available_space_in_bytes changes too frequently.
				ImportStateVerifyIgnore: []string{"available_space_in_bytes"},
			},
		},
	})
}

func testAccBlobstoreResourceFile(name string, bsType string, path string, quotaLimit int, quotaType string) string {
	return fmt.Sprintf(`
resource "nexus_blobstore" "acceptance" {
	name = "%s"
	path = "%s"
	type = "%s"

	soft_quota {
		limit = %d
		type  = "%s"
	}
}`, name, path, bsType, quotaLimit, quotaType)
}

func TestAccResourceBlobstoreS3(t *testing.T) {
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
