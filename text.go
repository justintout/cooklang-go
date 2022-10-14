package cooklang

import "encoding/json"

// Text is text in a recipe
type Text struct {
	Type  string `json:"type"`
	Value string `json:"value"`

	stepPos int
}

// NewText returns new recipe text
func NewText(value string) *Text {
	return &Text{
		Type:  "text",
		Value: value,
	}
}

// String implements Stringer for Text
func (t Text) String() string {
	return t.Value
}

// DirectionItem outputs a new direction item for the Text
func (t Text) DirectionItem() DirectionItem {
	return DirectionItem{
		Type:  "text",
		Value: t.Value,
	}
}

// MarshalJSON implements json.Marshaler
func (t Text) MarshalJSON() ([]byte, error) {
	m := map[string]string{
		"type":  "text",
		"value": t.Value,
	}
	return json.Marshal(m)
}
