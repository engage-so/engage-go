package engage

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	defaultTimeout = time.Second * 60

	// base URL for the Engage API
	baseURL = "https://api.engage.so"

	// version of the engage-go library
	version = "0.1.0"

	// user agent
	userAgent = "engage-go:" + version
)

var (
	errNoAPIKey = errors.New("API key not set")
)

// HTTPClient interface
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// resource represents a shared resource struct
type resource struct {
	client *Client
}

// Client represents the Engage Client
type Client struct {
	client HTTPClient

	// The Engage API key
	key string

	// The Engage API secret
	secret string

	baseURL *url.URL

	resource resource

	// Engage client Resources
	User *UserResource
}

// HTTPResponse ...
type HTTPResponse struct {
	Code int
	Data []byte
}

// New creates a new Engage API Client with the given API key
// and secret pair
func New(key, secret string) (*Client, error) {
	if key == "" || secret == "" {
		return nil, errNoAPIKey
	}

	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	c := &Client{
		client:  &http.Client{Timeout: defaultTimeout},
		key:     key,
		secret:  secret,
		baseURL: u,
	}

	c.resource.client = c
	// initialize resources
	c.User = (*UserResource)(&c.resource)

	return c, nil
}

// ParseJSON ...
func (r *HTTPResponse) ParseJSON(v interface{}) error {
	return json.Unmarshal(r.Data, v)
}

// SetClient updates the HTTP client
func (c *Client) SetClient(client HTTPClient) {
	c.client = client
}

func (c *Client) postRequest(endpoint string, body interface{}) (*HTTPResponse, error) {
	return c.makeRequest("POST", endpoint, body)
}

func (c *Client) putRequest(endpoint string, body interface{}) (*HTTPResponse, error) {
	return c.makeRequest("PUT", endpoint, body)
}

// MakeRequest makes an HTTP request to the Engage API
func (c *Client) makeRequest(method, endpoint string, body interface{}) (*HTTPResponse, error) {
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	u, _ := c.baseURL.Parse(endpoint)
	req, err := http.NewRequest(method, u.String(), buf)

	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}
	req.Header.Add("User-Agent", userAgent)

	req.SetBasicAuth(c.key, c.secret)

	response := HTTPResponse{}

	resp, err := c.client.Do(req)
	if err != nil {
		if clErr, ok := err.(*url.Error); ok {
			if clErr.Err == io.EOF {
				return nil, fmt.Errorf("remote server prematurely closed connection: %v", err)
			}
		}
		return nil, fmt.Errorf("while making http request: %v", err)
	}

	if resp != nil {
		response.Code = resp.StatusCode
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	response.Data = respBody

	return &response, nil
}
