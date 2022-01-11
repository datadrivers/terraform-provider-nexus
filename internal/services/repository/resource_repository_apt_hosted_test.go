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
	resourceRepositoryAptHostedTemplateString = `
resource "nexus_repository_apt_hosted" "acceptance" {
	name   = "{{ .Name }}"
	online = {{ .Online }}

	distribution = "{{ .Apt.Distribution }}"
	signing {
		keypair = "{{ .AptSigning.Keypair }}"
{{- if .AptSigning.Passphrase }}
		passphrase = "{{ .AptSigning.Passphrase }}"
{{- end }}
	}

	storage {
		blob_store_name                = "{{ .Storage.BlobStoreName }}"
		strict_content_type_validation = {{ .Storage.StrictContentTypeValidation }}
		{{- if .Storage.WritePolicy }}
		write_policy                   = "{{ .Storage.WritePolicy }}"
		{{- end }}
	}

{{ if .Cleanup }}
	cleanup {
		policy_names = [
		{{- range $val := .Cleanup.PolicyNames }}
			"{{ $val }}",
		{{ end -}}
		]
	}
{{ end -}}
{{ if .Component }}
	component {
		proprietary_components = {{ .Component.ProprietaryComponents }}
	}
{{ end -}}
}
`
)

func testAccResourceRepositoryAptHosted() repository.AptHostedRepository {
	writePolicy := repository.StorageWritePolicyAllow

	return repository.AptHostedRepository{
		Name:   fmt.Sprintf("test-repo-%s", acctest.RandString(10)),
		Online: true,
		Apt: repository.AptHosted{
			Distribution: "buster",
		},
		AptSigning: repository.AptSigning{
			Keypair:    "test-keypair",
			Passphrase: tools.GetStringPointer("test-passphrase"),
		},
		Storage: repository.HostedStorage{
			BlobStoreName:               "default",
			StrictContentTypeValidation: true,
			WritePolicy:                 &writePolicy,
		},
		Cleanup: &repository.Cleanup{
			PolicyNames: []string{"cleanup-weekly"},
		},
		Component: &repository.Component{
			ProprietaryComponents: true,
		},
	}
}

func testAccResourceRepositoryAptHostedConfig(repo repository.AptHostedRepository) string {
	buf := &bytes.Buffer{}
	resourceRepositoryAptHostedTemplate := template.Must(template.New("AptHostedRepository").Funcs(acceptance.TemplateFuncMap).Parse(resourceRepositoryAptHostedTemplateString))
	if err := resourceRepositoryAptHostedTemplate.Execute(buf, repo); err != nil {
		panic(err)
	}
	return buf.String()
}

func TestAccResourceRepositoryAptHosted(t *testing.T) {
	repo := testAccResourceRepositoryAptHosted()
	resourceName := "nexus_repository_apt_hosted.acceptance"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceRepositoryAptHostedConfig(repo),
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
						resource.TestCheckResourceAttr(resourceName, "storage.0.write_policy", string(*repo.Storage.WritePolicy)),
						resource.TestCheckResourceAttr(resourceName, "cleanup.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "cleanup.0.policy_names.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "cleanup.0.policy_names.0", repo.Cleanup.PolicyNames[0]),
						resource.TestCheckResourceAttr(resourceName, "component.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "component.0.proprietary_components", strconv.FormatBool(repo.Component.ProprietaryComponents)),
					),
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "distribution", repo.Apt.Distribution),
						resource.TestCheckResourceAttr(resourceName, "signing.0.keypair", repo.AptSigning.Keypair),
						resource.TestCheckResourceAttr(resourceName, "signing.0.passphrase", string(*repo.AptSigning.Passphrase)),
					),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateId:     repo.Name,
				ImportState:       true,
				ImportStateVerify: true,
				// Signing block is not returned
				ImportStateVerifyIgnore: []string{"signing"},
			},
		},
	})
}
