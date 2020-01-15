package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	rolesAPIEndpoint = "service/rest/beta/security/roles"
)

// Role ...
type Role struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Privileges  []string `json:"privileges"`
	Roles       []string `json:"roles"`
}

func roleIOReader(role Role) (io.Reader, error) {
	b, err := json.Marshal(role)
	if err != nil {
		return nil, fmt.Errorf("could not marshal role data: %v", err)
	}

	return bytes.NewReader(b), nil
}

func (c client) RoleCreate(role Role) error {
	ioReader, err := roleIOReader(role)
	if err != nil {
		return err
	}

	body, resp, err := c.Post(rolesAPIEndpoint, ioReader)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s", string(body))
	}

	return nil
}

func (c client) RoleRead(id string) (*Role, error) {
	body, resp, err := c.Get(rolesAPIEndpoint, nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s", string(body))
	}

	var roles []Role
	if err := json.Unmarshal(body, roles); err != nil {
		return nil, fmt.Errorf("could not unmarshal roles: %v", err)
	}

	for _, role := range roles {
		if role.ID == id {
			return &role, nil
		}
	}

	return nil, fmt.Errorf("could not find role '%s'", id)
}

func (c client) RoleUpdate(id string, role Role) error {
	ioReader, err := roleIOReader(role)
	if err != nil {
		return err
	}

	body, resp, err := c.Put(fmt.Sprintf("%s/%s", rolesAPIEndpoint, id), ioReader)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("%s", string(body))
	}

	return nil
}

func (c client) RoleDelete(id string) error {
	body, resp, err := c.Delete(fmt.Sprintf("%s/%s", rolesAPIEndpoint, id))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("%s", string(body))
	}

	return nil
}
