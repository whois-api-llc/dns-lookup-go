package dnslookupapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

var ErrUnsupportedDNSType = errors.New("unknown DNS type")

// unmarshalString parses the JSON-encoded data and returns value as a string.
func unmarshalString(raw json.RawMessage) (string, error) {
	var val string
	err := json.Unmarshal(raw, &val)
	if err != nil {
		return "", err
	}
	return val, nil
}

// Time is a helper wrapper on time.Time.
type Time time.Time

var emptyTime Time

// UnmarshalJSON decodes time as DNS Lookup API does.
func (t *Time) UnmarshalJSON(b []byte) error {
	str, err := unmarshalString(b)
	if err != nil {
		return err
	}
	if str == "" {
		*t = emptyTime
		return nil
	}
	v, err := time.Parse("2006-01-02 15:04:05 MST", str)
	if err != nil {
		return err
	}
	*t = Time(v)
	return nil
}

// MarshalJSON encodes time as DNS Lookup API does.
func (t Time) MarshalJSON() ([]byte, error) {
	if t == emptyTime {
		return []byte(`""`), nil
	}
	return []byte(`"` + time.Time(t).Format("2006-01-02 15:04:05 MST") + `"`), nil
}

type commonFields struct {
	// Type is the DNS record type code.
	Type int `json:"type"`

	// DNSType is the DNS record type.
	DNSType string `json:"dnsType"`

	// Name is a domain name.
	Name string `json:"name"`

	// TTL is the time to live of DNS record.
	TTL int `json:"ttl"`

	// RRsetType is the resource record type ID.
	RRsetType int `json:"rRsetType"`

	// RawText is the raw text of DNS record.
	RawText string `json:"rawText"`
}

type ARecord struct {
	commonFields

	// Address is the IPv4 address.
	Address string `json:"address"`
}

type AAAARecord struct {
	commonFields

	// Address is the IPv6 address.
	Address string `json:"address"`
}

type NSRecord struct {
	commonFields

	// Target is the name server.
	Target string `json:"target"`
}

type MXRecord struct {
	commonFields

	// Target is the domain name of a mail server.
	Target string `json:"target"`

	// Priority is the priority field.
	Priority int `json:"priority"`
}

type MDRecord struct {
	commonFields

	// AdditionalName is a compressed domain name which specifies a host which has a mail agent for the domain.
	AdditionalName string `json:"additionalName"`

	// MailAgent is a compressed domain name which specifies a host which has a mail agent for the domain.
	MailAgent string `json:"mailAgent"`
}

type MFRecord struct {
	commonFields

	// AdditionalName is a compressed domain name which specifies a host which has a mail agent for the domain.
	AdditionalName string `json:"additionalName"`

	// MailAgent is a compressed domain name which specifies a host which has a mail agent for the domain.
	MailAgent string `json:"mailAgent"`
}

type MBRecord struct {
	commonFields

	// AdditionalName is a compressed domain name which specifies a host which has the specified mailbox.
	AdditionalName string `json:"additionalName"`

	// Mailbox is a compressed domain name which specifies a host which has the specified mailbox.
	Mailbox string `json:"mailbox"`
}

type SOARecord struct {
	commonFields

	// Admin is the email address of the administrator.
	Admin string `json:"admin"`

	// Host is the primary master name server.
	Host string `json:"host"`

	// Expire is the number of seconds after which secondary name servers should stop answering request
	// if the master does not respond.
	Expire int `json:"expire"`

	// Minimum is the negative response caching TTL.
	Minimum int `json:"minimum"`

	// Refresh is the number of seconds after which secondary name servers should query the master for the SOA record,
	// to detect zone changes.
	Refresh int `json:"refresh"`

	// Retry is the number of seconds after which secondary name servers should retry to request the serial number
	// from the master if the master does not respond.
	Retry int `json:"retry"`

	// Serial is the serial number.
	Serial int `json:"serial"`
}

type TXTRecord struct {
	commonFields

	// Strings is the slice of text strings as part of the TXT record.
	Strings []string `json:"strings"`
}

