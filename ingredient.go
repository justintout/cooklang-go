package cooklang

import (
	"encoding/json"
	"strings"
)

type Ingredient struct {
	Name string
	Quantity

	raw     string
	stepPos int
}

func NewIngredient(source string) *Ingredient {
	i := Ingredient{raw: source}
	source = strings.TrimPrefix(source, "@")
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

func (i Ingredient) DirectionItem() DirectionItem {
	return DirectionItem{
		Type:     "ingredient",
		Name:     i.Name,
		Quantity: i.S,
		Units:    i.Units,
	}
}

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
