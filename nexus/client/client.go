package client

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

// Client represents the Nexus API Client interface
type Client interface {
	Post(string, io.Reader) ([]byte, *http.Response, error)
	Put(string, io.Reader) ([]byte, *http.Response, error)
	Delete(string) ([]byte, *http.Response, error)
	UserCreate(User) error
	UserRead(string) (*User, error)
	UserUpdate(string, User) error
	UserDelete(string) error
}

type client struct {
	config Config
	client *http.Client
}

// NewClient returns an instance of client that implements the Client interface
func NewClient(config Config) Client {
	return &client{
		config: config,
		client: &http.Client{},
	}
}

func (c client) NewRequest(method string, endpoint string, body io.Reader) (req *http.Request, err error) {
	url := fmt.Sprintf("%s/%s", c.config.URL, endpoint)
	req, err = http.NewRequest(method, url, body)
	if err != nil {
		return req, err
	}

	req.SetBasicAuth(c.config.Username, c.config.Password)
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func (c client) execute(method string, endpoint string, payload io.Reader) ([]byte, *http.Response, error) {
	req, err := c.NewRequest(method, endpoint, payload)
	if err != nil {
		return nil, nil, fmt.Errorf("could not create equest: %v", err)
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("could not execute request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	return body, resp, err
}

func (c client) Get(endpoint string, payload io.Reader) ([]byte, *http.Response, error) {
	return c.execute(http.MethodGet, endpoint, payload)
}

func (c client) Post(endpoint string, payload io.Reader) ([]byte, *http.Response, error) {
	return c.execute(http.MethodPost, endpoint, payload)
}

func (c client) Put(endpoint string, payload io.Reader) ([]byte, *http.Response, error) {
	return c.execute(http.MethodPut, endpoint, payload)
}

func (c client) Delete(endpoint string) ([]byte, *http.Response, error) {
	return c.execute(http.MethodDelete, endpoint, nil)
}
