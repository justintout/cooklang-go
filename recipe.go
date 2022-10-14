package cooklang

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// Recipe is a Cooklang recipe
type Recipe struct {
	Name        string    `yaml:"-"`
	Metadata    *Metadata `json:"metadata"`
	Servings    Servings  `yaml:"-"` // duplicated in Metadata?
	Steps       []*Step   `json:"steps" yaml:"steps"`
	Ingredients map[string][]string
	// Ingredients []*Ingredient `json:"ingredients"`
	// Cookware    []*Cookware   `json:"cookware"`
	// Timers      []*Timer      `json:"timers"`
	Cookware map[string][]string `json:"cookware"`
	Timers   map[string]string   `json:"timers"`

	filename string
}

// NewRecipe creates a new, empty recipe with the given name
func NewRecipe(name string) Recipe {
	return Recipe{
		Name:        name,
		Metadata:    &Metadata{},
		Steps:       make([]*Step, 0),
		Ingredients: make(map[string][]string),
		Cookware:    make(map[string][]string),
		Timers:      make(map[string]string),
	}
}

// AddStep appends a new step to the recipe
func (r *Recipe) AddStep(s *Step) {
	s.Number = len(r.Steps)
	r.Steps = append(r.Steps, s)
	// TODO: add em up. need conversions...
	for _, i := range s.Ingredients {
		if ri, ok := r.Ingredients[i.Name]; ok {
			r.Ingredients[i.Name] = append(ri, i.Quantity.String())
			continue
		}
		r.Ingredients[i.Name] = []string{i.Quantity.String()}
	}
	for _, c := range s.Cookware {
		if rc, ok := r.Cookware[c.Name]; ok {
			r.Cookware[c.Name] = append(rc, c.Quantity.String())
			continue
		}
		r.Cookware[c.Name] = []string{c.Quantity.String()}
	}
	for i, t := range s.Timers {
		if t.Name == "" {
			t.Name = fmt.Sprintf("timer:%d:%d", s.Number, i)
		}
		r.Timers[t.Name] = t.Quantity.String()
	}
}

// MarshalYAML implements yaml.Marshaler for Recipe, returning the
// Cooklang canonical test format for the recipe
func (r *Recipe) MarshalYAML() (interface{}, error) {
	cr := &canonicalRecipe{
		Steps:    make([]canonicalStep, 0, len(r.Steps)),
		Metadata: *r.Metadata,
	}
	for _, s := range r.Steps {
		cs := make(canonicalStep, 0, len(s.order))
		for _, di := range s.order {
			switch di := di.(type) {
			case *Text:
				cs = append(cs, newCanonicalText(*di))
			case *Ingredient:
				cs = append(cs, newCanonicalIngredient(*di))
			case *Cookware:
				cs = append(cs, newCanonicalCookware(*di))
			case *Timer:
				cs = append(cs, newCanonicalTimer(*di))
			default:
				return nil, fmt.Errorf("unexpected item type: %v %T", di, di)
			}
		}
		fmt.Printf("- %+v\n", cs)
		cr.Steps = append(cr.Steps, cs)
	}
	y, err := yaml.Marshal(cr)
	return string(y), err
}

type canonicalRecipe struct {
	Steps    []canonicalStep   `yaml:"steps"`
	Metadata map[string]string `yaml:"metadata"`
}

type canonicalStep []directionItemTyper

type canonicalDirectionItem struct {
	Typ      string `yaml:"type"`
	Name     string `yaml:"name"`
	Quantity string `yaml:"quantity"`
	Units    string `yaml:"units"`
}

func (cdi canonicalDirectionItem) Type() string {
	return cdi.Typ
}

func newCanonicalIngredient(i Ingredient) canonicalDirectionItem {
	return canonicalDirectionItem{
		Typ:      "ingredient",
		Name:     i.Name,
		Quantity: i.Quantity.S,
		Units:    i.Quantity.Units,
	}
}

func newCanonicalCookware(c Cookware) canonicalDirectionItem {
	return canonicalDirectionItem{
		Typ:      "cookware",
		Name:     c.Name,
		Quantity: c.Quantity.S,
	}
}

func newCanonicalTimer(t Timer) canonicalDirectionItem {
	return canonicalDirectionItem{
		Typ:      "timer",
		Name:     t.Name,
		Quantity: t.Quantity.S,
		Units:    t.Quantity.Units,
	}
}

type canonicalText struct {
	Typ   string `yaml:"type"`
	Value string `yaml:"value"`
}

func (ct canonicalText) Type() string {
	return "text"
}

func newCanonicalText(t Text) canonicalText {
	return canonicalText{
		Typ:   "text",
		Value: t.Value,
	}
}

type directionItemTyper interface {
	Type() string
}
