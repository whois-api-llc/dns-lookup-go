package dnslookupapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
)

const (
	pathDNSLookupResponseOK         = "/DNSLookup/ok"
	pathDNSLookupResponseError      = "/DNSLookup/error"
	pathDNSLookupResponse500        = "/DNSLookup/500"
	pathDNSLookupResponsePartial1   = "/DNSLookup/partial"
	pathDNSLookupResponsePartial2   = "/DNSLookup/partial2"
	pathDNSLookupResponseUnparsable = "/DNSLookup/unparsable"
)

const apiKey = "at_LoremIpsumDolorSitAmetConsect"

// dummyServer is the sample of the DNS Lookup API server for testing.
func dummyServer(resp, respUnparsable string, respErr string) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var response string

		response = resp

		switch req.URL.Path {
		case pathDNSLookupResponseOK:
		case pathDNSLookupResponseError:
			w.WriteHeader(499)
			response = respErr
		case pathDNSLookupResponse500:
			w.WriteHeader(500)
			response = respUnparsable
		case pathDNSLookupResponsePartial1:
			response = response[:len(response)-10]
		case pathDNSLookupResponsePartial2:
			w.Header().Set("Content-Length", strconv.Itoa(len(response)))
			response = response[:len(response)-10]
		case pathDNSLookupResponseUnparsable:
			response = respUnparsable
		default:
			panic(req.URL.Path)
		}
		_, err := w.Write([]byte(response))
		if err != nil {
			panic(err)
		}
	}))

	return server
}

// newAPI returns new DNS Lookup API client for testing.
func newAPI(apiServer *httptest.Server, link string) *Client {
	apiURL, err := url.Parse(apiServer.URL)
	if err != nil {
		panic(err)
	}

	apiURL.Path = link

	params := ClientParams{
		HTTPClient:       apiServer.Client(),
		DNSLookupBaseURL: apiURL,
	}

	return NewClient(apiKey, params)
}

// TestDNSLookupGet tests the Get function.
func TestDNSLookupGet(t *testing.T) {
	checkResultRec := func(res *DNSLookupResponse) bool {
		return res != nil
	}

	ctx := context.Background()

	const resp = ` {"DNSData": {
  "domainName": "whoisxmlapi.com",
  "types": [1],
  "dnsTypes": "A",
  "audit": {"createdDate": "2022-07-12 11:46:25 UTC","updatedDate": "2022-07-12 11:46:25 UTC"},
  "dnsRecords": [
    {
      "type": 1,
      "dnsType": "A",
      "name": "whoisxmlapi.com.",
      "ttl": 300,
      "rRsetType": 1,
      "rawText": "whoisxmlapi.com.\u0009300\u0009IN\u0009A\u0009104.26.13.210",
      "address": "104.26.13.210"
    }
]
}}`

	const respUnparsable = `<?xml version="1.0" encoding="utf-8"?><>`

	const errResp = `{"ErrorMessage":{"errorCode":"TEST_CODE","msg":"test error message"}}`

	server := dummyServer(resp, respUnparsable, errResp)
	defer server.Close()

	type args struct {
		ctx     context.Context
		options string
	}

	tests := []struct {
		name    string
		path    string
		args    args
		want    bool
		wantErr string
	}{
		{
			name: "successful request",
			path: pathDNSLookupResponseOK,
			args: args{
				ctx:     ctx,
				options: "whoisxmlapi.com",
			},
			want:    true,
			wantErr: "",
		},
		{
			name: "non 200 status code",
			path: pathDNSLookupResponse500,
			args: args{
				ctx:     ctx,
				options: "whoisxmlapi.com",
			},
			want:    false,
			wantErr: "cannot parse response: invalid character '<' looking for beginning of value",
		},
		{
			name: "partial response 1",
			path: pathDNSLookupResponsePartial1,
			args: args{
				ctx:     ctx,
				options: "whoisxmlapi.com",
			},
			want:    false,
			wantErr: "cannot parse response: unexpected EOF",
		},
		{
			name: "partial response 2",
			path: pathDNSLookupResponsePartial2,
			args: args{
				ctx:     ctx,
				options: "whoisxmlapi.com",
			},
			want:    false,
			wantErr: "cannot read response: unexpected EOF",
		},
		{
			name: "could not process request",
			path: pathDNSLookupResponseError,
			args: args{
				ctx:     ctx,
				options: "whoisxmlapi.com",
			},
			want:    false,
			wantErr: "API error: [TEST_CODE] test error message",
		},
		{
			name: "unparsable response",
			path: pathDNSLookupResponseUnparsable,
			args: args{
				ctx:     ctx,
				options: "whoisxmlapi.com",
			},
			want:    false,
			wantErr: "cannot parse response: invalid character '<' looking for beginning of value",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newAPI(server, tt.path)

			gotRec, _, err := api.Get(tt.args.ctx, tt.args.options)
			if (err != nil || tt.wantErr != "") && (err == nil || err.Error() != tt.wantErr) {
				t.Errorf("DNSLookup.Get() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if tt.want {
				if !checkResultRec(gotRec) {
					t.Errorf("DNSLookup.Get() got = %v, expected something else", gotRec)
				}
			} else {
				if gotRec != nil {
					t.Errorf("DNSLookup.Get() got = %v, expected nil", gotRec)
				}
			}
		})
	}
}

