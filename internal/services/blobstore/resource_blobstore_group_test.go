package blobstore_test

import (
	"bytes"
	"fmt"
	"strconv"
	"testing"
	"text/template"

	"github.com/datadrivers/terraform-provider-nexus/internal/acceptance"
	"github.com/datadrivers/terraform-provider-nexus/internal/tools"
	"github.com/gcroucher/go-nexus-client/nexus3/schema/blobstore"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	resourceBlobstoreGroupTemplateString = `
resource "nexus_blobstore_group" "acceptance" {
	name        = "{{ .Name }}"
    fill_policy = "{{ .FillPolicy }}"
	members = [
		nexus_blobstore_file.acceptance.name
	]

{{- if .SoftQuota }}
	soft_quota {
		limit = {{ .SoftQuota.Limit }}
		type  = "{{ .SoftQuota.Type }}"
	}
{{- end }}
}`
)

func testAccResourceBlobstoreGroupConfig(bs blobstore.Group) string {
	buf := &bytes.Buffer{}
	resourceTemplate := template.Must(template.New("BlobstoreGroup").Funcs(acceptance.TemplateFuncMap).Parse(resourceBlobstoreGroupTemplateString))
	if err := resourceTemplate.Execute(buf, bs); err != nil {
		panic(err)
	}
	return buf.String()
}

func TestAccResourceBlobstoreGroup(t *testing.T) {
	if tools.GetEnv("SKIP_PRO_TESTS", "false") == "true" {
		t.Skip("Skipping Nexus Pro tests")
	}

	resourceName := "nexus_blobstore_group.acceptance"

	memberBlobStore := blobstore.File{
		Name: fmt.Sprintf("test_file_%s", acctest.RandString(5)),
		Path: fmt.Sprintf("/nexus-data/test-file-%s", acctest.RandString(5)),
		SoftQuota: &blobstore.SoftQuota{
			Limit: int64(acctest.RandIntRange(100, 300) * 1000000),
			Type:  "spaceRemainingQuota",
		},
	}
	bs := blobstore.Group{
		Name: fmt.Sprintf("test-blobstore-%s", acctest.RandString(5)),
		SoftQuota: &blobstore.SoftQuota{
			Limit: int64(acctest.RandIntRange(100, 300) * 1000000),
			Type:  "spaceRemainingQuota",
		},
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
				Config: testAccResourceBlobstoreFileConfig(memberBlobStore) + testAccResourceBlobstoreGroupConfig(bs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", bs.Name),
					resource.TestCheckResourceAttr(resourceName, "name", bs.Name),
					resource.TestCheckResourceAttr(resourceName, "fill_policy", bs.FillPolicy),
					resource.TestCheckResourceAttr(resourceName, "members.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "members.0", bs.Members[0]),
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