type CAARecord struct {
	commonFields

	// Flags is the flag byte.
	Flags int `json:"flags"`

	// Tag is the property identifier.
	Tag string `json:"tag"`

	// Value is a sequence of octets representing the property value.
	Value string `json:"value"`
}

type CNAMERecord struct {
	commonFields

	// Alias is an alias for a domain name.
	Alias string `json:"alias"`

	// Target is the target domain name.
	Target string `json:"target"`
}

type DNAMERecord struct {
	commonFields

	// Alias is an alias for a domain name.
	Alias string `json:"alias"`

	// Target is the target domain name.
	Target string `json:"target"`
}

type DNSKEYRecord struct {
	commonFields

	// Algorithm is the public key's cryptographic algorithm.
	Algorithm int `json:"algorithm"`

	// Flags is the Zone Key flag.
	Flags int `json:"flags"`

	// Footprint is the key ID/tag/footprint.
	Footprint int `json:"footprint"`

	// Key holds the public key material.
	Key []string `json:"key"`

	// Protocol is the protocol identifier.
	Protocol int `json:"protocol"`

	// PublicKey is the public key description.
	PublicKey string `json:"publicKey"`
}

type NSEC3PARAMRecord struct {
	commonFields

	// Flags are 8 one-bit flags.
	Flags int `json:"flags"`

	// HashAlgorithm is the cryptographic hash algorithm used to construct the hash-value.
	HashAlgorithm int `json:"hashAlgorithm"`

	// Iterations defines the number of additional times the hash function has been performed.
	Iterations int `json:"iterations"`

	// Salt is a value which appended to the original owner name before hashing.
	Salt []string `json:"salt"`
}

type DSRecord struct {
	commonFields

	// Algorithm lists the algorithm number of the DNSKEY RR.
	Algorithm int `json:"algorithm"`

	// Digest is the digest of a DNSKEY RR.
	Digest []string `json:"digest"`

	// DigestID identifies the algorithm used to construct the digest.
	DigestID int `json:"digestID"`

	// Footprint lists the key tag of the DNSKEY RR.
	Footprint int `json:"footprint"`
}

type NSECRecord struct {
	commonFields

	// Next contains the next hashed owner name in hash order.
	Next string `json:"next"`

	// Types is the type bit maps.
	Types []int `json:"types"`
}

type PTRRecord struct {
	commonFields

	// Target is a domain name.
	Target string `json:"target"`
}

type SRVRecord struct {
	commonFields

	// Port is the port on the target host of the service.
	Port int `json:"port"`

	// Priority is the priority of the target host.
	Priority int `json:"priority"`

	// Target is the domain name of the target host.
	Target string `json:"target"`

	// Weight is a server selection mechanism.
	Weight int `json:"weight"`
}

type LOCRecord struct {
	commonFields

	// Altitude is the altitude of the center of the sphere described by the Size field.
	Altitude float64 `json:"altitude"`

	// HPrecision is the horizontal precision of the data, in centimeters.
	HPrecision float64 `json:"hPrecision"`

	// Latitude is the latitude of the center of the sphere described by the Size field.
	Latitude float64 `json:"latitude"`

	// Longitude is the longitude of the center of the sphere described by the Size field.
	Longitude float64 `json:"longitude"`

	// Size is the diameter of a sphere enclosing the described entity.
	Size float64 `json:"size"`

	// VPrecision is the vertical precision of the data, in centimeters.
	VPrecision float64 `json:"vPrecision"`
}

type NAPTRRecord struct {
	commonFields

	// Flags are flags to control aspects of the rewriting and interpretation of the fields in the record
	// as part of NAPTR record.
	Flags string `json:"flags"`

	// Order is a 16-bit unsigned integer specifying the order in which the NAPTR records MUST be processed.
	Order int `json:"order"`

	// Preference is a 16-bit unsigned integer that specifies the order in which NAPTR records with equal Order values
	// SHOULD be processed.
	Preference int `json:"preference"`

	// Regexp contains a substitution expression that is applied to the original string held by the client.
	Regexp string `json:"regexp"`

	// Replacement is a domain name which is the next domain-name to query for
	// depending on the potential values found in the flags field.
	Replacement string `json:"replacement"`

	// Service specifies the Service Parameters applicable to the delegation path.
	Service string `json:"service"`
}

