package cooklang

import (
	"encoding/json"
	"strings"
)

// Ingredient represents an ingredient used in a recipe
type Ingredient struct {
	Name string
	Quantity
	Optional    bool
	Preparation string
	Note        string

	raw string
}

// NewIngredient creates a new Ingredient from an ingredient definition
func NewIngredient(source string) *Ingredient {
	i := Ingredient{raw: source}
	source = strings.TrimSpace(strings.TrimPrefix(source, leftIngredient))
	if strings.HasPrefix(source, optionalMarker) {
		i.Optional = true
		source = strings.TrimPrefix(source, optionalMarker)
	}
	if strings.HasSuffix(source, rightPreparation) {
		if qe := strings.IndexRune(source, '}'); qe >= 0 && qe+1 < len(source) && strings.HasPrefix(source[qe+1:], leftPreparation) {
			i.Preparation = strings.TrimSpace(source[qe+2 : len(source)-1])
			i.Note = i.Preparation
			source = source[:qe+1]
		} else if qe < 0 {
			if ps := strings.IndexRune(source, '('); ps >= 0 {
				i.Preparation = strings.TrimSpace(source[ps+1 : len(source)-1])
				i.Note = i.Preparation
				source = source[:ps]
			}
		}
	}
	if !strings.HasSuffix(source, "}") {
		source += "{}"
	}
	qs := strings.IndexRune(source, '{')
	i.Quantity = parseQuantity(source[qs:], "some", -1)
	i.Name = strings.TrimSpace(source[:qs])
	return &i
}

// IsReference reports whether the ingredient references another Cooklang
// recipe via a relative path ("./" or "../").
func (i Ingredient) IsReference() bool {
	return strings.HasPrefix(i.Name, "./") || strings.HasPrefix(i.Name, "../")
}

// ReferencePath returns the relative path for a recipe reference ingredient,
// or an empty string if the ingredient is not a recipe reference.
func (i Ingredient) ReferencePath() string {
	if !i.IsReference() {
		return ""
	}
	return i.Name
}

func (i Ingredient) String() string {
	return i.raw
}

// DirectionItem creates a new DirectionItem from the Ingredient
func (i Ingredient) DirectionItem() DirectionItem {
	return DirectionItem{
		Type:        "ingredient",
		Name:        i.Name,
		Quantity:    i.Quantity.Canonical(),
		Units:       i.Units,
		Optional:    i.Optional,
		Preparation: i.Preparation,
	}
}

// MarshalJSON implements json.Marshaler for Ingredient
func (i Ingredient) MarshalJSON() ([]byte, error) {
	ii := struct {
		Type        string `json:"type"`
		Name        string `json:"name"`
		Quantity    string `json:"quantity"`
		Units       string `json:"units"`
		Optional    bool   `json:"optional,omitempty"`
		Preparation string `json:"preparation,omitempty"`
	}{
		Type:        "ingredient",
		Name:        i.Name,
		Quantity:    i.S,
		Units:       i.Units,
		Optional:    i.Optional,
		Preparation: i.Preparation,
	}

	return json.Marshal(ii)
}
