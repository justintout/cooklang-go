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
			"@whole milk{1 1/2%cup}",
			Ingredient{Name: "whole milk"},
		},
	}
	for _, tt := range tests {
		i := NewIngredient(tt.source, nil, nil, nil)
		if i.Name != tt.expected.Name {
			t.Errorf("wrong name: got: %v, want: %v", i.Name, tt.expected.Name)
		}
	}
}
