package cooklang

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

type Step struct {
	Number int `json:"stepNumber"`

	Ingredients []*Ingredient `json:"ingredients"`
	Cookware    []*Cookware   `json:"cookware"`
	Timers      []*Timer      `json:"timers"`
	Text        []*Text       `json:"text"`

	DirectionItems []DirectionItemer
	order          []json.Marshaler
	pos            int
	recipe         *Recipe
}

func (s Step) String() string {
	return fmt.Sprintf("#%d: I(%d): %+v, C(%d): %+v, T(%d): %+v, TXT(%d): %+v",
		s.Number, len(s.Ingredients), s.Ingredients, len(s.Cookware), s.Cookware, len(s.Timers), s.Timers, len(s.Text), s.Text)
}

func (s Step) Zero() bool {
	return len(s.DirectionItems) == 0 && len(s.order) == 0 && len(s.Text) == 0
}

func (s *Step) AddIngredient(i *Ingredient) {
	i.stepPos = s.pos
	s.Ingredients = append(s.Ingredients, i)
	s.order = append(s.order, i)
	s.DirectionItems = append(s.DirectionItems, i)
	s.pos++
}

func (s *Step) AddCookware(c *Cookware) {
	c.stepPos = s.pos
	s.Cookware = append(s.Cookware, c)
	s.DirectionItems = append(s.DirectionItems, c)
	s.order = append(s.order, c)
	s.pos++
}

func (s *Step) AddTimer(t *Timer) {
	t.stepPos = s.pos
	s.Timers = append(s.Timers, t)
	s.DirectionItems = append(s.DirectionItems, t)
	s.order = append(s.order, t)
	s.pos++
}

func (s *Step) AddText(t *Text) {
	t.stepPos = s.pos
	s.Text = append(s.Text, t)
	s.DirectionItems = append(s.DirectionItems, t)
	s.order = append(s.order, t)
	s.pos++
}

func (s Step) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.DirectionItems)
}

func (s Step) MarshalYAML() (interface{}, error) {
	var dd []DirectionItem
	for _, d := range s.DirectionItems {
		dd = append(dd, d.DirectionItem())
	}
	b, err := yaml.Marshal(dd)
	return string(b), err
}
