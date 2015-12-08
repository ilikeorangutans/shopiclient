package shopify

import (
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
	RemoteJSONResource
}

func (o *Orders) Get(id int) (*Order, error) {
	req, err := http.NewRequest("GET", o.BuildURL(fmt.Sprintf("orders/%d", id)), nil)
	if err != nil {
		log.Fatal(err)
	}

	var order *Order
	err = o.RequestAndDecode(req, "order", &order)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (o *Orders) List() ([]*Order, error) {
	req, err := http.NewRequest("GET", o.BuildURL("orders"), nil)
	if err != nil {
		return nil, err
	}

	var orders []*Order
	err = o.RequestAndDecode(req, "orders", &orders)
	if err != nil {
		return nil, err
	}

	return orders, nil

}
