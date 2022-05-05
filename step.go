package cooklang

import "encoding/json"

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
	return json.Marshal(s.order)
}
