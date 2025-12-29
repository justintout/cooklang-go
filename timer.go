package cooklang

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Timer represents a timer used in a recipe
type Timer struct {
	Name string `json:"name"`
	Quantity

	raw     string
	stepPos int
}

// NewTimer creates a new Timer from a timer definition
func NewTimer(source string) *Timer {
	t := Timer{raw: source}
	ns := strings.IndexRune(source, '~') + 1
	qs := strings.IndexRune(source, '{')
	t.Name = source[ns:qs]
	var err error
	if t.Quantity, err = strictParseQuantity(source[qs:]); err != nil {
		// TODO: how to handle a parse error here?
		// preferably the lexer shouldn't emit as a timer so
		// maybe unnecessary to do much
		fmt.Printf("invalid quantity for timer %q: %v\n", source, err)
	}
	return &t
}

// String implements Stringer for Timer
func (t Timer) String() string {
	return t.raw
}

// DirectionItem creates a new direction item from the Timer
func (t Timer) DirectionItem() DirectionItem {
	return DirectionItem{
		Type:     "timer",
		Name:     t.Name,
		Quantity: t.S,
		Units:    t.Units,
	}
}

// MarshalJSON implements json.Marshaler for Timer
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