type HINFORecord struct {
	commonFields

	// CPU specifies the CPU type.
	CPU string `json:"cpu"`

	// OS specifies the operating system type.
	OS string `json:"os"`
}

type RPRecord struct {
	commonFields

	// Mailbox is a domain name that specifies the mailbox for the responsible person.
	Mailbox string `json:"mailbox"`

	// TextDomain is a domain name for which TXT RR's exist.
	TextDomain string `json:"textDomain"`
}

type DLVRecord struct {
	commonFields

	// Algorithm lists the algorithm number of the DNSKEY RR.
	Algorithm int `json:"algorithm"`

	// Digest is the digest of a DNSKEY RR.
	Digest []string `json:"digest"`

	// DigestID identifies the algorithm used to construct the digest.
	DigestID int `json:"digestID"`

	// Footprint lists the key tag of the DNSKEY RR.
	Footprint int `json:"footprint"`
}

type SSHFPRecord struct {
	commonFields

	// Algorithm describes the algorithm of the public key.
	Algorithm int `json:"algorithm"`

	// DigestType describes the message-digest algorithm used to calculate the fingerprint of the public key.
	DigestType int `json:"digestType"`

	// FingerPrint is calculated over the public key blob.
	FingerPrint []string `json:"fingerPrint"`
}

type DHCIDRecord struct {
	commonFields

	// Data is several octets of binary data.
	Data []string `json:"data"`
}
type TLSARecord struct {
	commonFields

	// CertificateAssociationData specifies the "certificate association data" to be matched.
	CertificateAssociationData []string `json:"certificateAssociationData"`

	// CertificateUsage specifies the provided association that will be used to match the certificate
	// presented in the TLS handshake.
	CertificateUsage int `json:"certificateUsage"`

	// MatchingType specifies how the certificate association is presented.
	MatchingType int `json:"matchingType"`

	// Selector specifies which part of the TLS certificate presented by the server will be matched against
	// the association data.
	Selector int `json:"selector"`
}

type NSAPRecord struct {
	commonFields

	// Address is a variable length string of octets containing the NSAP.
	Address string `json:"address"`
}

type NULLRecord struct {
	commonFields

	// Data is anything, so long as it is 65535 octets or less.
	Data []string `json:"data"`
}

type DNSRecord struct {
	CommonFields commonFields

	// Raw is a not parsed DNS record.
	Raw json.RawMessage `json:"raw"`

	// ParseError is the error that occurred during parsing.
	ParseError error `json:"parseError"`
}

// DNSRecords is the struct where returned DNS records are stored.
type DNSRecords struct {
	// All is a slice of all parsed DNS records.
	All []DNSRecord

	// A is a slice of the parsed A records.
	A []ARecord

	// AAAA is a slice of the parsed AAAA records.
	AAAA []AAAARecord

	// NS is a slice of the parsed NS records.
	NS []NSRecord

	// MX is a slice of the parsed MX records.
	MX []MXRecord

	// MD is a slice of the parsed MD records.
	MD []MDRecord

	// MF is a slice of the parsed MF records.
	MF []MFRecord

	// MB is a slice of the parsed MF records.
	MB []MBRecord

	// SOA is a slice of the parsed SOA records.
	SOA []SOARecord

	// TXT is a slice of the parsed TXT records.
	TXT []TXTRecord

	// CAA is a slice of the parsed CAA records.
	CAA []CAARecord

	// CNAME is a slice of the parsed CNAME records.
	CNAME []CNAMERecord

	// DNAME is a slice of the parsed CNAME records.
	DNAME []DNAMERecord

	// DNSKEY is a slice of the parsed DNSKEY records.
	DNSKEY []DNSKEYRecord

	// NSEC3PARAM is a slice of the parsed NSEC3PARAM records.
	NSEC3PARAM []NSEC3PARAMRecord

	// NSEC is a slice of the parsed NSEC records.
	NSEC []NSECRecord

	// DS is a slice of the parsed DS records.
	DS []DSRecord

	// PTR is a slice of the parsed PTR records.
	PTR []PTRRecord

	// SRV is a slice of the parsed SRV records.
	SRV []SRVRecord

	// LOC is a slice of the parsed LOC records.
	LOC []LOCRecord

	// NAPTR is a slice of the parsed NAPTR records.
	NAPTR []NAPTRRecord

	// HINFO is a slice of the parsed HINFO records.
	HINFO []HINFORecord

	// RP is a slice of the parsed RP records.
	RP []RPRecord

	// DLV is a slice of the parsed DLV records.
	DLV []DLVRecord

	// SSHFP is a slice of the parsed SSHFP records.
	SSHFP []SSHFPRecord

	// TLSA is a slice of the parsed TLSA records.
	TLSA []TLSARecord

	// DHCID is a slice of the parsed DHCID records.
	DHCID []DHCIDRecord

	// NSAP is a slice of the parsed NSAP records.
	NSAP []NSAPRecord

	// NULL is a slice of the parsed NULL records.
	NULL []NULLRecord
}

