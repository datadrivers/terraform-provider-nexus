package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	nexus "github.com/datadrivers/go-nexus-client/nexus3"
	"github.com/datadrivers/go-nexus-client/nexus3/pkg/client"
	"github.com/datadrivers/go-nexus-client/nexus3/schema/repository"
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	repositorySchema "github.com/datadrivers/terraform-provider-nexus/internal/schema/repository"
	"github.com/datadrivers/terraform-provider-nexus/internal/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"net/http"
	"net/url"
	"sync"
)

// Hugging Face uses the common proxy fields without format-specific attributes.
type huggingfaceProxyRepository struct {
	Name            string                   `json:"name"`
	Online          bool                     `json:"online"`
	Storage         repository.Storage       `json:"storage"`
	Proxy           repository.Proxy         `json:"proxy"`
	NegativeCache   repository.NegativeCache `json:"negativeCache"`
	HTTPClient      repository.HTTPClient    `json:"httpClient"`
	RoutingRule     *string                  `json:"routingRule,omitempty"`
	RoutingRuleName *string                  `json:"routingRuleName,omitempty"`
	Cleanup         *repository.Cleanup      `json:"cleanup,omitempty"`
}

func ResourceRepositoryHuggingfaceProxy() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to create an huggingface proxy repository.",

		Create: resourceHuggingfaceProxyRepositoryCreate,
		Delete: resourceHuggingfaceProxyRepositoryDelete,
		Exists: resourceHuggingfaceProxyRepositoryExists,
		Read:   resourceHuggingfaceProxyRepositoryRead,
		Update: resourceHuggingfaceProxyRepositoryUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			// Common schemas
			"id":     common.ResourceID,
			"name":   repositorySchema.ResourceName,
			"online": repositorySchema.ResourceOnline,
			// Proxy schemas
			"cleanup":        repositorySchema.ResourceCleanup,
			"http_client":    huggingfaceHTTPClientSchema(),
			"negative_cache": repositorySchema.ResourceNegativeCache,
			"proxy":          repositorySchema.ResourceProxy,
			"routing_rule":   repositorySchema.ResourceRoutingRule,
			"storage":        repositorySchema.ResourceStorage,
		},
	}
}

// Hugging Face supports bearer-token authentication. Keep the extra field
// scoped to this resource so other proxy resources cannot accept a value
// their encoders do not send to Nexus.
func huggingfaceHTTPClientSchema() *schema.Schema {
	httpClient := *repositorySchema.ResourceHTTPClient
	httpClientResource := *repositorySchema.ResourceHTTPClient.Elem.(*schema.Resource)
	httpClientResource.Schema = make(map[string]*schema.Schema, len(httpClientResource.Schema))
	for name, field := range repositorySchema.ResourceHTTPClient.Elem.(*schema.Resource).Schema {
		httpClientResource.Schema[name] = field
	}

	authentication := *repositorySchema.ResourceHTTPClientAuthentication
	authenticationResource := *repositorySchema.ResourceHTTPClientAuthentication.Elem.(*schema.Resource)
	authenticationResource.Schema = make(map[string]*schema.Schema, len(authenticationResource.Schema)+1)
	for name, field := range repositorySchema.ResourceHTTPClientAuthentication.Elem.(*schema.Resource).Schema {
		authenticationResource.Schema[name] = field
	}
	authenticationType := *authenticationResource.Schema["type"]
	authenticationType.Description = "Authentication type. Possible values: `ntlm`, `username` or `bearerToken`"
	authenticationType.ValidateDiagFunc = validation.ToDiagFunc(validation.StringInSlice([]string{"ntlm", "username", "bearerToken"}, false))
	authenticationResource.Schema["type"] = &authenticationType
	authenticationResource.Schema["bearer_token"] = &schema.Schema{
		Description: "The bearer token used for bearerToken authentication",
		Optional:    true,
		Sensitive:   true,
		Type:        schema.TypeString,
	}
	authentication.Elem = &authenticationResource
	httpClientResource.Schema["authentication"] = &authentication
	httpClient.Elem = &httpClientResource
	return &httpClient
}

