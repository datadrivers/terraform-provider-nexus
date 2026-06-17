package security

import (
	"encoding/json"
	"fmt"
	"net/http"

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

// configuredOIDCService is populated by ConfigureOIDC during providerConfigure.
// Package-scoped because go-nexus-client v1.20.0 does not expose an OIDC
// service via NexusClient and the embedded `*client.Client` is unexported.
var configuredOIDCService *oidcService

func ConfigureOIDC(c *client.Client) {
	configuredOIDCService = &oidcService{c: c}
}

func oidc() (*oidcService, error) {
	if configuredOIDCService == nil {
		return nil, fmt.Errorf("nexus OIDC client not configured")
	}
	return configuredOIDCService, nil
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
