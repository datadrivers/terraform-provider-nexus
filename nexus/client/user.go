package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	usersAPIEndpoint = "service/rest/beta/security/users"
)

// User ..
type User struct {
	UserID       string   `json:"userId"`
	FirstName    string   `json:"firstName"`
	LastName     string   `json:"lastName"`
	EmailAddress string   `json:"emailAddress"`
	Password     string   `json:"password"`
	Status       string   `json:"status"`
	Source       string   `json:"source"`
	Roles        []string `json:"roles"`
}

func userIOReader(user User) (io.Reader, error) {
	b, err := json.Marshal(user)
	if err != nil {
		return nil, fmt.Errorf("could not marshal user data: %v", err)
	}

	return bytes.NewReader(b), nil
}

func (c client) UserCreate(user User) error {
	ioReader, err := userIOReader(user)
	if err != nil {
		return err
	}

	body, resp, err := c.Post(usersAPIEndpoint, ioReader)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s", string(body))
	}

	return nil
}

func (c client) UserRead(id string) (*User, error) {
	body, resp, err := c.Get(usersAPIEndpoint, nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s", string(body))
	}

	var users []User
	if err := json.Unmarshal(body, users); err != nil {
		return nil, fmt.Errorf("could not unmarschal users: %v", err)
	}

	for _, user := range users {
		if user.UserID == id {
			return &user, nil
		}
	}

	return nil, fmt.Errorf("could not find user '%s'", id)
}

func (c client) UserUpdate(id string, user User) error {
	// Not sure what this is and why is required to update a user
	if user.Source == "" {
		user.Source = "default"
	}

	ioReader, err := userIOReader(user)
	if err != nil {
		return err
	}

	body, resp, err := c.Put(fmt.Sprintf("%s/%s", usersAPIEndpoint, id), ioReader)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("%s", string(body))
	}

	return nil
}

func (c client) UserDelete(id string) error {
	body, resp, err := c.Delete(fmt.Sprintf("%s/%s", usersAPIEndpoint, id))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("%s", string(body))
	}
	return err
}

// UserChangePassword  must be send with content-type text/plain :-/
func (c client) UserChangePassword(id string, password string) error {
	data := bytes.NewReader([]byte(password))
	body, resp, err := c.Put(fmt.Sprintf("%s/%s/change-password", usersAPIEndpoint, id), data)
	if err != nil {
		return fmt.Errorf("could not change password of user '%s': %v", id, err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("could not change password of user '%s':  HTTP: %d, %s ", id, resp.StatusCode, string(body))
	}
	return nil
}