func getHuggingfaceProxyRepositoryFromResourceData(resourceData *schema.ResourceData) huggingfaceProxyRepository {
	httpClientConfig := resourceData.Get("http_client").([]interface{})[0].(map[string]interface{})
	negativeCacheConfig := resourceData.Get("negative_cache").([]interface{})[0].(map[string]interface{})
	proxyConfig := resourceData.Get("proxy").([]interface{})[0].(map[string]interface{})
	storageConfig := resourceData.Get("storage").([]interface{})[0].(map[string]interface{})

	repo := huggingfaceProxyRepository{
		Name:   resourceData.Get("name").(string),
		Online: resourceData.Get("online").(bool),
		Storage: repository.Storage{
			BlobStoreName:               storageConfig["blob_store_name"].(string),
			StrictContentTypeValidation: storageConfig["strict_content_type_validation"].(bool),
		},
		HTTPClient: repository.HTTPClient{
			AutoBlock: httpClientConfig["auto_block"].(bool),
			Blocked:   httpClientConfig["blocked"].(bool),
		},
		NegativeCache: repository.NegativeCache{
			Enabled: negativeCacheConfig["enabled"].(bool),
			TTL:     negativeCacheConfig["ttl"].(int),
		},
		Proxy: repository.Proxy{
			ContentMaxAge:  proxyConfig["content_max_age"].(int),
			MetadataMaxAge: proxyConfig["metadata_max_age"].(int),
			RemoteURL:      proxyConfig["remote_url"].(string),
		},
	}

	if routingRule, ok := resourceData.GetOk("routing_rule"); ok {
		repo.RoutingRule = tools.GetStringPointer(routingRule.(string))
		repo.RoutingRuleName = tools.GetStringPointer(routingRule.(string))
	}

	cleanupList := resourceData.Get("cleanup").([]interface{})
	if len(cleanupList) > 0 && cleanupList[0] != nil {
		cleanupConfig := cleanupList[0].(map[string]interface{})
		if len(cleanupConfig) > 0 {
			policy_names, ok := cleanupConfig["policy_names"]
			if ok {
				repo.Cleanup = &repository.Cleanup{
					PolicyNames: tools.InterfaceSliceToStringSlice(policy_names.(*schema.Set).List()),
				}
			}
		}
	}

	if v, ok := httpClientConfig["authentication"]; ok {
		authList := v.([]interface{})
		if len(authList) == 1 && authList[0] != nil {
			authConfig := authList[0].(map[string]interface{})

			repo.HTTPClient.Authentication = &repository.HTTPClientAuthentication{
				NTLMDomain:  authConfig["ntlm_domain"].(string),
				NTLMHost:    authConfig["ntlm_host"].(string),
				Type:        repository.HTTPClientAuthenticationType(authConfig["type"].(string)),
				Username:    authConfig["username"].(string),
				Password:    authConfig["password"].(string),
				BearerToken: authConfig["bearer_token"].(string),
			}
		}
	}

	if v, ok := httpClientConfig["connection"]; ok {
		repo.HTTPClient.Connection = getHTTPClientConnection(v.([]interface{}))
	}

	return repo
}

func setHuggingfaceProxyRepositoryToResourceData(repo *huggingfaceProxyRepository, resourceData *schema.ResourceData) error {
	resourceData.SetId(repo.Name)
	resourceData.Set("name", repo.Name)
	resourceData.Set("online", repo.Online)

	if repo.RoutingRuleName != nil {
		resourceData.Set("routing_rule", repo.RoutingRuleName)
	} else if repo.RoutingRule != nil {
		resourceData.Set("routing_rule", repo.RoutingRule)
	} else if repo.RoutingRuleName == nil && repo.RoutingRule == nil {
		resourceData.Set("routing_rule", nil)
	}

	if err := resourceData.Set("storage", flattenStorage(&repo.Storage)); err != nil {
		return err
	}

	if err := resourceData.Set("http_client", flattenHuggingfaceHTTPClient(&repo.HTTPClient, resourceData)); err != nil {
		return err
	}

	if err := resourceData.Set("negative_cache", flattenNegativeCache(&repo.NegativeCache)); err != nil {
		return err
	}

	if err := resourceData.Set("proxy", flattenProxy(&repo.Proxy)); err != nil {
		return err
	}

	if err := resourceData.Set("cleanup", flattenCleanup(repo.Cleanup)); err != nil {
		return err
	}
	return nil
}

// Nexus omits bearer-token values from repository reads. Preserve the
// configured value in Terraform state so refreshes and updates keep it.
func flattenHuggingfaceHTTPClient(httpClient *repository.HTTPClient, d *schema.ResourceData) []map[string]interface{} {
	flattened := flattenHTTPClient(httpClient, d)
	if httpClient == nil || httpClient.Authentication == nil || httpClient.Authentication.Type != repository.HTTPClientAuthenticationTypeBearerToken || len(flattened) == 0 {
		return flattened
	}

	authentication, _ := flattened[0]["authentication"].([]map[string]interface{})
	if len(authentication) > 0 {
		authentication[0]["bearer_token"] = d.Get("http_client.0.authentication.0.bearer_token").(string)
		flattened[0]["authentication"] = authentication
	}
	return flattened
}

