package repository

import (
	"github.com/datadrivers/go-nexus-client/nexus3/schema/repository"
	"github.com/datadrivers/terraform-provider-nexus/internal/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func flattenCleanup(cleanup *repository.Cleanup) []map[string]interface{} {
	if cleanup == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"policy_names": tools.StringSliceToInterfaceSlice(cleanup.PolicyNames),
		},
	}
}

func flattenDocker(docker *repository.Docker) []map[string]interface{} {
	data := map[string]interface{}{
		"force_basic_auth": docker.ForceBasicAuth,
		"v1_enabled":       docker.V1Enabled,
	}

	if docker.HTTPPort != nil {
		data["http_port"] = *docker.HTTPPort
	}
	if docker.HTTPSPort != nil {
		data["https_port"] = *docker.HTTPSPort
	}
	if docker.Subdomain != nil {
		data["subdomain"] = *docker.Subdomain
	}

	return []map[string]interface{}{data}
}

func flattenDockerProxy(dockerProxy *repository.DockerProxy) []map[string]interface{} {
	data := map[string]interface{}{
		"index_type": string(dockerProxy.IndexType),
	}

	if dockerProxy.IndexURL != nil {
		data["index_url"] = *dockerProxy.IndexURL
	}

	return []map[string]interface{}{data}
}

func flattenComponent(component *repository.Component) []map[string]interface{} {
	if component == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"proprietary_components": component.ProprietaryComponents,
		},
	}
}

func flattenGroup(group *repository.Group) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"member_names": tools.StringSliceToInterfaceSlice(group.MemberNames),
		},
	}
}

func flattenGroupDeploy(group *repository.GroupDeploy) []map[string]interface{} {
	data := map[string]interface{}{
		"member_names": tools.StringSliceToInterfaceSlice(group.MemberNames),
	}
	if group.WritableMember != nil {
		data["writable_member"] = *group.WritableMember
	}
	return []map[string]interface{}{data}
}

func flattenHTTPClient(httpClient *repository.HTTPClient, d *schema.ResourceData) []map[string]interface{} {
	if httpClient == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"authentication": flattenHTTPClientAuthentication(httpClient.Authentication, d),
			"auto_block":     httpClient.AutoBlock,
			"blocked":        httpClient.Blocked,
			"connection":     flattenHTTPClientConnection(httpClient.Connection),
		},
	}
}

func flattenHTTPClientWithPreemptiveAuth(httpClient *repository.HTTPClientWithPreemptiveAuth, d *schema.ResourceData) []map[string]interface{} {
	if httpClient == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"authentication": flattenHTTPClientAuthenticationWithPreemptive(httpClient.Authentication, d),
			"auto_block":     httpClient.AutoBlock,
			"blocked":        httpClient.Blocked,
			"connection":     flattenHTTPClientConnection(httpClient.Connection),
		},
	}
}

func flattenHTTPClientAuthentication(auth *repository.HTTPClientAuthentication, d *schema.ResourceData) []map[string]interface{} {
	if auth == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"ntlm_domain": auth.NTLMDomain,
			"ntlm_host":   auth.NTLMHost,
			"type":        auth.Type,
			"username":    auth.Username,
			"password":    d.Get("http_client.0.authentication.0.password").(string),
		},
	}
}

func flattenHTTPClientAuthenticationWithPreemptive(auth *repository.HTTPClientAuthenticationWithPreemptive, d *schema.ResourceData) []map[string]interface{} {
	if auth == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"ntlm_domain": auth.NTLMDomain,
			"ntlm_host":   auth.NTLMHost,
			"type":        auth.Type,
			"username":    auth.Username,
			"password":    d.Get("http_client.0.authentication.0.password").(string),
			"preemptive":  auth.Preemptive,
		},
	}
}

func flattenHTTPClientConnection(conn *repository.HTTPClientConnection) []map[string]interface{} {
	if conn == nil {
		return nil
	}
	data := map[string]interface{}{
		"user_agent_suffix": conn.UserAgentSuffix,
	}
	if conn.EnableCircularRedirects != nil {
		data["enable_circular_redirects"] = *conn.EnableCircularRedirects
	}
	if conn.EnableCookies != nil {
		data["enable_cookies"] = *conn.EnableCookies
	}
	if conn.Retries != nil {
		data["retries"] = *conn.Retries
	}
	if conn.Timeout != nil {
		data["timeout"] = *conn.Timeout
	}
	if conn.UseTrustStore != nil {
		data["use_trust_store"] = *conn.UseTrustStore
	}
	return []map[string]interface{}{data}
}

func flattenNegativeCache(negativeCache *repository.NegativeCache) []map[string]interface{} {
	if negativeCache == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"enabled": negativeCache.Enabled,
			"ttl":     negativeCache.TTL,
		},
	}
}

func flattenProxy(proxy *repository.Proxy) []map[string]interface{} {
	if proxy == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"content_max_age":  proxy.ContentMaxAge,
			"metadata_max_age": proxy.MetadataMaxAge,
			"remote_url":       proxy.RemoteURL,
		},
	}
}

func flattenStorage(storage *repository.Storage) []map[string]interface{} {
	if storage == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"blob_store_name":                storage.BlobStoreName,
			"strict_content_type_validation": storage.StrictContentTypeValidation,
		},
	}
}

func flattenHostedStorage(storage *repository.HostedStorage) []map[string]interface{} {
	if storage == nil {
		return nil
	}
	data := map[string]interface{}{
		"blob_store_name":                storage.BlobStoreName,
		"strict_content_type_validation": storage.StrictContentTypeValidation,
	}
	if storage.WritePolicy != nil {
		data["write_policy"] = storage.WritePolicy
	}
	return []map[string]interface{}{data}
}

func flattenDockerHostedStorage(dockerStorage *repository.DockerHostedStorage) []map[string]interface{} {
	if dockerStorage == nil {
		return nil
	}
	data := map[string]interface{}{
		"blob_store_name":                dockerStorage.BlobStoreName,
		"strict_content_type_validation": dockerStorage.StrictContentTypeValidation,
		"write_policy":                   dockerStorage.WritePolicy,
	}
	if dockerStorage.LatestPolicy != nil {
		data["latest_policy"] = dockerStorage.LatestPolicy
	}
	return []map[string]interface{}{data}
}

func flattenMaven(maven *repository.Maven) []map[string]interface{} {
	data := map[string]interface{}{
		"version_policy": maven.VersionPolicy,
		"layout_policy":  maven.LayoutPolicy,
	}
	if maven.ContentDisposition != nil {
		data["content_disposition"] = string(*maven.ContentDisposition)
	}

	return []map[string]interface{}{data}
}

func flattenAPTHostedSigningConfig(signingData repository.AptSigning) []interface{} {

	m := make(map[string]interface{})

	if signingData.Keypair != "" {
		m["keypair"] = signingData.Keypair
	}

	if signingData.Passphrase != nil {
		m["passphrase"] = *signingData.Passphrase
	}

	return []interface{}{m}
}
