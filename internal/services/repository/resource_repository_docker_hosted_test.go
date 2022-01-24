package repository_test

import (
	"bytes"
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"text/template"

	"github.com/datadrivers/go-nexus-client/nexus3/schema/repository"
	"github.com/datadrivers/terraform-provider-nexus/internal/acceptance"
	"github.com/datadrivers/terraform-provider-nexus/internal/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	resourceRepositoryDockerHostedTemplateString = `
resource "nexus_repository_docker_hosted" "acceptance" {
	name   = "{{ .Name }}"
	online = {{ .Online }}

	docker {
		force_basic_auth = "{{ .Docker.ForceBasicAuth }}"
{{- if .Docker.HTTPPort }}
		http_port = "{{ .Docker.HTTPPort }}"
{{- end }}
{{- if .Docker.HTTPSPort }}
		https_port = "{{ .Docker.HTTPSPort }}"
{{- end }}
		v1_enabled = "{{ .Docker.V1Enabled }}"
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

func testAccResourceRepositoryDockerHosted() repository.DockerHostedRepository {
	writePolicy := repository.StorageWritePolicyAllow

	return repository.DockerHostedRepository{
		Name:   fmt.Sprintf("test-repo-%s", acctest.RandString(10)),
		Online: true,
		Docker: repository.Docker{
			ForceBasicAuth: true,
			HTTPPort:       tools.GetIntPointer(rand.Intn(999) + 32000),
			HTTPSPort:      tools.GetIntPointer(rand.Intn(999) + 33000),
			V1Enabled:      false,
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

func testAccResourceRepositoryDockerHostedConfig(repo repository.DockerHostedRepository) string {
	buf := &bytes.Buffer{}
	resourceRepositoryDockerHostedTemplate := template.Must(template.New("DockerHostedRepository").Funcs(acceptance.TemplateFuncMap).Parse(resourceRepositoryDockerHostedTemplateString))
	if err := resourceRepositoryDockerHostedTemplate.Execute(buf, repo); err != nil {
		panic(err)
	}
	return buf.String()
}

func TestAccResourceRepositoryDockerHosted(t *testing.T) {
	repo := testAccResourceRepositoryDockerHosted()
	resourceName := "nexus_repository_docker_hosted.acceptance"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceRepositoryDockerHostedConfig(repo),
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
						resource.TestCheckResourceAttr(resourceName, "docker.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "docker.0.force_basic_auth", strconv.FormatBool(repo.Docker.ForceBasicAuth)),
						resource.TestCheckResourceAttr(resourceName, "docker.0.http_port", strconv.Itoa(*repo.Docker.HTTPPort)),
						resource.TestCheckResourceAttr(resourceName, "docker.0.https_port", strconv.Itoa(*repo.Docker.HTTPSPort)),
						resource.TestCheckResourceAttr(resourceName, "docker.0.v1_enabled", strconv.FormatBool(repo.Docker.V1Enabled)),
					),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateId:     repo.Name,
				ImportState:       true,
				ImportStateVerify: true,
				// Signing block is not returned
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}
