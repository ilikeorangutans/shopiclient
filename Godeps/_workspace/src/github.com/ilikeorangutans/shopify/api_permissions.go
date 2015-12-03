package shopify

import (
	"encoding/json"
	"log"
	"net/http"
)

type APIPermission struct {
	ID int `json:"id"`
}

type APIPermissions struct {
	requester  Requester
	urlBuilder URLBuilder
}

func (ap *APIPermissions) List() []*APIPermission {
	req, err := http.NewRequest("GET", ap.urlBuilder("/admin/apps.json"), nil)
	if err != nil {
		log.Fatal(err)
	}

	d, err := ap.requester(req)
	if err != nil {
		log.Fatal(err)
	}

	var permissions []*APIPermission
	json.Unmarshal(d["api_permissions"], &permissions)

	return permissions
}
