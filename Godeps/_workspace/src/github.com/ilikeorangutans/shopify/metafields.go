package shopify

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

type Metafield struct {
	ID          int       `json:"id"`
	Namespace   string    `json:"namespace"`
	Key         string    `json:"key"`
	Description string    `json:"description"`
	Resource    string    `json:"owner_resource"`
	Value       string    `json:"value"`
	Type        string    `json:"value_type"`
	ResourceID  int       `json:"owner_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type metafieldsResponse struct {
	Metafields []*Metafield `json:"metafields"`
}

type Metafields struct {
	requester  Requester
	urlBuilder URLBuilder
}

func (mf *Metafields) Create(under, namespace, key, value, valueType string) *Metafield {
	payload := fmt.Sprintf("{\"metafield\":{\"namespace\": \"%s\", \"key\":\"%s\", \"value\":\"%s\", \"value_type\": \"%s\", \"created_at\":Mon, 13 Apr 2015 16:22:16 -0400}}", namespace, key, value, valueType)

	req, err := http.NewRequest("POST", mf.urlBuilder("/admin", under, "metafields.json"), strings.NewReader(payload))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	body, _ := httputil.DumpRequest(req, true)
	log.Printf("%s", body)

	resp, err := mf.requester(req)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Respnose", resp)

	return nil
}

func (mf *Metafields) List(under string) []*Metafield {
	req, err := http.NewRequest("GET", mf.urlBuilder("/admin", under, "metafields.json"), nil)
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		log.Fatal(err)
	}

	resp, err := mf.requester(req)
	if err != nil {
		log.Fatal(err)
	}

	var metafields []*Metafield
	json.Unmarshal(resp["metafields"], &metafields)

	return metafields
}
