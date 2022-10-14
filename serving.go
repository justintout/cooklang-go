package cooklang

// Serving is a recipe serving
type Serving struct {
	S   string
	N   int
	Idx int
}

// Servings is a list of recipe servings
type Servings []Serving
