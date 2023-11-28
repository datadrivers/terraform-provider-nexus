package repository_test

import (
	"bytes"
	"fmt"
	"strconv"
	"testing"
	"text/template"

	"github.com/dre2004/go-nexus-client/nexus3/schema"
	"github.com/dre2004/go-nexus-client/nexus3/schema/repository"
	"github.com/dre2004/terraform-provider-nexus/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func testAccResourceRepositoryHelmProxy() repository.HelmProxyRepository {
	enableCircularRedirects := true
	enableCookies := true
	retries := 3
	timeout := 15
	useTrustStore := true

	return repository.HelmProxyRepository{
		Name:   fmt.Sprintf("test-repo-%s", acctest.RandString(10)),
		Online: true,
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
			RemoteURL:      "https://kubernetes-charts.storage.googleapis.com/",
		},
	}
}

func testAccResourceRepositoryHelmProxyConfig(repo repository.HelmProxyRepository) string {
	buf := &bytes.Buffer{}
	resourceRepositoryHelmProxyTemplate := template.Must(template.New("HelmProxyRepository").Funcs(acceptance.TemplateFuncMap).Parse(acceptance.TemplateStringRepositoryHelmProxy))
	if err := resourceRepositoryHelmProxyTemplate.Execute(buf, repo); err != nil {
		panic(err)
	}
	return buf.String()
}

func TestAccResourceRepositoryHelmProxy(t *testing.T) {
	routingRule := schema.RoutingRule{
		Name:        acctest.RandString(10),
		Description: "acceptance test",
		Mode:        schema.RoutingRuleModeAllow,
		Matchers: []string{
			"/",
		},
	}
	repo := testAccResourceRepositoryHelmProxy()
	repo.RoutingRule = &routingRule.Name
	resourceName := "nexus_repository_helm_proxy.acceptance"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceRoutingRuleConfig(routingRule) + testAccResourceRepositoryHelmProxyConfig(repo),
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
