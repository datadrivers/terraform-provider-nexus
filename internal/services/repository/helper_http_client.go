package repository

import (
	"github.com/datadrivers/terraform-provider-nexus/internal/tools"
	"github.com/gcroucher/go-nexus-client/nexus3/schema/repository"
)

func getHTTPClientConnection(connectionList []interface{}) *repository.HTTPClientConnection {
	if len(connectionList) == 1 && connectionList[0] != nil {
		connectionConfig := connectionList[0].(map[string]interface{})
		connection := repository.HTTPClientConnection{
			EnableCircularRedirects: tools.GetBoolPointer(connectionConfig["enable_circular_redirects"].(bool)),
			EnableCookies:           tools.GetBoolPointer(connectionConfig["enable_cookies"].(bool)),
			Retries:                 tools.GetIntPointer(connectionConfig["retries"].(int)),
			UserAgentSuffix:         connectionConfig["user_agent_suffix"].(string),
			UseTrustStore:           tools.GetBoolPointer(connectionConfig["use_trust_store"].(bool)),
		}
		timeout, ok := connectionConfig["timeout"]
		if ok && timeout.(int) != 0 {
			connection.Timeout = tools.GetIntPointer(timeout.(int))
		}
		return &connection
	}
	return nil
}
