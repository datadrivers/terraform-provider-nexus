package other_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	nexus "github.com/datadrivers/go-nexus-client/nexus3"
	"github.com/datadrivers/go-nexus-client/nexus3/pkg/client"
	"github.com/datadrivers/terraform-provider-nexus/internal/services/other"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newMailConfigServer returns a client for a fake Nexus that records the body of
// the last PUT to the mail config endpoint.
func newMailConfigServer(t *testing.T, sent *map[string]interface{}) *nexus.NexusClient {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			body, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			*sent = map[string]interface{}{}
			require.NoError(t, json.Unmarshal(body, sent))
			w.WriteHeader(http.StatusNoContent)
		default:
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"host":"smtp.example.com","port":587,"fromAddress":"noreply@example.com","username":"smtp-user"}`))
		}
	}))
	t.Cleanup(server.Close)
	return nexus.NewClient(client.Config{URL: server.URL, Username: "admin", Password: "admin"})
}

func TestResourceMailConfigSendsPassword(t *testing.T) {
	var sent map[string]interface{}
	c := newMailConfigServer(t, &sent)
	resource := other.ResourceMailConfig()

	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		"host":         "smtp.example.com",
		"port":         587,
		"from_address": "noreply@example.com",
		"username":     "smtp-user",
		"password":     "smtp-secret",
	})

	require.NoError(t, resource.Create(d, c))
	assert.Equal(t, "smtp-secret", sent["password"])

	require.NoError(t, resource.Update(d, c))
	assert.Equal(t, "smtp-secret", sent["password"])
}

func TestResourceMailConfigOmitsEmptyPassword(t *testing.T) {
	var sent map[string]interface{}
	c := newMailConfigServer(t, &sent)
	resource := other.ResourceMailConfig()

	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		"host":         "smtp.example.com",
		"port":         587,
		"from_address": "noreply@example.com",
	})

	require.NoError(t, resource.Create(d, c))
	assert.NotContains(t, sent, "password")
}
