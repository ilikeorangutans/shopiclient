package shopify

import (
	"net/http"
)

type APIPermission struct {
	ID int `json:"id"`
}

type APIPermissions struct {
	RemoteJSONResource
}

func (ap *APIPermissions) List() ([]*APIPermission, error) {
	req, err := http.NewRequest("GET", ap.BuildURL("apps"), nil)
	if err != nil {
		return nil, err
	}

	var permissions []*APIPermission
	if err = ap.RequestAndDecode(req, "api_permissions", &permissions); err != nil {
		return nil, err
	}

	return permissions, nil
}
