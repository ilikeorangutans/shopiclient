package shopify

import (
	"fmt"
	"net/http"
	"time"
)

const DefaultTimeout = time.Duration(10 * time.Second)

// Client is the facade for all API connections to Shopify. Obtain a new instance with the NewClient function.
type Client struct {
	RemoteJSONResource
	Settings ClientSettings
	client   *http.Client
	Verbose  bool
}

// NewClient returns a new client with default settings and the given host, username and password.
func NewClient(host, username, password string) *Client {
	settings := ClientSettings{host: host, username: username, password: password, timeout: DefaultTimeout, dumpRequestURLs: true}
	return NewClientWithSettings(settings)
}

// NewClientWithSettings creates a new client with the given settings.
func NewClientWithSettings(settings ClientSettings) *Client {
	client := &http.Client{
		Timeout: settings.timeout,
	}
	return &Client{
		client:   client,
		Settings: settings,
		RemoteJSONResource: &ShopifyRemoteJSONResource{
			URLBuilder:     &ShopifyAdminURLBuilder{baseURL: settings.ShopURL()},
			RemoteResource: NewRemoteResource(settings),
		},
		Verbose: true,
	}
}

func NewClientWithSettingsAndRemoteResource(settings ClientSettings, remote RemoteJSONResource) *Client {
	return &Client{
		Settings:           settings,
		RemoteJSONResource: remote,
	}
}

// Connect attempts a connection to the configured server. If the server responds with a 4xx or 5xx
// status code, this method will return an error. It's recommended to use this function at least
// once after creating a new client to ensure valid settings and credentials.
func (c *Client) Connect() error {
	req, err := http.NewRequest("HEAD", fmt.Sprintf("%s/admin/shop", c.Settings.ShopURL()), nil)
	if err != nil {
		return err
	}
	_, err = c.Request(req)
	if err != nil {
		return fmt.Errorf("Error connecting to server: \"%s\"", err.Error())
	}

	return nil
}

func (c *Client) Webhooks() *Webhooks {
	return &Webhooks{RemoteJSONResource: c}
}

func (c *Client) Apps() *APIPermissions {
	return &APIPermissions{RemoteJSONResource: c}
}

func (c *Client) Metafields() *Metafields {
	return &Metafields{RemoteJSONResource: c}
}

func (c *Client) FullfillmentServices() *FulfillmentServices {
	return &FulfillmentServices{RemoteJSONResource: c}
}

func (c *Client) Orders() *Orders {
	return &Orders{RemoteJSONResource: c}
}

func (c *Client) Transactions() *Transactions {
	return &Transactions{RemoteJSONResource: c}
}

func (c *Client) Themes() *Themes {
	return &Themes{RemoteJSONResource: c}
}

func (c *Client) Products() *Products {
	return &Products{RemoteJSONResource: c}
}

func (c *Client) Assets(theme *Theme) *Assets {
	return &Assets{
		RemoteJSONResource: c,
		Theme:              theme,
	}
}
