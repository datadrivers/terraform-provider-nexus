package nexus

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	nexus "github.com/datadrivers/go-nexus-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccResourceBlobstoreS3(t *testing.T) {
	if strings.ToLower(os.Getenv("SKIP_GOOGLE_TESTS")) != "false" {
		t.Skip("Skipping Google tests")
	}

	resName := "nexus_blobstore.acceptance"
	googleCrednentialFilePath := getEnv("GOOGLE_CREDENTIAL_FILE_PATH", "")

	bs := nexus.Blobstore{
		Name: fmt.Sprintf("test-blobstore-google-%d", acctest.RandIntRange(0, 99)),
		Type: nexus.BlobstoreTypeGoogle,
		BucketName: getEnv("GOOGLE_BUCKET_NAME", "terraform-provider-nexus-google-test"),
		Region: getEnv("GOOGLE_DEFAULT_REGION", "us-central1").(string),
		CredentialFilePath: googleCrednentialFilePath,,
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceBlobstoreTypeGoogleConfig(bs, googleCrednentialFilePath),
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
				ImportStateVerifyIgnore: []string{"available_space_in_bytes"},
			},
		},
	})
}

func testAccResourceBlobstoreTypeGoogleConfig(bs nexus.Blobstore, googleCrednentialFilePath string) string {
	return fmt.Sprintf(`
resource "nexus_blobstore" "acceptance" {
	name = "%s"
	type = "%s"

	bucket_name = "%s"
	region = "%s"
	credential_file_path = "%s"
}`, bs.Name, bs.Type, bs.BucketName, bs.Region, googleCrednentialFilePath)
}
