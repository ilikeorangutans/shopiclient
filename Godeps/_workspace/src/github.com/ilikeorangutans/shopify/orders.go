package shopify

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Order struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Number int    `json:"number"`
}

type Orders struct {
	requester  Requester
	urlBuilder URLBuilder
}

func (o *Orders) Get(id int) (*Order, error) {
	req, err := http.NewRequest("GET", o.urlBuilder(fmt.Sprintf("/admin/orders/%d.json", id)), nil)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := o.requester(req)
	if err != nil {
		return nil, err
	}
	var order *Order
	json.Unmarshal(resp["order"], &order)
	return order, nil
}

func (o *Orders) List() []*Order {
	req, err := http.NewRequest("GET", o.urlBuilder("/admin/orders.json"), nil)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := o.requester(req)
	if err != nil {
		log.Fatal(err)
	}

	return decodeOrdersList(resp["orders"])
}

func decodeOrdersList(body []byte) []*Order {
	var orders []*Order
	json.Unmarshal(body, &orders)
	return orders
}
