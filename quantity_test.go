package cooklang

import (
	"testing"
)

func TestParseQuantity(t *testing.T) {
	tests := []struct {
		source string
		ds     string
		dn     float32

		n    float32
		s    string
		unit string
	}{
		{"{}", "some", -1, -1, "some", ""},
		{"", "some", -1, -1, "some", ""},
		{"{3}", "some", -1, 3, "3", ""},
		{"{3.5%cups}", "some", -1, 3.5, "3.5", "cups"},
		{"{1 1/2%oz.}", "some", -1, 1.5, "1 1/2", "oz."},
		{"{a few%sprigs{", "some", -1, -1, "a few", "sprigs"},
		{"{%pounds}", "some", -1, -1, "some", "pounds"},
		{"{10%minutes}", "some", -1, 10, "10", "minutes"},
	}
	for _, tt := range tests {
		q := parseQuantity(tt.source, tt.ds, tt.dn)
		if q.N != tt.n {
			t.Errorf("numeric quantity for %q incorrect: got: %.2f, want: %.2f", tt.source, q.N, tt.n)
		}
		if q.S != tt.s {
			t.Errorf("string quantity for %q incorrect: got: %q, want: %q", tt.source, q.S, tt.s)
		}
		if q.Units != tt.unit {
			t.Errorf("unit for %q incorrect: got: %q, want: %q", tt.source, q.Units, tt.unit)
		}
	}
}
