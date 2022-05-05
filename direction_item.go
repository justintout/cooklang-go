package cooklang

type DirectionItemer interface {
	DirectionItem() DirectionItem
}

type DirectionItem struct {
	Type     string
	Name     string `yaml:"name,omitempty"`
	Quantity string `yaml:"quantity,omitempty"`
	Units    string `yaml:"units,omitempty"`
	Value    string `yaml:"value,omitempty"`
}
