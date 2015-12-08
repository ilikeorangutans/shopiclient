package shopify

import (
	"fmt"
	"net/http"
)

type Transactions struct {
	RemoteJSONResource
}

type Transaction struct {
	ID       int     `json:"id"`
	Amount   float64 `json:"amount,string"`
	Currency string  `json:"currency"`
	Status   string  `json:"status"`
	Test     bool    `json:"test"`
}

func (t *Transactions) List(order *Order) ([]*Transaction, error) {
	req, err := http.NewRequest("GET", t.BuildURL(fmt.Sprintf("orders/%d/transactions", order.ID)), nil)
	if err != nil {
		return nil, err
	}

	var transactions []*Transaction
	err = t.RequestAndDecode(req, "transactions", &transactions)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}
