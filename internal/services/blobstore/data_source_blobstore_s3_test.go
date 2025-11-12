package blobstore_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/datadrivers/terraform-provider-nexus/internal/acceptance"
	"github.com/datadrivers/terraform-provider-nexus/internal/tools"
	"github.com/gcroucher/go-nexus-client/nexus3/schema/blobstore"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceBlobstoreS3(t *testing.T) {
	if tools.GetEnv("SKIP_S3_TESTS", "false") == "true" {
		t.Skip("Skipping S3 tests")
	}

	dataSourceName := "data.nexus_blobstore_s3.acceptance"
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
				Config: testAccResourceBlobstoreTypeS3Config(bs, awsAccessKeyID, awsSecretAccessKey) + testAccDataSourceBlobstoreTypeS3Config(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", bs.Name),
					resource.TestCheckResourceAttr(dataSourceName, "bucket_configuration.0.bucket.0.name", bs.BucketConfiguration.Bucket.Name),
					resource.TestCheckResourceAttr(dataSourceName, "bucket_configuration.0.bucket.0.region", bs.BucketConfiguration.Bucket.Region),
					resource.TestCheckResourceAttr(dataSourceName, "bucket_configuration.0.bucket.0.expiration", strconv.FormatInt(int64(bs.BucketConfiguration.Bucket.Expiration), 10)),
					resource.TestCheckResourceAttr(dataSourceName, "bucket_configuration.0.advanced_bucket_connection.0.endpoint", bs.BucketConfiguration.AdvancedBucketConnection.Endpoint),
					resource.TestCheckResourceAttr(dataSourceName, "bucket_configuration.0.advanced_bucket_connection.0.force_path_style", strconv.FormatBool(*bs.BucketConfiguration.AdvancedBucketConnection.ForcePathStyle)),
					resource.TestCheckResourceAttrSet(dataSourceName, "blob_count"),
					resource.TestCheckResourceAttrSet(dataSourceName, "total_size_in_bytes"),
				),
			},
		},
	})
}

func testAccDataSourceBlobstoreTypeS3Config() string {
	return `
data "nexus_blobstore_s3" "acceptance" {
	name = nexus_blobstore_s3.acceptance.name
}`
}
