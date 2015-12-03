package shopify

import (
	"encoding/json"
	"fmt"
)

func AssetsForTheme(t *Theme) *Assets {
	return &Assets{}
}

type Assets struct {
	buildURL URLBuilder
	request  Requester
	Theme    Theme
}

// List downloads metadata for all assets associated with the Theme set on the instance.
func (a *Assets) List() ([]*Asset, error) {
	a.buildURL(a.themeBaseURL(), "assets.json")
	return nil, nil
}

func (a *Assets) themeBaseURL() string {
	return fmt.Sprintf("/admin/themes/%d", a.Theme.ID)
}

type Asset struct {
	CommonFields

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

func decodeAssetsList(body []byte) ([]*Asset, error) {
	var assets []*Asset
	err := json.Unmarshal(body, &assets)
	if err != nil {
		return nil, err
	}
	return assets, nil
}
