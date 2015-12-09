package shopify

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
)

// RemoteResource abstracts http requests.
type RemoteResource interface {
	// Request performs the given HTTP request and returns either an io.ReadCloser for the
	// response body or an error.
	Request(*http.Request) (io.ReadCloser, error)
}

type httpRemoteResource struct {
	client              *http.Client
	authenticateRequest AuthenticateRequest
	settings            ClientSettings
}

func (rr *httpRemoteResource) Request(req *http.Request) (io.ReadCloser, error) {
	if rr.settings.DumpRequestURLs {
		log.Printf("%s %s\n", req.Method, req.URL.String())
	}
	rr.authenticateRequest(req)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	if rr.settings.DumpRequests {
		rr.dumpRequest(req)
	}
	resp, err := rr.client.Do(req)
	if err != nil {
		return nil, err
	}
	if rr.settings.DumpResponses {
		rr.dumpResponse(resp)
	}

	if err := FromResponse(resp); err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func (rr *httpRemoteResource) dumpRequest(req *http.Request) {
	dump, err := httputil.DumpRequest(req, true)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%s", dump)
}

func (rr *httpRemoteResource) dumpResponse(resp *http.Response) {
	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%s", dump)
}

func NewRemoteResource(settings ClientSettings) RemoteResource {
	client := &http.Client{
		Timeout: settings.Timeout,
	}
	return &httpRemoteResource{
		client:              client,
		authenticateRequest: settings.AuthenticateRequest,
		settings:            settings,
	}
}

type RemoteJSONResource interface {
	URLBuilder
	RemoteResource
	RequestAndDecode(*http.Request, string, interface{}) error
}

type ShopifyRemoteJSONResource struct {
	URLBuilder
	RemoteResource
}

func (sr *ShopifyRemoteJSONResource) RequestAndDecode(req *http.Request, element string, v interface{}) error {
	reader, err := sr.Request(req)
	if err != nil {
		return err
	}
	defer reader.Close()

	decoder := json.NewDecoder(reader)
	var raw map[string]json.RawMessage
	if err = decoder.Decode(&raw); err != nil {
		return err
	}

	data, found := raw[element]
	if !found {
		return fmt.Errorf("Element \"%s\" could not be found in response from server.", element)
	}

	return json.Unmarshal(data, v)
}
