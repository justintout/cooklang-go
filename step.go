package cooklang

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

// Step represents a step of a recipe
type Step struct {
	Number int `json:"stepNumber"`

	// Ingredients are the ingredients used in this step
	Ingredients []*Ingredient `json:"ingredients"`
	// Cookware is the cookware used in this step
	Cookware []*Cookware `json:"cookware"`
	// Timers are the timers used in this step
	Timers []*Timer `json:"timers"`
	// Text is the text used in this step
	Text []*Text `json:"text"`

	// DirectionItems are the direction items in this step
	DirectionItems []DirectionItemer
	order          []json.Marshaler
	pos            int
	recipe         *Recipe
}

func (s Step) String() string {
	return fmt.Sprintf("#%d: I(%d): %+v, C(%d): %+v, T(%d): %+v, TXT(%d): %+v",
		s.Number, len(s.Ingredients), s.Ingredients, len(s.Cookware), s.Cookware, len(s.Timers), s.Timers, len(s.Text), s.Text)
}

// Zero returns whether the Step is equal to the Step zero value
func (s Step) Zero() bool {
	return len(s.DirectionItems) == 0 && len(s.order) == 0 && len(s.Text) == 0
}

// AddIngredient adds a new ingredient to the step
func (s *Step) AddIngredient(i *Ingredient) {
	i.stepPos = s.pos
	s.Ingredients = append(s.Ingredients, i)
	s.order = append(s.order, i)
	s.DirectionItems = append(s.DirectionItems, i)
	s.pos++
}

// AddCookware adds a new piece of cookware to the step
func (s *Step) AddCookware(c *Cookware) {
	c.stepPos = s.pos
	s.Cookware = append(s.Cookware, c)
	s.DirectionItems = append(s.DirectionItems, c)
	s.order = append(s.order, c)
	s.pos++
}

// AddTimer adds a new timer to the step
func (s *Step) AddTimer(t *Timer) {
	t.stepPos = s.pos
	s.Timers = append(s.Timers, t)
	s.DirectionItems = append(s.DirectionItems, t)
	s.order = append(s.order, t)
	s.pos++
}

// AddText adds new text to the step
func (s *Step) AddText(t *Text) {
	t.stepPos = s.pos
	s.Text = append(s.Text, t)
	s.DirectionItems = append(s.DirectionItems, t)
	s.order = append(s.order, t)
	s.pos++
}

// MarshalJSON implements json.Marshaler for Step
func (s Step) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.DirectionItems)
}

// MarshalYAML impelements yaml.Marshaler for the step, outputting Cooklang canonical test format
func (s Step) MarshalYAML() (interface{}, error) {
	var dd []DirectionItem
	for _, d := range s.DirectionItems {
		dd = append(dd, d.DirectionItem())
	}
	b, err := yaml.Marshal(dd)
	return string(b), err
}
