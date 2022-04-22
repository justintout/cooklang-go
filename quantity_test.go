package cooklang

import (
	"testing"
)

// func TestQuantityMarshalJSON(t *testing.T) {
// 	tests := []struct {
// 		q        Quantity
// 		expected string
// 	}{
// 		{
// 			Quantity{S: "some"},
// 			`{"quantity":"some","units":""}`,
// 		},
// 		{
// 			Quantity{S: "1 1/2", N: 1.5, Units: ""},
// 			`{"quantity":"1 1/2","units":""}`,
// 		},
// 		{
// 			Quantity{S: "4.32", N: 4.32, Units: "kg"},
// 			`{"quantity":"4.32","units":"kg"}`,
// 		},
// 	}
// 	for _, tt := range tests {
// 		b, err := json.Marshal(tt.q)
// 		if err != nil {
// 			t.Error(err)
// 		}
// 		if string(b) != tt.expected {
// 			t.Errorf("got: %s, want: %s", string(b), tt.expected)
// 		}
// 	}
// }

func TestParseQuantity(t *testing.T) {
	tests := []struct {
		source string

		n    float32
		s    string
		unit string
	}{
		{"{1}", 1, "1", ""},
		{"{3.5%cups}", 3.5, "3.5", "cups"},
		{"{1 1/2%oz.}", 1.5, "1 1/2", "oz."},
		{"{a few%sprigs{", -1, "a few", "sprigs"},
		{"{%pounds}", -1, "some", "pounds"},
		{"{10%minutes}", 10, "10", "minutes"},
	}
	for _, tt := range tests {
		q := parseQuantity(tt.source)
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
