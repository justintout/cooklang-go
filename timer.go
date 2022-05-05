package cooklang

import (
	"encoding/json"
	"fmt"
	"strings"
)

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