// UnmarshalJSON decodes DNS records and returns them as a DNSRecords struct.
func (r *DNSRecords) UnmarshalJSON(data []byte) error {
	// this just splits up the JSON array into the raw JSON for each object
	var raw []json.RawMessage
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}

	for _, record := range raw {
		r.All = append(r.All, r.parseRecord(record))
	}
	return nil
}

func (r *DNSRecords) parseRecord(record json.RawMessage) DNSRecord {
	var obj struct {
		commonFields
	}

	if err := json.Unmarshal(record, &obj); err != nil {
		return DNSRecord{
			CommonFields: commonFields{},
			Raw:          record,
			ParseError:   err,
		}
	}

	dnsRecord := DNSRecord{
		CommonFields: obj.commonFields,
		Raw:          record,
		ParseError:   nil,
	}

	// unmarshal again into the correct type
	actual := actualDNSType(obj.DNSType)
	if actual == nil {
		dnsRecord.ParseError = ErrUnsupportedDNSType
		return dnsRecord
	}

	if err := json.Unmarshal(record, actual); err != nil {
		dnsRecord.ParseError = err
		return dnsRecord
	}

	switch obj.DNSType {
	case "A":
		r.A = append(r.A, *actual.(*ARecord))
	case "AAAA":
		r.AAAA = append(r.AAAA, *actual.(*AAAARecord))
	case "NS":
		r.NS = append(r.NS, *actual.(*NSRecord))
	case "MX":
		r.MX = append(r.MX, *actual.(*MXRecord))
	case "MD":
		r.MD = append(r.MD, *actual.(*MDRecord))
	case "MF":
		r.MF = append(r.MF, *actual.(*MFRecord))
	case "MB":
		r.MB = append(r.MB, *actual.(*MBRecord))
	case "SOA":
		r.SOA = append(r.SOA, *actual.(*SOARecord))
	case "TXT":
		r.TXT = append(r.TXT, *actual.(*TXTRecord))
	case "CAA":
		r.CAA = append(r.CAA, *actual.(*CAARecord))
	case "CNAME":
		r.CNAME = append(r.CNAME, *actual.(*CNAMERecord))
	case "DNAME":
		r.DNAME = append(r.DNAME, *actual.(*DNAMERecord))
	case "DNSKEY":
		r.DNSKEY = append(r.DNSKEY, *actual.(*DNSKEYRecord))
	case "NSEC":
		r.NSEC = append(r.NSEC, *actual.(*NSECRecord))
	case "NSEC3PARAM":
		r.NSEC3PARAM = append(r.NSEC3PARAM, *actual.(*NSEC3PARAMRecord))
	case "DS":
		r.DS = append(r.DS, *actual.(*DSRecord))
	case "PTR":
		r.PTR = append(r.PTR, *actual.(*PTRRecord))
	case "SRV":
		r.SRV = append(r.SRV, *actual.(*SRVRecord))
	case "LOC":
		r.LOC = append(r.LOC, *actual.(*LOCRecord))
	case "NAPTR":
		r.NAPTR = append(r.NAPTR, *actual.(*NAPTRRecord))
	case "HINFO":
		r.HINFO = append(r.HINFO, *actual.(*HINFORecord))
	case "RP":
		r.RP = append(r.RP, *actual.(*RPRecord))
	case "DLV":
		r.DLV = append(r.DLV, *actual.(*DLVRecord))
	case "SSHFP":
		r.SSHFP = append(r.SSHFP, *actual.(*SSHFPRecord))
	case "DHCID":
		r.DHCID = append(r.DHCID, *actual.(*DHCIDRecord))
	case "TLSA":
		r.TLSA = append(r.TLSA, *actual.(*TLSARecord))
	case "NSAP":
		r.NSAP = append(r.NSAP, *actual.(*NSAPRecord))
	case "NULL":
		r.NULL = append(r.NULL, *actual.(*NULLRecord))
	}

	return dnsRecord
}

