package blobstore_test

import (
	"bytes"
	"fmt"
	"strconv"
	"testing"
	"text/template"

	"github.com/datadrivers/terraform-provider-nexus/internal/acceptance"
	"github.com/gcroucher/go-nexus-client/nexus3/schema/blobstore"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	resourceBlobstoreFileTemplateString = `
resource "nexus_blobstore_file" "acceptance" {
	name = "{{ .Name }}"
	path = "{{ .Path }}"
{{- if .SoftQuota }}
	soft_quota {
		limit = {{ .SoftQuota.Limit }}
		type  = "{{ .SoftQuota.Type }}"
	}
{{- end }}
}`
)

func testAccResourceBlobstoreFileConfig(bs blobstore.File) string {
	buf := &bytes.Buffer{}
	resourceTemplate := template.Must(template.New("BlobstoreFile").Funcs(acceptance.TemplateFuncMap).Parse(resourceBlobstoreFileTemplateString))
	if err := resourceTemplate.Execute(buf, bs); err != nil {
		panic(err)
	}
	return buf.String()
}

func TestAccResourceBlobstoreFile(t *testing.T) {
	resourceName := "nexus_blobstore_file.acceptance"

	bs := blobstore.File{
		Name: fmt.Sprintf("test-blobstore-%s", acctest.RandString(5)),
		Path: "/nexus-data/acceptance",
		SoftQuota: &blobstore.SoftQuota{
			Limit: int64(acctest.RandIntRange(100, 300) * 1000000),
			Type:  "spaceRemainingQuota",
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceBlobstoreFileConfig(bs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", bs.Name),
					resource.TestCheckResourceAttr(resourceName, "name", bs.Name),
					resource.TestCheckResourceAttr(resourceName, "path", bs.Path),
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "soft_quota.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "soft_quota.0.limit", strconv.FormatInt(bs.SoftQuota.Limit, 10)),
						resource.TestCheckResourceAttr(resourceName, "soft_quota.0.type", bs.SoftQuota.Type),
					),
					resource.TestCheckResourceAttrSet(resourceName, "blob_count"),
					resource.TestCheckResourceAttrSet(resourceName, "total_size_in_bytes"),
					resource.TestCheckResourceAttrSet(resourceName, "available_space_in_bytes"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateId:           bs.Name,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"available_space_in_bytes"},
			},
		},
	})
}
