package security

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	nexus "github.com/datadrivers/go-nexus-client/nexus3"
	"github.com/datadrivers/go-nexus-client/nexus3/pkg/client"
	"github.com/datadrivers/go-nexus-client/nexus3/pkg/tools"
)

// REST endpoint introduced in Sonatype Nexus Repository Pro 3.93.0
// (see https://help.sonatype.com/en/openid-connect.html).
//
// JSON field tags below were verified against a live Nexus Pro 3.93.0
// instance via `GET /service/rest/v1/security/oauth2`.
const oidcAPIEndpoint = client.BasePath + "v1/security/oauth2"

// OIDC is the request/response body of the Nexus OAuth2/OIDC config endpoint.
type OIDC struct {
	ClientID                  string            `json:"clientId"`
	ClientSecret              string            `json:"clientSecret"`
	IdpAuthorizationURL       string            `json:"idpAuthorizationUrl"`
	IdpTokenURL               string            `json:"idpTokenUrl"`
	IdpJwksURL                string            `json:"idpJwksUrl"`
	IdpJwsAlgorithm           string            `json:"idpJwsAlgorithm"`
	UsernameClaim             string            `json:"usernameClaim"`
	GroupsClaim               string            `json:"groupsClaim"`
	IdpLogoutURL              string            `json:"idpLogoutUrl"`
	IdpJwks                   string            `json:"idpJwks"`
	FirstNameClaim            string            `json:"firstNameClaim"`
	LastNameClaim             string            `json:"lastNameClaim"`
	EmailClaim                string            `json:"emailClaim"`
	UseTrustStore             bool              `json:"useTrustStore"`
	ExactMatchClaims          map[string]string `json:"exactMatchClaims"`
	AuthorizationCustomParams map[string]string `json:"authorizationCustomParams"`
	TokenRequestCustomParams  map[string]string `json:"tokenRequestCustomParams"`
}

type oidcService struct {
	c *client.Client
}

// oidcServices maps a provider instance's *nexus.NexusClient to its low-level
// OIDC client, populated by ConfigureOIDC during providerConfigure. Keyed
// per-instance (rather than a single package-level value) so that multiple
// aliased `nexus` provider configurations in the same run each talk to their
// own Nexus instance instead of racing on a shared global. Package-scoped
// because go-nexus-client v1.20.0 does not expose an OIDC service via
// NexusClient and the embedded `*client.Client` is unexported.
var (
	oidcServicesMu sync.Mutex
	oidcServices   = map[*nexus.NexusClient]*oidcService{}
)

func ConfigureOIDC(nc *nexus.NexusClient, c *client.Client) {
	oidcServicesMu.Lock()
	defer oidcServicesMu.Unlock()
	oidcServices[nc] = &oidcService{c: c}
}

func oidc(nc *nexus.NexusClient) (*oidcService, error) {
	oidcServicesMu.Lock()
	defer oidcServicesMu.Unlock()
	svc, ok := oidcServices[nc]
	if !ok {
		return nil, fmt.Errorf("nexus OIDC client not configured")
	}
	return svc, nil
}

func (s *oidcService) Apply(o OIDC) error {
	if o.ExactMatchClaims == nil {
		o.ExactMatchClaims = map[string]string{}
	}
	if o.AuthorizationCustomParams == nil {
		o.AuthorizationCustomParams = map[string]string{}
	}
	if o.TokenRequestCustomParams == nil {
		o.TokenRequestCustomParams = map[string]string{}
	}

	body, err := tools.JsonMarshalInterfaceToIOReader(o)
	if err != nil {
		return err
	}

	respBody, resp, err := s.c.Put(oidcAPIEndpoint, body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK &&
		resp.StatusCode != http.StatusCreated &&
		resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("could not create/update OIDC configuration: HTTP %d: %s", resp.StatusCode, string(respBody))
	}
	return nil
}

func (s *oidcService) Read() (*OIDC, error) {
	respBody, resp, err := s.c.Get(oidcAPIEndpoint, nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("could not read OIDC configuration: HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	out := &OIDC{}
	if err := json.Unmarshal(respBody, out); err != nil {
		return nil, fmt.Errorf("could not unmarshal OIDC configuration: %w", err)
	}
	return out, nil
}

func (s *oidcService) Delete() error {
	respBody, resp, err := s.c.Delete(oidcAPIEndpoint)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK &&
		resp.StatusCode != http.StatusNoContent &&
		resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("could not delete OIDC configuration: HTTP %d: %s", resp.StatusCode, string(respBody))
	}
	return nil
}
