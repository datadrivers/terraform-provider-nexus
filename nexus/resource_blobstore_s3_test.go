package nexus

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/datadrivers/go-nexus-client/nexus3/schema/blobstore"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceBlobstoreS3(t *testing.T) {
	if os.Getenv("SKIP_S3_TESTS") != "" {
		t.Skip("Skipping S3 tests")
	}

	resName := "nexus_blobstore.acceptance"
	awsAccessKeyID := getEnv("AWS_ACCESS_KEY_ID", "")
	awsSecretAccessKey := getEnv("AWS_SECRET_ACCESS_KEY", "")

	bs := blobstore.Legacy{
		Name: fmt.Sprintf("test-blobstore-s3-%d", acctest.RandIntRange(0, 99)),
		Type: blobstore.BlobstoreTypeS3,
		S3BucketConfiguration: &blobstore.S3BucketConfiguration{
			Bucket: blobstore.S3Bucket{
				Name:   getEnv("AWS_BUCKET_NAME", "terraform-provider-nexus-s3-test"),
				Region: getEnv("AWS_DEFAULT_REGION", "eu-central-1"),
			},
			AdvancedBucketConnection: &blobstore.S3AdvancedBucketConnection{
				Endpoint:       getEnv("AWS_ENDPOINT", ""),
				ForcePathStyle: true,
			},
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceBlobstoreTypeS3Config(bs, awsAccessKeyID, awsSecretAccessKey),
				// FIXME: Increase test coverage
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", bs.Name),
					resource.TestCheckResourceAttr(resName, "type", bs.Type),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateId:     bs.Name,
				ImportStateVerify: true,
				// available_space_in_bytes changes too frequently.
				ImportStateVerifyIgnore: []string{"available_space_in_bytes", "bucket_configuration.0.bucket_security.0.secret_access_key"},
			},
		},
	})
}

func testAccResourceBlobstoreTypeS3Config(bs blobstore.Legacy, awsAccessKeyID string, awsSecretAccessKey string) string {
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

		advanced_bucket_connection {
 		  endpoint			= "%s"
		  force_path_style	= %s
		}
	}
}`, bs.Name, bs.Type, bs.S3BucketConfiguration.Bucket.Name, bs.S3BucketConfiguration.Bucket.Region, awsAccessKeyID, awsSecretAccessKey, bs.S3BucketConfiguration.AdvancedBucketConnection.Endpoint, strconv.FormatBool(bs.S3BucketConfiguration.AdvancedBucketConnection.ForcePathStyle))
}
