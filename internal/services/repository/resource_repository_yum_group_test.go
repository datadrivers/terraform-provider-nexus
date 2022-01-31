package repository_test

import (
	"bytes"
	"fmt"
	"strconv"
	"testing"
	"text/template"

	"github.com/datadrivers/go-nexus-client/nexus3/pkg/tools"
	"github.com/datadrivers/go-nexus-client/nexus3/schema/repository"
	"github.com/datadrivers/terraform-provider-nexus/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	resourceRepositoryYumGroupTemplateString = `
resource "nexus_repository_yum_group" "acceptance" {
	name   = "{{ .Name }}"
	online = {{ .Online }}

	storage {
		blob_store_name                = "{{ .Storage.BlobStoreName }}"
		strict_content_type_validation = {{ .Storage.StrictContentTypeValidation }}
	}

	group {
		member_names = [
		{{- range $val := .Group.MemberNames }}
			"{{ $val }}",
		{{ end -}}
		]
	}

{{ if .YumSigning }}
	yum_signing {
		keypair = "{{ .YumSigning.Keypair }}"
{{ if .YumSigning.Passphrase }}
		passphrase = "{{ .YumSigning.Passphrase }}"
{{ end -}}
	}
{{ end -}}

	depends_on = [
		nexus_repository_yum_hosted.acceptance
	]
}
`
)

func testAccResourceRepositoryYumGroup() repository.YumGroupRepository {
	return repository.YumGroupRepository{
		Name:   fmt.Sprintf("test-repo-%s", acctest.RandString(10)),
		Online: true,
		YumSigning: &repository.YumSigning{
			Keypair:    tools.GetStringPointer("test-keypair"),
			Passphrase: tools.GetStringPointer("test-keypair"),
		},
		Storage: repository.Storage{
			BlobStoreName:               "default",
			StrictContentTypeValidation: true,
		},
		Group: repository.Group{
			MemberNames: []string{},
		},
	}
}

func testAccResourceRepositoryYumGroupConfig(repo repository.YumGroupRepository) string {
	buf := &bytes.Buffer{}
	resourceRepositoryYumGroupTemplate := template.Must(template.New("YumGroupRepository").Funcs(acceptance.TemplateFuncMap).Parse(resourceRepositoryYumGroupTemplateString))
	if err := resourceRepositoryYumGroupTemplate.Execute(buf, repo); err != nil {
		panic(err)
	}
	return buf.String()
}

func TestAccResourceRepositoryYumGroup(t *testing.T) {
	repoHosted := testAccResourceRepositoryYumHosted()
	repo := testAccResourceRepositoryYumGroup()
	repo.Group.MemberNames = append(repo.Group.MemberNames, repoHosted.Name)
	resourceName := "nexus_repository_yum_group.acceptance"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceRepositoryYumHostedConfig(repoHosted) + testAccResourceRepositoryYumGroupConfig(repo),
				Check: resource.ComposeTestCheckFunc(
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "id", repo.Name),
						resource.TestCheckResourceAttr(resourceName, "name", repo.Name),
						resource.TestCheckResourceAttr(resourceName, "online", strconv.FormatBool(repo.Online)),
					),
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "storage.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "storage.0.blob_store_name", repo.Storage.BlobStoreName),
						resource.TestCheckResourceAttr(resourceName, "storage.0.strict_content_type_validation", strconv.FormatBool(repo.Storage.StrictContentTypeValidation)),
						resource.TestCheckResourceAttr(resourceName, "group.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "group.0.member_names.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "group.0.member_names.0", repo.Group.MemberNames[0]),
					),
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "yum_signing.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "yum_signing.0.keypair", *repo.YumSigning.Keypair),
						resource.TestCheckResourceAttr(resourceName, "yum_signing.0.passphrase", *repo.YumSigning.Passphrase),
					),
				),
			},
			{
				ResourceName:            resourceName,
				ImportStateId:           repo.Name,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"yum_signing"},
			},
		},
	})
}
