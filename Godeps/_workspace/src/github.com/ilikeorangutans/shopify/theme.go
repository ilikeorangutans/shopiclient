package shopify

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Theme struct {
	CommonFields

	Name         string `json:"name"`
	Role         string `json:"role"`
	ThemeStoreID *int64 `json:"theme_store_id"`
	Previewable  bool   `json:"previewable"`
	Processing   bool   `json:"processing"`
}

type Themes struct {
	requester       Requester
	urlBuilder      URLBuilder
	requestAndParse RequestAndParse
}

func (t *Themes) List() ([]*Theme, error) {
	req, err := http.NewRequest("GET", t.urlBuilder("/admin/themes"), nil)
	if err != nil {
		return nil, err
	}

	if themes, err := t.requestAndParse(req, "themes", decodeThemesList); err != nil {
		return nil, err
	} else {
		return themes.([]*Theme), nil
	}
}

func (t *Themes) Get(id int64) (*Theme, error) {
	req, err := http.NewRequest("GET", t.urlBuilder(fmt.Sprintf("/admin/themes/%d.json", id)), nil)
	if err != nil {
		return nil, err
	}

	if theme, err := t.requestAndParse(req, "theme", decodeTheme); err != nil {
		return nil, err
	} else {
		return theme.(*Theme), nil
	}
}

func (t *Theme) String() string {
	return fmt.Sprintf("Theme{id: %d, name: %s, role: %s, theme_store_id: %d, processing: %t, previewable: %t}", t.ID, t.Name, t.Role, t.ThemeStoreID, t.Processing, t.Previewable)
}

func (t *Theme) Assets() *Assets {
	return &Assets{}
}

func decodeThemesList(body []byte) (interface{}, error) {
	var themes []*Theme
	if err := json.Unmarshal(body, &themes); err != nil {
		return nil, err
	}
	return themes, nil
}

func decodeTheme(body []byte) (interface{}, error) {
	var theme *Theme
	if err := json.Unmarshal(body, &theme); err != nil {
		return nil, err
	}
	return theme, nil
}
