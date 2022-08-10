package dnslookupapi

import (
	"net/url"
	"strings"
)

// Option adds parameters to the query.
type Option func(v url.Values)

var _ = []Option{
	OptionOutputFormat("JSON"),
	OptionType("A"),
	OptionCallback("func"),
}

// OptionOutputFormat sets Response output format JSON | XML. Default: JSON.
func OptionOutputFormat(outputFormat string) Option {
	return func(v url.Values) {
		v.Set("outputFormat", strings.ToUpper(outputFormat))
	}
}

// OptionType sets types of DNS records that should be returned. DNS type: A, NS, SOA, MX, etc.
// You can specify multiple comma-separated values, e.g., A,SOA,TXT;
// all records can be retrieved with type=_all (by default).
func OptionType(value string) Option {
	return func(v url.Values) {
		v.Set("type", strings.ToUpper(value))
	}
}

// OptionCallback sets a javascript function used when outputFormat is JSON;
// this is an implementation known as JSONP which invokes the callback on the returned response.
func OptionCallback(value string) Option {
	return func(v url.Values) {
		v.Set("callback", value)
	}
}
