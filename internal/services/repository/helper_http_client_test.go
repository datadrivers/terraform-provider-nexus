package repository

import (
	"testing"

	"github.com/gcroucher/go-nexus-client/nexus3/schema/repository"
	"github.com/stretchr/testify/assert"
)

func TestGetHTTPClientConnection(t *testing.T) {
	emptyConnectionList := []interface{}{}
	result := getHTTPClientConnection(emptyConnectionList)
	assert.Nil(t, result)

	enableCircularRedirects := true
	enableCookies := true
	retries := 3
	timeout := 15
	useTrustStore := true

	connectionList := []interface{}{
		map[string]interface{}{
			"enable_circular_redirects": enableCircularRedirects,
			"enable_cookies":            enableCookies,
			"retries":                   retries,
			"timeout":                   timeout,
			"user_agent_suffix":         "",
			"use_trust_store":           useTrustStore,
		},
	}
	expected := repository.HTTPClientConnection{
		EnableCircularRedirects: &enableCircularRedirects,
		EnableCookies:           &enableCookies,
		Retries:                 &retries,
		Timeout:                 &timeout,
		UserAgentSuffix:         "",
		UseTrustStore:           &useTrustStore,
	}
	result = getHTTPClientConnection(connectionList)
	assert.Equal(t, *expected.EnableCircularRedirects, *result.EnableCircularRedirects)
	assert.Equal(t, *expected.EnableCookies, *result.EnableCookies)
	assert.Equal(t, *expected.Retries, *result.Retries)
	assert.Equal(t, *expected.Timeout, *result.Timeout)
	assert.Equal(t, expected.UserAgentSuffix, result.UserAgentSuffix)
	assert.Equal(t, *expected.UseTrustStore, *result.UseTrustStore)

	connectionListWithoutTimeout := []interface{}{
		map[string]interface{}{
			"enable_circular_redirects": enableCircularRedirects,
			"enable_cookies":            enableCookies,
			"retries":                   retries,
			"user_agent_suffix":         "",
			"use_trust_store":           useTrustStore,
		},
	}
	resultWithoutTimeout := getHTTPClientConnection(connectionListWithoutTimeout)
	assert.Nil(t, resultWithoutTimeout.Timeout)
}
