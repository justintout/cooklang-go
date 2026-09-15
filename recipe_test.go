package cooklang

import (
	"testing"
)

// A symbol at the start of a line introduces its construct wherever the line
// falls in the file, not only on the first line, and a file without a
// trailing newline keeps its last step.
func TestSymbolsAtLineStart(t *testing.T) {
	tests := []struct {
		name   string
		source string
		check  func(Recipe) bool
	}{
		{
			"ingredient",
			"some text\n@salt{1%tsp}\n",
			func(r Recipe) bool { return len(r.Ingredients["salt"]) == 1 },
		},
		{
			"cookware",
			"some text\n#bowl{2}\n",
			func(r Recipe) bool { return len(r.Cookware["bowl"]) == 1 },
		},
		{
			"timer",
			"some text\n~{5%minutes}\n",
			func(r Recipe) bool { return len(r.Timers) == 1 },
		},
		{
			"metadata",
			"some text\n>> sourced: babooshka\n",
			func(r Recipe) bool { return r.Metadata["sourced"] == "babooshka" },
		},
		{
			"line comment",
			"some text\n-- a comment\n",
			func(r Recipe) bool { return len(r.Steps) == 1 },
		},
		{
			"block comment",
			"some text\n[- a comment -]\n",
			func(r Recipe) bool { return len(r.Steps) == 1 },
		},
		{
			"last step without a trailing newline",
			"first\n\nAdd @salt.",
			func(r Recipe) bool { return len(r.Steps) == 2 },
		},
	}
	for _, tt := range tests {
		if !tt.check(MustParse(tt.source)) {
			t.Errorf("%s: %q did not parse as expected: %+v", tt.name, tt.source, MustParse(tt.source))
		}
	}
}

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
			Cookware{Name: "frying pan", Quantity: Quantity{N: 1, S: "1"}},
		},
		{
			"#Oven",
			Cookware{Name: "Oven", Quantity: Quantity{N: 1, S: "1"}},
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
