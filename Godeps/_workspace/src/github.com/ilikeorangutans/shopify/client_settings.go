package shopify

import (
	"fmt"
	"net/http"
	"regexp"
	"time"
)

type AuthenticateRequest func(*http.Request)

type ClientSettings struct {
	Host, Username, Password    string
	Timeout                     time.Duration
	DumpRequestURLs             bool
	DumpRequests, DumpResponses bool
}

func (cs ClientSettings) ShopURL() string {
	pattern := regexp.MustCompile("https?://")
	if pattern.MatchString(cs.Host) {
		return cs.Host
	} else {
		return fmt.Sprintf("https://%s", cs.Host)
	}
}

func (cs ClientSettings) AuthenticateRequest(req *http.Request) {
	if len(cs.Username) > 0 && len(cs.Password) > 0 {
		req.SetBasicAuth(cs.Username, cs.Password)
	}
}

type RequestAuthenticator interface {
	AuthenticateRequest(*http.Request)
}
