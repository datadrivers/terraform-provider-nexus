package repository

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	nexus "github.com/datadrivers/go-nexus-client/nexus3"
	"github.com/datadrivers/go-nexus-client/nexus3/pkg/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestHuggingfaceProxyLifecycle(t *testing.T) {
	var saved huggingfaceProxyRepository
	exists := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, password, ok := r.BasicAuth()
		require.True(t, ok)
		require.Equal(t, "admin", user)
		require.Equal(t, "password", password)
		switch r.Method {
		case http.MethodPost, http.MethodPut:
			require.Contains(t, r.URL.Path, "/service/rest/v1/repositories/huggingface/proxy")
			require.NoError(t, json.NewDecoder(r.Body).Decode(&saved))
			require.Equal(t, "upstream-secret", saved.HTTPClient.Authentication.BearerToken)
			exists = true
			w.WriteHeader(http.StatusNoContent)
		case http.MethodGet:
			require.Equal(t, "/service/rest/v1/repositories/huggingface/proxy/huggingface-proxy", r.URL.Path)
			if !exists {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			// Nexus does not return the upstream credential.
			reply := saved
			auth := *saved.HTTPClient.Authentication
			auth.BearerToken = ""
			reply.HTTPClient.Authentication = &auth
			require.NoError(t, json.NewEncoder(w).Encode(reply))
		case http.MethodDelete:
			require.Equal(t, "/service/rest/v1/repositories/huggingface-proxy", r.URL.Path)
			exists = false
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()
	config := client.Config{URL: server.URL, Username: "admin", Password: "password"}
	nc := nexus.NewClient(config)
	ConfigureHuggingface(nc, client.NewClient(config))
	defer huggingfaceClients.Delete(nc)
	resource := ResourceRepositoryHuggingfaceProxy()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		"name": "huggingface-proxy", "online": true,
		"storage":        []interface{}{map[string]interface{}{"blob_store_name": "models", "strict_content_type_validation": true}},
		"proxy":          []interface{}{map[string]interface{}{"remote_url": "https://huggingface.co", "content_max_age": -1, "metadata_max_age": 1440}},
		"negative_cache": []interface{}{map[string]interface{}{"enabled": true, "ttl": 15}},
		"http_client": []interface{}{map[string]interface{}{"blocked": false, "auto_block": false,
			"authentication": []interface{}{map[string]interface{}{"type": "bearerToken", "bearer_token": "upstream-secret"}},
			"connection":     []interface{}{map[string]interface{}{"timeout": 120, "retries": 3}},
		}},
	})
	require.NoError(t, resource.Create(d, nc))
	require.Equal(t, "huggingface-proxy", d.Id())
	require.Equal(t, "upstream-secret", d.Get("http_client.0.authentication.0.bearer_token"))
	require.Equal(t, -1, saved.Proxy.ContentMaxAge)
	require.NoError(t, resource.Update(d, nc))
	require.NoError(t, resource.Delete(d, nc))
	require.NoError(t, resource.Read(d, nc))
	require.Empty(t, d.Id())
}

func TestHuggingfaceProxyDoesNotLeakCredentials(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("upstream-secret"))
	}))
	defer server.Close()
	c := &huggingfaceHTTPClient{client: client.NewClient(client.Config{URL: server.URL})}
	err := c.Create(huggingfaceProxyRepository{Name: "models"})
	require.Error(t, err)
	require.NotContains(t, err.Error(), "upstream-secret")
}
