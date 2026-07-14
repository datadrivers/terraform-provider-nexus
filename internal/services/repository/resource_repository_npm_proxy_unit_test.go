package repository

import (
	"testing"

	"github.com/datadrivers/go-nexus-client/nexus3/schema/repository"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TestSetNpmProxyRepositoryToResourceData_NilNpm reproduces
// https://github.com/datadrivers/terraform-provider-nexus/issues/596.
//
// Nexus Repository 3.94.0 stopped returning the `npm` sub-object in its GET
// response for npm proxy repositories, so the anonymously embedded *Npm pointer
// on NpmProxyRepository is nil. Before the nil-guard, reading the promoted
// RemoveQuarantined / RemoveNonCataloged fields dereferenced that nil pointer
// and panicked while refreshing state (terraform plan). This test asserts the
// read succeeds and both attributes default to false.
func TestSetNpmProxyRepositoryToResourceData_NilNpm(t *testing.T) {
	repo := &repository.NpmProxyRepository{
		Name:   "npm-official-registry",
		Online: true,
		Npm:    nil, // Nexus >= 3.94.0 no longer returns this sub-object
	}

	resourceData := schema.TestResourceDataRaw(t, ResourceRepositoryNpmProxy().Schema, map[string]interface{}{})

	if err := setNpmProxyRepositoryToResourceData(repo, resourceData); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if got := resourceData.Get("remove_quarantined").(bool); got {
		t.Errorf("remove_quarantined = %v, want false", got)
	}
	if got := resourceData.Get("remove_non_cataloged").(bool); got {
		t.Errorf("remove_non_cataloged = %v, want false", got)
	}
}

// TestSetNpmProxyRepositoryToResourceData_WithNpm ensures the values from the
// npm sub-object are still propagated when Nexus does return it (< 3.94.0).
func TestSetNpmProxyRepositoryToResourceData_WithNpm(t *testing.T) {
	repo := &repository.NpmProxyRepository{
		Name:   "npm-official-registry",
		Online: true,
		Npm: &repository.Npm{
			RemoveQuarantined:  true,
			RemoveNonCataloged: true,
		},
	}

	resourceData := schema.TestResourceDataRaw(t, ResourceRepositoryNpmProxy().Schema, map[string]interface{}{})

	if err := setNpmProxyRepositoryToResourceData(repo, resourceData); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if got := resourceData.Get("remove_quarantined").(bool); !got {
		t.Errorf("remove_quarantined = %v, want true", got)
	}
	if got := resourceData.Get("remove_non_cataloged").(bool); !got {
		t.Errorf("remove_non_cataloged = %v, want true", got)
	}
}