func resourceHuggingfaceProxyRepositoryCreate(resourceData *schema.ResourceData, m interface{}) error {
	client := huggingfaceClient(m)

	repo := getHuggingfaceProxyRepositoryFromResourceData(resourceData)

	if err := client.Create(repo); err != nil {
		return err
	}
	resourceData.SetId(repo.Name)

	return retryRepositoryRead(resourceHuggingfaceProxyRepositoryRead, resourceData, m)
}

func resourceHuggingfaceProxyRepositoryRead(resourceData *schema.ResourceData, m interface{}) error {
	client := huggingfaceClient(m)

	repo, err := client.Get(resourceData.Id())
	if err != nil {
		return err
	}

	if repo == nil {
		resourceData.SetId("")
		return nil
	}

	return setHuggingfaceProxyRepositoryToResourceData(repo, resourceData)
}

func resourceHuggingfaceProxyRepositoryUpdate(resourceData *schema.ResourceData, m interface{}) error {
	client := huggingfaceClient(m)

	repoName := resourceData.Id()
	repo := getHuggingfaceProxyRepositoryFromResourceData(resourceData)

	if err := client.Update(repoName, repo); err != nil {
		return err
	}

	return retryRepositoryRead(resourceHuggingfaceProxyRepositoryRead, resourceData, m)
}

func resourceHuggingfaceProxyRepositoryDelete(resourceData *schema.ResourceData, m interface{}) error {
	client := huggingfaceClient(m)
	return client.Delete(resourceData.Id())
}

func resourceHuggingfaceProxyRepositoryExists(resourceData *schema.ResourceData, m interface{}) (bool, error) {
	client := huggingfaceClient(m)

	repo, err := client.Get(resourceData.Id())
	return repo != nil, err
}

// Keep clients scoped to provider instances, including aliased providers.
var huggingfaceClients sync.Map

type huggingfaceHTTPClient struct{ client *client.Client }

func ConfigureHuggingface(nc *nexus.NexusClient, c *client.Client) {
	huggingfaceClients.Store(nc, &huggingfaceHTTPClient{client: c})
}

func huggingfaceClient(m interface{}) *huggingfaceHTTPClient {
	c, _ := huggingfaceClients.Load(m.(*nexus.NexusClient))
	return c.(*huggingfaceHTTPClient)
}

const huggingfaceEndpoint = client.BasePath + "v1/repositories/huggingface/proxy"

func (c *huggingfaceHTTPClient) write(name string, repo huggingfaceProxyRepository, update bool) error {
	payload, err := json.Marshal(repo)
	if err != nil {
		return err
	}
	var response *http.Response
	if update {
		_, response, err = c.client.Put(huggingfaceEndpoint+"/"+url.PathEscape(name), bytes.NewReader(payload))
	} else {
		_, response, err = c.client.Post(huggingfaceEndpoint, bytes.NewReader(payload))
	}
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated && response.StatusCode != http.StatusNoContent {
		// Do not include payloads or response bodies: they may contain credentials.
		return fmt.Errorf("write Hugging Face repository %q: HTTP %d", name, response.StatusCode)
	}
	return nil
}
func (c *huggingfaceHTTPClient) Create(repo huggingfaceProxyRepository) error {
	return c.write(repo.Name, repo, false)
}
func (c *huggingfaceHTTPClient) Update(name string, repo huggingfaceProxyRepository) error {
	return c.write(name, repo, true)
}
func (c *huggingfaceHTTPClient) Get(name string) (*huggingfaceProxyRepository, error) {
	body, response, err := c.client.Get(huggingfaceEndpoint+"/"+url.PathEscape(name), nil)
	if err != nil {
		return nil, err
	}
	if response.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("read Hugging Face repository %q: HTTP %d", name, response.StatusCode)
	}
	var repo huggingfaceProxyRepository
	if err := json.Unmarshal(body, &repo); err != nil {
		return nil, err
	}
	return &repo, nil
}
func (c *huggingfaceHTTPClient) Delete(name string) error {
	_, response, err := c.client.Delete(client.BasePath + "v1/repositories/" + url.PathEscape(name))
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusNotFound {
		return fmt.Errorf("delete Hugging Face repository %q: HTTP %d", name, response.StatusCode)
	}
	return nil
}
