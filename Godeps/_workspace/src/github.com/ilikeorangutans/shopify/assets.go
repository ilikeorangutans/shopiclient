package shopify

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Assets struct {
	RemoteJSONResource
	Theme *Theme
}

// List downloads metadata for all assets associated with the Theme set on the instance.
func (a *Assets) List() ([]*Asset, error) {
	req, err := http.NewRequest("GET", a.BuildURL(a.themeBaseURL(), "assets"), nil)
	if err != nil {
		return nil, err
	}

	var assets []*Asset
	if err := a.RequestAndDecode(req, "assets", &assets); err != nil {
		return nil, err
	}
	return assets, nil
}

// DownloadAll downloads all assets including their attachments. This can cause large requests!
func (a *Assets) DownloadAll() ([]*Asset, error) {
	req, err := http.NewRequest("GET", a.BuildURL(a.themeBaseURL(), "assets?fields=key,value,attachment"), nil)
	if err != nil {
		return nil, err
	}
	var assets []*Asset
	if err := a.RequestAndDecode(req, "assets", &assets); err != nil {
		return nil, err
	}

	for i := range assets {
		asset := assets[i]
		asset.DecodingComplete = make(chan struct{})
		go asset.decodeAttachment()
	}

	return assets, nil
}

// Download downloads a single Asset identified by the given key with all its data.
func (a *Assets) Download(key string) (*Asset, error) {
	req, err := http.NewRequest("GET", a.BuildURL(a.themeBaseURL(), fmt.Sprintf("assets?asset[key]=%s", key)), nil)
	if err != nil {
		return nil, err
	}
	var asset *Asset
	if err := a.RequestAndDecode(req, "asset", &asset); err != nil {
		return nil, err
	}
	asset.DecodingComplete = make(chan struct{})
	if err := asset.decodeAttachment(); err != nil {
		return nil, err
	}

	return asset, nil
}

func (a *Assets) Upload(asset *Asset) (*Asset, error) {
	envelope := make(map[string]interface{})
	envelope["asset"] = asset
	payload, err := json.Marshal(envelope)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", a.BuildURL(a.themeBaseURL(), "assets.json"), bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}

	var response *Asset
	err = a.RequestAndDecode(req, "asset", &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (a *Assets) themeBaseURL() string {
	return fmt.Sprintf("themes/%d", a.Theme.ID)
}

// Asset is a single asset in a Shopify theme. Assets are uniquely identified in a theme by their Key field
// and can have either a (string) value or a binary attachment. Attachments are retrieved BASE64 encoded and
// have to be decoded prior to usage. To check if either decoding or encoding is complete one has to read
// from the DecodingComplete or EncodingComplete channels. Reads from either channel will block until the
// operations are complete.
type Asset struct {
	Timestamps

	Key         string `json:"key"`
	ContentType string `json:"content_type,omitempty"`
	PublicURL   string `json:"public_url,omitempty"`
	Size        int    `json:"size,omitempty"`
	ThemeID     int64  `json:"theme_id,omitempty"`
	Value       string `json:"value,omitempty"`
	// Attachment holds the binary attachment of this asset, if available. Note that you should check the
	// DecodingComplete channel on this asset to ensure decoding is complete.
	Attachment []byte `json:"-"`
	// EncodedAttachment holds a base64 encoded representation of the attachment.
	EncodedAttachment string `json:"attachment,omitempty"`
	// DecodingComplete is a channel that blocks until decoding of this asset's attachment is complete.
	DecodingComplete chan struct{} `json:"-"`
	EncodingComplete chan struct{} `json:"-"`
	Src              string        `json:"src,omitempty"`
}

func (a *Asset) HasAttachment() bool {
	return len(a.Attachment) > 0 || len(a.EncodedAttachment) > 0
}

func (a *Asset) String() string {
	return fmt.Sprintf("Asset{key: %s, content_type: %s, size: %d}", a.Key, a.ContentType, a.Size)
}

func (a *Asset) decodeAttachment() error {
	defer close(a.DecodingComplete)
	if len(a.EncodedAttachment) == 0 {
		return nil
	}
	b, err := base64.StdEncoding.DecodeString(a.EncodedAttachment)
	if err != nil {
		return err
	}
	if a.Size != len(b) {
		return fmt.Errorf("Attachment length does not match expected value, expected %d bytes but got %d", a.Size, len(b))
	}
	a.Attachment = b
	return nil
}

func (a *Asset) encodeAttachment() {
	encodedAttachment := base64.StdEncoding.EncodeToString(a.Attachment)
	a.EncodedAttachment = encodedAttachment
}

func NewAssetWithAttachment(key string, reader io.Reader) (*Asset, error) {
	attachment, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return &Asset{
		Key:        key,
		Attachment: attachment,
	}, nil
}

func NewAssetWithValue(key, value string) (*Asset, error) {
	return &Asset{
		Key:   key,
		Value: value}, nil
}

func NewAssetWithURL(key, src string) (*Asset, error) {
	srcURL, err := url.Parse(src)
	if err != nil {
		return nil, err
	}
	return &Asset{
		Key: key,
		Src: srcURL.String(),
	}, nil
}
