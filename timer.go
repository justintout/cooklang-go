package cooklang

import (
	"encoding/json"
	"strings"
)

// Timer represents a timer used in a recipe
type Timer struct {
	Name string `json:"name"`
	Quantity

	raw     string
	stepPos int
}

// NewTimer creates a new Timer from a timer definition. A timer without a
// duration, such as "~rest", has an empty quantity.
func NewTimer(source string) *Timer {
	t := Timer{raw: source}
	ns := strings.IndexRune(source, '~') + 1
	qs := strings.IndexRune(source, '{')
	if qs == -1 {
		t.Name = strings.TrimSpace(source[ns:])
		t.Quantity = Quantity{N: -1}
		return &t
	}
	t.Name = strings.TrimSpace(source[ns:qs])
	if q, err := strictParseQuantity(source[qs:]); err == nil {
		t.Quantity = q
	} else {
		// A malformed duration is kept as written rather than dropped, and
		// a library must not print.
		t.Quantity = parseQuantity(source[qs:], "", -1)
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
		Quantity: t.Quantity.Canonical(),
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
