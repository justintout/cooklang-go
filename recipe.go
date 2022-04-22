package cooklang

import (
	"encoding/json"
	"strings"
)

type Serving struct {
	S   string
	N   int
	Idx int
}

type Servings []Serving

type Ingredient struct {
	Name string
	Quantity

	raw      string
	step     *Step
	recipe   *Recipe
	category *Category
}

func NewIngredient(source string, recipe *Recipe, step *Step, category *Category) Ingredient {
	i := Ingredient{
		raw:      source,
		recipe:   recipe,
		step:     step,
		category: category,
	}
	source = strings.TrimPrefix(source, "@")
	if !strings.HasSuffix(source, "}") {
		source += "{}"
	}
	qs := strings.IndexRune(source, '{')
	i.Quantity = parseQuantity(source[qs:])
	i.Name = source[:qs]
	return i
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
		Quantity: i.Quantity.S,
		Units:    i.Quantity.Units,
	}

	return json.Marshal(ii)
}

type Cookware struct {
	Type     string `json:"type"`
	Name     string `json:"name"`
	Quantity string `json:"quantity"`
	id       string
	step     *Step
}

type Timer struct {
	Name string  `json:"name"`
	Time float32 `json:"`

	id   string
	step *Step
}

type Step struct {
	Number int
	Line   string

	Ingredients []Ingredient
	Cookware    []Cookware
	Timers      []Timer

	id string
}

type Recipe struct {
	Metadata    map[string]string
	Servings    Servings // duplicated in Metadata?
	Ingredients []Ingredient
	Cookware    []Cookware
	Steps       []*Step
	timers      map[string]Timer
}
