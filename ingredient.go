package cooklang

import (
	"encoding/json"
	"strings"
)

// Ingredient represents an ingredient used in a recipe
type Ingredient struct {
	Name string
	Quantity

	raw     string
	stepPos int
}

// NewIngredient creates a new Ingredient from an ingredient definition
func NewIngredient(source string) *Ingredient {
	i := Ingredient{raw: source}
	source = strings.TrimSpace(strings.TrimPrefix(source, "@"))
	if !strings.HasSuffix(source, "}") {
		source += "{}"
	}
	qs := strings.IndexRune(source, '{')
	i.Quantity = parseQuantity(source[qs:], "some", -1)
	i.Name = source[:qs]
	return &i
}

func (i Ingredient) String() string {
	return i.raw
}

// DirectionItem creates a new DirectionItem from the Ingredient
func (i Ingredient) DirectionItem() DirectionItem {
	return DirectionItem{
		Type:     "ingredient",
		Name:     i.Name,
		Quantity: i.S,
		Units:    i.Units,
	}
}

// MarshalJSON implements json.Marshaler for Ingredient
func (i Ingredient) MarshalJSON() ([]byte, error) {
	ii := struct {
		Type     string `json:"type"`
		Name     string `json:"name"`
		Quantity string `json:"quantity"`
		Units    string `json:"units"`
	}{
		Type:     "ingredient",
		Name:     i.Name,
		Quantity: i.S,
		Units:    i.Units,
	}

	return json.Marshal(ii)
}
