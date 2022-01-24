package deprecated_test

import (
	"bytes"
	"fmt"
	"strconv"
	"text/template"

	"github.com/datadrivers/go-nexus-client/nexus3/schema/repository"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	resourceRepositoryTemplateString = `
resource "nexus_repository" "{{ .Name }}" {
	format = "{{ .Format }}"
	type   = "{{ .Type }}"
	name   = "{{ .Name }}"
	online = {{ .Online }}

{{ if .Apt }}
	apt {
		distribution = "{{ .Apt.Distribution }}"
		flat         = {{ .Apt.Flat }}
	}
{{ end -}}

{{ if .AptSigning }}
	apt_signing {
		keypair    = "{{ .AptSigning.Keypair }}"
		passphrase = "{{ .AptSigning.Passphrase }}"
	}
{{ end -}}

{{ if .Bower }}
	bower {
		rewrite_package_urls = {{ .Bower.RewritePackageUrls }}
	}
{{  end -}}

{{ if .Cleanup }}
	cleanup {
		policy_names = [
		{{- range $val := .Cleanup.PolicyNames }}
			{{ $val }},
		{{ end -}}
		]
	}
{{ end -}}

{{ if .Docker }}
	docker {
		force_basic_auth = {{ .Docker.ForceBasicAuth }}
		{{ if .Docker.HTTPPort -}}
		http_port        = {{ deref .Docker.HTTPPort }}
		{{ end -}}
		{{  if .Docker.HTTPSPort -}}
		https_port       = {{ deref .Docker.HTTPSPort }}
		{{ end -}}
		v1enabled        = {{ .Docker.V1Enabled }}
	}
{{ end -}}

{{ if .DockerProxy }}
	docker_proxy {
		index_type = "{{ .DockerProxy.IndexType }}"
		index_url  = "{{ deref .DockerProxy.IndexURL }}"
	}
{{ end -}}

{{ if .Group }}
	group {
		member_names = [
		{{- range $val := .Group.MemberNames }}
			{{ $val }},
		{{ end -}}
		]
	}
{{ end -}}

{{ if .HTTPClient }}
	http_client {
		{{ if .HTTPClient.Authentication -}}
		authentication {
			ntlm_domain = "{{ .HTTPClient.Authentication.NTLMDomain }}"
			ntlm_host   = "{{ .HTTPClient.Authentication.NTLMHost }}"
			{{ if .HTTPClient.Authentication.Password -}}
			password    = "{{ deref .HTTPClient.Authentication.Password }}"
			{{ end -}}
			type        = "{{ .HTTPClient.Authentication.Type }}"
			{{ if .HTTPClient.Authentication.Username -}}
			username    = "{{ deref .HTTPClient.Authentication.Username }}"
			{{ end -}}
		}
		{{ end -}}
	}
{{ end -}}

{{ if .Maven }}
	maven {
		layout_policy  = "{{ .Maven.LayoutPolicy }}"
		version_policy = "{{ .Maven.VersionPolicy }}"
	}
{{ end -}}

{{ if .NegativeCache }}
	negative_cache {

	}
{{ end -}}

{{ if .NugetProxy }}
	nuget_proxy {
		query_cache_item_max_age = {{ .NugetProxy.QueryCacheItemMaxAge }}
	}
{{ end -}}

{{ if .Proxy }}
	proxy {
		remote_url = "{{ .Proxy.RemoteURL }}"
	}
{{ end -}}

{{ if .Storage }}
	storage {
		blob_store_name                = "{{ .Storage.BlobStoreName }}"
		strict_content_type_validation = {{ .Storage.StrictContentTypeValidation }}
		{{- if eq .Type "hosted" }}
		write_policy                   = "{{ .Storage.WritePolicy }}"
		{{- end }}
	}
{{ end -}}

{{ if .Yum }}
	yum {
		deploy_policy  = "{{ .Yum.DeployPolicy }}"
		repodata_depth = {{ .Yum.RepodataDepth }}
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
				return *v
			case *int:
				return fmt.Sprintf("%d", *v)
			default:
				return fmt.Sprintf("%v", v)
			}
		},
	}
	resourceRepositoryTemplate = template.Must(template.New("repository").Funcs(resourceRepositoryTemplateFuncMap).Parse(resourceRepositoryTemplateString))
)

func testAccResourceRepositoryName(repo repository.LegacyRepository) string {
	return fmt.Sprintf("nexus_repository.%s", repo.Name)
}

func testAccResourceRepositoryConfig(repo repository.LegacyRepository) string {
	buf := &bytes.Buffer{}
	if err := resourceRepositoryTemplate.Execute(buf, repo); err != nil {
		panic(err)
	}
	return buf.String()
}

func testAccResourceRepositoryGroup(format string) repository.LegacyRepository {
	return repository.LegacyRepository{
		Format: format,
		Name:   fmt.Sprintf("test-repo-%s", acctest.RandString(10)),
		Online: true,
		Type:   repository.RepositoryTypeGroup,

		Group: &repository.Group{},

		Storage: &repository.HostedStorage{
			BlobStoreName:               "default",
			StrictContentTypeValidation: true,
		},
	}
}

func testAccResourceRepositoryHosted(format string) repository.LegacyRepository {
	writePolicy := repository.StorageWritePolicyAllow
	return repository.LegacyRepository{
		Format: format,
		Name:   fmt.Sprintf("test-repo-%s", acctest.RandString(10)),
		Online: true,
		Type:   repository.RepositoryTypeHosted,

		Cleanup: &repository.Cleanup{
			PolicyNames: []string{"\"cleanup-weekly\""},
		},

		Storage: &repository.HostedStorage{
			BlobStoreName:               "default",
			StrictContentTypeValidation: true,
			WritePolicy:                 &writePolicy,
		},
	}
}

func testAccResourceRepositoryProxy(format string) repository.LegacyRepository {
	return repository.LegacyRepository{
		Format: format,
		Name:   fmt.Sprintf("test-repo-%s", acctest.RandString(10)),
		Online: true,
		Type:   repository.RepositoryTypeProxy,

		Cleanup: &repository.Cleanup{
			PolicyNames: []string{"\"cleanup-weekly\""},
		},

		HTTPClient: &repository.HTTPClient{
			Authentication: &repository.HTTPClientAuthentication{
				Password: "t0ps3cr3t",
				Type:     "username",
				Username: "4dm1n",
			},
			AutoBlock: true,
			Blocked:   false,
		},

		NegativeCache: &repository.NegativeCache{},

		Proxy: &repository.Proxy{
			ContentMaxAge:  1440,
			MetadataMaxAge: 1440,
			RemoteURL:      "https://proxy.example.com",
		},

		Storage: &repository.HostedStorage{
			BlobStoreName:               "default",
			StrictContentTypeValidation: true,
		},
	}
}

func resourceRepositoryTestCheckFunc(repo repository.LegacyRepository) resource.TestCheckFunc {
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
			resource.TestCheckResourceAttr(resName, "storage.0.blob_store_name", repo.Storage.BlobStoreName),
			resource.TestCheckResourceAttr(resName, "storage.0.strict_content_type_validation", strconv.FormatBool(repo.Storage.StrictContentTypeValidation)),
		),
	)
}

func resourceRepositoryTypeGroupTestCheckFunc(repo repository.LegacyRepository) resource.TestCheckFunc {
	resName := testAccResourceRepositoryName(repo)
	return resource.ComposeAggregateTestCheckFunc(
		resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr(resName, "group.#", "1"),
			resource.TestCheckResourceAttr(resName, "group.0.member_names.#", strconv.Itoa(len(repo.Group.MemberNames))),
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

func resourceRepositoryTypeHostedTestCheckFunc(repo repository.LegacyRepository) resource.TestCheckFunc {
	resName := testAccResourceRepositoryName(repo)
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

func resourceRepositoryTypeProxyTestCheckFunc(repo repository.LegacyRepository) resource.TestCheckFunc {
	resName := testAccResourceRepositoryName(repo)
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resName, "http_client.#", "1"),
		resource.TestCheckResourceAttr(resName, "group.#", "0"),
		resource.TestCheckResourceAttr(resName, "negative_cache.#", "1"),
		resource.TestCheckResourceAttr(resName, "proxy.#", "1"),
		resource.TestCheckResourceAttr(resName, "proxy.0.content_max_age", strconv.Itoa(repo.Proxy.ContentMaxAge)),
		resource.TestCheckResourceAttr(resName, "proxy.0.metadata_max_age", strconv.Itoa(repo.Proxy.MetadataMaxAge)),
	)
}
