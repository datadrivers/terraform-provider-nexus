package security

import (
	"fmt"
	"strings"
)

// isRoleNotFoundError reports whether err is the error go-nexus-client returns
// when Nexus answers a role lookup with 404. The client only passes the response
// body on, so match the message Nexus sends for a missing role.
func isRoleNotFoundError(err error, roleID string) bool {
	return err != nil && strings.Contains(err.Error(), fmt.Sprintf("Role '%s' not found", roleID))
}

// isLDAPNotFoundError reports whether err is the error go-nexus-client returns
// when Nexus answers an LDAP server lookup with 404.
func isLDAPNotFoundError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "HTTP: 404")
}
