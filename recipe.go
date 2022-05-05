package cooklang

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

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

func (r *Recipe) MarshalYAML() (interface{}, error) {
	cr := &CanonicalRecipe{
		Steps:    make([]CanonicalStep, 0, len(r.Steps)),
		Metadata: *r.Metadata,
	}
	for _, s := range r.Steps {
		cs := make(CanonicalStep, 0, len(s.order))
		for _, di := range s.order {
			switch di.(type) {
			case *Text:
				cs = append(cs, NewCanonicalText(*di.(*Text)))
			case *Ingredient:
				cs = append(cs, NewCanonicalIngredient(*di.(*Ingredient)))
			case *Cookware:
				cs = append(cs, NewCanonicalCookware(*di.(*Cookware)))
			case *Timer:
				cs = append(cs, NewCanonicalTimer(*di.(*Timer)))
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

type CanonicalRecipe struct {
	Steps    []CanonicalStep   `yaml:"steps"`
	Metadata map[string]string `yaml:"metadata"`
}

type CanonicalStep []DirectionItemTyper

type CanonicalDirectionItem struct {
	Typ      string `yaml:"type"`
	Name     string `yaml:"name"`
	Quantity string `yaml:"quantity"`
	Units    string `yaml:"units"`
}

func (cdi CanonicalDirectionItem) Type() string {
	return cdi.Typ
}

func NewCanonicalIngredient(i Ingredient) CanonicalDirectionItem {
	return CanonicalDirectionItem{
		Typ:      "ingredient",
		Name:     i.Name,
		Quantity: i.Quantity.S,
		Units:    i.Quantity.Units,
	}
}

func NewCanonicalCookware(c Cookware) CanonicalDirectionItem {
	return CanonicalDirectionItem{
		Typ:      "cookware",
		Name:     c.Name,
		Quantity: c.Quantity.S,
	}
}

func NewCanonicalTimer(t Timer) CanonicalDirectionItem {
	return CanonicalDirectionItem{
		Typ:      "timer",
		Name:     t.Name,
		Quantity: t.Quantity.S,
		Units:    t.Quantity.Units,
	}
}

type CanonicalText struct {
	Typ   string `yaml:"type"`
	Value string `yaml:"value"`
}

func (ct CanonicalText) Type() string {
	return "text"
}

func NewCanonicalText(t Text) CanonicalText {
	return CanonicalText{
		Typ:   "text",
		Value: t.Value,
	}
}

type DirectionItemTyper interface {
	Type() string
}
