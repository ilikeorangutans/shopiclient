package shopify

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Assets struct {
	buildURL        URLBuilder
	requestAndParse RequestAndParse
	Theme           *Theme
}

// List downloads metadata for all assets associated with the Theme set on the instance.
func (a *Assets) List() ([]*Asset, error) {

	req, err := http.NewRequest("GET", a.buildURL(a.themeBaseURL(), "assets.json"), nil)
	if err != nil {
		return nil, err
	}
	assets, err := a.requestAndParse(req, "assets", decodeAssetsList)
	if err != nil {
		return nil, err
	}
	return assets.([]*Asset), nil
}

func (a *Assets) themeBaseURL() string {
	return fmt.Sprintf("/admin/themes/%d", a.Theme.ID)
}

type Asset struct {
	Timestamps

	Key                string          `json:"key"`
	ContentType        string          `json:"content_type"`
	PublicURL          string          `json:"public_url"`
	Size               int             `json:"size"`
	ThemeID            int64           `json:"theme_id"`
	unparsedAttachment json.RawMessage `json:"attachment"`
	Value              string          `json:"value"`
	Attachment         []byte          `json:"-"`
	metadataOnly       bool            `json:"-"`
}

func (a *Asset) String() string {
	return fmt.Sprintf("Asset{key: %s, content_type: %s, size: %d}", a.Key, a.ContentType, a.Size)
}

func decodeAssetsList(body []byte) (interface{}, error) {
	var assets []*Asset
	err := json.Unmarshal(body, &assets)
	if err != nil {
		return nil, err
	}
	return assets, nil
}