// MarshalJSON encodes DNSRecords.
func (r *DNSRecords) MarshalJSON() ([]byte, error) {
	if len(r.All) == 0 {
		return []byte(`[]`), nil
	}
	result, err := json.Marshal(r.All)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func actualDNSType(dnsType string) interface{} {
	switch dnsType {
	case "A":
		return &ARecord{}
	case "AAAA":
		return &AAAARecord{}
	case "NS":
		return &NSRecord{}
	case "MX":
		return &MXRecord{}
	case "MD":
		return &MDRecord{}
	case "MF":
		return &MFRecord{}
	case "MB":
		return &MBRecord{}
	case "SOA":
		return &SOARecord{}
	case "TXT":
		return &TXTRecord{}
	case "CAA":
		return &CAARecord{}
	case "CNAME":
		return &CNAMERecord{}
	case "DNAME":
		return &DNAMERecord{}
	case "DNSKEY":
		return &DNSKEYRecord{}
	case "NSEC3PARAM":
		return &NSEC3PARAMRecord{}
	case "NSEC":
		return &NSECRecord{}
	case "DS":
		return &DSRecord{}
	case "PTR":
		return &PTRRecord{}
	case "SRV":
		return &SRVRecord{}
	case "LOC":
		return &LOCRecord{}
	case "NAPTR":
		return &NAPTRRecord{}
	case "HINFO":
		return &HINFORecord{}
	case "RP":
		return &RPRecord{}
	case "DLV":
		return &DLVRecord{}
	case "SSHFP":
		return &SSHFPRecord{}
	case "DHCID":
		return &DHCIDRecord{}
	case "TLSA":
		return &TLSARecord{}
	case "NSAP":
		return &NSAPRecord{}
	case "NULL":
		return &NULLRecord{}
	}
	return nil
}

// Audit is a part of the DNS Lookup API response
// It represents dates when Whois record was collected and updated in our database
type Audit struct {
	// CreatedDate is the date the DNS records are collected on whoisxmlapi.com
	CreatedDate Time `json:"createdDate"`

	// UpdatedDate is the date the DNS records updated on whoisxmlapi.com
	UpdatedDate Time `json:"updatedDate"`
}

// DNSLookupResponse is a response of DNS Lookup API.
type DNSLookupResponse struct {
	// DomainName is a domain name.
	DomainName string `json:"domainName"`

	// Types are codes of the requested DNS record types.
	Types []int `json:"types"`

	// DNSTypes is the comma-separated list of DNS record types.
	DNSTypes string `json:"dnsTypes"`

	// Audit is a part of the DNS Lookup API response
	// It represents dates when Whois record was collected and updated in our database
	Audit Audit `json:"audit"`

	// DNSRecords is the struct where returned DNS records are stored.
	DNSRecords DNSRecords `json:"dnsRecords"`
}

// ErrorMessage is an error message.
type ErrorMessage struct {
	Code    string `json:"errorCode"`
	Message string `json:"msg"`
}

// Error returns error message as a string.
func (e *ErrorMessage) Error() string {
	return fmt.Sprintf("API error: [%s] %s", e.Code, e.Message)
}
