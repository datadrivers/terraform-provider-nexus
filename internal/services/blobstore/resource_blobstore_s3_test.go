package blobstore_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/datadrivers/go-nexus-client/nexus3/schema/blobstore"
	"github.com/datadrivers/terraform-provider-nexus/internal/acceptance"
	"github.com/datadrivers/terraform-provider-nexus/internal/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceBlobstoreS3(t *testing.T) {
	if tools.GetEnv("SKIP_S3_TESTS", "false") == "true" {
		t.Skip("Skipping S3 tests")
	}

	resourceName := "nexus_blobstore_s3.acceptance"
	awsAccessKeyID := tools.GetEnv("AWS_ACCESS_KEY_ID", "")
	awsSecretAccessKey := tools.GetEnv("AWS_SECRET_ACCESS_KEY", "")
	forcePathStyle := true

	bs := blobstore.S3{
		Name: fmt.Sprintf("test-blobstore-s3-%s", acctest.RandString(5)),
		BucketConfiguration: blobstore.S3BucketConfiguration{
			Bucket: blobstore.S3Bucket{
				Name:       tools.GetEnv("AWS_BUCKET_NAME", "terraform-provider-nexus-s3-test"),
				Region:     tools.GetEnv("AWS_DEFAULT_REGION", "eu-central-1"),
				Expiration: 0,
			},
			AdvancedBucketConnection: &blobstore.S3AdvancedBucketConnection{
				Endpoint:       tools.GetEnv("AWS_ENDPOINT", ""),
				ForcePathStyle: &forcePathStyle,
			},
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceBlobstoreTypeS3Config(bs, awsAccessKeyID, awsSecretAccessKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", bs.Name),
					resource.TestCheckResourceAttrSet(resourceName, "blob_count"),
					resource.TestCheckResourceAttrSet(resourceName, "total_size_in_bytes"),
					resource.TestCheckResourceAttr(resourceName, "bucket_configuration.0.bucket.0.name", bs.BucketConfiguration.Bucket.Name),
					resource.TestCheckResourceAttr(resourceName, "bucket_configuration.0.bucket.0.region", bs.BucketConfiguration.Bucket.Region),
					resource.TestCheckResourceAttr(resourceName, "bucket_configuration.0.bucket.0.expiration", strconv.FormatInt(int64(bs.BucketConfiguration.Bucket.Expiration), 10)),
					resource.TestCheckResourceAttr(resourceName, "bucket_configuration.0.bucket_security.0.access_key_id", awsAccessKeyID),
					resource.TestCheckResourceAttr(resourceName, "bucket_configuration.0.bucket_security.0.secret_access_key", awsSecretAccessKey),
					resource.TestCheckResourceAttr(resourceName, "bucket_configuration.0.advanced_bucket_connection.0.endpoint", bs.BucketConfiguration.AdvancedBucketConnection.Endpoint),
					resource.TestCheckResourceAttr(resourceName, "bucket_configuration.0.advanced_bucket_connection.0.force_path_style", strconv.FormatBool(*bs.BucketConfiguration.AdvancedBucketConnection.ForcePathStyle)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     bs.Name,
				ImportStateVerify: true,
				// available_space_in_bytes changes too frequently.
				ImportStateVerifyIgnore: []string{"available_space_in_bytes", "bucket_configuration.0.bucket_security.0.secret_access_key"},
			},
		},
	})
}

func testAccResourceBlobstoreTypeS3Config(bs blobstore.S3, awsAccessKeyID string, awsSecretAccessKey string) string {
	return fmt.Sprintf(`
resource "nexus_blobstore_s3" "acceptance" {
	name = "%s"

	bucket_configuration {
		bucket {
		  name       = "%s"
		  region     = "%s"
		  expiration = %d
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
}`, bs.Name, bs.BucketConfiguration.Bucket.Name, bs.BucketConfiguration.Bucket.Region, bs.BucketConfiguration.Bucket.Expiration, awsAccessKeyID, awsSecretAccessKey, bs.BucketConfiguration.AdvancedBucketConnection.Endpoint, strconv.FormatBool(*bs.BucketConfiguration.AdvancedBucketConnection.ForcePathStyle))
}
