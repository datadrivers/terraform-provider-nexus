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

	result := bytes.NewReader(b)
	return result, nil
}

func (c client) UserCreate(user User) error {
	ioReader, err := userIOReader(user)
	if err != nil {
		return err
	}

	body, resp, err := c.Post(usersAPIEndpoint, ioReader)
	if err != nil {
		return fmt.Errorf("could not create user '%s': %v", user.UserID, err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("could not create user '%s': HTTP: %d, %s", user.UserID, resp.StatusCode, string(body))
	}

	return nil
}

func (c client) UserRead(userId string) (*User, error) {
	body, resp, err := c.Get(usersAPIEndpoint, nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("could not get users: HTTP: %d, %s", resp.StatusCode, string(body))
	}

	var users []User
	if err := json.Unmarshal(body, users); err != nil {
		return nil, fmt.Errorf("could not unmarschal user '%s': %v", userId, err)
	}

	for _, user := range users {
		if user.UserID == userId {
			return &user, nil
		}
	}

	return nil, nil
}

func (c client) UserUpdate(userId string, user User) error {
	if user.Source == "" {
		user.Source = "default"
	}
	ioReader, err := userIOReader(user)
	if err != nil {
		return err
	}
	body, resp, err := c.Put(fmt.Sprintf("%s/%s", usersAPIEndpoint, userId), ioReader)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("could not update user '%s': HTTP: %d, %s", userId, resp.StatusCode, string(body))
	}

	return nil
}

func (c client) UserDelete(userId string) error {
	body, resp, err := c.Delete(fmt.Sprintf("%s/%s", usersAPIEndpoint, userId))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("could not delete user '%s': HTTP: %d, %s", userId, resp.StatusCode, string(body))
	}
	return err
}
