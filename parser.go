package cooklang

import (
	"fmt"
	"os"
)

func MustParseFile(path string) Recipe {
	r, err := ParseFile(path)
	if err != nil {
		panic(err)
	}
	return r
}

func ParseFile(path string) (Recipe, error) {
	c, err := os.ReadFile(path)
	if err != nil {
		return Recipe{}, fmt.Errorf("failed to parse %q: %v", path, err)
	}
	return parse(string(c))
}

func MustParse(input string) Recipe {
	r, err := Parse(input)
	if err != nil {
		panic(err)
	}
	return r
}

func Parse(input string) (Recipe, error) {
	return parse(input)
}

func parse(input string) (Recipe, error) {
	_, items := lex("recipe", input)
	recipe := NewRecipe("recipe")
	step := &Step{}
	for item := range items {
		switch item.typ {
		case itemMetadata:
			recipe.Metadata.Add(item.val)
		case itemComment:
		case itemText:
			step.AddText(NewText(item.val))
		case itemIngredient:
			step.AddIngredient(NewIngredient(item.val))
		case itemCookware:
			step.AddCookware(NewCookware(item.val))
		case itemTimer:
			step.AddTimer(NewTimer(item.val))
		case itemStep:
			// BAD
			if !step.Zero() {
				recipe.AddStep(step)
			}
			step = &Step{}
		}
	}

	return recipe, nil
}