// TestDNSLookupGetRaw tests the GetRaw function.
func TestDNSLookupGetRaw(t *testing.T) {
	checkResultRaw := func(res []byte) bool {
		return len(res) != 0
	}

	ctx := context.Background()

	const resp = ` {"DNSData": {
  "domainName": "whoisxmlapi.com",
  "types": [1],
  "dnsTypes": "A",
  "audit": {"createdDate": "2022-07-12 11:46:25 UTC","updatedDate": "2022-07-12 11:46:25 UTC"},
  "dnsRecords": [
    {
      "type": 1,
      "dnsType": "A",
      "name": "whoisxmlapi.com.",
      "ttl": 300,
      "rRsetType": 1,
      "rawText": "whoisxmlapi.com.\u0009300\u0009IN\u0009A\u0009104.26.13.210",
      "address": "104.26.13.210"
    }
]
}}`

	const respUnparsable = `<?xml version="1.0" encoding="utf-8"?><>`

	const errResp = `{"ErrorMessage":{"errorCode":"TEST_CODE","msg":"test error message"}}`

	server := dummyServer(resp, respUnparsable, errResp)
	defer server.Close()

	type args struct {
		ctx     context.Context
		options string
	}

	tests := []struct {
		name    string
		path    string
		args    args
		wantErr string
	}{
		{
			name: "successful request",
			path: pathDNSLookupResponseOK,
			args: args{
				ctx:     ctx,
				options: "whoisxmlapi.com",
			},
			wantErr: "",
		},
		{
			name: "non 200 status code",
			path: pathDNSLookupResponse500,
			args: args{
				ctx:     ctx,
				options: "whoisxmlapi.com",
			},
			wantErr: "API failed with status code: 500",
		},
		{
			name: "partial response 1",
			path: pathDNSLookupResponsePartial1,
			args: args{
				ctx:     ctx,
				options: "whoisxmlapi.com",
			},
			wantErr: "",
		},
		{
			name: "partial response 2",
			path: pathDNSLookupResponsePartial2,
			args: args{
				ctx:     ctx,
				options: "whoisxmlapi.com",
			},
			wantErr: "cannot read response: unexpected EOF",
		},
		{
			name: "unparsable response",
			path: pathDNSLookupResponseUnparsable,
			args: args{
				ctx:     ctx,
				options: "whoisxmlapi.com",
			},
			wantErr: "",
		},
		{
			name: "could not process request",
			path: pathDNSLookupResponseError,
			args: args{
				ctx:     ctx,
				options: "whoisxmlapi.com",
			},
			wantErr: "API failed with status code: 499",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newAPI(server, tt.path)

			resp, err := api.GetRaw(tt.args.ctx, tt.args.options)
			if (err != nil || tt.wantErr != "") && (err == nil || err.Error() != tt.wantErr) {
				t.Errorf("DNSLookup.Get() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !checkResultRaw(resp.Body) {
				t.Errorf("DNSLookup.Get() got = %v, expected something else", string(resp.Body))
			}
		})
	}
}
