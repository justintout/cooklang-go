package cooklang

import (
	"reflect"
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

func TestSumQuantities(t *testing.T) {
	tests := []struct {
		name string
		in   []Quantity
		want []string
	}{
		{
			"identical units are summed",
			[]Quantity{{N: 2, S: "2", Units: "cups"}, {N: 3, S: "3", Units: "cups"}},
			[]string{"5 cups"},
		},
		{
			"unitless numbers are summed",
			[]Quantity{{N: 2, S: "2"}, {N: 1, S: "1"}},
			[]string{"3"},
		},
		{
			"convertible units in a family collapse to the first unit",
			[]Quantity{{N: 1, S: "1", Units: "tbsp"}, {N: 3, S: "3", Units: "tsp"}},
			[]string{"2 tbsp"},
		},
		{
			"grams and kilograms combine",
			[]Quantity{{N: 500, S: "500", Units: "g"}, {N: 1, S: "1", Units: "kg"}},
			[]string{"1500 g"},
		},
		{
			"incompatible units stay separate, first unit seen leads",
			[]Quantity{{N: 2, S: "2", Units: "cups"}, {N: 100, S: "100", Units: "g"}},
			[]string{"2 cups", "100 g"},
		},
		{
			"non-numeric quantities are kept verbatim",
			[]Quantity{{N: -1, S: "some"}, {N: 2, S: "2", Units: "cups"}},
			[]string{"some", "2 cups"},
		},
		{
			"case and pluralization do not block summation",
			[]Quantity{{N: 1, S: "1", Units: "Cup"}, {N: 2, S: "2", Units: "cups"}},
			[]string{"3 Cup"},
		},
	}
	for _, tt := range tests {
		got := sumQuantities(tt.in)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%s: got: %v, want: %v", tt.name, got, tt.want)
		}
	}
}
