package shopify

import (
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
	RemoteJSONResource
}

func (t *Themes) List() ([]*Theme, error) {
	req, err := http.NewRequest("GET", t.BuildURL("themes"), nil)
	if err != nil {
		return nil, err
	}

	var themes []*Theme
	if err = t.RequestAndDecode(req, "themes", &themes); err != nil {
		return nil, err
	}

	return themes, nil
}

func (t *Themes) Get(id int64) (*Theme, error) {
	req, err := http.NewRequest("GET", t.BuildURL(fmt.Sprintf("themes/%d.json", id)), nil)
	if err != nil {
		return nil, err
	}

	var theme *Theme
	if err := t.RequestAndDecode(req, "theme", &theme); err != nil {
		return nil, err
	} else {
		return theme, nil
	}
}

func (t *Theme) String() string {
	return fmt.Sprintf("Theme{id: %d, name: %s, role: %s, theme_store_id: %d, processing: %t, previewable: %t}", t.ID, t.Name, t.Role, t.ThemeStoreID, t.Processing, t.Previewable)
}
