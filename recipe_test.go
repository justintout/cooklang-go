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
		{
			"section header at line start",
			"first\n== Filling ==\nAdd @salt.",
			func(r Recipe) bool {
				return len(r.Sections) == 1 && r.Sections[0].Name == "Filling" && len(r.Steps) == 2
			},
		},
		{
			"note at line start",
			"first\n\n> Keep warm.\n\nAdd @salt.",
			func(r Recipe) bool {
				return len(r.Notes) == 1 && r.Notes[0].Value == "Keep warm." && len(r.Steps) == 2
			},
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
			Ingredient{Name: "whole milk", Quantity: Quantity{S: "1 1/2", N: 1.5, Units: "cup"}},
		},
		{
			"@?chives",
			Ingredient{Name: "chives", Optional: true, Quantity: Quantity{S: "some", N: -1}},
		},
		{
			"@?chilli flakes{1%pinch}",
			Ingredient{Name: "chilli flakes", Optional: true, Quantity: Quantity{S: "1", N: 1, Units: "pinch"}},
		},
		{
			"@onion{1}(peeled and finely chopped)",
			Ingredient{Name: "onion", Preparation: "peeled and finely chopped", Quantity: Quantity{S: "1", N: 1}},
		},
		{
			"@garlic{}(minced)",
			Ingredient{Name: "garlic", Preparation: "minced", Quantity: Quantity{S: "some", N: -1}},
		},
		{
			"@?./sauces/tomato sauce{200%g}",
			Ingredient{Name: "./sauces/tomato sauce", Optional: true, Quantity: Quantity{S: "200", N: 200, Units: "g"}},
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
		if i.Optional != tt.expected.Optional {
			t.Errorf("wrong optional flag for %q: got: %v, want: %v", tt.source, i.Optional, tt.expected.Optional)
		}
		if i.Preparation != tt.expected.Preparation {
			t.Errorf("wrong preparation for %q: got: %q, want: %q", tt.source, i.Preparation, tt.expected.Preparation)
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
		{
			"#?splatter guard{}",
			Cookware{Name: "splatter guard", Optional: true, Quantity: Quantity{N: 1, S: "1"}},
		},
		{
			"#?lid",
			Cookware{Name: "lid", Optional: true, Quantity: Quantity{N: 1, S: "1"}},
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
		if c.Optional != tt.expected.Optional {
			t.Errorf("wrong optional flag for %q: got: %v, want: %v", tt.source, c.Optional, tt.expected.Optional)
		}
	}
}
