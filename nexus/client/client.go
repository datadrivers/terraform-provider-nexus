package client

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type Client interface {
	Post(string, io.Reader) ([]byte, *http.Response, error)
	Put(string, io.Reader) ([]byte, *http.Response, error)
	Delete(string) (*http.Response, error)
}

type client struct {
	config Config
}

// NewClient returns an instance of client that implements the Client interface
func NewClient(config Config) Client {
	return &client{
		config: config,
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
		return nil, nil, fmt.Errorf("could not create new request: %v", err)
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, nil, fmt.Errorf("could not execute request: %v", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	return body, resp, err
}

func (c client) Post(endpoint string, payload io.Reader) ([]byte, *http.Response, error) {
	return c.execute(http.MethodPost, endpoint, payload)
}

func (c client) Put(endpoint string, payload io.Reader) ([]byte, *http.Response, error) {
	return c.execute(http.MethodPut, endpoint, payload)
}

func (c client) Delete(endpoint string) (*http.Response, error) {
	_, resp, err := c.execute(http.MethodDelete, endpoint, nil)
	return resp, err
}
