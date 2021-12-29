package nexus

import (
	"bytes"
	"fmt"
	"strconv"
	"testing"
	"text/template"

	"github.com/datadrivers/go-nexus-client/nexus3/schema/repository"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	resourceRepositoryYumHostedTemplateString = `
resource "nexus_repository_yum_hosted" "{{ .Name }}" {
	name   = "{{ .Name }}"
	online = {{ .Online }}

	{{- if .Yum.DeployPolicy }}
	deploy_policy  = "{{ .Yum.DeployPolicy }}"
	{{- end }}
	repodata_depth = {{ .Yum.RepodataDepth }}

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
			{{ $val }},
		{{ end -}}
		]
	}
{{ end -}}
}
`
)

var ()

func testAccResourceRepositoryYumHosted() repository.YumHostedRepository {
	writePolicy := repository.StorageWritePolicyAllow
	deployPolicy := repository.YumDeployPolicyPermissive
	return repository.YumHostedRepository{
		Name:   fmt.Sprintf("test-repo-%s", acctest.RandString(10)),
		Online: true,
		Yum: repository.Yum{
			DeployPolicy:  &deployPolicy,
			RepodataDepth: 0,
		},
		Storage: repository.HostedStorage{
			BlobStoreName:               "default",
			StrictContentTypeValidation: true,
			WritePolicy:                 &writePolicy,
		},
		Cleanup: &repository.Cleanup{
			PolicyNames: []string{"\"cleanup-weekly\""},
		},
	}
}

func resourceYumRepositoryTestCheckFunc(repo repository.YumHostedRepository) resource.TestCheckFunc {
	resName := fmt.Sprintf("nexus_repository_yum_hosted.%s", repo.Name)
	return resource.ComposeAggregateTestCheckFunc(
		resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr(resName, "id", repo.Name),
			resource.TestCheckResourceAttr(resName, "name", repo.Name),
			resource.TestCheckResourceAttr(resName, "online", strconv.FormatBool(repo.Online)),
		),
		resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr(resName, "storage.#", "1"),
			resource.TestCheckResourceAttr(resName, "storage.0.blob_store_name", repo.Storage.BlobStoreName),
			resource.TestCheckResourceAttr(resName, "storage.0.strict_content_type_validation", strconv.FormatBool(repo.Storage.StrictContentTypeValidation)),
		),
	)
}

func resourceYumRepositoryTypeHostedTestCheckFunc(repo repository.YumHostedRepository) resource.TestCheckFunc {
	resName := fmt.Sprintf("nexus_repository_yum_hosted.%s", repo.Name)
	return resource.ComposeAggregateTestCheckFunc(
		resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr(resName, "http_client.#", "0"),
			resource.TestCheckResourceAttr(resName, "group.#", "0"),
			resource.TestCheckResourceAttr(resName, "negative_cache.#", "0"),
			resource.TestCheckResourceAttr(resName, "proxy.#", "0"),
		),
		resource.TestCheckResourceAttr(resName, "storage.0.write_policy", string(*repo.Storage.WritePolicy)),
	)
}

func testAccResourceRepositoryYumHostedConfig(repo repository.YumHostedRepository) string {
	buf := &bytes.Buffer{}
	resourceRepositoryYumHostedTemplate := template.Must(template.New("repository").Funcs(resourceRepositoryTemplateFuncMap).Parse(resourceRepositoryYumHostedTemplateString))
	if err := resourceRepositoryYumHostedTemplate.Execute(buf, repo); err != nil {
		panic(err)
	}
	return buf.String()
}

func TestAccResourceRepositoryYumHosted(t *testing.T) {
	repo := testAccResourceRepositoryYumHosted()
	resName := fmt.Sprintf("nexus_repository_yum_hosted.%s", repo.Name)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceRepositoryYumHostedConfig(repo),
				Check: resource.ComposeTestCheckFunc(
					resourceYumRepositoryTestCheckFunc(repo),
					resourceYumRepositoryTypeHostedTestCheckFunc(repo),
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resName, "deploy_policy", string(*repo.DeployPolicy)),
						resource.TestCheckResourceAttr(resName, "repodata_depth", strconv.Itoa(repo.RepodataDepth)),
					),
				),
			},
			{
				ResourceName:      resName,
				ImportStateId:     repo.Name,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
