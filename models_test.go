package dnslookupapi

import (
	"encoding/json"
	"testing"
)

// TestTime tests JSON encoding/parsing functions for the time values
func TestTime(t *testing.T) {
	tests := []struct {
		name   string
		decErr string
		encErr string
	}{
		{
			name:   `"2006-01-02 15:04:05 EST"`,
			decErr: "",
			encErr: "",
		},
		{
			name:   `"2006-01-02 12:04:05 UTC"`,
			decErr: "",
			encErr: "",
		},
		{
			name:   `"2006-01-02T15:04:05-07:00"`,
			decErr: `parsing time "2006-01-02T15:04:05-07:00" as "2006-01-02 15:04:05 MST": cannot parse "T15:04:05-07:00" as " "`,
			encErr: "",
		},
		{
			name:   `""`,
			decErr: "",
			encErr: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var v Time

			err := json.Unmarshal([]byte(tt.name), &v)
			checkErr(t, err, tt.decErr)
			if tt.decErr != "" {
				return
			}

			bb, err := json.Marshal(v)
			checkErr(t, err, tt.encErr)
			if tt.encErr != "" {
				return
			}

			if string(bb) != tt.name {
				t.Errorf("got = %v, want %v", string(bb), tt.name)
			}
		})
	}
}

// TestDNSRecords tests JSON encoding/parsing functions for DNSRecords
func TestDNSRecords(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output string
		decErr string
		encErr string
	}{
		{
			name:   `test-1`,
			input:  `[]`,
			output: `[]`,
			decErr: "",
			encErr: "",
		},
		{
			name: `test-2`,
			input: `[
{
      "type": 1,
      "dnsType": "A",
      "name": "whoisxmlapi.com.",
      "ttl": 300,
      "rRsetType": 1,
      "rawText": "whoisxmlapi.com.\u0009300\u0009IN\u0009A\u0009172.67.71.123",
      "address": "172.67.71.123"
    },
    {
      "type": 2,
      "dnsType": "NS",
      "name": "whoisxmlapi.com.",
      "additionalName": "elle.ns.cloudflare.com.",
      "ttl": 21600,
      "rRsetType": 2,
      "rawText": "whoisxmlapi.com.\u000921600\u0009IN\u0009NS\u0009elle.ns.cloudflare.com.",
      "target": "elle.ns.cloudflare.com."
    }
]`,
			output: `[{"CommonFields":{"type":1,"dnsType":"A","name":"whoisxmlapi.com.","ttl":300,"rRsetType":1,"rawText":"whoisxmlapi.com.\t300\tIN\tA\t172.67.71.123"},"raw":{"type":1,"dnsType":"A","name":"whoisxmlapi.com.","ttl":300,"rRsetType":1,"rawText":"whoisxmlapi.com.\u0009300\u0009IN\u0009A\u0009172.67.71.123","address":"172.67.71.123"},"parseError":null},{"CommonFields":{"type":2,"dnsType":"NS","name":"whoisxmlapi.com.","ttl":21600,"rRsetType":2,"rawText":"whoisxmlapi.com.\t21600\tIN\tNS\telle.ns.cloudflare.com."},"raw":{"type":2,"dnsType":"NS","name":"whoisxmlapi.com.","additionalName":"elle.ns.cloudflare.com.","ttl":21600,"rRsetType":2,"rawText":"whoisxmlapi.com.\u000921600\u0009IN\u0009NS\u0009elle.ns.cloudflare.com.","target":"elle.ns.cloudflare.com."},"parseError":null}]`,
			decErr: "",
			encErr: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var v *DNSRecords

			err := json.Unmarshal([]byte(tt.input), &v)
			checkErr(t, err, tt.decErr)
			if tt.decErr != "" {
				return
			}

			bb, err := json.Marshal(v)
			checkErr(t, err, tt.encErr)
			if tt.encErr != "" {
				return
			}

			if string(bb) != tt.output {
				t.Errorf("got  = %v", string(bb))
				t.Errorf("want = %v", tt.output)
			}
		})
	}
}

// checkErr checks for an error.
func checkErr(t *testing.T, err error, want string) {
	if (err != nil || want != "") && (err == nil || err.Error() != want) {
		t.Errorf("error = %v, wantErr %v", err, want)
	}
}
