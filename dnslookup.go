package dnslookupapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// DNSLookupService is an interface for DNS Lookup API.
type DNSLookupService interface {
	// Get returns parsed DNS Lookup API response
	Get(ctx context.Context, domainName string, opts ...Option) (*DNSLookupResponse, *Response, error)

	// GetRaw returns raw DNS Lookup API response as Response struct with Body saved as a byte slice
	GetRaw(ctx context.Context, domainName string, opts ...Option) (*Response, error)
}

// Response is the http.Response wrapper with Body saved as a byte slice.
type Response struct {
	*http.Response

	// Body is the byte slice representation of http.Response Body
	Body []byte
}

// dnsLookupServiceOp is the type implementing the DNSLookupService interface.
type dnsLookupServiceOp struct {
	client  *Client
	baseURL *url.URL
}

var _ DNSLookupService = &dnsLookupServiceOp{}

// newRequest creates the API request with default parameters and the specified apiKey.
func (service *dnsLookupServiceOp) newRequest() (*http.Request, error) {
	req, err := service.client.NewRequest(http.MethodGet, service.baseURL, nil)
	if err != nil {
		return nil, err
	}

	query := url.Values{}
	query.Set("apiKey", service.client.apiKey)

	req.URL.RawQuery = query.Encode()

	return req, nil
}

// apiResponse is used for parsing DNS Lookup API response as a model instance.
type apiResponse struct {
	DNSLookupResponse `json:"DNSData"`
	ErrorMessage      `json:"ErrorMessage"`
}

// request returns intermediate API response for further actions.
func (service *dnsLookupServiceOp) request(ctx context.Context, domainName string, opts ...Option) (*Response, error) {
	req, err := service.newRequest()
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Set("domainName", domainName)
	q.Set("type", "_all")

	for _, opt := range opts {
		opt(q)
	}

	req.URL.RawQuery = q.Encode()

	var b bytes.Buffer

	resp, err := service.client.Do(ctx, req, &b)
	if err != nil {
		return &Response{
			Response: resp,
			Body:     b.Bytes(),
		}, err
	}

	return &Response{
		Response: resp,
		Body:     b.Bytes(),
	}, nil
}

// parse parses raw DNS Lookup API response.
func parse(raw []byte) (*apiResponse, error) {
	var response apiResponse

	err := json.NewDecoder(bytes.NewReader(raw)).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("cannot parse response: %w", err)
	}

	return &response, nil
}

// Get returns parsed DNS Lookup API response.
func (service dnsLookupServiceOp) Get(
	ctx context.Context,
	domainName string,
	opts ...Option,
) (dnsLookupResponse *DNSLookupResponse, resp *Response, err error) {
	optsJSON := make([]Option, 0, len(opts)+1)
	optsJSON = append(optsJSON, opts...)
	optsJSON = append(optsJSON, OptionOutputFormat("JSON"))

	resp, err = service.request(ctx, domainName, optsJSON...)
	if err != nil {
		return nil, resp, err
	}

	dnsLookupResp, err := parse(resp.Body)
	if err != nil {
		return nil, resp, err
	}

	if dnsLookupResp.Message != "" || dnsLookupResp.Code != "" {
		return nil, nil, &ErrorMessage{
			Code:    dnsLookupResp.Code,
			Message: dnsLookupResp.Message,
		}
	}

	return &dnsLookupResp.DNSLookupResponse, resp, nil
}

// GetRaw returns raw DNS Lookup API response as Response struct with Body saved as a byte slice.
func (service dnsLookupServiceOp) GetRaw(
	ctx context.Context,
	domainName string,
	opts ...Option,
) (resp *Response, err error) {
	resp, err = service.request(ctx, domainName, opts...)
	if err != nil {
		return resp, err
	}

	if respErr := checkResponse(resp.Response); respErr != nil {
		return resp, respErr
	}

	return resp, nil
}

// ArgError is the argument error.
type ArgError struct {
	Name    string
	Message string
}

// Error returns error message as a string.
func (a *ArgError) Error() string {
	return `invalid argument: "` + a.Name + `" ` + a.Message
}
