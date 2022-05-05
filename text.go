package cooklang

import "encoding/json"

type Text struct {
	Type  string `json:"type"`
	Value string `json:"value"`

	stepPos int
}

func NewText(value string) *Text {
	return &Text{
		Type:  "text",
		Value: value,
	}
}

func (t Text) String() string {
	return t.Value
}

func (t Text) MarshalJSON() ([]byte, error) {
	m := map[string]string{
		"type":  "text",
		"value": t.Value,
	}
	return json.Marshal(m)
}
