package repository_test

import (
	"bytes"
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"text/template"

	"github.com/datadrivers/go-nexus-client/nexus3/schema"
	"github.com/datadrivers/go-nexus-client/nexus3/schema/repository"
	"github.com/datadrivers/terraform-provider-nexus/internal/acceptance"
	"github.com/datadrivers/terraform-provider-nexus/internal/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	resourceRepositoryDockerProxyTemplateString = `
resource "nexus_repository_docker_proxy" "acceptance" {
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

	docker_proxy {
		index_type = "{{ .DockerProxy.IndexType }}"
{{- if .DockerProxy.IndexURL }}
		index_url = "{{ .DockerProxy.IndexURL }}"
{{- end }}
	}

{{ if .RoutingRule }}
	routing_rule = nexus_routing_rule.acceptance.name
{{ end -}}

{{ if .Storage }}
	storage {
		blob_store_name                = "{{ .Storage.BlobStoreName }}"
		strict_content_type_validation = {{ .Storage.StrictContentTypeValidation }}
	}
{{ end -}}

	proxy {
		remote_url = "{{ .Proxy.RemoteURL }}"
		{{ if .Proxy.ContentMaxAge }}
			content_max_age = "{{ .Proxy.ContentMaxAge }}"
		{{ end -}}
		{{ if .Proxy.MetadataMaxAge }}
			metadata_max_age = "{{ .Proxy.MetadataMaxAge }}"
		{{ end -}}
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

{{ if .HTTPClient }}
	http_client {
		auto_block = {{ .HTTPClient.AutoBlock }}
		blocked    = {{ .HTTPClient.Blocked }}

		{{ if .HTTPClient.Authentication -}}
		authentication {
			ntlm_domain = "{{ .HTTPClient.Authentication.NTLMDomain }}"
			ntlm_host   = "{{ .HTTPClient.Authentication.NTLMHost }}"
			{{ if .HTTPClient.Authentication.Password -}}
			password    = "{{ .HTTPClient.Authentication.Password }}"
			{{ end -}}
			type        = "{{ .HTTPClient.Authentication.Type }}"
			{{ if .HTTPClient.Authentication.Username -}}
			username    = "{{ .HTTPClient.Authentication.Username }}"
			{{ end -}}
		}
		{{ end -}}

		{{ if .HTTPClient.Connection -}}
		connection {
			{{ if .HTTPClient.Connection.EnableCircularRedirects -}}
			enable_circular_redirects = {{ .HTTPClient.Connection.EnableCircularRedirects }}
			{{ end -}}
			{{ if .HTTPClient.Connection.EnableCookies -}}
			enable_cookies = {{ .HTTPClient.Connection.EnableCookies }}
			{{ end -}}
			{{ if .HTTPClient.Connection.Retries -}}
			retries = {{ .HTTPClient.Connection.Retries }}
			{{ end -}}
			{{ if .HTTPClient.Connection.Timeout -}}
			timeout = {{ .HTTPClient.Connection.Timeout }}
			{{ end -}}
			{{ if .HTTPClient.Connection.UserAgentSuffix -}}
			user_agent_suffix = "{{ .HTTPClient.Connection.UserAgentSuffix }}"
			{{ end -}}
			{{ if .HTTPClient.Connection.UseTrustStore -}}
			use_trust_store = {{ .HTTPClient.Connection.UseTrustStore }}
			{{ end -}}
		}
		{{ end -}}
	}
{{ end -}}

{{ if .NegativeCache }}
	negative_cache {
		{{ if .NegativeCache.Enabled }}
		enabled = {{ .NegativeCache.Enabled }}
		{{ end -}}
		{{ if .NegativeCache.TTL }}
		ttl = {{ .NegativeCache.TTL }}
		{{ end }}
	}
{{ end -}}
}
`
)

func testAccResourceRepositoryDockerProxy() repository.DockerProxyRepository {
	enableCircularRedirects := true
	enableCookies := true
	retries := 3
	timeout := 15
	useTrustStore := true

	return repository.DockerProxyRepository{
		Name:   fmt.Sprintf("test-repo-%s", acctest.RandString(10)),
		Online: true,
		DockerProxy: repository.DockerProxy{
			IndexType: repository.DockerProxyIndexTypeRegistry,
			IndexURL:  tools.GetStringPointer("https://docker.elastic.co/index.json"),
		},
		Docker: repository.Docker{
			ForceBasicAuth: false,
			HTTPPort:       tools.GetIntPointer(rand.Intn(999) + 34000),
			HTTPSPort:      tools.GetIntPointer(rand.Intn(999) + 35000),
			V1Enabled:      true,
		},
		Storage: repository.Storage{
			BlobStoreName:               "default",
			StrictContentTypeValidation: true,
		},
		Cleanup: &repository.Cleanup{
			PolicyNames: []string{"cleanup-weekly"},
		},
		HTTPClient: repository.HTTPClient{
			AutoBlock: true,
			Blocked:   false,
			Authentication: &repository.HTTPClientAuthentication{
				Password: "acceptance-password",
				Type:     repository.HTTPClientAuthenticationTypeUsername,
				Username: "acceptance-user",
			},
			Connection: &repository.HTTPClientConnection{
				EnableCircularRedirects: &enableCircularRedirects,
				EnableCookies:           &enableCookies,
				Retries:                 &retries,
				Timeout:                 &timeout,
				UserAgentSuffix:         "acceptance-test",
				UseTrustStore:           &useTrustStore,
			},
		},
		NegativeCache: repository.NegativeCache{
			Enabled: true,
			TTL:     5,
		},
		Proxy: repository.Proxy{
			ContentMaxAge:  770,
			MetadataMaxAge: 770,
			RemoteURL:      "https://docker.elastic.co",
		},
	}
}

func testAccResourceRepositoryDockerProxyConfig(repo repository.DockerProxyRepository) string {
	buf := &bytes.Buffer{}
	resourceRepositoryDockerProxyTemplate := template.Must(template.New("DockerProxyRepository").Funcs(acceptance.TemplateFuncMap).Parse(resourceRepositoryDockerProxyTemplateString))
	if err := resourceRepositoryDockerProxyTemplate.Execute(buf, repo); err != nil {
		panic(err)
	}
	return buf.String()
}

func TestAccResourceRepositoryDockerProxy(t *testing.T) {
	routingRule := schema.RoutingRule{
		Name:        acctest.RandString(10),
		Description: "acceptance test",
		Mode:        schema.RoutingRuleModeAllow,
		Matchers: []string{
			"/",
		},
	}
	repo := testAccResourceRepositoryDockerProxy()
	repo.RoutingRule = &routingRule.Name
	resourceName := "nexus_repository_docker_proxy.acceptance"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceRoutingRuleConfig(routingRule) + testAccResourceRepositoryDockerProxyConfig(repo),
				Check: resource.ComposeTestCheckFunc(
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "id", repo.Name),
						resource.TestCheckResourceAttr(resourceName, "name", repo.Name),
						resource.TestCheckResourceAttr(resourceName, "online", strconv.FormatBool(repo.Online)),
					),
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "http_client.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "http_client.0.auto_block", strconv.FormatBool(repo.HTTPClient.AutoBlock)),
						resource.TestCheckResourceAttr(resourceName, "http_client.0.blocked", strconv.FormatBool(repo.HTTPClient.Blocked)),
						resource.TestCheckResourceAttr(resourceName, "http_client.0.authentication.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "http_client.0.authentication.0.type", string(repo.HTTPClient.Authentication.Type)),
						resource.TestCheckResourceAttr(resourceName, "http_client.0.authentication.0.username", repo.HTTPClient.Authentication.Username),
						resource.TestCheckResourceAttr(resourceName, "http_client.0.authentication.0.password", repo.HTTPClient.Authentication.Password),
						resource.TestCheckResourceAttr(resourceName, "http_client.0.connection.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "http_client.0.connection.0.enable_circular_redirects", strconv.FormatBool(*repo.HTTPClient.Connection.EnableCircularRedirects)),
						resource.TestCheckResourceAttr(resourceName, "http_client.0.connection.0.enable_cookies", strconv.FormatBool(*repo.HTTPClient.Connection.EnableCookies)),
						resource.TestCheckResourceAttr(resourceName, "http_client.0.connection.0.retries", strconv.Itoa(*repo.HTTPClient.Connection.Retries)),
						resource.TestCheckResourceAttr(resourceName, "http_client.0.connection.0.timeout", strconv.Itoa(*repo.HTTPClient.Connection.Timeout)),
						resource.TestCheckResourceAttr(resourceName, "http_client.0.connection.0.user_agent_suffix", repo.HTTPClient.Connection.UserAgentSuffix),
						resource.TestCheckResourceAttr(resourceName, "http_client.0.connection.0.use_trust_store", strconv.FormatBool(*repo.HTTPClient.Connection.UseTrustStore)),
						resource.TestCheckResourceAttr(resourceName, "negative_cache.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "negative_cache.0.enabled", strconv.FormatBool(repo.NegativeCache.Enabled)),
						resource.TestCheckResourceAttr(resourceName, "negative_cache.0.ttl", strconv.Itoa(repo.NegativeCache.TTL)),
						resource.TestCheckResourceAttr(resourceName, "proxy.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "proxy.0.content_max_age", strconv.Itoa(repo.Proxy.ContentMaxAge)),
						resource.TestCheckResourceAttr(resourceName, "proxy.0.metadata_max_age", strconv.Itoa(repo.Proxy.MetadataMaxAge)),
						resource.TestCheckResourceAttr(resourceName, "proxy.0.remote_url", repo.Proxy.RemoteURL),
						resource.TestCheckResourceAttr(resourceName, "storage.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "storage.0.blob_store_name", repo.Storage.BlobStoreName),
						resource.TestCheckResourceAttr(resourceName, "storage.0.strict_content_type_validation", strconv.FormatBool(repo.Storage.StrictContentTypeValidation)),
						resource.TestCheckResourceAttr(resourceName, "cleanup.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "cleanup.0.policy_names.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "cleanup.0.policy_names.0", repo.Cleanup.PolicyNames[0]),
						resource.TestCheckResourceAttr(resourceName, "routing_rule", *repo.RoutingRule),
					),
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "docker.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "docker.0.force_basic_auth", strconv.FormatBool(repo.Docker.ForceBasicAuth)),
						resource.TestCheckResourceAttr(resourceName, "docker.0.http_port", strconv.Itoa(*repo.Docker.HTTPPort)),
						resource.TestCheckResourceAttr(resourceName, "docker.0.https_port", strconv.Itoa(*repo.Docker.HTTPSPort)),
						resource.TestCheckResourceAttr(resourceName, "docker.0.v1_enabled", strconv.FormatBool(repo.Docker.V1Enabled)),
						resource.TestCheckResourceAttr(resourceName, "docker_proxy.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "docker_proxy.0.index_type", string(repo.DockerProxy.IndexType)),
						resource.TestCheckResourceAttr(resourceName, "docker_proxy.0.index_url", *repo.DockerProxy.IndexURL),
					),
				),
			},
			{
				ResourceName:            resourceName,
				ImportStateId:           repo.Name,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"http_client.0.authentication.0.password"},
			},
		},
	})
}
