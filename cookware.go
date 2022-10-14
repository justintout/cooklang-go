package cooklang

import (
	"encoding/json"
	"strings"
)

// Cookware represents a tool used for a recipe
type Cookware struct {
	Name string `json:"name"`
	Quantity

	raw     string
	stepPos int
}

// NewCookware creates a cookware from a cookware definition. Default quantity is 1.
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

// String returns the cookware's raw string
func (c Cookware) String() string {
	return c.raw
}

// MarshalJSON implements json.Marshaler
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

// DirectionItem creates a new DirectionItem from the Cookware
func (c Cookware) DirectionItem() DirectionItem {
	return DirectionItem{
		Type:     "cookware",
		Name:     c.Name,
		Quantity: c.S,
	}
}
