package repository

import "github.com/datadrivers/go-nexus-client/nexus3/schema/repository"

// TerraformHostedRepository represents a Nexus Terraform hosted repository.
type TerraformHostedRepository struct {
	Name      string                    `json:"name"`
	Online    bool                      `json:"online"`
	Storage   repository.HostedStorage  `json:"storage"`
	Cleanup   *repository.Cleanup       `json:"cleanup,omitempty"`
	Component *repository.Component     `json:"component,omitempty"`
}

// TerraformProxyRepository represents a Nexus Terraform proxy repository.
type TerraformProxyRepository struct {
	Name            string                   `json:"name"`
	Online          bool                     `json:"online"`
	Storage         repository.Storage       `json:"storage"`
	Cleanup         *repository.Cleanup      `json:"cleanup,omitempty"`
	Proxy           repository.Proxy         `json:"proxy"`
	NegativeCache   repository.NegativeCache `json:"negativeCache"`
	HTTPClient      repository.HTTPClient    `json:"httpClient"`
	RoutingRule     *string                  `json:"routingRule,omitempty"`
	RoutingRuleName *string                  `json:"routingRuleName,omitempty"`
}

// TerraformGroupRepository represents a Nexus Terraform group repository.
type TerraformGroupRepository struct {
	Name    string             `json:"name"`
	Online  bool               `json:"online"`
	Storage repository.Storage `json:"storage"`
	Group   repository.Group   `json:"group"`
}
