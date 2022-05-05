package cooklang

import (
	"encoding/json"
	"strings"
)

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

func (c Cookware) DirectionItem() DirectionItem {
	return DirectionItem{
		Type:     "cookware",
		Name:     c.Name,
		Quantity: c.S,
	}
}
