package other

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	nexus "github.com/datadrivers/go-nexus-client/nexus3"
	"github.com/datadrivers/go-nexus-client/nexus3/pkg/client"
	"github.com/datadrivers/go-nexus-client/nexus3/pkg/tools"
)

// REST endpoint: https://help.sonatype.com/en/http-configuration-api.html
const httpClientAPIEndpoint = client.BasePath + "v1/http"

type HTTPAuthInfo struct {
	Enabled    bool   `json:"enabled"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	NtlmHost   string `json:"ntlmHost"`
	NtlmDomain string `json:"ntlmDomain"`
}

type HTTPProxy struct {
	Enabled  bool         `json:"enabled"`
	Host     string       `json:"host"`
	Port     interface{}  `json:"port"`
	AuthInfo HTTPAuthInfo `json:"authInfo"`
}

type HTTPClientConfig struct {
	NonProxyHosts []string  `json:"nonProxyHosts"`
	UserAgent     string    `json:"userAgent"`
	Timeout       int       `json:"timeout"`
	Retries       int       `json:"retries"`
	HTTPProxy     HTTPProxy `json:"httpProxy"`
	HTTPSProxy    HTTPProxy `json:"httpsProxy"`
}

func (p HTTPProxy) portInt() int {
	switch v := p.Port.(type) {
	case float64:
		return int(v)
	case int:
		return v
	case json.Number:
		n, _ := v.Int64()
		return int(n)
	case string:
		var n int
		fmt.Sscanf(v, "%d", &n)
		return n
	default:
		return 0
	}
}

type httpClientService struct {
	c *client.Client
}

var (
	httpClientServicesMu sync.Mutex
	httpClientServices   = map[*nexus.NexusClient]*httpClientService{}
)

func ConfigureHTTPClient(nc *nexus.NexusClient, c *client.Client) {
	httpClientServicesMu.Lock()
	defer httpClientServicesMu.Unlock()
	httpClientServices[nc] = &httpClientService{c: c}
}

func httpClient(nc *nexus.NexusClient) (*httpClientService, error) {
	httpClientServicesMu.Lock()
	defer httpClientServicesMu.Unlock()
	svc, ok := httpClientServices[nc]
	if !ok {
		return nil, fmt.Errorf("nexus HTTP client not configured")
	}
	return svc, nil
}

func (s *httpClientService) Apply(cfg HTTPClientConfig) error {
	if cfg.NonProxyHosts == nil {
		cfg.NonProxyHosts = []string{}
	}

	body, err := tools.JsonMarshalInterfaceToIOReader(cfg)
	if err != nil {
		return err
	}

	respBody, resp, err := s.c.Put(httpClientAPIEndpoint, body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK &&
		resp.StatusCode != http.StatusCreated &&
		resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("could not create/update HTTP configuration: HTTP %d: %s", resp.StatusCode, string(respBody))
	}
	return nil
}

func (s *httpClientService) Read() (*HTTPClientConfig, error) {
	respBody, resp, err := s.c.Get(httpClientAPIEndpoint, nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("could not read HTTP configuration: HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	out := &HTTPClientConfig{}
	if err := json.Unmarshal(respBody, out); err != nil {
		return nil, fmt.Errorf("could not unmarshal HTTP configuration: %w", err)
	}
	return out, nil
}

func (s *httpClientService) Delete() error {
	respBody, resp, err := s.c.Delete(httpClientAPIEndpoint)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK &&
		resp.StatusCode != http.StatusNoContent &&
		resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("could not delete HTTP configuration: HTTP %d: %s", resp.StatusCode, string(respBody))
	}
	return nil
}
