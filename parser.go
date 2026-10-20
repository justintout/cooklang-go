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

// ParseFile parses the file at the given path as a Cooklang recipe of the
// latest spec version
func ParseFile(path string) (Recipe, error) {
	return ParseFileSpec(path, latest)
}

// MustParseFileSpec calls ParseFileSpec and panics on error
func MustParseFileSpec(path string, spec Spec) Recipe {
	r, err := ParseFileSpec(path, spec)
	if err != nil {
		panic(err)
	}
	return r
}

// ParseFileSpec parses the file at the given path as a recipe of the given
// spec version
func ParseFileSpec(path string, spec Spec) (Recipe, error) {
	c, err := os.ReadFile(path)
	if err != nil {
		return Recipe{}, fmt.Errorf("failed to parse %q: %v", path, err)
	}
	return parse(string(c), spec)
}

// MustParse calls Parse and panics on error
func MustParse(input string) Recipe {
	r, err := Parse(input)
	if err != nil {
		panic(err)
	}
	return r
}

// Parse parses the input string as a Cooklang recipe of the latest spec
// version
func Parse(input string) (Recipe, error) {
	return parse(input, latest)
}

// MustParseSpec calls ParseSpec and panics on error
func MustParseSpec(input string, spec Spec) Recipe {
	r, err := ParseSpec(input, spec)
	if err != nil {
		panic(err)
	}
	return r
}

// ParseSpec parses the input string as a recipe of the given spec version
func ParseSpec(input string, spec Spec) (Recipe, error) {
	return parse(input, spec)
}

func parse(input string, spec Spec) (Recipe, error) {
	if !spec.known() {
		return Recipe{}, fmt.Errorf("unknown spec version %d", int(spec))
	}
	_, items := lex(input, spec)
	recipe := NewRecipe("recipe")
	step := &Step{}
	var note *Note
	// A step is a paragraph: a blank line ends it, while a single line break
	// inside it renders as a space and keeps the step going. Before version 7
	// every line was its own step, so there a newline ends the step instead.
	lineSteps := spec == SpecV5
	lineBreak := false
	forcedBreak := false

	flush := func() {
		if !step.Zero() {
			recipe.AddStep(step)
		}
		if note != nil {
			recipe.AddNote(note)
			note = nil
		}
		step = &Step{}
		lineBreak = false
		forcedBreak = false
	}

	for item := range items {
		switch item.typ {
		case itemMetadata:
			recipe.Metadata.Add(item.val)
		case itemComment:
		case itemSection:
			flush()
			recipe.AddSection(NewSection(item.val))
		case itemNote:
			if !step.Zero() {
				flush()
			}
			if note == nil {
				note = NewNote(item.val)
			} else if forcedBreak {
				note.Value += "\n" + item.val
			} else if lineBreak {
				note.Value += " " + item.val
			} else {
				note.Value += item.val
			}
			lineBreak = false
			forcedBreak = false
		case itemLineBreak:
			if forcedBreak {
				if note != nil {
					note.Value += "\n"
				} else if !step.Zero() {
					step.AddJoinedText(NewText("\n"))
				}
			}
			forcedBreak = true
			lineBreak = false
		case itemText:
			if note != nil {
				if forcedBreak {
					note.Value += "\n" + item.val
				} else if lineBreak {
					note.Value += " " + item.val
				} else {
					note.Value += item.val
				}
				lineBreak = false
				forcedBreak = false
				continue
			}
			text := NewText(item.val)
			if forcedBreak {
				text.Value = "\n" + text.Value
			} else if lineBreak {
				text.Value = " " + text.Value
			}
			step.AddJoinedText(text)
			lineBreak = false
			forcedBreak = false
		case itemIngredient:
			if note != nil {
				flush()
			}
			if forcedBreak {
				step.AddJoinedText(NewText("\n"))
				forcedBreak = false
			}
			step.AddIngredient(NewIngredient(item.val))
			lineBreak = false
		case itemCookware:
			if note != nil {
				flush()
			}
			if forcedBreak {
				step.AddJoinedText(NewText("\n"))
				forcedBreak = false
			}
			step.AddCookware(NewCookware(item.val))
			lineBreak = false
		case itemTimer:
			if note != nil {
				flush()
			}
			if forcedBreak {
				step.AddJoinedText(NewText("\n"))
				forcedBreak = false
			}
			step.AddTimer(NewTimer(item.val))
			lineBreak = false
		case itemStep:
			if step.Zero() && note == nil {
				// a blank line between steps, or before any content
				continue
			}
			if lineBreak || forcedBreak || lineSteps {
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
