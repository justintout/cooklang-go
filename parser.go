package cooklang

import (
	"fmt"
	"os"
)

// MustParseFile calls ParseFile and panics on error
func MustParseFile(path string) Recipe {
	r, err := ParseFile(path)
	if err != nil {
		panic(err)
	}
	return r
}

// ParseFile parses the file at the given path as a Cooklang recipe
func ParseFile(path string) (Recipe, error) {
	c, err := os.ReadFile(path)
	if err != nil {
		return Recipe{}, fmt.Errorf("failed to parse %q: %v", path, err)
	}
	return parse(string(c))
}

// MustParse calls Parse and panics on error
func MustParse(input string) Recipe {
	r, err := Parse(input)
	if err != nil {
		panic(err)
	}
	return r
}

// Parse parses the input string as a Cooklang recipe
func Parse(input string) (Recipe, error) {
	return parse(input)
}

func parse(input string) (Recipe, error) {
	_, items := lex("recipe", input)
	recipe := NewRecipe("recipe")
	step := &Step{}
	// A step is a paragraph: a blank line ends it, while a single line break
	// inside it renders as a space and keeps the step going.
	lineBreak := false

	flush := func() {
		if !step.Zero() {
			recipe.AddStep(step)
		}
		step = &Step{}
		lineBreak = false
	}

	for item := range items {
		switch item.typ {
		case itemMetadata:
			recipe.Metadata.Add(item.val)
		case itemComment:
		case itemText:
			text := NewText(item.val)
			if lineBreak {
				text.Value = " " + text.Value
			}
			step.AddJoinedText(text)
			lineBreak = false
		case itemIngredient:
			step.AddIngredient(NewIngredient(item.val))
			lineBreak = false
		case itemCookware:
			step.AddCookware(NewCookware(item.val))
			lineBreak = false
		case itemTimer:
			step.AddTimer(NewTimer(item.val))
			lineBreak = false
		case itemStep:
			if step.Zero() {
				// a blank line between steps, or before any content
				continue
			}
			if lineBreak {
				flush()
				continue
			}
			lineBreak = true
		}
	}
	// The last line of a file need not end with a newline, so the trailing
	// step has to be flushed here as well.
	flush()

	return recipe, nil
}
