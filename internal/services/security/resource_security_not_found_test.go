package security_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	nexus "github.com/datadrivers/go-nexus-client/nexus3"
	"github.com/datadrivers/go-nexus-client/nexus3/pkg/client"
	"github.com/datadrivers/terraform-provider-nexus/internal/services/security"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newNotFoundClient returns a client for a Nexus that answers every request with
// 404 and the given body, like Nexus does for a role or LDAP server that does not exist.
func newNotFoundClient(t *testing.T, body string) *nexus.NexusClient {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(server.Close)
	return nexus.NewClient(client.Config{URL: server.URL, Username: "admin", Password: "admin"})
}

func TestResourceSecurityRoleDeletedOutsideTerraform(t *testing.T) {
	c := newNotFoundClient(t, `{"id":"*","message":"\"Role 'missing-role' not found.\""}`)
	resource := security.ResourceSecurityRole()

	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{"roleid": "missing-role", "name": "missing-role"})
	d.SetId("missing-role")

	exists, err := resource.Exists(d, c)
	require.NoError(t, err)
	assert.False(t, exists)

	require.NoError(t, resource.Read(d, c))
	assert.Equal(t, "", d.Id(), "a role deleted outside Terraform must be removed from state")
}

func TestResourceSecurityRoleOtherErrorsAreReturned(t *testing.T) {
	c := newNotFoundClient(t, `{"id":"*","message":"\"Role 'other-role' not found.\""}`)
	resource := security.ResourceSecurityRole()

	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{"roleid": "missing-role", "name": "missing-role"})
	d.SetId("missing-role")

	assert.Error(t, resource.Read(d, c))
	assert.Equal(t, "missing-role", d.Id())
}

func TestDataSourceSecurityRoleNotFound(t *testing.T) {
	c := newNotFoundClient(t, `{"id":"*","message":"\"Role 'missing-role' not found.\""}`)
	dataSource := security.DataSourceSecurityRole()

	d := schema.TestResourceDataRaw(t, dataSource.Schema, map[string]interface{}{"roleid": "missing-role"})

	err := dataSource.Read(d, c)
	require.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), `role "missing-role" not found`), err.Error())
}

func TestResourceSecurityLDAPDeletedOutsideTerraform(t *testing.T) {
	c := newNotFoundClient(t, `{"id":"*","message":"\"LDAP server named 'test-ldap' was not found\""}`)
	resource := security.ResourceSecurityLDAP()

	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{"name": "test-ldap"})
	d.SetId("test-ldap")

	require.NoError(t, resource.Read(d, c))
	assert.Equal(t, "", d.Id(), "a missing LDAP server must be removed from state")
}
