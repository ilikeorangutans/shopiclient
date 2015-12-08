package shopify

import (
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
	RemoteJSONResource
}

func (mf *Metafields) Create(under, namespace, key, value, valueType string) *Metafield {
	payload := fmt.Sprintf("{\"metafield\":{\"namespace\": \"%s\", \"key\":\"%s\", \"value\":\"%s\", \"value_type\": \"%s\", \"created_at\":Mon, 13 Apr 2015 16:22:16 -0400}}", namespace, key, value, valueType)

	req, err := http.NewRequest("POST", mf.BuildURL(under, "metafields"), strings.NewReader(payload))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	body, _ := httputil.DumpRequest(req, true)
	log.Printf("%s", body)

	resp, err := mf.Request(req)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Respnose", resp)

	return nil
}

func (mf *Metafields) List(under string) ([]*Metafield, error) {
	req, err := http.NewRequest("GET", mf.BuildURL(under, "metafields"), nil)
	if err != nil {
		return nil, err
	}

	var metafields []*Metafield
	err = mf.RequestAndDecode(req, "metafields", &metafields)
	if err != nil {
		return nil, err
	}

	return metafields, nil
}
