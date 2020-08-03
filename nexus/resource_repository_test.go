package nexus

import (
	"bytes"
	"fmt"
	"strconv"
	"text/template"

	nexus "github.com/datadrivers/go-nexus-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	resourceRepositoryTemplateString = `
resource "nexus_repository" "{{ .Name }}" {
	format = "{{ .Format }}"
	name   = "{{ .Name }}"
	online = {{ .Online }}
	type   = "{{ .Type }}"

{{ if .RepositoryApt }}
	apt {
		distribution = "{{ .RepositoryApt.Distribution }}"
		flat         = {{ .RepositoryApt.Flat }}
	}
{{ end -}}

{{ if .RepositoryAptSigning }}
	apt_signing {
		keypair    = "{{ .RepositoryAptSigning.Keypair }}"
		passphrase = "{{ .RepositoryAptSigning.Passphrase }}"
	}
{{ end -}}

{{ if .RepositoryBower }}
	bower {
		rewrite_package_urls = {{ .RepositoryBower.RewritePackageUrls }}
	}
{{  end -}}

{{ if .RepositoryCleanup }}
	cleanup {
		policy_names = [
		{{- range $val := .RepositoryCleanup.PolicyNames }}
			{{ $val }},
		{{ end -}}
		]
	}
{{ end -}}

{{ if .RepositoryDocker }}
	docker {
		force_basic_auth = {{ .RepositoryDocker.ForceBasicAuth }}
		{{ if .RepositoryDocker.HTTPPort -}}
		http_port        = {{ deref .RepositoryDocker.HTTPPort }}
		{{ end -}}
		{{  if .RepositoryDocker.HTTPSPort -}}
		https_port       = {{ deref .RepositoryDocker.HTTPSPort }}
		{{ end -}}
		v1enabled        = {{ .RepositoryDocker.V1Enabled }}
	}
{{ end -}}

{{ if .RepositoryDockerProxy }}
	docker_proxy {
		index_type = "{{ .RepositoryDockerProxy.IndexType }}"
		index_url  = "{{ deref .RepositoryDockerProxy.IndexURL }}"
	}
{{ end -}}

{{ if .RepositoryGroup }}
	group {
		member_names = [
		{{- range $val := .RepositoryGroup.MemberNames }}
			{{ $val }},
		{{ end -}}
		]
	}
{{ end -}}

{{ if .RepositoryHTTPClient }}
	http_client {
		{{ if .RepositoryHTTPClient.Authentication -}}
		authentication {
			ntlm_domain = "{{ .RepositoryHTTPClient.Authentication.NTLMDomain }}"
			ntlm_host   = "{{ .RepositoryHTTPClient.Authentication.NTLMHost }}"
			{{ if .RepositoryHTTPClient.Authentication.Password -}}
			password    = "{{ deref .RepositoryHTTPClient.Authentication.Password }}"
			{{ end -}}
			type        = "{{ .RepositoryHTTPClient.Authentication.Type }}"
			{{ if .RepositoryHTTPClient.Authentication.Username -}}
			username    = "{{ deref .RepositoryHTTPClient.Authentication.Username }}"
			{{ end -}}
		}
		{{ end -}}
	}
{{ end -}}

{{ if .RepositoryMaven }}
	maven {
		layout_policy  = "{{ .RepositoryMaven.LayoutPolicy }}"
		version_policy = "{{ .RepositoryMaven.VersionPolicy }}"
	}
{{ end -}}

{{ if .RepositoryNegativeCache }}
	negative_cache {

	}
{{ end -}}

{{ if .RepositoryNugetProxy }}
	nuget_proxy {
		query_cache_item_max_age = {{ .RepositoryNugetProxy.QueryCacheItemMaxAge }}
	}
{{ end -}}

{{ if .RepositoryProxy }}
	proxy {
		remote_url = "{{ .RepositoryProxy.RemoteURL }}"
	}
{{ end -}}

{{ if .RepositoryStorage }}
	storage {
		blob_store_name                = "{{ .RepositoryStorage.BlobStoreName }}"
		strict_content_type_validation = {{ .RepositoryStorage.StrictContentTypeValidation }}
		{{- if eq .Type "hosted" }}
		write_policy                   = "{{ .RepositoryStorage.WritePolicy }}"
		{{- end }}
	}
{{ end -}}

{{ if .RepositoryYum }}
	yum {
		deploy_policy  = "{{ .RepositoryYum.DeployPolicy }}"
		repodata_depth = {{ .RepositoryYum.RepodataDepth }}
	}
{{ end -}}
}
`
)

var (
	resourceRepositoryTemplateFuncMap = template.FuncMap{
		"deref": func(data interface{}) string {
			switch v := data.(type) {
			case *string:
				return fmt.Sprintf("%s", *v)
			case *int:
				return fmt.Sprintf("%d", *v)
			default:
				return fmt.Sprintf("%v", v)
			}
		},
	}
	resourceRepositoryTemplate = template.Must(template.New("repository").Funcs(resourceRepositoryTemplateFuncMap).Parse(resourceRepositoryTemplateString))
)

func testAccResourceRepositoryName(repo nexus.Repository) string {
	return fmt.Sprintf("nexus_repository.%s", repo.Name)
}

func testAccResourceRepositoryConfig(repo nexus.Repository) string {
	buf := &bytes.Buffer{}
	if err := resourceRepositoryTemplate.Execute(buf, repo); err != nil {
		panic(err)
	}
	return buf.String()
}

func testAccResourceRepositoryGroup(format string) nexus.Repository {
	return nexus.Repository{
		Format: format,
		Name:   fmt.Sprintf("test-repo-%s", acctest.RandString(10)),
		Online: true,
		Type:   nexus.RepositoryTypeGroup,

		RepositoryGroup: &nexus.RepositoryGroup{},

		RepositoryStorage: &nexus.RepositoryStorage{
			BlobStoreName:               "default",
			StrictContentTypeValidation: true,
		},
	}
}

func testAccResourceRepositoryHosted(format string) nexus.Repository {
	writePolicy := "ALLOW"
	return nexus.Repository{
		Format: format,
		Name:   fmt.Sprintf("test-repo-%s", acctest.RandString(10)),
		Online: true,
		Type:   nexus.RepositoryTypeHosted,

		RepositoryCleanup: &nexus.RepositoryCleanup{
			PolicyNames: []string{"\"cleanup-weekly\""},
		},

		RepositoryStorage: &nexus.RepositoryStorage{
			BlobStoreName:               "default",
			StrictContentTypeValidation: true,
			WritePolicy:                 &writePolicy,
		},
	}
}

func testAccResourceRepositoryProxy(format string) nexus.Repository {
	return nexus.Repository{
		Format: format,
		Name:   fmt.Sprintf("test-repo-%s", acctest.RandString(10)),
		Online: true,
		Type:   nexus.RepositoryTypeProxy,

		RepositoryCleanup: &nexus.RepositoryCleanup{
			PolicyNames: []string{"\"cleanup-weekly\""},
		},

		RepositoryHTTPClient: &nexus.RepositoryHTTPClient{
			Authentication: &nexus.RepositoryHTTPClientAuthentication{
				Password: "t0ps3cr3t",
				Type:     "username",
				Username: "4dm1n",
			},
			AutoBlock: true,
			Blocked:   false,
		},

		RepositoryNegativeCache: &nexus.RepositoryNegativeCache{},

		RepositoryProxy: &nexus.RepositoryProxy{
			ContentMaxAge:  1440,
			MetadataMaxAge: 1440,
			RemoteURL:      "https://proxy.example.com",
		},

		RepositoryStorage: &nexus.RepositoryStorage{
			BlobStoreName:               "default",
			StrictContentTypeValidation: true,
		},
	}
}

func resourceRepositoryTestCheckFunc(repo nexus.Repository) resource.TestCheckFunc {
	resName := testAccResourceRepositoryName(repo)
	return resource.ComposeAggregateTestCheckFunc(
		resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr(resName, "format", repo.Format),
			resource.TestCheckResourceAttr(resName, "id", repo.Name),
			resource.TestCheckResourceAttr(resName, "name", repo.Name),
			resource.TestCheckResourceAttr(resName, "online", strconv.FormatBool(repo.Online)),
			resource.TestCheckResourceAttr(resName, "type", repo.Type),
		),
		resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr(resName, "storage.#", "1"),
			resource.TestCheckResourceAttr(resName, "storage.0.blob_store_name", repo.RepositoryStorage.BlobStoreName),
			resource.TestCheckResourceAttr(resName, "storage.0.strict_content_type_validation", strconv.FormatBool(repo.RepositoryStorage.StrictContentTypeValidation)),
		),
	)
}

func resourceRepositoryTypeGroupTestCheckFunc(repo nexus.Repository) resource.TestCheckFunc {
	resName := testAccResourceRepositoryName(repo)
	return resource.ComposeAggregateTestCheckFunc(
		resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr(resName, "group.#", "1"),
			resource.TestCheckResourceAttr(resName, "group.0.member_names.#", strconv.Itoa(len(repo.RepositoryGroup.MemberNames))),
			// FIXME: (BUG) Incorrect member_names state representation.
			// For some reasons, 1st ans 2nd elements in array are not stored as group.0.member_names.0, but instead they're stored
			// as group.0.member_names.2126137474 where 2126137474 is a "random" number.
			// This number changes from test run to test run.
			// It may be a pointer to int instead of int itself, but it's not clear and requires additional research.
			// resource.TestCheckResourceAttr(resName, "group.0.member_names.2126137474", memberRepoName),
		),
		resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr(resName, "http_client.#", "0"),
			resource.TestCheckResourceAttr(resName, "negative_cache.#", "0"),
			resource.TestCheckResourceAttr(resName, "proxy.#", "0"),
		),
	)
}

func resourceRepositoryTypeHostedTestCheckFunc(repo nexus.Repository) resource.TestCheckFunc {
	resName := testAccResourceRepositoryName(repo)
	return resource.ComposeAggregateTestCheckFunc(
		resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr(resName, "http_client.#", "0"),
			resource.TestCheckResourceAttr(resName, "group.#", "0"),
			resource.TestCheckResourceAttr(resName, "negative_cache.#", "0"),
			resource.TestCheckResourceAttr(resName, "proxy.#", "0"),
		),
		resource.TestCheckResourceAttr(resName, "storage.0.write_policy", *repo.RepositoryStorage.WritePolicy),
	)
}

func resourceRepositoryTypeProxyTestCheckFunc(repo nexus.Repository) resource.TestCheckFunc {
	resName := testAccResourceRepositoryName(repo)
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resName, "http_client.#", "1"),
		resource.TestCheckResourceAttr(resName, "group.#", "0"),
		resource.TestCheckResourceAttr(resName, "negative_cache.#", "1"),
		resource.TestCheckResourceAttr(resName, "proxy.#", "1"),
		resource.TestCheckResourceAttr(resName, "proxy.0.content_max_age", strconv.Itoa(repo.RepositoryProxy.ContentMaxAge)),
		resource.TestCheckResourceAttr(resName, "proxy.0.metadata_max_age", strconv.Itoa(repo.RepositoryProxy.MetadataMaxAge)),
	)
}
