package shopify

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Transactions struct {
	requester  Requester
	urlBuilder URLBuilder
}

type Transaction struct {
	ID       int     `json:"id"`
	Amount   float64 `json:"amount,string"`
	Currency string  `json:"currency"`
	Status   string  `json:"status"`
	Test     bool    `json:"test"`
}

func (t *Transactions) List(order *Order) []*Transaction {
	req, err := http.NewRequest("GET", t.urlBuilder(fmt.Sprintf("/admin/orders/%d/transactions.json", order.ID)), nil)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := t.requester(req)
	if err != nil {
		log.Fatal(err)
	}

	return decodeTransactionsList(resp["transactions"])
}

func decodeTransactionsList(body []byte) []*Transaction {
	var transactions []*Transaction
	json.Unmarshal(body, &transactions)
	return transactions
}
