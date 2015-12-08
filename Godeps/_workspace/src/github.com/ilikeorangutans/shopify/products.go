package shopify

import (
	"net/http"
)

type Products struct {
	RemoteJSONResource
}

type Product struct {
	CommonFields
	Title string `json:"title"`
}

func (p *Products) List() ([]*Product, error) {
	url := p.BuildURL("products")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var products []*Product
	err = p.RequestAndDecode(req, "products", &products)
	if err != nil {
		return nil, err
	}

	return products, nil
}
