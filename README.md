[![dns-lookup-go license](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)
[![dns-lookup-go made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](https://pkg.go.dev/github.com/whois-api-llc/dns-lookup-go)
[![dns-lookup-go test](https://github.com/whois-api-llc/dns-lookup-go/workflows/Test/badge.svg)](https://github.com/whois-api-llc/dns-lookup-go/actions/)

# Overview

The client library for
[DNS Lookup API](https://dns-lookup.whoisxmlapi.com/)
in Go language.

The minimum go version is 1.17.

# Installation

The library is distributed as a Go module

```bash
go get github.com/whois-api-llc/dns-lookup-go
```

# Examples

Full API documentation available [here](https://dns-lookup.whoisxmlapi.com/api/documentation/making-requests)

You can find all examples in `example` directory.

## Create a new client

To start making requests you need the API Key. 
You can find it on your profile page on [whoisxmlapi.com](https://whoisxmlapi.com/).
Using the API Key you can create Client.

Most users will be fine with `NewBasicClient` function. 
```go
client := dnslookupapi.NewBasicClient(apiKey)
```

If you want to set custom `http.Client` to use proxy then you can use `NewClient` function.
```go
transport := &http.Transport{Proxy: http.ProxyURL(proxyUrl)}

client := dnslookupapi.NewClient(apiKey, dnslookupapi.ClientParams{
    HTTPClient: &http.Client{
        Transport: transport,
        Timeout:   20 * time.Second,
    },
})
```

## Make basic requests

DNS Lookup API lets you get well-structured a domain’s corresponding IP address from its A record as well as the domain’s mail server (MX record), nameserver (NS record), SPF (TXT record), and more records.

```go

// Make request to get all parsed DNS records for the domain name
dnsLookupResp, resp, err := client.Get(ctx, "whoisxmlapi.com")
if err != nil {
    log.Fatal(err)
}

for _, record := range dnsLookupResp.DNSRecords.A {
    log.Println(record.Name)
    log.Println(record.Address)
}

// Make request to get raw DNS Lookup API data
resp, err := client.GetRaw(context.Background(), "whoisxmlapi.com")
if err != nil {
    log.Fatal(err)
}

log.Println(string(resp.Body))


```
