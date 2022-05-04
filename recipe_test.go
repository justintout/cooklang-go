package cooklang

import (
	"testing"
)

func TestNewIngredient(t *testing.T) {
	tests := []struct {
		source   string
		expected Ingredient
	}{
		{
			"@eggs",
			Ingredient{Name: "eggs", Quantity: Quantity{S: "some", N: -1}},
		},
		{
			"@whole milk{1 1/2%cup}",
			Ingredient{Name: "whole milk", Quantity: Quantity{S: "1 1/2", N: 1.5}},
		},
	}
	for _, tt := range tests {
		i := NewIngredient(tt.source)
		if i.Name != tt.expected.Name {
			t.Errorf("wrong name: got: %v, want: %v", i.Name, tt.expected.Name)
		}
		if i.S != tt.expected.S {
			t.Errorf("wrong string quantity: got: %v, want: %v", i.S, tt.expected.S)
		}
		if i.N != tt.expected.N {
			t.Errorf("wrong number quantity: got: %v, want: %v", i.N, tt.expected.N)
		}
	}
}

func TestNewCookware(t *testing.T) {
	tests := []struct {
		source   string
		expected Cookware
	}{
		{
			"#frying pan{}",
			Cookware{Name: "frying pan", Quantity: Quantity{N: 1, S: ""}},
		},
		{
			"#Oven",
			Cookware{Name: "Oven", Quantity: Quantity{N: 1, S: ""}},
		},
		{
			"#bowls{3}",
			Cookware{Name: "bowls", Quantity: Quantity{N: 3, S: "3"}},
		},
	}
	for _, tt := range tests {
		c := NewCookware(tt.source)
		if c.Name != tt.expected.Name {
			t.Errorf("wrong name: got: %v, want: %v", c.Name, tt.expected.Name)
		}
		if c.S != tt.expected.S {
			t.Errorf("wrong string quantity: got: %v, want: %v", c.S, tt.expected.S)
		}
		if c.N != tt.expected.N {
			t.Errorf("wrong number quantity: got: %v, want: %v", c.N, tt.expected.N)
		}
	}
}
