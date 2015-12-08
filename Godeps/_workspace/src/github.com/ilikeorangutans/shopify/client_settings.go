package shopify

import (
	"fmt"
	"net/http"
	"regexp"
	"time"
)

type AuthenticateRequest func(*http.Request)

type ClientSettings struct {
	host, username, password    string
	timeout                     time.Duration
	dumpRequestURLs             bool
	dumpRequests, dumpResponses bool
}

func (cs ClientSettings) ShopURL() string {
	pattern := regexp.MustCompile("https?://")
	if pattern.MatchString(cs.host) {
		return cs.host
	} else {
		return fmt.Sprintf("https://%s", cs.host)
	}
}

func (cs ClientSettings) AuthenticateRequest(req *http.Request) {
	if len(cs.username) > 0 && len(cs.password) > 0 {
		req.SetBasicAuth(cs.username, cs.password)
	}
}

type RequestAuthenticator interface {
	AuthenticateRequest(*http.Request)
}
