package shopify

import (
	"log"
	"net/url"
	"strings"
)

type URLBuilder interface {
	BuildURL(segments ...string) string
}

type ShopifyAdminURLBuilder struct {
	baseURL string
}

func (sa *ShopifyAdminURLBuilder) BuildURL(segments ...string) string {
	url, err := url.Parse(sa.baseURL + "/admin/" + strings.Join(segments, "/"))
	if err != nil {
		log.Fatal(err)
	}

	return url.String()
}
