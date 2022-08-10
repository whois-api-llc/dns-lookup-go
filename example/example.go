package example

import (
	"context"
	"errors"
	dnslookupapi "github.com/whois-api-llc/dns-lookup-go"
	"log"
)

func GetData(apikey string) {
	client := dnslookupapi.NewBasicClient(apikey)

	// Get parsed DNS Lookup API response as a model instance
	dnsLookupResp, resp, err := client.Get(context.Background(),
		"whoisxmlapi.com",
		// this option is ignored, as the inner parser works with JSON only
		dnslookupapi.OptionOutputFormat("XML"))

	if err != nil {
		// Handle error message returned by server
		var apiErr *dnslookupapi.ErrorMessage
		if errors.As(err, &apiErr) {
			log.Println(apiErr.Code)
			log.Println(apiErr.Message)
		}
		log.Fatal(err)
	}

	// Then print some values from each returned DNS record.
	for _, obj := range dnsLookupResp.DNSRecords.All {
		log.Printf("DomainName: %s, DNSType: %s, RawText: %s\n",
			obj.CommonFields.Name,
			obj.CommonFields.DNSType,
			obj.CommonFields.RawText,
		)
	}

	// Or just print IP addresses from A records.
	for _, record := range dnsLookupResp.DNSRecords.A {
		log.Println(record.Address)
	}

	log.Println("raw response is always in JSON format. Most likely you don't need it.")
	log.Printf("raw response: %s\n", string(resp.Body))
}

func GetRawData(apikey string) {
	client := dnslookupapi.NewBasicClient(apikey)

	// Get raw API response
	resp, err := client.GetRaw(context.Background(),
		"whoisxmlapi.com",
		dnslookupapi.OptionOutputFormat("JSON"),
		// this option causes the only NS and MX records to be returned
		dnslookupapi.OptionType("NS,MX"))

	if err != nil {
		// Handle error message returned by server
		log.Fatal(err)
	}

	log.Println(string(resp.Body))
}
