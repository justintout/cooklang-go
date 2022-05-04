package cooklang

import (
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

type Serving struct {
	S   string
	N   int
	Idx int
}

type Servings []Serving

type Text struct {
	Type  string `json:"type"`
	Value string `json:"value"`

	stepPos int
}

func NewText(value string) *Text {
	return &Text{
		Type:  "text",
		Value: value,
	}
}

func (t Text) String() string {
	return t.Value
}

func (t Text) MarshalJSON() ([]byte, error) {
	m := map[string]string{
		"type":  "text",
		"value": t.Value,
	}
	return json.Marshal(m)
}

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

type Cookware struct {
	Name string `json:"name"`
	Quantity

	raw     string
	stepPos int
}

func NewCookware(source string) *Cookware {
	c := Cookware{raw: source}
	ns := strings.IndexRune(source, '#') + 1
	qs := strings.IndexRune(source, '{')
	if qs == -1 {
		c.Name = source[ns:]
		c.Quantity = Quantity{N: 1}
		return &c
	}
	c.Name = source[ns:qs]
	c.Quantity = parseQuantity(source[qs:], "", 1)
	return &c
}

func (c Cookware) String() string {
	return c.raw
}

func (c Cookware) MarshalJSON() ([]byte, error) {
	cc := struct {
		Type     string `json:"type"`
		Name     string `json:"name"`
		Quantity string `json:"quantity"`
	}{
		Type:     "cookware",
		Name:     c.Name,
		Quantity: c.S,
	}
	return json.Marshal(cc)
}

type Timer struct {
	Name string `json:"name"`
	Quantity

	raw     string
	stepPos int
}

func NewTimer(source string) *Timer {
	t := Timer{raw: source}
	ns := strings.IndexRune(source, '~') + 1
	qs := strings.IndexRune(source, '{')
	t.Name = source[ns:qs]
	var err error
	if t.Quantity, err = strictParseQuantity(source[qs : len(source)-2]); err != nil {
		// TODO: how to handle a parse error here?
		// preferably the lexer shouldn't emit as a timer so
		// maybe unnecessary to do much
		fmt.Printf("invalid quantity for timer %q: %v\n", source, err)
	}
	return &t
}

func (t Timer) String() string {
	return t.raw
}

func (t Timer) MarshalJSON() ([]byte, error) {
	tt := struct {
		Name     string `json:"name,omitempty"`
		Quantity string `json:"quantity"`
		Units    string `json:"units"`
	}{
		Name:     t.Name,
		Quantity: t.Quantity.S,
		Units:    t.Units,
	}
	return json.Marshal(tt)
}

type Step struct {
	Number int `json:"stepNumber"`

	Ingredients []*Ingredient `json:"ingredients"`
	Cookware    []*Cookware   `json:"cookware"`
	Timers      []*Timer      `json:"timers"`
	Text        []*Text       `json:"text"`

	order  []json.Marshaler
	pos    int
	recipe *Recipe
}

func (s *Step) AddIngredient(i *Ingredient) {
	i.stepPos = s.pos
	s.Ingredients = append(s.Ingredients, i)
	s.order = append(s.order, i)
	s.pos++
}

func (s *Step) AddCookware(c *Cookware) {
	c.stepPos = s.pos
	s.Cookware = append(s.Cookware, c)
	s.order = append(s.order, c)
	s.pos++
}

func (s *Step) AddTimer(t *Timer) {
	t.stepPos = s.pos
	s.Timers = append(s.Timers, t)
	s.order = append(s.order, t)
	s.pos++
}

func (s *Step) AddText(t *Text) {
	t.stepPos = s.pos
	s.Text = append(s.Text, t)
	s.order = append(s.order, t)
	s.pos++
}

func (s Step) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.order)
}

type Metadata map[string]string

func (m Metadata) Add(input string) {
	input = strings.TrimSpace(strings.TrimPrefix(input, ">>"))
	s := strings.SplitN(input, ":", 2)
	if len(s) == 1 {
		s = append(s, "")
	}
	s[0], s[1] = strings.TrimSpace(s[0]), strings.TrimSpace(s[1])
	m[s[0]] = s[1]
}

type Recipe struct {
	Name        string
	Metadata    *Metadata `json:"metadata"`
	Servings    Servings  // duplicated in Metadata?
	Steps       []*Step   `json:"steps"`
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
