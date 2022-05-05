package cooklang

type DirectionItemer interface {
	DirectionItem() DirectionItem
}

type DirectionItem struct {
	Type     string
	Name     string `yaml:"omitempty"`
	Quantity string `yaml:"omitempty"`
	Units    string
	Value    string `yaml:"omitempty"`
}
