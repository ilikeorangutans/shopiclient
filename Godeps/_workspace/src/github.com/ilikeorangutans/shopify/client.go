package shopify

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"
	"time"
)

const DefaultTimeout = time.Duration(10 * time.Second)

type Requester func(req *http.Request) (map[string]json.RawMessage, error)

type URLBuilder func(...string) string

type ClientSettings struct {
	host, username, password string
	timeout                  time.Duration
}

func (cs ClientSettings) ShopURL() string {
	pattern := regexp.MustCompile("https?://")
	if pattern.MatchString(cs.host) {
		return cs.host
	} else {
		return fmt.Sprintf("https://%s", cs.host)
	}
}

func (cs ClientSettings) AuthenticateRequest(req *http.Request) {
	if len(cs.username) > 0 && len(cs.password) > 0 {
		req.SetBasicAuth(cs.username, cs.password)
	}
}

type RequestAuthenticator interface {
	AuthenticateRequest(*http.Request)
}

// Client is the facade for all API connections to Shopify. Obtain a new instance with the NewClient function.
type Client struct {
	Settings ClientSettings
	client   *http.Client
	Verbose  bool
}

// NewClient returns a new client with default settings and the given host, username and password.
func NewClient(host, username, password string) *Client {
	settings := ClientSettings{host: host, username: username, password: password, timeout: DefaultTimeout}
	return NewClientWithSettings(settings)
}

// NewClientWithSettings creates a new client with the given settings.
func NewClientWithSettings(settings ClientSettings) *Client {
	client := &http.Client{
		Timeout: settings.timeout,
	}
	return &Client{client: client, Settings: settings}
}

// Connect attempts a connection to the configured server. If the server responds with a 4xx or 5xx
// status code, this method will return an error. It's recommended to use this function at least
// once after creating a new client to ensure valid settings and credentials.
func (c *Client) Connect() error {

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/admin/", c.Settings.ShopURL()), nil)
	if err != nil {
		return err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("Error connecting to server: \"%s\"", err.Error())
	} else if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		return fmt.Errorf("Server responded with \"%s\", check credentials", resp.Status)
	} else if resp.StatusCode >= 500 {
		return fmt.Errorf("Error connecting to server: %s", resp.Status)
	}

	return nil
}

func (c *Client) Webhooks() *Webhooks {
	return &Webhooks{requester: c.doRequest, urlBuilder: c.buildURL}
}

func (c *Client) Apps() *APIPermissions {
	return &APIPermissions{requester: c.doRequest, urlBuilder: c.buildURL}
}

func (c *Client) Metafields() *Metafields {
	return &Metafields{requester: c.doRequest, urlBuilder: c.buildURL}
}

func (c *Client) FullfillmentServices() *FulfillmentServices {
	return &FulfillmentServices{requester: c.doRequest, urlBuilder: c.buildURL}
}

func (c *Client) Orders() *Orders {
	return &Orders{requester: c.doRequest, urlBuilder: c.buildURL}
}

func (c *Client) Transactions() *Transactions {
	return &Transactions{requester: c.doRequest, urlBuilder: c.buildURL}
}

func (c *Client) Themes() *Themes {
	return &Themes{requestAndParse: c.requestAndParseJSON, urlBuilder: c.buildURL}
}

func (c *Client) Assets(theme *Theme) *Assets {
	return AssetsForTheme(theme)
}

func (c *Client) debug(msg string) {
	if c.Verbose {
		log.Println(msg)
	}
}

func (c *Client) doRequest(req *http.Request) (map[string]json.RawMessage, error) {
	c.debug(fmt.Sprintf("%s: %s \n", req.Method, req.URL))
	c.Settings.AuthenticateRequest(req)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	b, _ := httputil.DumpResponse(resp, true)
	c.debug(fmt.Sprintf("Response: \n%s", b))

	if err := FromResponse(resp); err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(resp.Body)

	var d map[string]json.RawMessage

	decoder.Decode(&d)

	return d, nil
}

func (c *Client) buildURL(input ...string) string {
	url, err := url.Parse("https://" + c.Settings.host + strings.Join(input, "/"))
	if err != nil {
		log.Fatal(err)
	}

	return url.String()
}

type JSONResourceParser func(body []byte) (interface{}, error)

type RequestAndParse func(*http.Request, string, JSONResourceParser) (interface{}, error)

func (c *Client) requestAndParseJSON(req *http.Request, element string, parser JSONResourceParser) (interface{}, error) {
	raw, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	data, found := raw[element]
	if !found {
		return nil, fmt.Errorf("Element \"%s\" could not be found in response from server.", element)
	}

	return parser(data)
}
