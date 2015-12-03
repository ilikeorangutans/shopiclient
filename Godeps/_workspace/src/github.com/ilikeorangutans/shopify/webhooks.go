package shopify

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Webhook struct {
	CommonFields

	Format  string `json:"format"`
	Topic   string `json:"topic"`
	Address string `json:"address"`
}

type Webhooks struct {
	requester  Requester
	urlBuilder URLBuilder
}

func (webhooks *Webhooks) Create(topic string, address *url.URL, format string) (*Webhook, error) {

	payload := fmt.Sprintf("{\"webhook\":{\"topic\":\"%s\", \"address\":\"%s\", \"format\": \"%s\"}}", topic, address.String(), format)
	req, err := http.NewRequest("POST", webhooks.urlBuilder("/admin/webhooks.json"), strings.NewReader(payload))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := webhooks.requester(req)
	if err != nil {
		log.Fatal(err)
	}

	var webhook *Webhook
	json.Unmarshal(resp["webhook"], &webhook)

	return webhook, nil
}

func (webhooks *Webhooks) Delete(id int64) {
	webhook, err := webhooks.Get(id)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Deleting webhook [%d] %s %s %s", webhook.ID, webhook.Topic, webhook.Format, webhook.Address)

	req, err := http.NewRequest("POST", webhooks.urlBuilder(fmt.Sprintf("/admin/webhooks/%d.json", id)), strings.NewReader("_method=delete"))
	if err != nil {
		log.Fatal(err)
	}

	_, err = webhooks.requester(req)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Webhook %d deleted\n", id)
}

func (ws *Webhooks) List() []*Webhook {
	req, err := http.NewRequest("GET", ws.urlBuilder("/admin/webhooks.json"), nil)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := ws.requester(req)
	if err != nil {
		log.Fatal(err)
	}

	var webhooks []*Webhook
	json.Unmarshal(resp["webhooks"], &webhooks)

	return webhooks
}

func (webhooks *Webhooks) Get(id int64) (*Webhook, error) {
	req, err := http.NewRequest("GET", webhooks.urlBuilder(fmt.Sprintf("/admin/webhooks/%d.json", id)), nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := webhooks.requester(req)
	if err != nil {
		log.Fatal(err)
	}
	var webhook *Webhook
	json.Unmarshal(resp["webhook"], &webhook)

	return webhook, nil
}
