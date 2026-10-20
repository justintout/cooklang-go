package cooklang

import (
	"encoding/json"
	"strings"
)

// Cookware represents a tool used for a recipe
type Cookware struct {
	Name string `json:"name"`
	Quantity
	Optional bool `json:"optional,omitempty"`

	raw string
}

// NewCookware creates a cookware from a cookware definition. Default quantity is 1.
func NewCookware(source string) *Cookware {
	c := Cookware{raw: source}
	rest := strings.TrimSpace(strings.TrimPrefix(source, leftCookware))
	if strings.HasPrefix(rest, optionalMarker) {
		c.Optional = true
		rest = strings.TrimPrefix(rest, optionalMarker)
	}
	qs := strings.IndexRune(rest, '{')
	if qs == -1 {
		c.Name = strings.TrimSpace(rest)
		c.Quantity = Quantity{N: 1, S: "1"}
		return &c
	}
	c.Name = strings.TrimSpace(rest[:qs])
	c.Quantity = parseQuantity(rest[qs:], "1", 1)
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
		Optional bool   `json:"optional,omitempty"`
	}{
		Type:     "cookware",
		Name:     c.Name,
		Quantity: c.S,
		Optional: c.Optional,
	}
	return json.Marshal(cc)
}

// DirectionItem creates a new DirectionItem from the Cookware
func (c Cookware) DirectionItem() DirectionItem {
	return DirectionItem{
		Type:     "cookware",
		Name:     c.Name,
		Quantity: c.Quantity.Canonical(),
		Optional: c.Optional,
	}
}
