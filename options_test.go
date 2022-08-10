package dnslookupapi

import (
	"net/url"
	"reflect"
	"testing"
)

// TestOptions tests the Options functions.
func TestOptions(t *testing.T) {
	tests := []struct {
		name   string
		values url.Values
		option Option
		want   string
	}{
		{
			name:   "output format",
			values: url.Values{},
			option: OptionOutputFormat("JSON"),
			want:   "outputFormat=JSON",
		},
		{
			name:   "type",
			values: url.Values{},
			option: OptionType("A,AAAA,NS,MX,CAA,SOA,TXT"),
			want:   "type=A%2CAAAA%2CNS%2CMX%2CCAA%2CSOA%2CTXT",
		},
		{
			name:   "callback",
			values: url.Values{},
			option: OptionCallback("func"),
			want:   "callback=func",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.option(tt.values)
			if got := tt.values.Encode(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Option() = %v, want %v", got, tt.want)
			}
		})
	}
}
